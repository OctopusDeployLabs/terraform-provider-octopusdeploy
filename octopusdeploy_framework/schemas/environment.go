package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	EnvironmentResourceDescription                             = "environment"
	EnvironmentSortOrder                                       = "sort_order"
	EnvironmentAllowDynamicInfrastructure                      = "allow_dynamic_infrastructure"
	EnvironmentUseGuidedFailure                                = "use_guided_failure"
	EnvironmentJiraExtensionSettings                           = "jira_extension_settings"
	EnvironmentJiraServiceManagementExtensionSettings          = "jira_service_management_extension_settings"
	EnvironmentServiceNowExtensionSettings                     = "servicenow_extension_settings"
	EnvironmentJiraExtensionSettingsEnvironmentType            = "environment_type"
	EnvironmentJiraServiceManagementExtensionSettingsIsEnabled = "is_enabled"
	EnvironmentServiceNowExtensionSettingsIsEnabled            = "is_enabled"
)

func GetEnvironmentDatasourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":                 GetIdDatasourceSchema(true),
		"slug":               util.GetSlugDatasourceSchema(EnvironmentResourceDescription, true),
		"name":               GetReadonlyNameDatasourceSchema(),
		"description":        util.GetDescriptionDatasourceSchema(EnvironmentResourceDescription),
		EnvironmentSortOrder: util.GetSortOrderDataSourceSchema(EnvironmentResourceDescription),
		EnvironmentAllowDynamicInfrastructure: datasourceSchema.BoolAttribute{
			Computed: true,
		},
		EnvironmentUseGuidedFailure: datasourceSchema.BoolAttribute{
			Computed: true,
		},
		"space_id": GetSpaceIdDatasourceSchema(EnvironmentResourceDescription, true),
		EnvironmentJiraExtensionSettings: datasourceSchema.ListNestedAttribute{
			Description: "Provides extension settings for the Jira integration for this environment.",
			Computed:    true,
			NestedObject: datasourceSchema.NestedAttributeObject{
				Attributes: map[string]datasourceSchema.Attribute{
					EnvironmentJiraExtensionSettingsEnvironmentType: datasourceSchema.StringAttribute{
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
			Computed:    true,
			NestedObject: datasourceSchema.NestedAttributeObject{
				Attributes: map[string]datasourceSchema.Attribute{
					EnvironmentJiraServiceManagementExtensionSettingsIsEnabled: datasourceSchema.BoolAttribute{Computed: true},
				},
			},
		},
		EnvironmentServiceNowExtensionSettings: datasourceSchema.ListNestedAttribute{
			Description: "Provides extension settings for the ServiceNow integration for this environment.",
			Computed:    true,
			NestedObject: datasourceSchema.NestedAttributeObject{
				Attributes: map[string]datasourceSchema.Attribute{
					EnvironmentJiraServiceManagementExtensionSettingsIsEnabled: datasourceSchema.BoolAttribute{Computed: true},
				},
			},
		},
	}
}

func GetEnvironmentResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":                 util.GetIdResourceSchema(),
			"slug":               util.GetSlugResourceSchema(EnvironmentResourceDescription),
			"name":               util.GetNameResourceSchema(true),
			"description":        util.GetDescriptionResourceSchema(EnvironmentResourceDescription),
			EnvironmentSortOrder: util.GetSortOrderResourceSchema(EnvironmentResourceDescription),
			EnvironmentAllowDynamicInfrastructure: resourceSchema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			EnvironmentUseGuidedFailure: resourceSchema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"space_id": util.GetSpaceIdResourceSchema(EnvironmentResourceDescription),
		},
		Blocks: map[string]resourceSchema.Block{
			EnvironmentJiraExtensionSettings: resourceSchema.ListNestedBlock{
				Description: "Provides extension settings for the Jira integration for this environment.",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"environment_type": resourceSchema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"development",
									"production",
									"staging",
									"testing",
									"unmapped",
								),
							},
						},
					},
				},
			},
			EnvironmentJiraServiceManagementExtensionSettings: resourceSchema.ListNestedBlock{
				Description: "Provides extension settings for the Jira Service Management (JSM) integration for this environment.",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"is_enabled": resourceSchema.BoolAttribute{Optional: true},
					},
				},
			},
			EnvironmentServiceNowExtensionSettings: resourceSchema.ListNestedBlock{
				Description: "Provides extension settings for the ServiceNow integration for this environment.",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"is_enabled": resourceSchema.BoolAttribute{Optional: true},
					},
				},
			},
		},
	}
}

func JiraExtensionSettingsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"environment_type": types.StringType,
	}
}

func MapJiraExtensionSettings(jiraExtensionSettings *environments.JiraExtensionSettings) attr.Value {
	return types.ObjectValueMust(JiraExtensionSettingsObjectType(), map[string]attr.Value{
		"environment_type": types.StringValue(jiraExtensionSettings.JiraEnvironmentType),
	})
}

func JiraServiceManagementExtensionSettingsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"is_enabled": types.BoolType,
	}
}

func MapJiraServiceManagementExtensionSettings(jiraServiceManagementExtensionSettings *environments.JiraServiceManagementExtensionSettings) attr.Value {
	return types.ObjectValueMust(JiraServiceManagementExtensionSettingsObjectType(), map[string]attr.Value{
		"is_enabled": types.BoolValue(jiraServiceManagementExtensionSettings.IsChangeControlled()),
	})
}

func ServiceNowExtensionSettingsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"is_enabled": types.BoolType,
	}
}

func MapServiceNowExtensionSettings(serviceNowExtensionSettings *environments.ServiceNowExtensionSettings) attr.Value {
	return types.ObjectValueMust(ServiceNowExtensionSettingsObjectType(), map[string]attr.Value{
		"is_enabled": types.BoolValue(serviceNowExtensionSettings.IsChangeControlled()),
	})
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
