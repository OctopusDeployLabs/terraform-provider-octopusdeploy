package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ProjectResourceName = "project"
const ProjectDataSourceName = "projects"

func GetProjectResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages projects in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":                                   util.GetIdResourceSchema(),
			"space_id":                             util.GetSpaceIdResourceSchema(ProjectResourceName),
			"name":                                 util.GetNameResourceSchema(true),
			"description":                          util.GetDescriptionResourceSchema(ProjectResourceName),
			"auto_create_release":                  util.Bool().Optional().Computed().Build(),
			"cloned_from_project_id":               util.String().Optional().Description("The ID of the project this project was cloned from.").Build(),
			"default_guided_failure_mode":          util.String().Optional().Computed().Build(),
			"default_to_skip_if_already_installed": util.Bool().Optional().Computed().Build(),
			"deployment_changes_template":          util.String().Optional().Computed().Build(),
			"discrete_channel_release":             util.Bool().Optional().Computed().Description("Treats releases of different channels to the same environment as a separate deployment dimension").Build(),
			"is_disabled":                          util.Bool().Optional().Computed().Build(),
			"is_discrete_channel_release":          util.Bool().Optional().Computed().Description("Treats releases of different channels to the same environment as a separate deployment dimension").Build(),
			"is_version_controlled":                util.Bool().Optional().Computed().Build(),
			"lifecycle_id":                         util.String().Required().Description("The lifecycle ID associated with this project.").Build(),
			"project_group_id":                     util.String().Required().Description("The project group ID associated with this project.").Build(),
			"tenanted_deployment_participation":    util.String().Optional().Computed().Build(),
			"included_library_variable_sets":       util.List(types.StringType).Optional().Computed().Description("The list of included library variable set IDs.").Build(),
			"release_notes_template":               util.String().Optional().Computed().Build(),
			"slug":                                 util.String().Optional().Computed().Build(),
			"deployment_process_id":                util.String().Optional().Computed().Build(),
			"variable_set_id":                      util.String().Optional().Computed().Build(),
		},
		Blocks: map[string]resourceSchema.Block{
			"auto_deploy_release_overrides": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"environment_id": util.String().Optional().Build(),
						"release_id":     util.String().Optional().Build(),
						"tenant_id":      util.String().Optional().Build(),
					},
				},
			},
			"connectivity_policy": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"allow_deployments_to_no_targets": util.Bool().Optional().Computed().Build(),
						"exclude_unhealthy_targets":       util.Bool().Optional().Computed().Build(),
						"skip_machine_behavior":           util.String().Optional().Build(),
						"target_roles":                    util.List(types.StringType).Optional().Computed().Build(),
					},
				},
			},
			"git_anonymous_persistence_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"url":                util.String().Optional().Description("The URL associated with these version control settings.").Build(),
						"base_path":          util.String().Optional().Description("The base path associated with these version control settings.").Build(),
						"default_branch":     util.String().Optional().Description("The default branch associated with these version control settings.").Build(),
						"protected_branches": util.Set(types.StringType).Optional().Description("A list of protected branch patterns.").Build(),
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"git_library_persistence_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"git_credential_id":  util.String().Optional().Build(),
						"url":                util.String().Optional().Description("The URL associated with these version control settings.").Build(),
						"base_path":          util.String().Optional().Description("The base path associated with these version control settings.").Build(),
						"default_branch":     util.String().Optional().Description("The default branch associated with these version control settings.").Build(),
						"protected_branches": util.Set(types.StringType).Optional().Description("A list of protected branch patterns.").Build(),
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"git_username_password_persistence_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"url":                util.String().Optional().Description("The URL associated with these version control settings.").Build(),
						"username":           util.String().Optional().Description("The username for the Git credential.").Build(),
						"password":           util.GetPasswordResourceSchema(false),
						"base_path":          util.String().Optional().Description("The base path associated with these version control settings.").Build(),
						"default_branch":     util.String().Optional().Description("The default branch associated with these version control settings.").Build(),
						"protected_branches": util.Set(types.StringType).Optional().Description("A list of protected branch patterns.").Build(),
					},
				},
				Description: "Provides Git-related persistence settings for a version-controlled project.",
			},
			"jira_service_management_extension_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"connection_id":             util.String().Optional().Description("The connection identifier associated with the extension settings.").Build(),
						"is_enabled":                util.Bool().Optional().Computed().Description("Specifies whether or not this extension is enabled for this project.").Build(),
						"service_desk_project_name": util.String().Optional().Description("The project name associated with this extension.").Build(),
					},
				},
				Description: "Provides extension settings for the Jira Service Management (JSM) integration for this project.",
			},
			"servicenow_extension_settings": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"connection_id":                       util.String().Optional().Description("The connection identifier associated with the extension settings.").Build(),
						"is_enabled":                          util.Bool().Optional().Computed().Description("Specifies whether or not this extension is enabled for this project.").Build(),
						"is_state_automatically_transitioned": util.Bool().Optional().Computed().Description("Specifies whether or not this extension will automatically transition the state of a deployment for this project.").Build(),
						"standard_change_template_name":       util.String().Optional().Description("The name of the standard change template associated with this extension.").Build(),
					},
				},
				Description: "Provides extension settings for the ServiceNow integration for this project.",
			},
			"template": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"id":            util.String().Optional().Computed().Description("The ID of the template parameter.").Build(),
						"name":          util.String().Required().Description("The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods.").Build(),
						"label":         util.String().Optional().Description("The label shown beside the parameter when presented in the deployment process.").Build(),
						"help_text":     util.String().Optional().Description("The help presented alongside the parameter input.").Build(),
						"default_value": util.String().Optional().Description("A default value for the parameter, if applicable.").Build(),
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
						"donor_package_step_id": util.String().Optional().Build(),
						"template":              util.String().Optional().Build(),
					},
					Blocks: map[string]resourceSchema.Block{
						"donor_package": resourceSchema.ListNestedBlock{
							NestedObject: resourceSchema.NestedBlockObject{
								Attributes: map[string]resourceSchema.Attribute{
									"deployment_action": util.String().Optional().Build(),
									"package_reference": util.String().Optional().Build(),
								},
							},
						},
					},
				},
			},
			"release_creation_strategy": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: map[string]resourceSchema.Attribute{
						"channel_id":                       util.String().Optional().Build(),
						"release_creation_package_step_id": util.String().Optional().Build(),
					},
					Blocks: map[string]resourceSchema.Block{
						"release_creation_package": resourceSchema.ListNestedBlock{
							NestedObject: resourceSchema.NestedBlockObject{
								Attributes: map[string]resourceSchema.Attribute{
									"deployment_action": util.String().Optional().Build(),
									"package_reference": util.String().Optional().Build(),
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
			"cloned_from_project_id": util.DataSourceString().Optional().Description("The ID of the project this project was cloned from.").Build(),
			"ids":                    util.GetQueryIDsDatasourceSchema(),
			"is_clone":               util.DataSourceBool().Optional().Description("If set, only return projects that are clones.").Build(),
			"name":                   util.GetNameDatasourceSchema(false),
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
				"auto_create_release":                        util.DataSourceBool().Computed().Description("Whether to automatically create a release when a package is pushed to a trigger.").Build(),
				"auto_deploy_release_overrides":              getAutoDeployReleaseOverrides(),
				"cloned_from_project_id":                     util.DataSourceString().Computed().Description("The ID of the project this project was cloned from.").Build(),
				"default_guided_failure_mode":                util.DataSourceString().Computed().Description("The default guided failure mode setting for the project.").Build(),
				"default_to_skip_if_already_installed":       util.DataSourceBool().Computed().Description("Whether deployment steps should be skipped if the relevant package is already installed.").Build(),
				"deployment_changes_template":                util.DataSourceString().Computed().Description("The template to use for deployment change details.").Build(),
				"deployment_process_id":                      util.DataSourceString().Computed().Description("The ID of the deployment process associated with this project.").Build(),
				"description":                                util.DataSourceString().Computed().Description("The description of this project.").Build(),
				"discrete_channel_release":                   util.DataSourceBool().Computed().Description("Treats releases of different channels to the same environment as a separate deployment dimension.").Build(),
				"id":                                         util.DataSourceString().Computed().Description("The unique ID of the project.").Build(),
				"included_library_variable_sets":             util.DataSourceList(types.StringType).Computed().Description("The list of included library variable set IDs.").Build(),
				"is_disabled":                                util.DataSourceBool().Computed().Description("Whether the project is disabled.").Build(),
				"is_discrete_channel_release":                util.DataSourceBool().Computed().Description("Treats releases of different channels to the same environment as a separate deployment dimension.").Build(),
				"is_version_controlled":                      util.DataSourceBool().Computed().Description("Whether the project is version controlled.").Build(),
				"lifecycle_id":                               util.DataSourceString().Computed().Description("The lifecycle ID associated with this project.").Build(),
				"name":                                       util.DataSourceString().Computed().Description("The name of the project.").Build(),
				"project_group_id":                           util.DataSourceString().Computed().Description("The project group ID associated with this project.").Build(),
				"release_notes_template":                     util.DataSourceString().Computed().Description("The template to use for release notes.").Build(),
				"slug":                                       util.DataSourceString().Computed().Description("A human-readable, unique identifier, used to identify a project.").Build(),
				"space_id":                                   util.DataSourceString().Computed().Description("The space ID associated with this project.").Build(),
				"tenanted_deployment_participation":          util.DataSourceString().Computed().Description("The tenanted deployment mode of the project.").Build(),
				"variable_set_id":                            util.DataSourceString().Computed().Description("The ID of the variable set associated with this project.").Build(),
				"connectivity_policy":                        getDataSourceConnectivityPolicyAttribute(),
				"git_library_persistence_settings":           getDataSourceGitPersistenceSettingsAttribute("library"),
				"git_username_password_persistence_settings": getDataSourceGitPersistenceSettingsAttribute("username_password"),
				"git_anonymous_persistence_settings":         getDataSourceGitPersistenceSettingsAttribute("anonymous"),
				"jira_service_management_extension_settings": getDataSourceJSMExtensionSettingsAttribute(),
				"servicenow_extension_settings":              getDataSourceServiceNowExtensionSettingsAttribute(),
				"versioning_strategy":                        getDataSourceVersioningStrategyAttribute(),
				"release_creation_strategy":                  getDataSourceReleaseCreationStrategyAttribute(),
				"template":                                   getDataSourceTemplateAttribute(),
			},
		},
	}
}

func getAutoDeployReleaseOverrides() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "A list of release overrides for auto deployments.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"environment_id": util.DataSourceString().Computed().Description("The environment ID for the auto deploy release override.").Build(),
				"release_id":     util.DataSourceString().Computed().Description("The release ID for the auto deploy release override.").Build(),
				"tenant_id":      util.DataSourceString().Computed().Description("The tenant ID for the auto deploy release override.").Build(),
			},
		},
	}
}

func getDataSourceConnectivityPolicyAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "Defines the connectivity policy for deployments.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"allow_deployments_to_no_targets": util.DataSourceBool().Computed().Description("Allow deployments to be created when there are no targets.").Build(),
				"exclude_unhealthy_targets":       util.DataSourceBool().Computed().Description("Exclude unhealthy targets from deployments.").Build(),
				"skip_machine_behavior":           util.DataSourceString().Computed().Description("The behavior when a machine is skipped.").Build(),
				"target_roles":                    util.DataSourceList(types.StringType).Computed().Description("The target roles for the connectivity policy.").Build(),
			},
		},
	}
}

func getDataSourceGitPersistenceSettingsAttribute(settingType string) datasourceSchema.ListNestedAttribute {
	attributes := map[string]datasourceSchema.Attribute{
		"base_path":          util.DataSourceString().Computed().Description("The base path associated with these version control settings.").Build(),
		"default_branch":     util.DataSourceString().Computed().Description("The default branch associated with these version control settings.").Build(),
		"protected_branches": util.DataSourceSet(types.StringType).Computed().Description("A list of protected branch patterns.").Build(),
		"url":                util.DataSourceString().Computed().Description("The URL associated with these version control settings.").Build(),
	}

	switch settingType {
	case "library":
		attributes["git_credential_id"] = util.DataSourceString().Computed().Description("The ID of the Git credential.").Build()
	case "username_password":
		attributes["username"] = util.DataSourceString().Computed().Description("The username for the Git credential.").Build()
		attributes["password"] = util.DataSourceString().Computed().Sensitive().Description("The password for the Git credential.").Build()
	case "anonymous":
		// No additional attributes for anonymous
	}

	return datasourceSchema.ListNestedAttribute{
		Description: "Git-related persistence settings for a version-controlled project using " + settingType + " authentication.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: attributes,
		},
	}
}

func getDataSourceJSMExtensionSettingsAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "Extension settings for the Jira Service Management (JSM) integration.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"connection_id":             util.DataSourceString().Computed().Description("The connection identifier for JSM.").Build(),
				"is_enabled":                util.DataSourceBool().Computed().Description("Whether the JSM extension is enabled.").Build(),
				"service_desk_project_name": util.DataSourceString().Computed().Description("The JSM service desk project name.").Build(),
			},
		},
	}
}

func getDataSourceServiceNowExtensionSettingsAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "Extension settings for the ServiceNow integration.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"connection_id":                       util.DataSourceString().Computed().Description("The connection identifier for ServiceNow.").Build(),
				"is_enabled":                          util.DataSourceBool().Computed().Description("Whether the ServiceNow extension is enabled.").Build(),
				"is_state_automatically_transitioned": util.DataSourceBool().Computed().Description("Whether state is automatically transitioned in ServiceNow.").Build(),
				"standard_change_template_name":       util.DataSourceString().Computed().Description("The name of the standard change template in ServiceNow.").Build(),
			},
		},
	}
}

func getDataSourceVersioningStrategyAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "The versioning strategy for the project.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"donor_package_step_id": util.DataSourceString().Computed().Description("The ID of the step containing the donor package.").Build(),
				"donor_package": datasourceSchema.ListNestedAttribute{
					Computed: true,
					NestedObject: datasourceSchema.NestedAttributeObject{
						Attributes: map[string]datasourceSchema.Attribute{
							"deployment_action": util.DataSourceString().Computed().Description("The deployment action for the donor package.").Build(),
							"package_reference": util.DataSourceString().Computed().Description("The package reference for the donor package.").Build(),
						},
					},
				},
				"template": util.DataSourceString().Computed().Description("The template to use for version numbers.").Build(),
			},
		},
	}
}

func getDataSourceReleaseCreationStrategyAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "The release creation strategy for the project.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"channel_id": util.DataSourceString().Computed().Description("The ID of the channel to use for release creation.").Build(),
				"release_creation_package": datasourceSchema.ListNestedAttribute{
					Description: "Details of the package used for release creation.",
					Computed:    true,
					NestedObject: datasourceSchema.NestedAttributeObject{
						Attributes: map[string]datasourceSchema.Attribute{
							"deployment_action": util.DataSourceString().Computed().Description("The deployment action for the release creation package.").Build(),
							"package_reference": util.DataSourceString().Computed().Description("The package reference for the release creation package.").Build(),
						},
					},
				},
				"release_creation_package_step_id": util.DataSourceString().Computed().Description("The ID of the step containing the package for release creation.").Build(),
			},
		},
	}
}

func getDataSourceTemplateAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "Template parameters for the project.",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"id":               util.DataSourceString().Computed().Description("The ID of the template parameter.").Build(),
				"name":             util.DataSourceString().Computed().Description("The name of the variable set by the parameter.").Build(),
				"label":            util.DataSourceString().Computed().Description("The label shown beside the parameter.").Build(),
				"help_text":        util.DataSourceString().Computed().Description("The help text for the parameter.").Build(),
				"default_value":    util.DataSourceString().Computed().Description("The default value for the parameter.").Build(),
				"display_settings": util.DataSourceMap(types.StringType).Computed().Description("The display settings for the parameter.").Build(),
			},
		},
	}
}
