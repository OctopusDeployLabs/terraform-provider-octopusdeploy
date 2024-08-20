package schemas

import (
	"context"
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/extensions"
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

type EnvironmentSchema struct{}

var _ EntitySchema = EnvironmentSchema{}

var jiraEnvironmentTypeNames = struct {
	Development string
	Production  string
	Testing     string
	Staging     string
	Unmapped    string
}{
	Development: "development",
	Production:  "production",
	Testing:     "testing",
	Staging:     "staging",
	Unmapped:    "unmapped",
}

var jiraEnvironmentTypes = []string{
	jiraEnvironmentTypeNames.Development,
	jiraEnvironmentTypeNames.Production,
	jiraEnvironmentTypeNames.Staging,
	jiraEnvironmentTypeNames.Testing,
	jiraEnvironmentTypeNames.Unmapped,
}

func EnvironmentObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                  types.StringType,
		"name":                                types.StringType,
		"slug":                                types.StringType,
		"description":                         types.StringType,
		EnvironmentAllowDynamicInfrastructure: types.BoolType,
		EnvironmentSortOrder:                  types.Int64Type,
		EnvironmentUseGuidedFailure:           types.BoolType,
		"space_id":                            types.StringType,
		EnvironmentJiraExtensionSettings: types.ListType{
			ElemType: types.ObjectType{AttrTypes: JiraExtensionSettingsObjectType()},
		},
		EnvironmentJiraServiceManagementExtensionSettings: types.ListType{
			ElemType: types.ObjectType{AttrTypes: JiraServiceManagementExtensionSettingsObjectType()},
		},
		EnvironmentServiceNowExtensionSettings: types.ListType{
			ElemType: types.ObjectType{AttrTypes: ServiceNowExtensionSettingsObjectType()},
		},
	}
}

func (e EnvironmentSchema) GetDatasourceSchemaAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":                 GetIdDatasourceSchema(true),
		"slug":               GetSlugDatasourceSchema(EnvironmentResourceDescription, true),
		"name":               GetReadonlyNameDatasourceSchema(),
		"description":        GetDescriptionDatasourceSchema(EnvironmentResourceDescription),
		EnvironmentSortOrder: GetSortOrderDatasourceSchema(EnvironmentResourceDescription),
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
								jiraEnvironmentTypes...,
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

func (e EnvironmentSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(EnvironmentResourceDescription),
		Attributes: map[string]resourceSchema.Attribute{
			"id":                 GetIdResourceSchema(),
			"slug":               GetSlugResourceSchema(EnvironmentResourceDescription),
			"name":               GetNameResourceSchema(true),
			"description":        GetDescriptionResourceSchema(EnvironmentResourceDescription),
			EnvironmentSortOrder: GetSortOrderResourceSchema(EnvironmentResourceDescription),
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
			"space_id": GetSpaceIdResourceSchema(EnvironmentResourceDescription),
		},
		Blocks: map[string]resourceSchema.Block{
			EnvironmentJiraExtensionSettings: resourceSchema.ListNestedBlock{
				Description: "Provides extension settings for the Jira integration for this environment.",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"environment_type": resourceSchema.StringAttribute{
							Description: fmt.Sprintf("The Jira environment type of this Octopus deployment environment. Valid values are %s.", strings.Join(util.Map(jiraEnvironmentTypes, func(item string) string { return fmt.Sprintf("`\"%s\"`", item) }), ", ")),
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									jiraEnvironmentTypes...,
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
						"is_enabled": resourceSchema.BoolAttribute{
							Description: "Specifies whether or not this extension is enabled for this project.",
							Optional:    true,
						},
					},
				},
			},
			EnvironmentServiceNowExtensionSettings: resourceSchema.ListNestedBlock{
				Description: "Provides extension settings for the ServiceNow integration for this environment.",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"is_enabled": resourceSchema.BoolAttribute{
							Description: "Specifies whether or not this extension is enabled for this project.",
							Optional:    true,
						},
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

func MapFromEnvironment(ctx context.Context, environment *environments.Environment) EnvironmentTypeResourceModel {
	var env EnvironmentTypeResourceModel
	env.ID = types.StringValue(environment.ID)
	env.SpaceID = types.StringValue(environment.SpaceID)
	env.Slug = types.StringValue(environment.Slug)
	env.Name = types.StringValue(environment.Name)
	env.Description = types.StringValue(environment.Description)
	env.AllowDynamicInfrastructure = types.BoolValue(environment.AllowDynamicInfrastructure)
	env.SortOrder = types.Int64Value(int64(environment.SortOrder))
	env.UseGuidedFailure = types.BoolValue(environment.UseGuidedFailure)
	env.JiraExtensionSettings, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: JiraExtensionSettingsObjectType()}, []any{})
	env.JiraServiceManagementExtensionSettings, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: JiraServiceManagementExtensionSettingsObjectType()}, []any{})
	env.ServiceNowExtensionSettings, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: ServiceNowExtensionSettingsObjectType()}, []any{})

	for _, extensionSetting := range environment.ExtensionSettings {
		switch extensionSetting.ExtensionID() {
		case extensions.JiraExtensionID:
			if jiraExtension, ok := extensionSetting.(*environments.JiraExtensionSettings); ok {
				env.JiraExtensionSettings, _ = types.ListValueFrom(
					ctx,
					types.ObjectType{AttrTypes: JiraExtensionSettingsObjectType()},
					[]any{MapJiraExtensionSettings(jiraExtension)},
				)
			}
		case extensions.JiraServiceManagementExtensionID:
			if jiraServiceManagementExtensionSettings, ok := extensionSetting.(*environments.JiraServiceManagementExtensionSettings); ok {
				env.JiraServiceManagementExtensionSettings, _ = types.ListValueFrom(
					ctx,
					types.ObjectType{AttrTypes: JiraServiceManagementExtensionSettingsObjectType()},
					[]any{MapJiraServiceManagementExtensionSettings(jiraServiceManagementExtensionSettings)},
				)
			}
		case extensions.ServiceNowExtensionID:
			if serviceNowExtensionSettings, ok := extensionSetting.(*environments.ServiceNowExtensionSettings); ok {
				env.ServiceNowExtensionSettings, _ = types.ListValueFrom(
					ctx,
					types.ObjectType{AttrTypes: ServiceNowExtensionSettingsObjectType()},
					[]any{MapServiceNowExtensionSettings(serviceNowExtensionSettings)},
				)
			}
		}
	}
	return env
}

type EnvironmentTypeResourceModel struct {
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

	ResourceModel
}
