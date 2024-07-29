package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ProjectResourceName = "project"

func getOptionalComputedStringAttribute() resourceSchema.StringAttribute {
	return resourceSchema.StringAttribute{
		Optional: true,
		Computed: true,
	}
}

func getOptionalComputedBoolAttribute(description string) resourceSchema.BoolAttribute {
	attr := resourceSchema.BoolAttribute{
		Optional: true,
		Computed: true,
	}
	if description != "" {
		attr.Description = description
	}
	return attr
}

func getOptionalStringAttribute(description string) resourceSchema.StringAttribute {
	return resourceSchema.StringAttribute{
		Optional:    true,
		Description: description,
	}
}

func getRequiredStringAttribute(description string) resourceSchema.StringAttribute {
	return resourceSchema.StringAttribute{
		Required:    true,
		Description: description,
	}
}

func getResourceOptionalStringListAttribute(description string) resourceSchema.ListAttribute {
	attr := resourceSchema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
	}
	if description != "" {
		attr.Description = description
	}

	return attr
}

func GetProjectResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages projects in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":                                   util.GetIdResourceSchema(),
			"space_id":                             util.GetSpaceIdResourceSchema(ProjectResourceName),
			"name":                                 util.GetNameResourceSchema(true),
			"description":                          util.GetDescriptionResourceSchema(ProjectResourceName),
			"auto_create_release":                  getOptionalComputedBoolAttribute(""),
			"cloned_from_project_id":               getOptionalStringAttribute(""),
			"default_guided_failure_mode":          getOptionalComputedStringAttribute(),
			"default_to_skip_if_already_installed": getOptionalComputedBoolAttribute(""),
			"deployment_changes_template":          getOptionalComputedStringAttribute(),
			"discrete_channel_release":             getOptionalComputedBoolAttribute("Treats releases of different channels to the same environment as a separate deployment dimension"),
			"is_disabled":                          getOptionalComputedBoolAttribute(""),
			"is_discrete_channel_release":          getOptionalComputedBoolAttribute("Treats releases of different channels to the same environment as a separate deployment dimension"),
			"is_version_controlled":                getOptionalComputedBoolAttribute(""),
			"lifecycle_id":                         getRequiredStringAttribute("The lifecycle ID associated with this project."),
			"project_group_id":                     getRequiredStringAttribute("The project group ID associated with this project."),
			"tenanted_deployment_participation":    getOptionalComputedStringAttribute(),
			"included_library_variable_sets":       getResourceOptionalStringListAttribute(""),
			"release_notes_template":               getOptionalComputedStringAttribute(),
			"slug":                                 getOptionalComputedStringAttribute(),
			"deployment_process_id":                getOptionalComputedStringAttribute(),
			"variable_set_id":                      getOptionalComputedStringAttribute(),
		},
		Blocks: map[string]resourceSchema.Block{
			"auto_deploy_release_overrides": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"environment_id": getOptionalStringAttribute(""),
						"release_id":     getOptionalStringAttribute(""),
						"tenant_id":      getOptionalStringAttribute(""),
					},
				},
			},
			"connectivity_policy": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"allow_deployments_to_no_targets": getOptionalComputedBoolAttribute(""),
						"exclude_unhealthy_targets":       getOptionalComputedBoolAttribute(""),
						"skip_machine_behavior":           getOptionalStringAttribute(""),
						"target_roles":                    getResourceOptionalStringListAttribute(""),
					},
				},
			},
			"git_anonymous_persistence_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"url":                getOptionalStringAttribute("The URL associated with these version control settings."),
						"base_path":          getOptionalStringAttribute("The base path associated with these version control settings."),
						"default_branch":     getOptionalStringAttribute("The default branch associated with these version control settings."),
						"protected_branches": getResourceOptionalStringListAttribute("A list of protected branch patterns."),
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"git_library_persistence_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"git_credential_id":  getOptionalStringAttribute(""),
						"url":                getOptionalStringAttribute("The URL associated with these version control settings."),
						"base_path":          getOptionalStringAttribute("The base path associated with these version control settings."),
						"default_branch":     getOptionalStringAttribute("The default branch associated with these version control settings."),
						"protected_branches": getResourceOptionalStringListAttribute("A list of protected branch patterns."),
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"git_username_password_persistence_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"url":                getOptionalStringAttribute("The URL associated with these version control settings."),
						"username":           getOptionalStringAttribute("The username for the Git credential."),
						"password":           util.GetPasswordResourceSchema(false),
						"base_path":          getOptionalStringAttribute("The base path associated with these version control settings."),
						"default_branch":     getOptionalStringAttribute("The default branch associated with these version control settings."),
						"protected_branches": getResourceOptionalStringListAttribute("A list of protected branch patterns."),
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"jira_service_management_extension_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"connection_id":             getOptionalStringAttribute("The connection identifier associated with the extension settings."),
						"is_enabled":                getOptionalComputedBoolAttribute("Specifies whether or not this extension is enabled for this project."),
						"service_desk_project_name": getOptionalStringAttribute("The project name associated with this extension."),
					},
				},
				Description: "Provides extension settings for the Jira Service Management (JSM) integration for this project.",
			},
			"servicenow_extension_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"connection_id":                       getOptionalStringAttribute("The connection identifier associated with the extension settings."),
						"is_enabled":                          getOptionalComputedBoolAttribute("Specifies whether or not this extension is enabled for this project."),
						"is_state_automatically_transitioned": getOptionalComputedBoolAttribute("Specifies whether or not this extension will automatically transition the state of a deployment for this project."),
						"standard_change_template_name":       getOptionalStringAttribute("The name of the standard change template associated with this extension."),
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
						"name":          getRequiredStringAttribute("The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods."),
						"label":         getOptionalStringAttribute("The label shown beside the parameter when presented in the deployment process."),
						"help_text":     getOptionalStringAttribute("The help presented alongside the parameter input."),
						"default_value": getOptionalStringAttribute("A default value for the parameter, if applicable."),
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
						"donor_package_step_id": getOptionalStringAttribute(""),
						"template":              getOptionalStringAttribute(""),
					},
					Blocks: map[string]resourceSchema.Block{
						"donor_package": resourceSchema.ListNestedBlock{
							NestedObject: resourceSchema.NestedBlockObject{
								Attributes: map[string]resourceSchema.Attribute{
									"deployment_action": getOptionalStringAttribute(""),
									"package_reference": getOptionalStringAttribute(""),
								},
							},
						},
					},
				},
			},
			"release_creation_strategy": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"channel_id":                       getOptionalStringAttribute(""),
						"release_creation_package_step_id": getOptionalStringAttribute(""),
					},
					Blocks: map[string]resourceSchema.Block{
						"release_creation_package": resourceSchema.ListNestedBlock{
							NestedObject: resourceSchema.NestedBlockObject{
								Attributes: map[string]resourceSchema.Attribute{
									"deployment_action": getOptionalStringAttribute(""),
									"package_reference": getOptionalStringAttribute(""),
								},
							},
						},
					},
				},
			},
		},
	}
}

func GetProjectDataSourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing Octopus Deploy projects.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":                     util.GetIdDatasourceSchema(),
			"cloned_from_project_id": getDataSourceStringAttribute("The ID of the project this project was cloned from.", true),
			"ids":                    util.GetQueryIDsDatasourceSchema(),
			"is_clone":               getDataSourceBoolAttribute("If set, only return projects that are clones.", true),
			"name":                   util.GetNameDatasourceSchema(true),
			"partial_name":           util.GetQueryPartialNameDatasourceSchema(),
			"skip":                   util.GetQuerySkipDatasourceSchema(),
			"space_id":               util.GetSpaceIdDatasourceSchema(ProjectResourceName),
			"take":                   util.GetQueryTakeDatasourceSchema(),
			"projects":               getProjectsDataSourceAttribute(),
		},
	}
}

func getProjectsDataSourceAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "A list of projects that match the filter(s).",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"allow_deployments_to_no_targets":            getDataSourceBoolAttribute("Deprecated: Whether deployments can be created to no targets.", false),
				"auto_create_release":                        getDataSourceBoolAttribute("Whether to automatically create a release when a package is pushed to a trigger.", false),
				"auto_deploy_release_overrides":              getDataSourceListAttribute("A list of release overrides for auto deployments.", types.StringType),
				"cloned_from_project_id":                     getDataSourceStringAttribute("The ID of the project this project was cloned from.", false),
				"default_guided_failure_mode":                getDataSourceStringAttribute("The default guided failure mode setting for the project.", false),
				"default_to_skip_if_already_installed":       getDataSourceBoolAttribute("Whether deployment steps should be skipped if the relevant package is already installed.", false),
				"deployment_changes_template":                getDataSourceStringAttribute("The template to use for deployment change details.", false),
				"deployment_process_id":                      getDataSourceStringAttribute("The ID of the deployment process associated with this project.", false),
				"description":                                getDataSourceStringAttribute("The description of this project.", false),
				"discrete_channel_release":                   getDataSourceBoolAttribute("Treats releases of different channels to the same environment as a separate deployment dimension.", false),
				"id":                                         getDataSourceStringAttribute("The unique ID of the project.", false),
				"included_library_variable_sets":             getDataSourceListAttribute("The list of included library variable set IDs.", types.StringType),
				"is_disabled":                                getDataSourceBoolAttribute("Whether the project is disabled.", false),
				"is_discrete_channel_release":                getDataSourceBoolAttribute("Treats releases of different channels to the same environment as a separate deployment dimension.", false),
				"is_version_controlled":                      getDataSourceBoolAttribute("Whether the project is version controlled.", false),
				"lifecycle_id":                               getDataSourceStringAttribute("The lifecycle ID associated with this project.", false),
				"name":                                       getDataSourceStringAttribute("The name of the project.", false),
				"project_group_id":                           getDataSourceStringAttribute("The project group ID associated with this project.", false),
				"release_notes_template":                     getDataSourceStringAttribute("The template to use for release notes.", false),
				"slug":                                       getDataSourceStringAttribute("A human-readable, unique identifier, used to identify a project.", false),
				"space_id":                                   getDataSourceStringAttribute("The space ID associated with this project.", false),
				"tenanted_deployment_participation":          getDataSourceStringAttribute("The tenanted deployment mode of the project.", false),
				"variable_set_id":                            getDataSourceStringAttribute("The ID of the variable set associated with this project.", false),
				"connectivity_policy":                        getDataSourceConnectivityPolicyAttribute(),
				"git_library_persistence_settings":           getDataSourceGitPersistenceSettingsAttribute("Git-related persistence settings for a version-controlled project.", true),
				"git_username_password_persistence_settings": getDataSourceGitPersistenceSettingsAttribute("Git-related persistence settings for a version-controlled project using username/password authentication.", false),
				"git_anonymous_persistence_settings":         getDataSourceGitPersistenceSettingsAttribute("Git-related persistence settings for a version-controlled project using anonymous authentication.", false),
				"jira_service_management_extension_settings": getDataSourceJSMExtensionSettingsAttribute(),
				"servicenow_extension_settings":              getDataSourceServiceNowExtensionSettingsAttribute(),
				"versioning_strategy":                        getDataSourceVersioningStrategyAttribute(),
				"release_creation_strategy":                  getDataSourceReleaseCreationStrategyAttribute(),
				"template":                                   getDataSourceTemplateAttribute(),
			},
		},
	}
}

func getDataSourceStringAttribute(description string, optional bool) datasourceSchema.StringAttribute {
	attribute := datasourceSchema.StringAttribute{
		Description: description,
		Computed:    true,
	}
	if optional {
		attribute.Optional = true
	}
	return attribute
}

func getDataSourceBoolAttribute(description string, optional bool) datasourceSchema.BoolAttribute {
	attribute := datasourceSchema.BoolAttribute{
		Description: description,
		Computed:    true,
	}
	if optional {
		attribute.Optional = true
	}
	return attribute
}

func getDataSourceListAttribute(description string, elementType attr.Type) datasourceSchema.ListAttribute {
	return datasourceSchema.ListAttribute{
		Description: description,
		Computed:    true,
		ElementType: elementType,
	}
}

func getDataSourceConnectivityPolicyAttribute() datasourceSchema.SingleNestedAttribute {
	return datasourceSchema.SingleNestedAttribute{
		Description: "Defines the connectivity policy for deployments.",
		Computed:    true,
		Attributes: map[string]datasourceSchema.Attribute{
			"allow_deployments_to_no_targets": getDataSourceBoolAttribute("Allow deployments to be created when there are no targets.", false),
			"exclude_unhealthy_targets":       getDataSourceBoolAttribute("Exclude unhealthy targets from deployments.", false),
			"skip_machine_behavior":           getDataSourceStringAttribute("The behavior when a machine is skipped.", false),
			"target_roles":                    getDataSourceListAttribute("The target roles for the connectivity policy.", types.StringType),
		},
	}
}

func getDataSourceGitPersistenceSettingsAttribute(description string, includeCredential bool) datasourceSchema.SingleNestedAttribute {
	attributes := map[string]datasourceSchema.Attribute{
		"base_path":      getDataSourceStringAttribute("The base path associated with these version control settings.", false),
		"default_branch": getDataSourceStringAttribute("The default branch associated with these version control settings.", false),
		"protected_branches": datasourceSchema.SetAttribute{
			Description: "A list of protected branch patterns.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"url": getDataSourceStringAttribute("The URL associated with these version control settings.", false),
	}

	if includeCredential {
		attributes["git_credential_id"] = getDataSourceStringAttribute("The ID of the Git credential.", false)
	} else {
		attributes["username"] = getDataSourceStringAttribute("The username for the Git credential.", false)
		attributes["password"] = datasourceSchema.StringAttribute{
			Description: "The password for the Git credential.",
			Computed:    true,
			Sensitive:   true,
		}
	}

	return datasourceSchema.SingleNestedAttribute{
		Description: description,
		Computed:    true,
		Attributes:  attributes,
	}
}

func getDataSourceJSMExtensionSettingsAttribute() datasourceSchema.SingleNestedAttribute {
	return datasourceSchema.SingleNestedAttribute{
		Description: "Extension settings for the Jira Service Management (JSM) integration.",
		Computed:    true,
		Attributes: map[string]datasourceSchema.Attribute{
			"connection_id":             getDataSourceStringAttribute("The connection identifier for JSM.", false),
			"is_enabled":                getDataSourceBoolAttribute("Whether the JSM extension is enabled.", false),
			"service_desk_project_name": getDataSourceStringAttribute("The JSM service desk project name.", false),
		},
	}
}

func getDataSourceServiceNowExtensionSettingsAttribute() datasourceSchema.SingleNestedAttribute {
	return datasourceSchema.SingleNestedAttribute{
		Description: "Extension settings for the ServiceNow integration.",
		Computed:    true,
		Attributes: map[string]datasourceSchema.Attribute{
			"connection_id":                       getDataSourceStringAttribute("The connection identifier for ServiceNow.", false),
			"is_enabled":                          getDataSourceBoolAttribute("Whether the ServiceNow extension is enabled.", false),
			"is_state_automatically_transitioned": getDataSourceBoolAttribute("Whether state is automatically transitioned in ServiceNow.", false),
			"standard_change_template_name":       getDataSourceStringAttribute("The name of the standard change template in ServiceNow.", false),
		},
	}
}

func getDataSourceVersioningStrategyAttribute() datasourceSchema.SingleNestedAttribute {
	return datasourceSchema.SingleNestedAttribute{
		Description: "The versioning strategy for the project.",
		Computed:    true,
		Attributes: map[string]datasourceSchema.Attribute{
			"donor_package_step_id": getDataSourceStringAttribute("The ID of the step containing the donor package.", false),
			"donor_package": datasourceSchema.SingleNestedAttribute{
				Description: "Details of the donor package.",
				Computed:    true,
				Attributes: map[string]datasourceSchema.Attribute{
					"deployment_action": getDataSourceStringAttribute("The deployment action for the donor package.", false),
					"package_reference": getDataSourceStringAttribute("The package reference for the donor package.", false),
				},
			},
			"template": getDataSourceStringAttribute("The template to use for version numbers.", false),
		},
	}
}

func getDataSourceReleaseCreationStrategyAttribute() datasourceSchema.SingleNestedAttribute {
	return datasourceSchema.SingleNestedAttribute{
		Description: "The release creation strategy for the project.",
		Computed:    true,
		Attributes: map[string]datasourceSchema.Attribute{
			"channel_id": getDataSourceStringAttribute("The ID of the channel to use for release creation.", false),
			"release_creation_package": datasourceSchema.SingleNestedAttribute{
				Description: "Details of the package used for release creation.",
				Computed:    true,
				Attributes: map[string]datasourceSchema.Attribute{
					"deployment_action": getDataSourceStringAttribute("The deployment action for the release creation package.", false),
					"package_reference": getDataSourceStringAttribute("The package reference for the release creation package.", false),
				},
			},
			"release_creation_package_step_id": getDataSourceStringAttribute("The ID of the step containing the package for release creation.", false),
		},
	}
}

func getDataSourceTemplateAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "Template parameters for the project.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"name":          getDataSourceStringAttribute("The name of the variable set by the parameter.", false),
				"label":         getDataSourceStringAttribute("The label shown beside the parameter.", false),
				"help_text":     getDataSourceStringAttribute("The help text for the parameter.", false),
				"default_value": getDataSourceStringAttribute("The default value for the parameter.", false),
				"display_settings": datasourceSchema.MapAttribute{
					Description: "The display settings for the parameter.",
					Computed:    true,
					ElementType: types.StringType,
				},
			},
		},
	}
}
