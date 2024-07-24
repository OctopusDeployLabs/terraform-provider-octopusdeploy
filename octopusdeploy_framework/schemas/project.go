package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ProjectResourceName = "project"

func GetProjectResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "This resource manages projects in Octopus Deploy.",
		Attributes: map[string]schema.Attribute{
			"id":               util.GetIdResourceSchema(),
			"space_id":         util.GetSpaceIdResourceSchema("project"),
			"name":             util.GetNameResourceSchema(true),
			"description":      util.GetDescriptionResourceSchema("project"),
			"lifecycle_id":     schema.StringAttribute{Description: "The lifecycle ID associated with this project.", Required: true},
			"project_group_id": schema.StringAttribute{Description: "The project group ID associated with this project.", Required: true},
			"is_disabled": schema.BoolAttribute{
				Description: "Indicates whether the project is disabled.",
				Optional:    true,
				Computed:    true,
			},
			"auto_create_release": schema.BoolAttribute{
				Description: "Indicates whether to automatically create a release when a package is pushed to a trigger.",
				Optional:    true,
				Computed:    true,
			},
			"default_guided_failure_mode": schema.StringAttribute{
				Description: "The default guided failure mode setting for the project.",
				Optional:    true,
				Computed:    true,
			},
			"default_to_skip_if_already_installed": schema.BoolAttribute{
				Description: "Indicates whether deployment steps should be skipped if the relevant package is already installed.",
				Optional:    true,
				Computed:    true,
			},
			"deployment_changes_template": schema.StringAttribute{
				Description: "The template to use for deployment change details.",
				Optional:    true,
				Computed:    true,
			},
			"deployment_process_id": schema.StringAttribute{
				Description: "The ID of the deployment process associated with this project.",
				Computed:    true,
			},
			"discrete_channel_release": schema.BoolAttribute{
				Description: "Treats releases of different channels to the same environment as a separate deployment dimension.",
				Optional:    true,
				Computed:    true,
			},
			"is_discrete_channel_release": schema.BoolAttribute{
				Description: "Treats releases of different channels to the same environment as a separate deployment dimension.",
				Optional:    true,
				Computed:    true,
			},
			"is_version_controlled": schema.BoolAttribute{
				Description: "Indicates whether the project is version controlled.",
				Optional:    true,
				Computed:    true,
			},
			"included_library_variable_sets": schema.ListAttribute{
				Description: "The list of included library variable set IDs.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"tenanted_deployment_participation": schema.StringAttribute{
				Description: "The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.",
				Optional:    true,
				Computed:    true,
			},
			"variable_set_id": schema.StringAttribute{
				Description: "The ID of the variable set associated with this project.",
				Computed:    true,
			},
			"release_notes_template": schema.StringAttribute{
				Description: "The template to use for release notes.",
				Optional:    true,
				Computed:    true,
			},
			"slug": util.GetSlugResourceSchema("project"),
		},
		Blocks: map[string]schema.Block{
			"connectivity_policy": schema.SingleNestedBlock{
				Description: "Defines the connectivity policy for deployments.",
				Attributes: map[string]schema.Attribute{
					"allow_deployments_to_no_targets": schema.BoolAttribute{
						Description: "Allow deployments to be created when there are no targets.",
						Optional:    true,
					},
					"exclude_unhealthy_targets": schema.BoolAttribute{
						Description: "Exclude unhealthy targets from deployments.",
						Optional:    true,
					},
					"skip_machine_behavior": schema.StringAttribute{
						Description: "The behavior when a machine is skipped.",
						Optional:    true,
					},
					"target_roles": schema.ListAttribute{
						Description: "The target roles for the connectivity policy.",
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
			"git_anonymous_persistence_settings": schema.SingleNestedBlock{
				Description: "Provides Git-related persistence settings for a version-controlled project.",
				Attributes: map[string]schema.Attribute{
					"url":                schema.StringAttribute{Optional: true},
					"base_path":          schema.StringAttribute{Description: "The base path associated with these version control settings.", Optional: true},
					"default_branch":     schema.StringAttribute{Description: "The default branch associated with these version control settings.", Optional: true},
					"protected_branches": schema.SetAttribute{Description: "A list of protected branch patterns.", ElementType: types.StringType, Optional: true},
				},
			},
			"git_library_persistence_settings": schema.SingleNestedBlock{
				Description: "Provides Git-related persistence settings for a version-controlled project.",
				Attributes: map[string]schema.Attribute{
					"git_credential_id":  schema.StringAttribute{Description: "The ID of the Git credential to use.", Optional: true},
					"url":                schema.StringAttribute{Optional: true},
					"base_path":          schema.StringAttribute{Description: "The base path associated with these version control settings.", Optional: true},
					"default_branch":     schema.StringAttribute{Description: "The default branch associated with these version control settings.", Optional: true},
					"protected_branches": schema.SetAttribute{Description: "A list of protected branch patterns.", ElementType: types.StringType, Optional: true},
				},
			},
			"git_username_password_persistence_settings": schema.SingleNestedBlock{
				Description: "Provides Git-related persistence settings for a version-controlled project.",
				Attributes: map[string]schema.Attribute{
					"url":                schema.StringAttribute{Optional: true},
					"username":           util.GetUsernameResourceSchema(false),
					"password":           util.GetPasswordResourceSchema(false),
					"base_path":          schema.StringAttribute{Description: "The base path associated with these version control settings.", Optional: true},
					"default_branch":     schema.StringAttribute{Description: "The default branch associated with these version control settings.", Optional: true},
					"protected_branches": schema.SetAttribute{Description: "A list of protected branch patterns.", ElementType: types.StringType, Optional: true},
				},
			},
			"jira_service_management_extension_settings": schema.SingleNestedBlock{
				Description: "Provides extension settings for the Jira Service Management (JSM) integration for this project.",
				Attributes: map[string]schema.Attribute{
					"connection_id":             schema.StringAttribute{Description: "The connection identifier associated with the extension settings.", Optional: true},
					"is_enabled":                schema.BoolAttribute{Description: "Specifies whether or not this extension is enabled for this project.", Optional: true},
					"service_desk_project_name": schema.StringAttribute{Description: "The project name associated with this extension.", Optional: true},
				},
			},
			"servicenow_extension_settings": schema.SingleNestedBlock{
				Description: "Provides extension settings for the ServiceNow integration for this project.",
				Attributes: map[string]schema.Attribute{
					"connection_id":                       schema.StringAttribute{Description: "The connection identifier associated with the extension settings.", Optional: true},
					"is_enabled":                          schema.BoolAttribute{Description: "Specifies whether or not this extension is enabled for this project.", Optional: true},
					"is_state_automatically_transitioned": schema.BoolAttribute{Description: "Specifies whether or not this extension will automatically transition the state of a deployment for this project.", Optional: true},
					"standard_change_template_name":       schema.StringAttribute{Description: "The name of the standard change template associated with this extension.", Optional: true},
				},
			},
			"versioning_strategy": schema.SingleNestedBlock{
				Description: "Defines the versioning strategy for the project.",
				Attributes: map[string]schema.Attribute{
					"donor_package_step_id": schema.StringAttribute{Description: "The ID of the step containing the donor package.", Optional: true},
					"template":              schema.StringAttribute{Description: "The template to use for version numbers.", Optional: true},
				},
				Blocks: map[string]schema.Block{
					"donor_package": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"deployment_action": schema.StringAttribute{Description: "The deployment action for the donor package.", Optional: true},
							"package_reference": schema.StringAttribute{Description: "The package reference for the donor package.", Optional: true},
						},
					},
				},
			},
			"release_creation_strategy": schema.SingleNestedBlock{
				Description: "Defines the release creation strategy for the project.",
				Attributes: map[string]schema.Attribute{
					"channel_id":                       schema.StringAttribute{Description: "The ID of the channel to use for release creation.", Optional: true},
					"release_creation_package_step_id": schema.StringAttribute{Description: "The ID of the step containing the package for release creation.", Optional: true},
				},
				Blocks: map[string]schema.Block{
					"release_creation_package": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"deployment_action": schema.StringAttribute{Description: "The deployment action for the release creation package.", Optional: true},
							"package_reference": schema.StringAttribute{Description: "The package reference for the release creation package.", Optional: true},
						},
					},
				},
			},
			"template": schema.ListNestedBlock{
				Description: "Defines template parameters for the project.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":             util.GetNameResourceSchema(false),
						"label":            schema.StringAttribute{Description: "The label shown beside the parameter when presented in the deployment process.", Optional: true},
						"help_text":        schema.StringAttribute{Description: "The help presented alongside the parameter input.", Optional: true},
						"default_value":    schema.StringAttribute{Description: "A default value for the parameter, if applicable.", Optional: true},
						"display_settings": schema.MapAttribute{Description: "The display settings for the parameter.", ElementType: types.StringType, Optional: true},
					},
				},
			},
		},
	}
}
