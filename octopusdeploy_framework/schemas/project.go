package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ProjectResourceName = "project"

func GetProjectResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages projects in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": resourceSchema.StringAttribute{
				Description: "The unique ID for this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"space_id": resourceSchema.StringAttribute{
				Description: "The space ID associated with this project.",
				Optional:    true,
				Computed:    true,
			},
			"name": resourceSchema.StringAttribute{
				Description: "The name of the project in Octopus Deploy. This name must be unique.",
				Required:    true,
			},
			"description": resourceSchema.StringAttribute{
				Description: "The description of this project.",
				Optional:    true,
			},
			"auto_create_release": resourceSchema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"cloned_from_project_id": resourceSchema.StringAttribute{
				Optional: true,
			},
			"default_guided_failure_mode": resourceSchema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"default_to_skip_if_already_installed": resourceSchema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"deployment_changes_template": resourceSchema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"discrete_channel_release": resourceSchema.BoolAttribute{
				Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
				Optional:    true,
				Computed:    true,
			},
			"is_disabled": resourceSchema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"is_discrete_channel_release": resourceSchema.BoolAttribute{
				Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
				Optional:    true,
				Computed:    true,
			},
			"is_version_controlled": resourceSchema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"lifecycle_id": resourceSchema.StringAttribute{
				Description: "The lifecycle ID associated with this project.",
				Required:    true,
			},
			"project_group_id": resourceSchema.StringAttribute{
				Description: "The project group ID associated with this project.",
				Required:    true,
			},
			"tenanted_deployment_participation": resourceSchema.StringAttribute{
				Description: "The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.",
				Optional:    true,
				Computed:    true,
			},
			"included_library_variable_sets": resourceSchema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"release_notes_template": resourceSchema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"slug": resourceSchema.StringAttribute{
				Description: "A human-readable, unique identifier, used to identify a project.",
				Optional:    true,
				Computed:    true,
			},
			"deployment_process_id": resourceSchema.StringAttribute{
				Computed: true,
			},
			"variable_set_id": resourceSchema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]resourceSchema.Block{
			"auto_deploy_release_overrides": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"environment_id": resourceSchema.StringAttribute{Optional: true},
						"release_id":     resourceSchema.StringAttribute{Optional: true},
						"tenant_id":      resourceSchema.StringAttribute{Optional: true},
					},
				},
			},
			"connectivity_policy": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"allow_deployments_to_no_targets": resourceSchema.BoolAttribute{
							Optional: true,
						},
						"exclude_unhealthy_targets": resourceSchema.BoolAttribute{
							Optional: true,
						},
						"skip_machine_behavior": resourceSchema.StringAttribute{
							Optional: true,
						},
						"target_roles": resourceSchema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			"git_anonymous_persistence_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"url": resourceSchema.StringAttribute{
							Description: "The URL associated with these version control settings.",
							Optional:    true,
						},
						"base_path": resourceSchema.StringAttribute{
							Description: "The base path associated with these version control settings.",
							Optional:    true,
						},
						"default_branch": resourceSchema.StringAttribute{
							Description: "The default branch associated with these version control settings.",
							Optional:    true,
						},
						"protected_branches": resourceSchema.SetAttribute{
							Description: "A list of protected branch patterns.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"git_library_persistence_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"git_credential_id": resourceSchema.StringAttribute{
							Optional: true,
						},
						"url": resourceSchema.StringAttribute{
							Description: "The URL associated with these version control settings.",
							Optional:    true,
						},
						"base_path": resourceSchema.StringAttribute{
							Description: "The base path associated with these version control settings.",
							Optional:    true,
						},
						"default_branch": resourceSchema.StringAttribute{
							Description: "The default branch associated with these version control settings.",
							Optional:    true,
						},
						"protected_branches": resourceSchema.SetAttribute{
							Description: "A list of protected branch patterns.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"git_username_password_persistence_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"url": resourceSchema.StringAttribute{
							Description: "The URL associated with these version control settings.",
							Optional:    true,
						},
						"username": resourceSchema.StringAttribute{
							Description: "The username for the Git credential.",
							Optional:    true,
						},
						"password": resourceSchema.StringAttribute{
							Description: "The password for the Git credential.",
							Optional:    true,
							Sensitive:   true,
						},
						"base_path": resourceSchema.StringAttribute{
							Description: "The base path associated with these version control settings.",
							Optional:    true,
						},
						"default_branch": resourceSchema.StringAttribute{
							Description: "The default branch associated with these version control settings.",
							Optional:    true,
						},
						"protected_branches": resourceSchema.SetAttribute{
							Description: "A list of protected branch patterns.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"jira_service_management_extension_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"connection_id": resourceSchema.StringAttribute{
							Description: "The connection identifier associated with the extension settings.",
							Optional:    true,
						},
						"is_enabled": resourceSchema.BoolAttribute{
							Description: "Specifies whether or not this extension is enabled for this project.",
							Optional:    true,
						},
						"service_desk_project_name": resourceSchema.StringAttribute{
							Description: "The project name associated with this extension.",
							Optional:    true,
						},
					},
				},
				Description: "Provides extension settings for the Jira Service Management (JSM) integration for this project.",
			},
			"servicenow_extension_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"connection_id": resourceSchema.StringAttribute{
							Description: "The connection identifier associated with the extension settings.",
							Optional:    true,
						},
						"is_enabled": resourceSchema.BoolAttribute{
							Description: "Specifies whether or not this extension is enabled for this project.",
							Optional:    true,
						},
						"is_state_automatically_transitioned": resourceSchema.BoolAttribute{
							Description: "Specifies whether or not this extension will automatically transition the state of a deployment for this project.",
							Optional:    true,
						},
						"standard_change_template_name": resourceSchema.StringAttribute{
							Description: "The name of the standard change template associated with this extension.",
							Optional:    true,
						},
					},
				},
				Description: "Provides extension settings for the ServiceNow integration for this project.",
			},
			"template": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"id": resourceSchema.StringAttribute{
							Description: "The ID of the template parameter.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": resourceSchema.StringAttribute{
							Description: "The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods.",
							Required:    true,
						},
						"label": resourceSchema.StringAttribute{
							Description: "The label shown beside the parameter when presented in the deployment process.",
							Optional:    true,
						},
						"help_text": resourceSchema.StringAttribute{
							Description: "The help presented alongside the parameter input.",
							Optional:    true,
						},
						"default_value": resourceSchema.StringAttribute{
							Description: "A default value for the parameter, if applicable.",
							Optional:    true,
						},
						"display_settings": resourceSchema.MapAttribute{
							Description: "The display settings for the parameter.",
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
			"versioning_strategy": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"donor_package_step_id": resourceSchema.StringAttribute{
							Optional: true,
						},
						"template": resourceSchema.StringAttribute{
							Optional: true,
						},
					},
					Blocks: map[string]resourceSchema.Block{
						"donor_package": resourceSchema.ListNestedBlock{
							NestedObject: resourceSchema.NestedBlockObject{
								Attributes: map[string]resourceSchema.Attribute{
									"deployment_action": resourceSchema.StringAttribute{
										Optional: true,
									},
									"package_reference": resourceSchema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"release_creation_strategy": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"channel_id": resourceSchema.StringAttribute{
							Optional: true,
						},
						"release_creation_package_step_id": resourceSchema.StringAttribute{
							Optional: true,
						},
					},
					Blocks: map[string]resourceSchema.Block{
						"release_creation_package": resourceSchema.ListNestedBlock{
							NestedObject: resourceSchema.NestedBlockObject{
								Attributes: map[string]resourceSchema.Attribute{
									"deployment_action": resourceSchema.StringAttribute{
										Optional: true,
									},
									"package_reference": resourceSchema.StringAttribute{
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

func GetProjectDataSourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id": datasourceSchema.StringAttribute{
			Computed:    true,
			Description: "An identifier for the data source, which includes the query timestamp.",
		},
		"space_id": datasourceSchema.StringAttribute{
			Optional:    true,
			Description: "The space ID to filter by.",
		},
		"cloned_from_project_id": datasourceSchema.StringAttribute{
			Optional:    true,
			Description: "The ID of the project this project was cloned from.",
		},
		"ids": datasourceSchema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "A list of project IDs to filter by.",
		},
		"is_clone": datasourceSchema.BoolAttribute{
			Optional:    true,
			Description: "If set, only return projects that are clones.",
		},
		"name": datasourceSchema.StringAttribute{
			Optional:    true,
			Description: "The name of the project to filter by.",
		},
		"partial_name": datasourceSchema.StringAttribute{
			Optional:    true,
			Description: "A partial name of the project to filter by.",
		},
		"skip": datasourceSchema.Int64Attribute{
			Optional:    true,
			Description: "Number of items to skip. Defaults to zero.",
		},
		"take": datasourceSchema.Int64Attribute{
			Optional:    true,
			Description: "Number of items to take. Defaults to 30.",
		},
		"projects": datasourceSchema.ListNestedAttribute{
			Computed: true,
			NestedObject: datasourceSchema.NestedAttributeObject{
				Attributes: map[string]datasourceSchema.Attribute{
					"id": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The ID of the project.",
					},
					"space_id": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The space ID of the project.",
					},
					"name": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The name of the project.",
					},
					"description": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The description of the project.",
					},
					"auto_create_release": datasourceSchema.BoolAttribute{
						Computed:    true,
						Description: "Whether to automatically create a release when a package is pushed to a trigger.",
					},
					"default_guided_failure_mode": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The default guided failure mode setting for the project.",
					},
					"default_to_skip_if_already_installed": datasourceSchema.BoolAttribute{
						Computed:    true,
						Description: "Whether deployment steps should be skipped if the relevant package is already installed.",
					},
					"deployment_changes_template": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The template to use for deployment change details.",
					},
					"deployment_process_id": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The ID of the deployment process associated with this project.",
					},
					"discrete_channel_release": datasourceSchema.BoolAttribute{
						Computed:    true,
						Description: "Treats releases of different channels to the same environment as a separate deployment dimension.",
					},
					"is_disabled": datasourceSchema.BoolAttribute{
						Computed:    true,
						Description: "Whether the project is disabled.",
					},
					"is_discrete_channel_release": datasourceSchema.BoolAttribute{
						Computed:    true,
						Description: "Treats releases of different channels to the same environment as a separate deployment dimension.",
					},
					"is_version_controlled": datasourceSchema.BoolAttribute{
						Computed:    true,
						Description: "Whether the project is version controlled.",
					},
					"lifecycle_id": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The lifecycle ID associated with this project.",
					},
					"project_group_id": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The project group ID associated with this project.",
					},
					"included_library_variable_sets": datasourceSchema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "The list of included library variable set IDs.",
					},
					"tenanted_deployment_participation": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The tenanted deployment mode of the project.",
					},
					"variable_set_id": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The ID of the variable set associated with this project.",
					},
					"release_notes_template": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "The template to use for release notes.",
					},
					"slug": datasourceSchema.StringAttribute{
						Computed:    true,
						Description: "A human-readable, unique identifier, used to identify a project.",
					},
					"connectivity_policy": datasourceSchema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]datasourceSchema.Attribute{
							"allow_deployments_to_no_targets": datasourceSchema.BoolAttribute{
								Computed:    true,
								Description: "Allow deployments to be created when there are no targets.",
							},
							"exclude_unhealthy_targets": datasourceSchema.BoolAttribute{
								Computed:    true,
								Description: "Exclude unhealthy targets from deployments.",
							},
							"skip_machine_behavior": datasourceSchema.StringAttribute{
								Computed:    true,
								Description: "The behavior when a machine is skipped.",
							},
							"target_roles": datasourceSchema.ListAttribute{
								Computed:    true,
								ElementType: types.StringType,
								Description: "The target roles for the connectivity policy.",
							},
						},
						Description: "Defines the connectivity policy for deployments.",
					},
					"git_library_persistence_settings": datasourceSchema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]datasourceSchema.Attribute{
							"url": datasourceSchema.StringAttribute{
								Computed:    true,
								Description: "The URL associated with these version control settings.",
							},
							"base_path": datasourceSchema.StringAttribute{
								Computed:    true,
								Description: "The base path associated with these version control settings.",
							},
							"default_branch": datasourceSchema.StringAttribute{
								Computed:    true,
								Description: "The default branch associated with these version control settings.",
							},
							"protected_branches": datasourceSchema.SetAttribute{
								Computed:    true,
								ElementType: types.StringType,
								Description: "A list of protected branch patterns.",
							},
						},
						Description: "Provides Git-related persistence settings for a version-controlled project.",
					},
					// Note: We're not including username/password or anonymous git settings here as they're typically not returned in read operations
					"versioning_strategy": datasourceSchema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]datasourceSchema.Attribute{
							"donor_package_step_id": datasourceSchema.StringAttribute{
								Computed:    true,
								Description: "The ID of the step containing the donor package.",
							},
							"template": datasourceSchema.StringAttribute{
								Computed:    true,
								Description: "The template to use for version numbers.",
							},
						},
						Description: "Defines the versioning strategy for the project.",
					},
				},
			},
			Description: "The list of projects that match the filter criteria.",
		},
	}
}
