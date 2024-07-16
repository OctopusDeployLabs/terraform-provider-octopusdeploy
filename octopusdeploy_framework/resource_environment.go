package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/extensions"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type environmentTypeResource struct {
	*Config
}

func NewEnvironmentResource() resource.Resource {
	return &environmentTypeResource{}
}

func (r *environmentTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ProviderTypeName + "_environment"
}

func (r *environmentTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetEnvironmentResourceSchema()
}

func (r *environmentTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *environmentTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.EnvironmentTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newEnvironment := environments.NewEnvironment(data.Name.ValueString())
	newEnvironment.SpaceID = data.SpaceID.ValueString()
	newEnvironment.Description = data.Description.ValueString()
	newEnvironment.AllowDynamicInfrastructure = data.AllowDynamicInfrastructure.ValueBool()
	newEnvironment.UseGuidedFailure = data.UseGuidedFailure.ValueBool()
	newEnvironment.SortOrder = int(data.SortOrder.ValueInt64())
	if len(data.JiraExtensionSettings.Elements()) > 0 {
		jiraExtensionSettings := mapJiraExtensionSettings(data.JiraExtensionSettings)
		if jiraExtensionSettings != nil {
			newEnvironment.ExtensionSettings = append(newEnvironment.ExtensionSettings, jiraExtensionSettings)
		}
	}
	if len(data.JiraServiceManagementExtensionSettings.Elements()) > 0 {
		jiraServiceManagementExtensionSettings := mapJiraServiceManagementExtensionSettings(data.JiraServiceManagementExtensionSettings)
		if jiraServiceManagementExtensionSettings != nil {
			newEnvironment.ExtensionSettings = append(newEnvironment.ExtensionSettings, jiraServiceManagementExtensionSettings)
		}
	}
	if len(data.ServiceNowExtensionSettings.Elements()) > 0 {
		serviceNowExtensionSettings := mapServiceNowExtensionSettings(data.ServiceNowExtensionSettings)
		if serviceNowExtensionSettings != nil {
			newEnvironment.ExtensionSettings = append(newEnvironment.ExtensionSettings, serviceNowExtensionSettings)
		}
	}

	env, err := environments.Add(r.Config.Client, newEnvironment)
	if err != nil {
		resp.Diagnostics.AddError("unable to create environment", err.Error())
		return
	}

	updateEnvironment(ctx, &data, env)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *environmentTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.EnvironmentTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environment, err := environments.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load environment", err.Error())
	}

	updateEnvironment(ctx, &data, environment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *environmentTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state schemas.EnvironmentTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	env, err := environments.GetByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load environment", err.Error())
		return
	}

	updatedEnv := environments.NewEnvironment(data.Name.ValueString())
	updatedEnv.ID = env.ID
	updatedEnv.SpaceID = env.SpaceID
	updatedEnv.Slug = env.Slug
	updatedEnv.Description = data.Description.ValueString()
	updatedEnv.AllowDynamicInfrastructure = data.AllowDynamicInfrastructure.ValueBool()
	updatedEnv.UseGuidedFailure = data.UseGuidedFailure.ValueBool()
	updatedEnv.SortOrder = int(data.SortOrder.ValueInt64())
	if len(data.JiraExtensionSettings.Elements()) > 0 {
		jiraExtensionSettings := mapJiraExtensionSettings(data.JiraExtensionSettings)
		if jiraExtensionSettings != nil {
			updatedEnv.ExtensionSettings = append(updatedEnv.ExtensionSettings, jiraExtensionSettings)
		}
	}
	if len(data.JiraServiceManagementExtensionSettings.Elements()) > 0 {
		jiraServiceManagementExtensionSettings := mapJiraServiceManagementExtensionSettings(data.JiraServiceManagementExtensionSettings)
		if jiraServiceManagementExtensionSettings != nil {
			updatedEnv.ExtensionSettings = append(updatedEnv.ExtensionSettings, jiraServiceManagementExtensionSettings)
		}
	}
	if len(data.ServiceNowExtensionSettings.Elements()) > 0 {
		serviceNowExtensionSettings := mapServiceNowExtensionSettings(data.ServiceNowExtensionSettings)
		if serviceNowExtensionSettings != nil {
			updatedEnv.ExtensionSettings = append(updatedEnv.ExtensionSettings, serviceNowExtensionSettings)
		}
	}

	updatedEnvironment, err := environments.Update(r.Config.Client, updatedEnv)
	if err != nil {
		resp.Diagnostics.AddError("unable to update environment", err.Error())
		return
	}

	updateEnvironment(ctx, &data, updatedEnvironment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *environmentTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.EnvironmentTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := environments.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete environment", err.Error())
		return
	}
}

func updateEnvironment(ctx context.Context, data *schemas.EnvironmentTypeResourceModel, environment *environments.Environment) {
	data.ID = types.StringValue(environment.ID)
	data.SpaceID = types.StringValue(environment.SpaceID)
	data.Slug = types.StringValue(environment.Slug)
	data.Name = types.StringValue(environment.Name)
	data.Description = types.StringValue(environment.Description)
	if !data.AllowDynamicInfrastructure.IsNull() {
		data.AllowDynamicInfrastructure = types.BoolValue(environment.AllowDynamicInfrastructure)
	}
	if !data.UseGuidedFailure.IsNull() {
		data.UseGuidedFailure = types.BoolValue(environment.UseGuidedFailure)
	}
	data.SortOrder = types.Int64Value(int64(environment.SortOrder))
	if len(environment.ExtensionSettings) != 0 {
		for _, extensionSettings := range environment.ExtensionSettings {
			switch extensionSettings.ExtensionID() {
			case extensions.JiraExtensionID:
				if jiraExtensionSettings, ok := extensionSettings.(*environments.JiraExtensionSettings); ok {
					data.JiraExtensionSettings, _ = types.ListValueFrom(
						ctx,
						types.ObjectType{AttrTypes: schemas.JiraExtensionSettingsObjectType()},
						[]any{schemas.MapJiraExtensionSettings(jiraExtensionSettings)},
					)
				}
			case extensions.JiraServiceManagementExtensionID:
				if jiraServiceManagementExtensionSettings, ok := extensionSettings.(*environments.JiraServiceManagementExtensionSettings); ok {
					data.JiraServiceManagementExtensionSettings, _ = types.ListValueFrom(
						ctx,
						types.ObjectType{AttrTypes: schemas.JiraServiceManagementExtensionSettingsObjectType()},
						[]any{schemas.MapJiraServiceManagementExtensionSettings(jiraServiceManagementExtensionSettings)},
					)
				}
			case extensions.ServiceNowExtensionID:
				if serviceNowExtensionSettings, ok := extensionSettings.(*environments.ServiceNowExtensionSettings); ok {
					data.ServiceNowExtensionSettings, _ = types.ListValueFrom(
						ctx,
						types.ObjectType{AttrTypes: schemas.ServiceNowExtensionSettingsObjectType()},
						[]any{schemas.MapServiceNowExtensionSettings(serviceNowExtensionSettings)},
					)
				}
			}
		}
	}
}

func mapJiraExtensionSettings(jiraExtensionSettings types.List) *environments.JiraExtensionSettings {
	obj := jiraExtensionSettings.Elements()[0].(types.Object)
	attrs := obj.Attributes()
	if environmentType, ok := attrs[schemas.EnvironmentJiraExtensionSettingsEnvironmentType].(types.String); ok && !environmentType.IsNull() {
		return environments.NewJiraExtensionSettings(environmentType.ValueString())
	}
	return nil
}

func mapJiraServiceManagementExtensionSettings(jiraServiceManagementExtensionSettings types.List) *environments.JiraServiceManagementExtensionSettings {
	obj := jiraServiceManagementExtensionSettings.Elements()[0].(types.Object)
	attrs := obj.Attributes()
	if isEnabled, ok := attrs[schemas.EnvironmentJiraServiceManagementExtensionSettingsIsEnabled].(types.Bool); ok && !isEnabled.IsNull() {
		return environments.NewJiraServiceManagementExtensionSettings(isEnabled.ValueBool())
	}
	return nil
}

func mapServiceNowExtensionSettings(serviceNowExtensionSettings types.List) *environments.ServiceNowExtensionSettings {
	obj := serviceNowExtensionSettings.Elements()[0].(types.Object)
	attrs := obj.Attributes()
	if isEnabled, ok := attrs[schemas.EnvironmentJiraServiceManagementExtensionSettingsIsEnabled].(types.Bool); ok && !isEnabled.IsNull() {
		return environments.NewServiceNowExtensionSettings(isEnabled.ValueBool())
	}
	return nil
}
