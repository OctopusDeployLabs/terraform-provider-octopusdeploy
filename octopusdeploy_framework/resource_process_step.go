package octopusdeploy_framework

import (
	"context"
	"fmt"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

var _ resource.ResourceWithImportState = &processStepResource{}

type processStepResource struct {
	*Config
}

func NewProcessStepResource() resource.Resource {
	return &processStepResource{}
}

func (r *processStepResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessStepResourceName)
}

func (r *processStepResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessStepSchema{}.GetResourceSchema()
}

func (r *processStepResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processStepResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	identifiers := strings.Split(request.ID, ":")

	if len(identifiers) != 2 {
		response.Diagnostics.AddError(
			"Incorrect Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ProcessId:StepId (e.g. deploymentprocess-Projects-123:00000000-0000-0000-0000-000000000001). Got: %q", request.ID),
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("process_id"), identifiers[0])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), identifiers[1])...)
}

func (r *processStepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("creating process step: %s", data.Name.ValueString()))

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	step := deployments.NewDeploymentStep(data.Name.ValueString())

	fromStateDiagnostics := mapProcessStepFromState(ctx, data, step)
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

	toStateDiagnostics := mapProcessStepToState(updatedProcess, createdStep, data)
	resp.Diagnostics.Append(toStateDiagnostics...)
	if toStateDiagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process step created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("reading process step (%s)", data.ID))

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.HasError()

	step, exists := process.FindStepByID(stepId)
	if !exists {
		// Remove from state when not found in the process, so terraform will try to recreate it
		tflog.Info(ctx, fmt.Sprintf("process step read (id: %s), but not found, removing from state ...", stepId))
		resp.State.RemoveResource(ctx)
		return
	}

	mapProcessStepToState(process, step, data)

	tflog.Info(ctx, fmt.Sprintf("process step read (%s)", step.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("updating process step (%s)", stepId))

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

	diagnostics := mapProcessStepFromState(ctx, data, step)
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

	mapProcessStepToState(updatedProcess, updatedStep, data)

	tflog.Info(ctx, fmt.Sprintf("process step updated (%s)", updatedStep.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessStepResourceModel
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

func mapProcessStepFromState(ctx context.Context, state *schemas.ProcessStepResourceModel, step *deployments.DeploymentStep) diag.Diagnostics {
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

	return mapProcessStepEmbeddedActionFromState(ctx, state, step)
}

func mapProcessStepEmbeddedActionFromState(ctx context.Context, state *schemas.ProcessStepResourceModel, step *deployments.DeploymentStep) diag.Diagnostics {
	actionType := state.Type.ValueString()
	name := state.Name.ValueString()

	if step.Actions == nil || len(step.Actions) == 0 {
		newAction := deployments.NewDeploymentAction(name, actionType)
		step.Actions = []*deployments.DeploymentAction{newAction}
	}

	if step.Actions[0] == nil {
		step.Actions[0] = deployments.NewDeploymentAction(name, actionType)
	}

	return mapProcessStepActionFromState(ctx, state, step.Actions[0])
}

func mapProcessStepActionFromState(ctx context.Context, state *schemas.ProcessStepResourceModel, action *deployments.DeploymentAction) diag.Diagnostics {
	action.Name = state.Name.ValueString()
	action.Slug = state.Slug.ValueString() // update only embedded action slug(step slug remains original), same as UI behaviour
	action.ActionType = state.Type.ValueString()
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
	var dependenciesMap map[string]types.Object
	diags = state.GitDependencies.ElementsAs(ctx, &dependenciesMap, false)
	if diags.HasError() {
		return diags
	}

	var gitDependencies = make([]*gitdependencies.GitDependency, 0)
	for key, dependencyObject := range dependenciesMap {
		if dependencyObject.IsNull() {
			continue
		}

		var dependencyState schemas.ProcessStepGitDependencyResourceModel
		diags = dependencyObject.As(ctx, &dependencyState, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return diags
		}

		gitDependency := &gitdependencies.GitDependency{
			Name:              key,
			RepositoryUri:     dependencyState.RepositoryUri.ValueString(),
			DefaultBranch:     dependencyState.DefaultBranch.ValueString(),
			GitCredentialType: dependencyState.GitCredentialType.ValueString(),
			GitCredentialId:   dependencyState.GitCredentialID.ValueString(),
		}

		if dependencyState.FilePathFilters.IsNull() {
			gitDependency.FilePathFilters = nil
		} else {
			gitDependency.FilePathFilters, diags = util.SetToStringArray(ctx, dependencyState.FilePathFilters)
			if diags.HasError() {
				return diags
			}
		}

		gitDependencies = append(gitDependencies, gitDependency)
	}

	action.GitDependencies = gitDependencies

	// Packages
	var packagesMap map[string]types.Object
	diags = state.Packages.ElementsAs(ctx, &packagesMap, false)
	if diags.HasError() {
		return diags
	}

	packageReferences := make([]*packages.PackageReference, 0)
	for key, packageObject := range packagesMap {
		if packageObject.IsNull() {
			continue
		}

		var packageState schemas.ProcessStepPackageReferenceResourceModel
		diags = packageObject.As(ctx, &packageState, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return diags
		}

		stateProperties := make(map[string]types.String, len(packageState.Properties.Elements()))
		diags = packageState.Properties.ElementsAs(ctx, &stateProperties, false)
		if diags.HasError() {
			return diags
		}

		packageProperties := make(map[string]string, len(stateProperties))
		for propertyKey, value := range stateProperties {
			if value.IsNull() {
				packageProperties[propertyKey] = ""
			} else {
				packageProperties[propertyKey] = value.ValueString()
			}
		}

		packageReference := &packages.PackageReference{
			ID:                  packageState.GetID(),
			Name:                key,
			PackageID:           packageState.PackageID.ValueString(),
			FeedID:              packageState.FeedID.ValueString(),
			AcquisitionLocation: packageState.AcquisitionLocation.ValueString(),
			Properties:          packageProperties,
		}
		packageReferences = append(packageReferences, packageReference)
	}

	action.Packages = packageReferences

	// Execution Properties
	stateProperties := make(map[string]types.String, len(state.ExecutionProperties.Elements()))
	propertiesDiags := state.ExecutionProperties.ElementsAs(ctx, &stateProperties, false)
	if propertiesDiags.HasError() {
		return propertiesDiags
	}

	properties := make(map[string]core.PropertyValue, len(stateProperties))
	for key, value := range stateProperties {
		if value.IsNull() {
			properties[key] = core.NewPropertyValue("", false)
		} else {
			properties[key] = core.NewPropertyValue(value.ValueString(), false)
		}
	}

	action.Properties = properties

	return diag.Diagnostics{}
}

func mapProcessStepToState(process processWrapper, step *deployments.DeploymentStep, state *schemas.ProcessStepResourceModel) diag.Diagnostics {
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
		return mapProcessStepActionToState(step.Actions[0], state)
	}

	return diag.Diagnostics{}
}

func mapProcessStepActionToState(action *deployments.DeploymentAction, state *schemas.ProcessStepResourceModel) diag.Diagnostics {
	state.Type = types.StringValue(action.ActionType)
	state.Slug = types.StringValue(action.Slug)
	state.IsRequired = types.BoolValue(action.IsRequired)
	state.IsDisabled = types.BoolValue(action.IsDisabled)
	state.Notes = types.StringValue(action.Notes)
	state.WorkerPoolID = types.StringValue(action.WorkerPool)
	state.WorkerPoolVariable = types.StringValue(action.WorkerPoolVariable)

	if action.Container == nil {
		state.Container = nil
	} else {
		state.Container = &schemas.ProcessStepActionContainerModel{
			FeedID: types.StringValue(action.Container.FeedID),
			Image:  types.StringValue(action.Container.Image),
		}
	}

	state.TenantTags = util.BuildStringSetOrEmpty(action.TenantTags)
	state.Environments = util.BuildStringSetOrEmpty(action.Environments)
	state.ExcludedEnvironments = util.BuildStringSetOrEmpty(action.ExcludedEnvironments)
	state.Channels = util.BuildStringSetOrEmpty(action.Channels)

	if action.GitDependencies == nil {
		state.GitDependencies = types.MapValueMust(schemas.ProcessStepGitDependencyObjectType(), map[string]attr.Value{})
	} else {
		stateDependencies := make(map[string]attr.Value, len(action.GitDependencies))
		for _, dependency := range action.GitDependencies {
			stateDependency := types.ObjectValueMust(
				schemas.ProcessStepGitDependencyAttributeTypes(),
				map[string]attr.Value{
					"repository_uri":      types.StringValue(dependency.RepositoryUri),
					"default_branch":      types.StringValue(dependency.DefaultBranch),
					"git_credential_type": types.StringValue(dependency.GitCredentialType),
					"file_path_filters":   types.SetValueMust(types.StringType, util.ToValueSlice(dependency.FilePathFilters)),
					"git_credential_id":   types.StringValue(dependency.GitCredentialId),
				},
			)

			stateDependencies[dependency.Name] = stateDependency
		}

		state.GitDependencies = types.MapValueMust(schemas.ProcessStepGitDependencyObjectType(), stateDependencies)
	}

	if action.Packages == nil {
		state.Packages = types.MapValueMust(schemas.ProcessStepPackageReferenceObjectType(), map[string]attr.Value{})
	} else {
		statePackages := make(map[string]attr.Value, len(action.Packages))
		for _, packageReference := range action.Packages {
			packageProperties := util.ConvertMapStringToMapAttrValue(packageReference.Properties)
			statePackage := types.ObjectValueMust(
				schemas.ProcessStepPackageReferenceAttributeTypes(),
				map[string]attr.Value{
					"id":                   types.StringValue(packageReference.ID),
					"package_id":           types.StringValue(packageReference.PackageID),
					"feed_id":              types.StringValue(packageReference.FeedID),
					"acquisition_location": types.StringValue(packageReference.AcquisitionLocation),
					"properties":           types.MapValueMust(types.StringType, packageProperties),
				},
			)

			statePackages[packageReference.Name] = statePackage
		}

		state.Packages = types.MapValueMust(schemas.ProcessStepPackageReferenceObjectType(), statePackages)
	}

	stateProperties, diags := util.ConvertPropertiesToAttributeValuesMap(action.Properties)
	if diags.HasError() {
		return diags
	}

	state.ExecutionProperties = stateProperties

	return diag.Diagnostics{}
}
