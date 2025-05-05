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
	_ resource.ResourceWithImportState = &processTemplatedStepResource{}
	_ resource.ResourceWithModifyPlan  = &processTemplatedStepResource{}
)

type processTemplatedStepResource struct {
	*Config
}

func NewProcessTemplatedStepResource() resource.Resource {
	return &processTemplatedStepResource{}
}

func (r *processTemplatedStepResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessTemplatedStepResourceName)
}

func (r *processTemplatedStepResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessTemplatedStepSchema{}.GetResourceSchema()
}

func (r *processTemplatedStepResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processTemplatedStepResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
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

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
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

	version, err := strconv.ParseInt(templateVersion.Value, 10, 32)
	if err != nil {
		response.Diagnostics.AddError("Unable to import process step", "Process step's template version is invalid")
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), stepId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("process_id"), processId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("template_id"), templateId.Value)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("template_version"), version)...)
}

func (r *processTemplatedStepResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	if request.Plan.Raw.IsNull() {
		return // When deleting
	}

	if request.State.Raw.IsNull() {
		return // When creating
	}

	var plan *schemas.ProcessTemplatedStepResourceModel
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

func (r *processTemplatedStepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessTemplatedStepResourceModel
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

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	step := deployments.NewDeploymentStep(data.Name.ValueString())

	fromStateDiagnostics := mapProcessTemplatedStepFromState(ctx, data, template, step)
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

	toStateDiagnostics := mapProcessTemplatedStepToState(ctx, updatedProcess, createdStep, template, data)
	resp.Diagnostics.Append(toStateDiagnostics...)
	if toStateDiagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process step with template created (step: %s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processTemplatedStepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessTemplatedStepResourceModel
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

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
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

	diags = mapProcessTemplatedStepToState(ctx, process, step, template, data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process step with template read (step: %s)", step.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processTemplatedStepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessTemplatedStepResourceModel
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

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	step, exists := process.FindStepByID(stepId)
	if !exists {
		resp.Diagnostics.AddError("unable to find process step '%s'", stepId)
		return
	}

	diagnostics := mapProcessTemplatedStepFromState(ctx, data, template, step)
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

	diags = mapProcessTemplatedStepToState(ctx, updatedProcess, updatedStep, template, data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process step with template updated (%s)", updatedStep.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processTemplatedStepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessTemplatedStepResourceModel
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

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
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

func mapProcessTemplatedStepFromState(ctx context.Context, state *schemas.ProcessTemplatedStepResourceModel, template *actiontemplates.ActionTemplate, step *deployments.DeploymentStep) diag.Diagnostics {
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

	return mapProcessTemplatedStepEmbeddedActionFromState(ctx, state, template, step)
}

func mapProcessTemplatedStepEmbeddedActionFromState(ctx context.Context, state *schemas.ProcessTemplatedStepResourceModel, template *actiontemplates.ActionTemplate, step *deployments.DeploymentStep) diag.Diagnostics {
	actionType := template.ActionType
	name := state.Name.ValueString()

	if step.Actions == nil || len(step.Actions) == 0 {
		newAction := deployments.NewDeploymentAction(name, actionType)
		step.Actions = []*deployments.DeploymentAction{newAction}
	}

	if step.Actions[0] == nil {
		step.Actions[0] = deployments.NewDeploymentAction(name, actionType)
	}

	return mapProcessTemplatedStepActionFromState(ctx, state, template, step.Actions[0])
}

func mapProcessTemplatedStepActionFromState(ctx context.Context, state *schemas.ProcessTemplatedStepResourceModel, template *actiontemplates.ActionTemplate, action *deployments.DeploymentAction) diag.Diagnostics {
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

	action.GitDependencies = copyActionTemplateGitDependencies(template)
	action.Packages = copyActionTemplatePackages(template)
	action.Properties, diags = mapTemplatedActionPropertiesFromState(ctx, template, state.Parameters, state.ExecutionProperties)
	if diags.HasError() {
		return diags
	}

	return diag.Diagnostics{}
}

func mapProcessTemplatedStepToState(ctx context.Context, process processWrapper, step *deployments.DeploymentStep, template *actiontemplates.ActionTemplate, state *schemas.ProcessTemplatedStepResourceModel) diag.Diagnostics {
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
		return mapProcessTemplatedStepActionToState(ctx, step.Actions[0], template, state)
	}

	return diag.Diagnostics{}
}

func mapProcessTemplatedStepActionToState(ctx context.Context, action *deployments.DeploymentAction, template *actiontemplates.ActionTemplate, state *schemas.ProcessTemplatedStepResourceModel) diag.Diagnostics {
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

func copyActionTemplatePackages(template *actiontemplates.ActionTemplate) []*packages.PackageReference {
	packageReferences := make([]*packages.PackageReference, len(template.Packages))
	for index := range template.Packages {
		packageReferences[index] = &template.Packages[index]
	}
	return packageReferences
}

func copyActionTemplateGitDependencies(template *actiontemplates.ActionTemplate) []*gitdependencies.GitDependency {
	gitDependencies := make([]*gitdependencies.GitDependency, len(template.GitDependencies))
	for index := range template.GitDependencies {
		gitDependencies[index] = &template.GitDependencies[index]
	}
	return gitDependencies
}

func mapTemplatedActionPropertiesFromState(ctx context.Context, template *actiontemplates.ActionTemplate, parameters types.Map, executionProperties types.Map) (map[string]core.PropertyValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	properties := make(map[string]core.PropertyValue)

	// Parameters
	var stateParameters map[string]types.String
	stateParameters, diags = util.ConvertMapToStringMap(ctx, parameters)
	if diags.HasError() {
		return nil, diags
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
	diags = util.MergePropertyValues(ctx, properties, executionProperties)
	if diags.HasError() {
		return nil, diags
	}

	properties["Octopus.Action.Template.Id"] = core.NewPropertyValue(template.ID, false)
	templateVersion := strconv.Itoa(int(template.Version))
	properties["Octopus.Action.Template.Version"] = core.NewPropertyValue(templateVersion, false)

	return properties, diags
}

func mapTemplatedActionPropertiesToState(ctx context.Context, template *actiontemplates.ActionTemplate, action *deployments.DeploymentAction, parameters types.Map, executionProperties types.Map) (*schemas.ProcessTemplatedStepGroupedPropertyValues, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	state := &schemas.ProcessTemplatedStepGroupedPropertyValues{}

	// Split properties into groups defined by the templated step schema
	parameterValues := make(map[string]attr.Value)
	unmanagedParameterValues := make(map[string]attr.Value)
	templatePropertyValues := make(map[string]attr.Value)
	executionPropertyValues := make(map[string]attr.Value)

	parametersLookup := make(map[string]actiontemplates.ActionTemplateParameter)
	for _, parameter := range template.Parameters {
		parametersLookup[parameter.Name] = parameter
	}

	var stateParameters map[string]types.String
	stateParameters, diags = util.ConvertMapToStringMap(ctx, parameters)
	if diags.HasError() {
		return nil, diags
	}

	var stateExecutionProperties map[string]types.String
	stateExecutionProperties, diags = util.ConvertMapToStringMap(ctx, executionProperties)
	if diags.HasError() {
		return nil, diags
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
			version, parsingError := strconv.ParseInt(property.Value, 10, 32)
			if parsingError != nil {
				diags.AddError("Cannot map template version into the configuration", fmt.Sprintf("Value '%s' is not valid version number(int)", property.Value))
				return nil, diags
			}
			state.TemplateVersion = types.Int32Value(int32(version))
			continue
		}

		executionPropertyValues[key] = value
	}

	state.Parameters, diags = types.MapValue(types.StringType, parameterValues)
	if diags.HasError() {
		return nil, diags
	}

	state.UnmanagedParameters, diags = types.MapValue(types.StringType, unmanagedParameterValues)
	if diags.HasError() {
		return nil, diags
	}

	state.TemplateProperties, diags = types.MapValue(types.StringType, templatePropertyValues)
	if diags.HasError() {
		return nil, diags
	}

	state.ExecutionProperties, diags = types.MapValue(types.StringType, executionPropertyValues)
	if diags.HasError() {
		return nil, diags
	}

	return state, diags
}

func loadActionTemplate(client *client.Client, spaceId string, id string, version int32) (*actiontemplates.ActionTemplate, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	versioned, err := actiontemplates.GetVersionByID(client, spaceId, id, version)
	if err != nil {
		diags.AddError("Unable to load template", err.Error())
		return nil, diags
	}

	// When version is not latest, id of returned template is the id of versioned template record
	// We want to keep original template id
	versioned.SetID(id)
	return versioned, diags
}
