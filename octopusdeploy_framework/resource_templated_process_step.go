package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
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
	_ resource.ResourceWithImportState = &templatedProcessStepResource{}
	_ resource.ResourceWithModifyPlan  = &templatedProcessStepResource{}
)

type templatedProcessStepResource struct {
	*Config
}

func NewProcessStepWithTemplateResource() resource.Resource {
	return &templatedProcessStepResource{}
}

func (r *templatedProcessStepResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.TemplatedProcessStepResourceName)
}

func (r *templatedProcessStepResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.TemplatedProcessStepSchema{}.GetResourceSchema()
}

func (r *templatedProcessStepResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *templatedProcessStepResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	identifiers := strings.Split(request.ID, ":")

	if len(identifiers) != 2 {
		response.Diagnostics.AddError(
			"Incorrect Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ProcessId:StepId (e.g. deploymentprocess-Projects-123:00000000-0000-0000-0000-000000000001). Got: %q", request.ID),
		)
		return
	}

	spaceId := "" // Client's space is used for imported resources
	processId := identifiers[0]
	stepId := identifiers[1]
	tflog.Info(ctx, fmt.Sprintf("importing templated process step (%s) from process (%s)", stepId, processId))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		response.Diagnostics.Append(diags...)
		return
	}

	step, exists := process.FindStepByID(stepId)
	if !exists {
		response.Diagnostics.AddError("Unable to import process step", fmt.Sprintf("Process step (%s) is not found in process (%s)", stepId, processId))
		return
	}

	if len(step.Actions) == 0 || step.Actions[0] == nil {
		response.Diagnostics.AddError("Unable to import process step", "Process step doesn't have actions")
		return
	}

	action := step.Actions[0]
	templateId, hasTemplateId := action.Properties["Octopus.Action.Template.Id"]
	if !hasTemplateId {
		response.Diagnostics.AddError("Unable to import process step", "Process step doesn't have template id")
		return
	}

	templateVersion, hasTemplateVersion := action.Properties["Octopus.Action.Template.Version"]
	if !hasTemplateVersion {
		response.Diagnostics.AddError("Unable to import process step", "Process step doesn't have template version")
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), stepId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("process_id"), processId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("template_id"), templateId.Value)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("template_version"), templateVersion.Value)...)
}

func (r *templatedProcessStepResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	if request.Plan.Raw.IsNull() {
		return
	}

	var plan *schemas.TemplatedProcessStepResourceModel
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

func (r *templatedProcessStepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.TemplatedProcessStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
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

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	step := deployments.NewDeploymentStep(data.Name.ValueString())

	fromStateDiagnostics := mapTemplatedProcessStepFromState(ctx, data, template, step)
	resp.Diagnostics.Append(fromStateDiagnostics...)
	if fromStateDiagnostics.HasError() {
		return
	}

	process.AppendStep(step)

	updatedProcess, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create process step", err.Error())
		return
	}

	createdStep, exists := updatedProcess.FindStepByName(step.Name)
	if !exists {
		resp.Diagnostics.AddError("Unable to create process step '%s'", step.Name)
		return
	}

	toStateDiagnostics := mapTemplatedProcessStepToState(ctx, updatedProcess, createdStep, template, data)
	resp.Diagnostics.Append(toStateDiagnostics...)
	if toStateDiagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process step with template created (step: %s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *templatedProcessStepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.TemplatedProcessStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()
	templateId := data.TemplateID.ValueString()
	templateVersion := data.TemplateVersion.ValueInt32()

	template, templateDiags := loadActionTemplate(r.Config.Client, spaceId, templateId, templateVersion)
	if templateDiags.HasError() {
		resp.Diagnostics.Append(templateDiags...)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading process step with template (%s)", data.ID))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	step, exists := process.FindStepByID(stepId)
	if !exists {
		// Remove from state when not found in the process, so terraform will try to recreate it
		tflog.Info(ctx, fmt.Sprintf("process step read (id: %s), but not found, removing from state ...", stepId))
		resp.State.RemoveResource(ctx)
		return
	}

	diags = mapTemplatedProcessStepToState(ctx, process, step, template, data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process step with template read (step: %s)", step.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *templatedProcessStepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.TemplatedProcessStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()
	templateId := data.TemplateID.ValueString()
	templateVersion := data.TemplateVersion.ValueInt32()

	template, templateDiags := loadActionTemplate(r.Config.Client, spaceId, templateId, templateVersion)
	if templateDiags.HasError() {
		resp.Diagnostics.Append(templateDiags...)
		return
	}

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("updating process step with template (step: %s)", stepId))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	step, exists := process.FindStepByID(stepId)
	if !exists {
		resp.Diagnostics.AddError("unable to find process step '%s'", stepId)
		return
	}

	diagnostics := mapTemplatedProcessStepFromState(ctx, data, template, step)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	updatedProcess, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("unable to update process step", err.Error())
		return
	}

	updatedStep, updatedExists := updatedProcess.FindStepByID(stepId)
	if !updatedExists {
		resp.Diagnostics.AddError("unable to find updated process step '%s'", stepId)
		return
	}

	diags = mapTemplatedProcessStepToState(ctx, updatedProcess, updatedStep, template, data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process step with template updated (%s)", updatedStep.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *templatedProcessStepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.TemplatedProcessStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("deleting process step (%s)", stepId))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	process.RemoveStep(stepId)

	_, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("unable to delete process step", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func mapTemplatedProcessStepFromState(ctx context.Context, state *schemas.TemplatedProcessStepResourceModel, template *actiontemplates.ActionTemplate, step *deployments.DeploymentStep) diag.Diagnostics {
	step.Name = state.Name.ValueString()
	step.StartTrigger = deployments.DeploymentStepStartTrigger(state.StartTrigger.ValueString())
	step.PackageRequirement = deployments.DeploymentStepPackageRequirement(state.PackageRequirement.ValueString())
	step.Condition = deployments.DeploymentStepConditionType(state.Condition.ValueString())

	if state.Properties.IsNull() {
		step.Properties = make(map[string]core.PropertyValue)
	} else {
		stateProperties := make(map[string]types.String, len(state.Properties.Elements()))
		diags := state.Properties.ElementsAs(ctx, &stateProperties, false)
		if diags.HasError() {
			return diags
		}

		properties := make(map[string]core.PropertyValue, len(stateProperties))
		for key, value := range stateProperties {
			if value.IsNull() {
				properties[key] = core.NewPropertyValue("", false)
			} else {
				properties[key] = core.NewPropertyValue(value.ValueString(), false)
			}
		}

		step.Properties = properties
	}

	return mapTemplatedProcessStepEmbeddedActionFromState(ctx, state, template, step)
}

func mapTemplatedProcessStepEmbeddedActionFromState(ctx context.Context, state *schemas.TemplatedProcessStepResourceModel, template *actiontemplates.ActionTemplate, step *deployments.DeploymentStep) diag.Diagnostics {
	actionType := template.ActionType
	name := state.Name.ValueString()

	if step.Actions == nil || len(step.Actions) == 0 {
		newAction := deployments.NewDeploymentAction(name, actionType)
		step.Actions = []*deployments.DeploymentAction{newAction}
	}

	if step.Actions[0] == nil {
		step.Actions[0] = deployments.NewDeploymentAction(name, actionType)
	}

	return mapTemplatedProcessStepActionFromState(ctx, state, template, step.Actions[0])
}

func mapTemplatedProcessStepActionFromState(ctx context.Context, state *schemas.TemplatedProcessStepResourceModel, template *actiontemplates.ActionTemplate, action *deployments.DeploymentAction) diag.Diagnostics {
	action.Name = state.Name.ValueString()
	action.Slug = state.Slug.ValueString() // update only embedded action slug(step slug remains original), same as UI behaviour
	action.ActionType = template.ActionType
	// action.Condition is not updated: replicates UI behaviour where condition of the first action of step is always a default value (Success)

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

	// Git Dependencies
	// Should be copied from the template, but ActionTemplate type missing GitDependencies attribute
	action.GitDependencies = make([]*gitdependencies.GitDependency, 0)

	// Packages
	// Always copied from the template
	packageReferences := make([]*packages.PackageReference, len(template.Packages))
	for index := range template.Packages {
		packageReferences[index] = &template.Packages[index]
	}
	action.Packages = packageReferences

	// Properties
	properties := make(map[string]core.PropertyValue)

	// Parameters
	stateParameters, paramDiags := util.ConvertMapToStringMap(ctx, state.Parameters)
	if paramDiags.HasError() {
		return paramDiags
	}

	for _, parameter := range template.Parameters {
		value, ok := stateParameters[parameter.Name]
		if ok {
			properties[parameter.Name] = util.ConvertToPropertyValue(value, false)
			continue
		}

		if defaultValue := parameter.DefaultValue; defaultValue != nil {
			if defaultValue.Value != "" {
				properties[parameter.Name] = *defaultValue
				continue
			}

			if defaultValue.IsSensitive && defaultValue.SensitiveValue.HasValue {
				properties[parameter.Name] = *defaultValue
				continue
			}
		}
	}

	// Template properties
	for key, value := range template.Properties {
		properties[key] = value
	}

	// Rest of the properties
	diags.Append(util.MergePropertyValues(ctx, properties, state.ExecutionProperties)...)
	if diags.HasError() {
		return diags
	}

	properties["Octopus.Action.Template.Id"] = core.NewPropertyValue(template.ID, false)
	templateVersion := strconv.Itoa(int(state.TemplateVersion.ValueInt32()))
	properties["Octopus.Action.Template.Version"] = core.NewPropertyValue(templateVersion, false)

	action.Properties = properties

	return diag.Diagnostics{}
}

func mapTemplatedProcessStepToState(ctx context.Context, process processWrapper, step *deployments.DeploymentStep, template *actiontemplates.ActionTemplate, state *schemas.TemplatedProcessStepResourceModel) diag.Diagnostics {
	state.ID = types.StringValue(step.GetID())
	state.SpaceID = types.StringValue(process.GetSpaceID())
	state.ProcessID = types.StringValue(process.GetID())
	state.Name = types.StringValue(step.Name)
	state.StartTrigger = types.StringValue(string(step.StartTrigger))
	state.PackageRequirement = types.StringValue(string(step.PackageRequirement))
	state.Condition = types.StringValue(string(step.Condition))

	stepProperties := make(map[string]attr.Value, len(step.Properties))
	for key, value := range step.Properties {
		stepProperties[key] = types.StringValue(value.Value)
	}

	stateProperties, diags := types.MapValue(types.StringType, stepProperties)
	if diags.HasError() {
		return diags
	}

	state.Properties = stateProperties

	if len(step.Actions) > 0 && step.Actions[0] != nil {
		return mapTemplatedProcessStepActionToState(ctx, step.Actions[0], template, state)
	}

	return diag.Diagnostics{}
}

func mapTemplatedProcessStepActionToState(ctx context.Context, action *deployments.DeploymentAction, template *actiontemplates.ActionTemplate, state *schemas.TemplatedProcessStepResourceModel) diag.Diagnostics {
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

	diags := diag.Diagnostics{}
	// Split properties into 3 groups (parameters, template properties and rest of the provided properties)
	parameterValues := make(map[string]attr.Value)
	unmanagedParameterValues := make(map[string]attr.Value)
	templatePropertyValues := make(map[string]attr.Value)
	executionPropertyValues := make(map[string]attr.Value)

	parametersLookup := make(map[string]actiontemplates.ActionTemplateParameter)
	for _, parameter := range template.Parameters {
		parametersLookup[parameter.Name] = parameter
	}

	var stateParameters map[string]types.String
	stateParameters, diags = util.ConvertMapToStringMap(ctx, state.Parameters)
	if diags.HasError() {
		return diags
	}

	var stateExecutionProperties map[string]types.String
	stateExecutionProperties, diags = util.ConvertMapToStringMap(ctx, state.ExecutionProperties)
	if diags.HasError() {
		return diags
	}

	for key, property := range action.Properties {
		value := types.StringValue(property.Value)

		if _, isParameter := parametersLookup[key]; isParameter {
			if _, configured := stateParameters[key]; configured {
				parameterValues[key] = value
			} else {
				unmanagedParameterValues[key] = value
			}

			continue
		}

		if _, isTemplateProperty := template.Properties[key]; isTemplateProperty {
			if _, overridden := stateExecutionProperties[key]; !overridden {
				templatePropertyValues[key] = value
				continue
			}
		}

		if key == "Octopus.Action.Template.Id" {
			state.TemplateID = value
			continue
		}

		if key == "Octopus.Action.Template.Version" {
			version, _ := strconv.Atoi(property.Value)
			state.TemplateVersion = types.Int32Value(int32(version))
			continue
		}

		executionPropertyValues[key] = value
	}

	state.Parameters, diags = types.MapValue(types.StringType, parameterValues)
	if diags.HasError() {
		return diags
	}

	state.UnmanagedParameters, diags = types.MapValue(types.StringType, unmanagedParameterValues)
	if diags.HasError() {
		return diags
	}

	state.TemplateProperties, diags = types.MapValue(types.StringType, templatePropertyValues)
	if diags.HasError() {
		return diags
	}

	state.ExecutionProperties, diags = types.MapValue(types.StringType, executionPropertyValues)
	if diags.HasError() {
		return diags
	}

	return diags
}

func loadActionTemplate(client *client.Client, spaceId string, id string, version int32) (*actiontemplates.ActionTemplate, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	versioned, err := actiontemplates.GetVersionByID(client, spaceId, id, version)
	if err != nil {
		diags.AddError("Unable to load template", err.Error())
		return nil, diags
	}

	// Versioned endpoint returns template with id of the version record which is different from id of actual template,
	// but we want to keep original template id
	versioned.SetID(id)
	return versioned, diags
}
