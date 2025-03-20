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

var _ resource.ResourceWithImportState = &processChildStepResource{}

type processChildStepResource struct {
	*Config
}

func NewProcessChildStepResource() resource.Resource {
	return &processChildStepResource{}
}

func (r *processChildStepResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessChildStepResourceName)
}

func (r *processChildStepResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessChildStepSchema{}.GetResourceSchema()
}

func (r *processChildStepResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processChildStepResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	identifiers := strings.Split(request.ID, ":")

	if len(identifiers) != 3 {
		response.Diagnostics.AddError(
			"Incorrect Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ProcessId:ParentStepId:ChildStepId (e.g. deploymentprocess-Projects-123:00000000-0000-0000-0000-000000000010:00000000-0000-0000-0000-000000000012). Got: %q", request.ID),
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("process_id"), identifiers[0])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("parent_id"), identifiers[1])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), identifiers[2])...)
}

func (r *processChildStepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessChildStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("creating process child step: %s", data.Name.ValueString()))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, ok := process.FindStepByID(parentId)
	if !ok {
		resp.Diagnostics.AddError("Error creating process child step", fmt.Sprintf("unable to find a parent step with id '%s'", parentId))
		return
	}

	action := deployments.NewDeploymentAction(data.Name.ValueString(), data.Type.ValueString())
	mapDiagnostics := mapProcessChildStepActionFromState(ctx, data, action)
	resp.Diagnostics.Append(mapDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	parent.Actions = append(parent.Actions, action)

	updatedProcess, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create process child step", err.Error())
		return
	}

	updatedStep, parentFound := updatedProcess.FindStepByID(parentId)
	if !parentFound {
		resp.Diagnostics.AddError("unable to create process child step, unable to find a parent step '%s'", parent.ID)
		return
	}

	createdAction, ok := findActionFromProcessStepByName(updatedStep, action.Name)
	if !ok {
		resp.Diagnostics.AddError("unable to create process child step", action.Name)
		return
	}

	mapDiagnostics = mapProcessChildStepActionToState(updatedProcess, updatedStep, createdAction, data)
	resp.Diagnostics.Append(mapDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process child step created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processChildStepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessChildStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading process child step (%s)", data.ID))

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()
	actionId := data.ID.ValueString()

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
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

	mapDiagnostics := mapProcessChildStepActionToState(process, parent, action, data)
	resp.Diagnostics.Append(mapDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process chidl step read (%s)", actionId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processChildStepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessChildStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()
	actionId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("updating process child step (%s)", actionId))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, parentFound := process.FindStepByID(parentId)
	if !parentFound {
		resp.Diagnostics.AddError("unable to find parent step '%s'", parentId)
		return
	}

	action, actionFound := findActionFromProcessStepByID(parent, actionId)
	if !actionFound {
		resp.Diagnostics.AddError("unable to find process child step", actionId)
		return
	}

	mapDiagnostics := mapProcessChildStepActionFromState(ctx, data, action)
	resp.Diagnostics.Append(mapDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProcess, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("unable to update process child step", err.Error())
		return
	}

	updatedStep, updatedParentFound := updatedProcess.FindStepByID(parentId)
	if !updatedParentFound {
		resp.Diagnostics.AddError("unable to update process child step, unable to find a parent step '%s'", parent.ID)
		return
	}

	updatedAction, updatedActionFound := findActionFromProcessStepByID(updatedStep, actionId)
	if !updatedActionFound {
		resp.Diagnostics.AddError("unable to update process child step", actionId)
		return
	}

	mapDiagnostics = mapProcessChildStepActionToState(updatedProcess, updatedStep, updatedAction, data)
	resp.Diagnostics.Append(mapDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process child step updated (%s)", actionId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processChildStepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessChildStepResourceModel
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

	tflog.Info(ctx, fmt.Sprintf("deleting process child step (%s)", data.ID))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, ok := process.FindStepByID(parentId)
	if !ok {
		resp.Diagnostics.AddError("unable to find parent step '%s'", parentId)
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
		resp.Diagnostics.AddError("unable to delete process child step", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func mapProcessChildStepActionFromState(ctx context.Context, state *schemas.ProcessChildStepResourceModel, action *deployments.DeploymentAction) diag.Diagnostics {
	action.Name = state.Name.ValueString()
	action.Slug = state.Slug.ValueString()
	action.ActionType = state.Type.ValueString()
	action.Condition = state.Condition.ValueString()

	action.IsRequired = state.IsRequired.ValueBool()
	action.IsDisabled = state.IsDisabled.ValueBool()
	action.Notes = state.Notes.ValueString()
	action.WorkerPool = state.WorkerPoolID.ValueString()
	action.WorkerPoolVariable = state.WorkerPoolVariable.ValueString()
	if state.Container == nil {
		action.Container = nil
	} else {
		action.Container = deployments.NewDeploymentActionContainer(state.Container.FeedID.ValueStringPointer(), state.Container.Image.ValueStringPointer())
	}

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

	var packageReferences = make([]*packages.PackageReference, 0)
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

func mapProcessChildStepActionToState(process processWrapper, step *deployments.DeploymentStep, action *deployments.DeploymentAction, state *schemas.ProcessChildStepResourceModel) diag.Diagnostics {
	state.ID = types.StringValue(action.GetID())
	state.SpaceID = types.StringValue(process.GetSpaceID())
	state.ProcessID = types.StringValue(process.GetID())
	state.ParentID = types.StringValue(step.GetID())
	state.Name = types.StringValue(action.Name)

	state.Type = types.StringValue(action.ActionType)
	state.Slug = types.StringValue(action.Slug)
	state.IsRequired = types.BoolValue(action.IsRequired)
	state.IsDisabled = types.BoolValue(action.IsDisabled)
	state.Condition = types.StringValue(action.Condition)
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

	// Git Dependencies
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

	// Packages
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

	// Execution Properties
	stateProperties, diags := util.ConvertPropertiesToAttributeValuesMap(action.Properties)
	if diags.HasError() {
		return diags
	}

	state.ExecutionProperties = stateProperties

	return diag.Diagnostics{}
}

func findActionFromProcessStepByID(step *deployments.DeploymentStep, actionId string) (*deployments.DeploymentAction, bool) {
	for _, action := range step.Actions {
		if action.ID == actionId {
			return action, true
		}
	}
	return nil, false
}

func findActionFromProcessStepByName(step *deployments.DeploymentStep, name string) (*deployments.DeploymentAction, bool) {
	for _, action := range step.Actions {
		if action.Name == name {
			return action, true
		}
	}
	return nil, false
}
