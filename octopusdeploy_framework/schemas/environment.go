package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	EnvironmentResourceDescription                    = "environment"
	EnvironmentSortOrder                              = "sort_order"
	EnvironmentAllowDynamicInfrastructure             = "allow_dynamic_infrastructure"
	EnvironmentUseGuidedFailure                       = "use_guided_failure"
	EnvironmentJiraExtensionSettings                  = "jira_extension_settings"
	EnvironmentJiraServiceManagementExtensionSettings = "jira_service_management_extension_settings"
	EnvironmentServiceNowExtensionSettings            = "servicenow_extension_settings"
)

func GetEnvironmentDatasourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":                 util.GetIdDatasourceSchema(),
		"slug":               util.GetSlugDatasourceSchema(EnvironmentResourceDescription),
		"name":               util.GetNameDatasourceWithMaxLengthSchema(true, 50),
		"description":        util.GetDescriptionDatasourceSchema(EnvironmentResourceDescription),
		EnvironmentSortOrder: util.GetSortOrderDataSourceSchema(EnvironmentResourceDescription),
		EnvironmentAllowDynamicInfrastructure: datasourceSchema.BoolAttribute{
			Optional: true,
		},
		EnvironmentUseGuidedFailure: datasourceSchema.BoolAttribute{
			Optional: true,
		},
		EnvironmentJiraExtensionSettings: datasourceSchema.ListNestedAttribute{
			Description: "Provides extension settings for the Jira integration for this environment.",
			Optional:    true,
			Computed:    true,
			NestedObject: datasourceSchema.NestedAttributeObject{
				Attributes: map[string]datasourceSchema.Attribute{
					"environment_type": datasourceSchema.StringAttribute{
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"development",
								"production",
								"testing",
								"staging",
								"unmapped",
							),
						},
					},
				},
			},
		},
		EnvironmentJiraServiceManagementExtensionSettings: datasourceSchema.ListNestedAttribute{
			Description: "Provides extension settings for the Jira Service Management (JSM) integration for this environment.",
			Optional:    true,
			Computed:    true,
			NestedObject: datasourceSchema.NestedAttributeObject{
				Attributes: map[string]datasourceSchema.Attribute{
					"is_enabled": datasourceSchema.BoolAttribute{Computed: true},
				},
			},
		},
		EnvironmentServiceNowExtensionSettings: datasourceSchema.ListNestedAttribute{
			Description: "Provides extension settings for the ServiceNow integration for this environment.",
			Optional:    true,
			Computed:    true,
			NestedObject: datasourceSchema.NestedAttributeObject{
				Attributes: map[string]datasourceSchema.Attribute{
					"is_enabled": datasourceSchema.BoolAttribute{Computed: true},
				},
			},
		},
		"space_id": util.GetSpaceIdDatasourceSchema(EnvironmentResourceDescription),
	}
}

type EnvironmentTypeResourceModel struct {
	ID                                     types.String `tfsdk:"id"`
	Slug                                   types.String `tfsdk:"slug"`
	Name                                   types.String `tfsdk:"name"`
	Description                            types.String `tfsdk:"description"`
	AllowDynamicInfrastructure             types.Bool   `tfsdk:"allow_dynamic_infrastructure"`
	SortOrder                              types.Int64  `tfsdk:"sort_order"`
	UseGuidedFailure                       types.Bool   `tfsdk:"use_guided_failure"`
	SpaceID                                types.String `tfsdk:"space_id"`
	JiraExtensionSettings                  types.List   `tfsdk:"jira_extension_settings"`
	JiraServiceManagementExtensionSettings types.List   `tfsdk:"jira_service_management_extension_settings"`
	ServiceNowExtensionSettings            types.List   `tfsdk:"servicenow_extension_settings"`
}
