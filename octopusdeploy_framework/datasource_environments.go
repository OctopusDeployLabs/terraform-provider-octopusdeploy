package octopusdeploy_framework

import (
	"context"
	"fmt"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/extensions"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type environmentDataSource struct {
	*Config
}

type environmentsDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	SpaceID      types.String `tfsdk:"space_id"`
	IDs          types.List   `tfsdk:"ids"`
	Name         types.String `tfsdk:"name"`
	PartialName  types.String `tfsdk:"partial_name"`
	Skip         types.Int64  `tfsdk:"skip"`
	Take         types.Int64  `tfsdk:"take"`
	Environments types.List   `tfsdk:"environments"`
}

func NewEnvironmentsDataSource() datasource.DataSource {
	return &environmentDataSource{}
}

func (*environmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("environments")
}

func (*environmentDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides information about existing environments.",
		Attributes: map[string]schema.Attribute{
			//request
			"ids":          util.GetQueryIDsDatasourceSchema(),
			"space_id":     schemas.GetSpaceIdDatasourceSchema(schemas.EnvironmentResourceDescription, false),
			"name":         util.GetQueryNameDatasourceSchema(),
			"partial_name": util.GetQueryPartialNameDatasourceSchema(),
			"skip":         util.GetQuerySkipDatasourceSchema(),
			"take":         util.GetQueryTakeDatasourceSchema(),

			//response
			"id": schemas.GetIdDatasourceSchema(false),
		},
		Blocks: map[string]schema.Block{
			"environments": schema.ListNestedBlock{
				Description: "Provides information about existing environments.",
				NestedObject: schema.NestedBlockObject{
					Attributes: schemas.GetEnvironmentDatasourceSchema(),
				},
			},
		},
	}
}

func (e *environmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	e.Config = DataSourceConfiguration(req, resp)
}

func (e *environmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data environmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := environments.EnvironmentsQuery{
		IDs:         util.GetIds(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Name:        data.Name.ValueString(),
		Skip:        util.GetNumber(data.Skip),
		Take:        util.GetNumber(data.Take),
	}

	existingEnvironments, err := environments.Get(e.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load environments", err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("environments returned from API: %#v", existingEnvironments))
	var mappedEnvironments []schemas.EnvironmentTypeResourceModel
	for _, environment := range existingEnvironments.Items {
		var env schemas.EnvironmentTypeResourceModel
		env.ID = types.StringValue(environment.ID)
		env.SpaceID = types.StringValue(environment.SpaceID)
		env.Slug = types.StringValue(environment.Slug)
		env.Name = types.StringValue(environment.Name)
		env.Description = types.StringValue(environment.Description)
		env.AllowDynamicInfrastructure = types.BoolValue(environment.AllowDynamicInfrastructure)
		env.SortOrder = types.Int64Value(int64(environment.SortOrder))
		env.UseGuidedFailure = types.BoolValue(environment.UseGuidedFailure)
		env.JiraExtensionSettings, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.JiraExtensionSettingsObjectType()}, []any{})
		env.JiraServiceManagementExtensionSettings, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.JiraServiceManagementExtensionSettingsObjectType()}, []any{})
		env.ServiceNowExtensionSettings, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.ServiceNowExtensionSettingsObjectType()}, []any{})

		for _, extensionSetting := range environment.ExtensionSettings {
			switch extensionSetting.ExtensionID() {
			case extensions.JiraExtensionID:
				if jiraExtension, ok := extensionSetting.(*environments.JiraExtensionSettings); ok {
					env.JiraExtensionSettings, _ = types.ListValueFrom(
						ctx,
						types.ObjectType{AttrTypes: schemas.JiraExtensionSettingsObjectType()},
						[]any{schemas.MapJiraExtensionSettings(jiraExtension)},
					)
				}
			case extensions.JiraServiceManagementExtensionID:
				if jiraServiceManagementExtensionSettings, ok := extensionSetting.(*environments.JiraServiceManagementExtensionSettings); ok {
					env.JiraServiceManagementExtensionSettings, _ = types.ListValueFrom(
						ctx,
						types.ObjectType{AttrTypes: schemas.JiraServiceManagementExtensionSettingsObjectType()},
						[]any{schemas.MapJiraServiceManagementExtensionSettings(jiraServiceManagementExtensionSettings)},
					)
				}
			case extensions.ServiceNowExtensionID:
				if serviceNowExtensionSettings, ok := extensionSetting.(*environments.ServiceNowExtensionSettings); ok {
					env.ServiceNowExtensionSettings, _ = types.ListValueFrom(
						ctx,
						types.ObjectType{AttrTypes: schemas.ServiceNowExtensionSettingsObjectType()},
						[]any{schemas.MapServiceNowExtensionSettings(serviceNowExtensionSettings)},
					)
				}
			}
		}

		mappedEnvironments = append(mappedEnvironments, env)
	}

	data.Environments, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: environmentObjectType()}, mappedEnvironments)
	data.ID = types.StringValue("Environments " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func environmentObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"slug":        types.StringType,
		"description": types.StringType,
		schemas.EnvironmentAllowDynamicInfrastructure: types.BoolType,
		schemas.EnvironmentSortOrder:                  types.Int64Type,
		schemas.EnvironmentUseGuidedFailure:           types.BoolType,
		"space_id":                                    types.StringType,
		schemas.EnvironmentJiraExtensionSettings: types.ListType{
			ElemType: types.ObjectType{AttrTypes: schemas.JiraExtensionSettingsObjectType()},
		},
		schemas.EnvironmentJiraServiceManagementExtensionSettings: types.ListType{
			ElemType: types.ObjectType{AttrTypes: schemas.JiraServiceManagementExtensionSettingsObjectType()},
		},
		schemas.EnvironmentServiceNowExtensionSettings: types.ListType{
			ElemType: types.ObjectType{AttrTypes: schemas.ServiceNowExtensionSettingsObjectType()},
		},
	}
}
