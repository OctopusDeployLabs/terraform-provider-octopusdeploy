package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"strings"
)

var (
	_ resource.ResourceWithImportState = &processTemplatedChildStepResource{}
	_ resource.ResourceWithModifyPlan  = &processTemplatedChildStepResource{}
)

type processTemplatedChildStepResource struct {
	*Config
}

func NewProcessTemplatedChildStepResource() resource.Resource {
	return &processTemplatedChildStepResource{}
}

func (r *processTemplatedChildStepResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessTemplatedChildStepResourceName)
}

func (r *processTemplatedChildStepResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessTemplatedChildStepSchema{}.GetResourceSchema()
}

func (r *processTemplatedChildStepResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processTemplatedChildStepResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	identifiers := strings.Split(request.ID, ":")

	if len(identifiers) != 3 {
		response.Diagnostics.AddError(
			"Incorrect Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ProcessId:ParentStepId:ChildStepId (e.g. deploymentprocess-Projects-123:00000000-0000-0000-0000-000000000010:00000000-0000-0000-0000-000000000012). Got: %q", request.ID),
		)
		return
	}

	spaceId := "" // Client's space is used for imported resources
	processId := identifiers[0]
	parentId := identifiers[1]
	actionId := identifiers[2]
	tflog.Info(ctx, fmt.Sprintf("importing templated process child step (%s) from parent (%s) and process (%s)", actionId, parentId, processId))

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		response.Diagnostics.Append(diags...)
		return
	}

	parent, stepExists := process.FindStepByID(parentId)
	if !stepExists {
		response.Diagnostics.AddError("Unable to import process child step", fmt.Sprintf("Parent step (%s) is not found in process (%s)", parentId, processId))
		return
	}

	action, actionExists := findActionFromProcessStepByID(parent, actionId)
	if !actionExists {
		response.Diagnostics.AddError("Unable to import process child step", fmt.Sprintf("Process child step (%s) is not found in parent (%s)", actionId, parentId))
		return
	}

	templateId, hasTemplateId := action.Properties["Octopus.Action.Template.Id"]
	if !hasTemplateId {
		response.Diagnostics.AddError("Unable to import process child step", "Process child step doesn't have template id")
		return
	}

	templateVersion, hasTemplateVersion := action.Properties["Octopus.Action.Template.Version"]
	if !hasTemplateVersion {
		response.Diagnostics.AddError("Unable to import process step", "Process child step doesn't have template version")
		return
	}

	version, err := strconv.ParseInt(templateVersion.Value, 10, 32)
	if err != nil {
		response.Diagnostics.AddError("Unable to import process step", "Process step's template version is invalid")
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), actionId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("process_id"), processId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("parent_id"), parentId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("template_id"), templateId.Value)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("template_version"), version)...)
}

func (r *processTemplatedChildStepResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	if request.Plan.Raw.IsNull() {
		return // When deleting
	}

	if request.State.Raw.IsNull() {
		return // When creating
	}

	var plan *schemas.ProcessTemplatedChildStepResourceModel
	diags := request.Plan.Get(ctx, &plan)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	spaceId := plan.SpaceID.ValueString()
	templateId := plan.TemplateID.ValueString()
	templateVersion := plan.TemplateVersion.ValueInt32()

	var template *actiontemplates.ActionTemplate
	template, diags = loadActionTemplate(r.Config.Client, spaceId, templateId, templateVersion)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	// Explicitly set computed attributes to avoid "state drift",
	// because terraform complains about differences between plan and state after apply

	// Set unmanaged parameters
	unmanagedParameters := make(map[string]attr.Value)
	var managedParameters map[string]types.String
	managedParameters, diags = util.ConvertMapToStringMap(ctx, plan.Parameters)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	for _, parameter := range template.Parameters {
		if _, configured := managedParameters[parameter.Name]; configured {
			continue
		}

		// Not configured - add to unmanaged only if default value is not empty
		if defaultValue := parameter.DefaultValue; defaultValue != nil {
			if defaultValue.Value != "" {
				unmanagedParameters[parameter.Name] = types.StringValue(defaultValue.Value)
				continue
			}

			if defaultValue.IsSensitive && defaultValue.SensitiveValue.HasValue {
				unmanagedParameters[parameter.Name] = types.StringValue(defaultValue.Value)
				continue
			}
		}
	}

	plan.UnmanagedParameters, diags = types.MapValue(types.StringType, unmanagedParameters)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	// Set template properties
	templateProperties := make(map[string]attr.Value)
	var executionProperties map[string]types.String
	executionProperties, diags = util.ConvertMapToStringMap(ctx, plan.ExecutionProperties)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	for key, property := range template.Properties {
		if _, overridden := executionProperties[key]; overridden {
			continue
		}

		templateProperties[key] = types.StringValue(property.Value)
	}

	plan.TemplateProperties, diags = types.MapValue(types.StringType, templateProperties)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	response.Diagnostics.Append(response.Plan.Set(ctx, &plan)...)
}

func (r *processTemplatedChildStepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessTemplatedChildStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()
	templateId := data.TemplateID.ValueString()
	templateVersion := data.TemplateVersion.ValueInt32()

	template, templateDiags := loadActionTemplate(r.Config.Client, spaceId, templateId, templateVersion)
	if templateDiags.HasError() {
		resp.Diagnostics.Append(templateDiags...)
		return
	}

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("creating process step with template: %s", data.Name.ValueString()))

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, ok := process.FindStepByID(parentId)
	if !ok {
		resp.Diagnostics.AddError("Error creating templated process child step", fmt.Sprintf("unable to find a parent step with id '%s'", parentId))
		return
	}

	action := deployments.NewDeploymentAction(data.Name.ValueString(), template.ActionType)

	diags = mapProcessTemplatedChildStepActionFromState(ctx, data, template, action)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent.Actions = append(parent.Actions, action)

	updatedProcess, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create process child step", err.Error())
		return
	}

	updatedParent, exists := updatedProcess.FindStepByID(parentId)
	if !exists {
		resp.Diagnostics.AddError("unable to create process child step", fmt.Sprintf("unable to find a parent step '%s'", parent.ID))
		return
	}

	createdAction, ok := findActionFromProcessStepByName(updatedParent, action.Name)
	if !ok {
		resp.Diagnostics.AddError("unable to create process child step", action.Name)
		return
	}

	diags = mapProcessTemplatedChildStepActionToState(ctx, updatedProcess, updatedParent, createdAction, template, data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process step with template created (step: %s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processTemplatedChildStepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessTemplatedChildStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()
	actionId := data.ID.ValueString()
	templateId := data.TemplateID.ValueString()
	templateVersion := data.TemplateVersion.ValueInt32()

	template, templateDiags := loadActionTemplate(r.Config.Client, spaceId, templateId, templateVersion)
	if templateDiags.HasError() {
		resp.Diagnostics.Append(templateDiags...)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading process step with template (%s)", data.ID))

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, ok := process.FindStepByID(parentId)
	if !ok {
		resp.Diagnostics.AddError("unable to find parent step '%s'", parentId)
		return
	}

	action, exists := findActionFromProcessStepByID(parent, actionId)
	if !exists {
		// Remove from state when action is not found in the step, so terraform will try to recreate it
		tflog.Info(ctx, fmt.Sprintf("reading process child step (id: %s), but not found, removing from state ...", actionId))
		resp.State.RemoveResource(ctx)
		return
	}

	diags = mapProcessTemplatedChildStepActionToState(ctx, process, parent, action, template, data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process step with template read (step: %s)", action.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processTemplatedChildStepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessTemplatedChildStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()
	actionId := data.ID.ValueString()
	templateId := data.TemplateID.ValueString()
	templateVersion := data.TemplateVersion.ValueInt32()

	template, templateDiags := loadActionTemplate(r.Config.Client, spaceId, templateId, templateVersion)
	if templateDiags.HasError() {
		resp.Diagnostics.Append(templateDiags...)
		return
	}

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("updating process step with template (step: %s)", actionId))

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, ok := process.FindStepByID(parentId)
	if !ok {
		resp.Diagnostics.AddError("unable to find parent step '%s'", parentId)
		return
	}

	action, actionFound := findActionFromProcessStepByID(parent, actionId)
	if !actionFound {
		resp.Diagnostics.AddError("unable to find process child step", actionId)
		return
	}

	diags = mapProcessTemplatedChildStepActionFromState(ctx, data, template, action)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	updatedProcess, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("unable to update process step", err.Error())
		return
	}

	updatedParent, updatedParentFound := updatedProcess.FindStepByID(parentId)
	if !updatedParentFound {
		resp.Diagnostics.AddError("unable to update process child step", fmt.Sprintf("unable to find a parent step '%s'", parent.ID))
		return
	}

	updatedAction, updatedActionFound := findActionFromProcessStepByID(updatedParent, actionId)
	if !updatedActionFound {
		resp.Diagnostics.AddError("unable to update process child step", actionId)
		return
	}

	diags = mapProcessTemplatedChildStepActionToState(ctx, updatedProcess, updatedParent, updatedAction, template, data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("templated process step updated (%s)", updatedParent.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processTemplatedChildStepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessTemplatedChildStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()
	actionId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("deleting process step (%s)", actionId))

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, ok := process.FindStepByID(parentId)
	if !ok {
		resp.Diagnostics.AddError("Unable to delete process step", fmt.Sprintf("unable to find parent step '%s'", parentId))
		return
	}

	var filteredActions []*deployments.DeploymentAction
	for _, action := range parent.Actions {
		if actionId != action.GetID() {
			filteredActions = append(filteredActions, action)
		}
	}
	parent.Actions = filteredActions

	_, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("unable to delete process step", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func mapProcessTemplatedChildStepActionFromState(ctx context.Context, state *schemas.ProcessTemplatedChildStepResourceModel, template *actiontemplates.ActionTemplate, action *deployments.DeploymentAction) diag.Diagnostics {
	action.Name = state.Name.ValueString()
	action.Slug = state.Slug.ValueString()
	action.ActionType = template.ActionType
	action.Condition = state.Condition.ValueString()

	action.IsRequired = state.IsRequired.ValueBool()
	action.IsDisabled = state.IsDisabled.ValueBool()
	action.Notes = state.Notes.ValueString()
	action.WorkerPool = state.WorkerPoolID.ValueString()
	action.WorkerPoolVariable = state.WorkerPoolVariable.ValueString()
	action.Container = deployments.NewDeploymentActionContainer(state.Container.FeedID.ValueStringPointer(), state.Container.Image.ValueStringPointer())

	diags := diag.Diagnostics{}

	action.TenantTags, diags = util.SetToStringArray(ctx, state.TenantTags)
	if diags.HasError() {
		return diags
	}

	action.Environments, diags = util.SetToStringArray(ctx, state.Environments)
	if diags.HasError() {
		return diags
	}

	action.ExcludedEnvironments, diags = util.SetToStringArray(ctx, state.ExcludedEnvironments)
	if diags.HasError() {
		return diags
	}

	action.Channels, diags = util.SetToStringArray(ctx, state.Channels)
	if diags.HasError() {
		return diags
	}

	action.GitDependencies = copyActionTemplateGitDependencies(template)
	action.Packages = copyActionTemplatePackages(template)
	action.Properties, diags = mapTemplatedActionPropertiesFromState(ctx, template, state.Parameters, state.ExecutionProperties)
	if diags.HasError() {
		return diags
	}

	return diag.Diagnostics{}
}

func mapProcessTemplatedChildStepActionToState(ctx context.Context, process processWrapper, parent *deployments.DeploymentStep, action *deployments.DeploymentAction, template *actiontemplates.ActionTemplate, state *schemas.ProcessTemplatedChildStepResourceModel) diag.Diagnostics {
	state.ID = types.StringValue(action.GetID())
	state.SpaceID = types.StringValue(process.GetSpaceID())
	state.ProcessID = types.StringValue(process.GetID())
	state.ParentID = types.StringValue(parent.GetID())
	state.Type = types.StringValue(action.ActionType)
	state.Slug = types.StringValue(action.Slug)
	state.IsRequired = types.BoolValue(action.IsRequired)
	state.IsDisabled = types.BoolValue(action.IsDisabled)
	state.Notes = types.StringValue(action.Notes)
	state.WorkerPoolID = types.StringValue(action.WorkerPool)
	state.WorkerPoolVariable = types.StringValue(action.WorkerPoolVariable)

	state.Container = mapDeploymentActionContainerToState(action.Container)
	state.TenantTags = util.BuildStringSetOrEmpty(action.TenantTags)
	state.Environments = util.BuildStringSetOrEmpty(action.Environments)
	state.ExcludedEnvironments = util.BuildStringSetOrEmpty(action.ExcludedEnvironments)
	state.Channels = util.BuildStringSetOrEmpty(action.Channels)

	state.GitDependencies = mapGitDependenciesToState(action.GitDependencies)
	state.Packages = mapPackageReferencesToState(action.Packages)

	properties, diags := mapTemplatedActionPropertiesToState(ctx, template, action, state.Parameters, state.ExecutionProperties)
	if diags.HasError() {
		return diags
	}
	state.TemplateID = properties.TemplateID
	state.TemplateVersion = properties.TemplateVersion
	state.Parameters = properties.Parameters
	state.UnmanagedParameters = properties.UnmanagedParameters
	state.TemplateProperties = properties.TemplateProperties
	state.ExecutionProperties = properties.ExecutionProperties

	return diags
}
