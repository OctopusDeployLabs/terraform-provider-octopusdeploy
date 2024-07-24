package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ProjectResourceName = "project"

func GetProjectResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "This resource manages projects in Octopus Deploy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique ID for this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"space_id": schema.StringAttribute{
				Description: "The space ID associated with this project.",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the project in Octopus Deploy. This name must be unique.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of this project.",
				Optional:    true,
			},
			"auto_create_release": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"cloned_from_project_id": schema.StringAttribute{
				Optional: true,
			},
			"default_guided_failure_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"default_to_skip_if_already_installed": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"deployment_changes_template": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"discrete_channel_release": schema.BoolAttribute{
				Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
				Optional:    true,
				Computed:    true,
			},
			"is_disabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"is_discrete_channel_release": schema.BoolAttribute{
				Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
				Optional:    true,
				Computed:    true,
			},
			"is_version_controlled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"lifecycle_id": schema.StringAttribute{
				Description: "The lifecycle ID associated with this project.",
				Required:    true,
			},
			"project_group_id": schema.StringAttribute{
				Description: "The project group ID associated with this project.",
				Required:    true,
			},
			"tenanted_deployment_participation": schema.StringAttribute{
				Description: "The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.",
				Optional:    true,
				Computed:    true,
			},
			"included_library_variable_sets": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"release_notes_template": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"slug": schema.StringAttribute{
				Description: "A human-readable, unique identifier, used to identify a project.",
				Optional:    true,
				Computed:    true,
			},
			"deployment_process_id": schema.StringAttribute{
				Computed: true,
			},
			"variable_set_id": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"auto_deploy_release_overrides": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"environment_id": schema.StringAttribute{Optional: true},
						"release_id":     schema.StringAttribute{Optional: true},
						"tenant_id":      schema.StringAttribute{Optional: true},
					},
				},
			},
			"connectivity_policy": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"allow_deployments_to_no_targets": schema.BoolAttribute{
							Optional: true,
						},
						"exclude_unhealthy_targets": schema.BoolAttribute{
							Optional: true,
						},
						"skip_machine_behavior": schema.StringAttribute{
							Optional: true,
						},
						"target_roles": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			"git_anonymous_persistence_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Description: "The URL associated with these version control settings.",
							Optional:    true,
						},
						"base_path": schema.StringAttribute{
							Description: "The base path associated with these version control settings.",
							Optional:    true,
						},
						"default_branch": schema.StringAttribute{
							Description: "The default branch associated with these version control settings.",
							Optional:    true,
						},
						"protected_branches": schema.SetAttribute{
							Description: "A list of protected branch patterns.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"git_library_persistence_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"git_credential_id": schema.StringAttribute{
							Optional: true,
						},
						"url": schema.StringAttribute{
							Description: "The URL associated with these version control settings.",
							Optional:    true,
						},
						"base_path": schema.StringAttribute{
							Description: "The base path associated with these version control settings.",
							Optional:    true,
						},
						"default_branch": schema.StringAttribute{
							Description: "The default branch associated with these version control settings.",
							Optional:    true,
						},
						"protected_branches": schema.SetAttribute{
							Description: "A list of protected branch patterns.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"git_username_password_persistence_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"url": schema.StringAttribute{
							Description: "The URL associated with these version control settings.",
							Optional:    true,
						},
						"username": schema.StringAttribute{
							Description: "The username for the Git credential.",
							Optional:    true,
						},
						"password": schema.StringAttribute{
							Description: "The password for the Git credential.",
							Optional:    true,
							Sensitive:   true,
						},
						"base_path": schema.StringAttribute{
							Description: "The base path associated with these version control settings.",
							Optional:    true,
						},
						"default_branch": schema.StringAttribute{
							Description: "The default branch associated with these version control settings.",
							Optional:    true,
						},
						"protected_branches": schema.SetAttribute{
							Description: "A list of protected branch patterns.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"jira_service_management_extension_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"connection_id": schema.StringAttribute{
							Description: "The connection identifier associated with the extension settings.",
							Optional:    true,
						},
						"is_enabled": schema.BoolAttribute{
							Description: "Specifies whether or not this extension is enabled for this project.",
							Optional:    true,
						},
						"service_desk_project_name": schema.StringAttribute{
							Description: "The project name associated with this extension.",
							Optional:    true,
						},
					},
				},
				Description: "Provides extension settings for the Jira Service Management (JSM) integration for this project.",
			},
			"servicenow_extension_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"connection_id": schema.StringAttribute{
							Description: "The connection identifier associated with the extension settings.",
							Optional:    true,
						},
						"is_enabled": schema.BoolAttribute{
							Description: "Specifies whether or not this extension is enabled for this project.",
							Optional:    true,
						},
						"is_state_automatically_transitioned": schema.BoolAttribute{
							Description: "Specifies whether or not this extension will automatically transition the state of a deployment for this project.",
							Optional:    true,
						},
						"standard_change_template_name": schema.StringAttribute{
							Description: "The name of the standard change template associated with this extension.",
							Optional:    true,
						},
					},
				},
				Description: "Provides extension settings for the ServiceNow integration for this project.",
			},
			"template": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the template parameter.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							Description: "The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods.",
							Required:    true,
						},
						"label": schema.StringAttribute{
							Description: "The label shown beside the parameter when presented in the deployment process.",
							Optional:    true,
						},
						"help_text": schema.StringAttribute{
							Description: "The help presented alongside the parameter input.",
							Optional:    true,
						},
						"default_value": schema.StringAttribute{
							Description: "A default value for the parameter, if applicable.",
							Optional:    true,
						},
						"display_settings": schema.MapAttribute{
							Description: "The display settings for the parameter.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			"versioning_strategy": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"donor_package_step_id": schema.StringAttribute{
							Optional: true,
						},
						"template": schema.StringAttribute{
							Optional: true,
						},
					},
					Blocks: map[string]schema.Block{
						"donor_package": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"deployment_action": schema.StringAttribute{
										Optional: true,
									},
									"package_reference": schema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"release_creation_strategy": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"channel_id": schema.StringAttribute{
							Optional: true,
						},
						"release_creation_package_step_id": schema.StringAttribute{
							Optional: true,
						},
					},
					Blocks: map[string]schema.Block{
						"release_creation_package": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"deployment_action": schema.StringAttribute{
										Optional: true,
									},
									"package_reference": schema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
