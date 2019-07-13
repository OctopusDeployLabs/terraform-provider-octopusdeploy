package octopusdeploy

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deployment_process_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_failure_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "EnvironmentDefault",
				ValidateFunc: validateValueFunc([]string{
					"EnvironmentDefault",
					"Off",
					"On",
				}),
			},
			"skip_machine_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "None",
				ValidateFunc: validateValueFunc([]string{
					"SkipUnavailableMachines",
					"None",
				}),
			},
			"allow_deployments_to_no_targets": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tenanted_deployment_mode": getTenantedDeploymentSchema(),
			"included_library_variable_sets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"discrete_channel_release": {
				Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"skip_package_steps_that_are_already_installed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"deployment_step_windows_service": getDeploymentStepWindowsServiceSchema(),
			"deployment_step_iis_website":     getDeploymentStepIISWebsiteSchema(),
			"deployment_step_inline_script":   getDeploymentStepInlineScriptSchema(),
			"deployment_step_kubernetes_helm": getDeploymentStepKubernetesHelmSchema(),
			"deployment_step_kubernetes_yaml": getDeploymentStepKubernetesYamlSchema(),
			"deployment_step_package_script":  getDeploymentStepPackageScriptSchema(),
			"deployment_step_apply_terraform": getDeploymentStepApplyTerraformSchema(),
		},
	}
}

// addFeedAndPackageDeploymentStepSchema adds schemas related packages and feeds
func addFeedAndPackageDeploymentStepSchema(schemaToAddToo interface{}) *schema.Resource {
	schemaResource := schemaToAddToo.(*schema.Resource)

	schemaResource.Schema["feed_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The ID of the feed a package will be found in.",
		Optional:    true,
		Default:     "feeds-builtin",
	}

	schemaResource.Schema["package"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "ID / Name of the package to be deployed.",
		Required:    true,
	}

	return schemaResource
}

// addConfigurationTransformDeploymentStepSchema adds schemas related to modifying configuration files
func addConfigurationTransformDeploymentStepSchema(schemaToAddToo interface{}) *schema.Resource {
	schemaResource := schemaToAddToo.(*schema.Resource)

	schemaResource.Schema["configuration_transforms"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Enables XML configuration transformations.",
		Optional:    true,
		Default:     true,
	}

	schemaResource.Schema["configuration_variables"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Enables replacing appSettings and connectionString entries in any .config file.",
		Optional:    true,
		Default:     true,
	}

	schemaResource.Schema["json_file_variable_replacement"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A comma-separated list of file names to replace settings in, relative to the package contents.",
	}

	schemaResource.Schema["variable_substitution_in_files"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A newline-separated list of file names to transform, relative to the package contents. Extended wildcard syntax is supported.",
	}

	return schemaResource
}

// addStandardDeploymentStepSchema adds the common schema for Octopus Deploy Steps
func addStandardDeploymentStepSchema(schemaToAddToo interface{}, requireRole bool) *schema.Resource {
	schemaResource := schemaToAddToo.(*schema.Resource)
	schemaResource.Schema["step_condition"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Limit when this step will run by setting this condition.",
		Optional:    true,
		ValidateFunc: validateValueFunc([]string{
			"success",
			"failure",
			"always",
			"variable",
		}),
		Default: "success",
	}

	schemaResource.Schema["step_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The name of the deployment step.",
		Required:    true,
	}

	schemaResource.Schema["step_start_trigger"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "StartAfterPrevious",
		Description: "Control whether the step waits for the previous step to complete, or runs parallel with it.",
		ValidateFunc: validateValueFunc([]string{
			"StartAfterPrevious",
			"StartWithPrevious",
		}),
	}

	schemaResource.Schema["target_roles"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: requireRole,
		Optional: !requireRole,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	return schemaResource
}

// addIISApplicationPoolSchema adds schema for Octopus Deploy Steps needing IIS AppPool configuration
func addIISApplicationPoolSchema(schemaToAddToo interface{}) *schema.Resource {
	schemaResource := schemaToAddToo.(*schema.Resource)

	schemaResource.Schema["application_pool_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Name of the application pool in IIS to create or reconfigure.",
		Required:    true,
	}

	schemaResource.Schema["application_pool_framework"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The version of the .NET common language runtime that this application pool will use. Choose v2.0 for applications built against .NET 2.0, 3.0 or 3.5. Choose v4.0 for .NET 4.0 or 4.5.",
		Default:     "v4.0",
		Optional:    true,
		ValidateFunc: validateValueFunc([]string{
			"v2.0",
			"v4.0",
		}),
	}

	schemaResource.Schema["application_pool_identity"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Which built-in account will the application pool run under.",
		Default:     "ApplicationPoolIdentity",
		ValidateFunc: validateValueFunc([]string{
			"ApplicationPoolIdentity",
			"LocalService",
			"LocalSystem",
			"NetworkService",
		}),
	}

	return schemaResource
}

func getDeploymentStepInlineScriptSchema() *schema.Schema {
	schemaToReturn := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"script_type": {
					Type:        schema.TypeString,
					Description: "The scripting language of the deployment step.",
					Required:    true,
					ValidateFunc: validateValueFunc([]string{
						"PowerShell",
						"CSharp",
						"Bash",
						"FSharp",
					}),
				},
				"script_body": {
					Type:        schema.TypeString,
					Description: "The script body.",
					Required:    true,
				},
				"run_on_server": {
					Type:        schema.TypeBool,
					Description: "Whether the script runs on the server (true) or target (false)",
					Optional:    true,
					Default:     false,
				},
			},
		},
	}

	schemaToReturn.Elem = addStandardDeploymentStepSchema(schemaToReturn.Elem, false)

	return schemaToReturn
}

func getDeploymentStepKubernetesHelmSchema() *schema.Schema {
	schemaToReturn := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"reset_values": {
					Type:        schema.TypeBool,
					Description: "Whether the Helm install can reset values.",
					Optional:    true,
					Default:     true,
				},
				"release_name": {
					Type:        schema.TypeString,
					Description: "The release name of the Helm chart.",
					Required:    true,
				},
				"namespace": {
					Type:        schema.TypeString,
					Description: "The namespace for the Helm chart.",
					Required:    true,
				},
				"yaml_values": {
					Type:        schema.TypeString,
					Description: "The YAML values to pass to the Helm chart.",
					Required:    true,
				},
				"tiller_namespace": {
					Type:        schema.TypeString,
					Description: "The tiller namespace for the Helm chart.",
					Required:    true,
				},
				"package_id": {
					Type:        schema.TypeString,
					Description: "The Package ID of the Helm chart.",
					Required:    true,
				},
				"feed_id": {
					Type:        schema.TypeString,
					Description: "The Feed ID of the Helm chart.",
					Required:    true,
				},
			},
		},
	}

	schemaToReturn.Elem = addStandardDeploymentStepSchema(schemaToReturn.Elem, false)

	return schemaToReturn
}

func getDeploymentStepKubernetesYamlSchema() *schema.Schema {
	schemaToReturn := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"yaml_values": {
					Type:        schema.TypeString,
					Description: "The YAML values to pass to the Helm chart.",
					Required:    true,
				},
				"run_on_server": {
					Type:        schema.TypeBool,
					Description: "Whether the script runs on the server (true) or target (false)",
					Optional:    true,
					Default:     false,
				},
			},
		},
	}

	schemaToReturn.Elem = addStandardDeploymentStepSchema(schemaToReturn.Elem, false)

	return schemaToReturn
}

func getDeploymentStepPackageScriptSchema() *schema.Schema {
	schemaToReturn := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"script_file_name": {
					Type:        schema.TypeString,
					Description: "The script file name in the package.",
					Required:    true,
				},
				"script_parameters": {
					Type:        schema.TypeString,
					Description: "Parameters expected by the script. Use platform specific calling convention. e.g. -Path #{VariableStoringPath} for PowerShell or -- #{VariableStoringPath} for ScriptCS.",
					Optional:    true,
				},
				"run_on_server": {
					Type:        schema.TypeBool,
					Description: "Whether the script runs on the server (true) or target (false)",
					Optional:    true,
					Default:     false,
				},
			},
		},
	}

	schemaToReturn.Elem = addConfigurationTransformDeploymentStepSchema(schemaToReturn.Elem)
	schemaToReturn.Elem = addFeedAndPackageDeploymentStepSchema(schemaToReturn.Elem)
	schemaToReturn.Elem = addStandardDeploymentStepSchema(schemaToReturn.Elem, false)

	return schemaToReturn
}

func getDeploymentStepApplyTerraformSchema() *schema.Schema {
	schemaToReturn := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"additional_init_params": {
					Type:        schema.TypeString,
					Description: "Additional parameters passed to the init command.",
					Optional:    true,
				},
				"run_on_server": {
					Type:        schema.TypeBool,
					Description: "Whether the script runs on the server (true) or target (false)",
					Optional:    true,
					Default:     false,
				},
				"terraform_file_variable_replacement": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}

	schemaToReturn.Elem = addFeedAndPackageDeploymentStepSchema(schemaToReturn.Elem)
	schemaToReturn.Elem = addStandardDeploymentStepSchema(schemaToReturn.Elem, false)

	return schemaToReturn
}

// getDeploymentStepIISWebsiteSchema returns schema for an IIS deployment step
func getDeploymentStepIISWebsiteSchema() *schema.Schema {
	schemaToReturn := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"anonymous_authentication": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Whether IIS should allow anonymous authentication.",
					Default:     false,
				},
				"basic_authentication": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Whether IIS should allow basic authentication with a 401 challenge.",
					Default:     false,
				},
				"website_name": {
					Type:        schema.TypeString,
					Description: "The name of the Website to be created",
					Required:    true,
				},
				"windows_authentication": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Whether IIS should allow integrated Windows authentication with a 401 challenge.",
					Default:     true,
				},
			},
		},
	}

	schemaToReturn.Elem = addConfigurationTransformDeploymentStepSchema(schemaToReturn.Elem)
	schemaToReturn.Elem = addStandardDeploymentStepSchema(schemaToReturn.Elem, true)
	schemaToReturn.Elem = addFeedAndPackageDeploymentStepSchema(schemaToReturn.Elem)
	schemaToReturn.Elem = addIISApplicationPoolSchema(schemaToReturn.Elem)

	return schemaToReturn
}

// getDeploymentStepWindowsServiceSchema returns schema for a Windows Service deployment step
func getDeploymentStepWindowsServiceSchema() *schema.Schema {
	schemaToReturn := &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"executable_path": {
					Type:     schema.TypeString,
					Required: true,
				},
				"service_account": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "LocalSystem",
				},
				"service_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"service_start_mode": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "auto",
					ValidateFunc: validateValueFunc([]string{
						"auto",
						"delayed-auto",
						"demand",
						"unchanged",
					}),
				},
			},
		},
	}

	schemaToReturn.Elem = addFeedAndPackageDeploymentStepSchema(schemaToReturn.Elem)
	schemaToReturn.Elem = addStandardDeploymentStepSchema(schemaToReturn.Elem, true)
	schemaToReturn.Elem = addConfigurationTransformDeploymentStepSchema(schemaToReturn.Elem)

	return schemaToReturn
}

func buildDeploymentProcess(d *schema.ResourceData, deploymentProcess *octopusdeploy.DeploymentProcess) *octopusdeploy.DeploymentProcess {
	deploymentProcess.Steps = nil // empty the steps

	if v, ok := d.GetOk("deployment_step_windows_service"); ok {
		steps := v.([]interface{})
		for _, raw := range steps {

			localStep := raw.(map[string]interface{})

			configurationTransforms := localStep["configuration_transforms"].(bool)
			configurationVariables := localStep["configuration_variables"].(bool)
			executablePath := localStep["executable_path"].(string)
			feedID := localStep["feed_id"].(string)
			jsonFileVariableReplacement := localStep["json_file_variable_replacement"].(string)
			variableSubstitutionInFiles := localStep["variable_substitution_in_files"].(string)
			packageID := localStep["package"].(string)
			serviceAccount := localStep["service_account"].(string)
			serviceName := localStep["service_name"].(string)
			serviceStartMode := localStep["service_start_mode"].(string)
			stepCondition := localStep["step_condition"].(string)
			stepName := localStep["step_name"].(string)
			stepStartTrigger := localStep["step_start_trigger"].(string)

			deploymentStep := &octopusdeploy.DeploymentStep{
				Name:               stepName,
				PackageRequirement: "LetOctopusDecide",
				Condition:          octopusdeploy.DeploymentStepCondition(stepCondition),
				StartTrigger:       octopusdeploy.DeploymentStepStartTrigger(stepStartTrigger),
				Actions: []octopusdeploy.DeploymentAction{
					{
						Name:       stepName,
						ActionType: "Octopus.WindowsService",
						Properties: map[string]string{
							"Octopus.Action.WindowsService.CreateOrUpdateService":                       "True",
							"Octopus.Action.WindowsService.ServiceAccount":                              serviceAccount,
							"Octopus.Action.WindowsService.StartMode":                                   serviceStartMode,
							"Octopus.Action.Package.AutomaticallyRunConfigurationTransformationFiles":   strconv.FormatBool(configurationTransforms),
							"Octopus.Action.Package.AutomaticallyUpdateAppSettingsAndConnectionStrings": strconv.FormatBool(configurationVariables),
							"Octopus.Action.EnabledFeatures":                                            "Octopus.Features.WindowsService,Octopus.Features.ConfigurationTransforms,Octopus.Features.ConfigurationVariables",
							"Octopus.Action.Package.FeedId":                                             feedID,
							"Octopus.Action.Package.PackageId":                                          packageID,
							"Octopus.Action.Package.DownloadOnTentacle":                                 "False",
							"Octopus.Action.WindowsService.ServiceName":                                 serviceName,
							"Octopus.Action.WindowsService.ExecutablePath":                              executablePath,
						},
					},
				},
			}

			if jsonFileVariableReplacement != "" {
				deploymentStep.Actions[0].Properties["Octopus.Action.Package.JsonConfigurationVariablesTargets"] = jsonFileVariableReplacement
				deploymentStep.Actions[0].Properties["Octopus.Action.Package.JsonConfigurationVariablesEnabled"] = "True"

				deploymentStep.Actions[0].Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.JsonConfigurationVariables"
			}

			if variableSubstitutionInFiles != "" {
				deploymentStep.Actions[0].Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] = variableSubstitutionInFiles
				deploymentStep.Actions[0].Properties["Octopus.Action.SubstituteInFiles.Enabled"] = "True"

				deploymentStep.Actions[0].Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.SubstituteInFiles"
			}

			if targetRolesInterface, ok := localStep["target_roles"]; ok {
				var targetRoleSlice []string

				targetRoles := targetRolesInterface.([]interface{})

				for _, role := range targetRoles {
					targetRoleSlice = append(targetRoleSlice, role.(string))
				}

				deploymentStep.Properties = map[string]string{"Octopus.Action.TargetRoles": strings.Join(targetRoleSlice, ",")}
			}

			deploymentProcess.Steps = append(deploymentProcess.Steps, *deploymentStep)
		}
	}

	if v, ok := d.GetOk("deployment_step_iis_website"); ok {
		steps := v.([]interface{})
		for _, raw := range steps {

			localStep := raw.(map[string]interface{})

			anonymousAuthentication := localStep["anonymous_authentication"].(bool)
			applicationPoolFramework := localStep["application_pool_framework"].(string)
			applicationPoolIdentity := localStep["application_pool_identity"].(string)
			applicationPoolName := localStep["application_pool_name"].(string)
			basicAuthentication := localStep["basic_authentication"].(bool)
			configurationTransforms := localStep["configuration_transforms"].(bool)
			configurationVariables := localStep["configuration_variables"].(bool)
			feedID := localStep["feed_id"].(string)
			jsonFileVariableReplacement := localStep["json_file_variable_replacement"].(string)
			variableSubstitutionInFiles := localStep["variable_substitution_in_files"].(string)
			packageID := localStep["package"].(string)
			stepCondition := localStep["step_condition"].(string)
			stepName := localStep["step_name"].(string)
			stepStartTrigger := localStep["step_start_trigger"].(string)
			websiteName := localStep["website_name"].(string)
			windowsAuthentication := localStep["windows_authentication"].(bool)

			deploymentStep := &octopusdeploy.DeploymentStep{
				Name:               stepName,
				PackageRequirement: "LetOctopusDecide",
				Condition:          octopusdeploy.DeploymentStepCondition(stepCondition),
				StartTrigger:       octopusdeploy.DeploymentStepStartTrigger(stepStartTrigger),
				Actions: []octopusdeploy.DeploymentAction{
					{
						Name:       stepName,
						ActionType: "Octopus.IIS",
						Properties: map[string]string{
							"Octopus.Action.IISWebSite.DeploymentType":                                  "webSite",
							"Octopus.Action.IISWebSite.CreateOrUpdateWebSite":                           "True",
							"Octopus.Action.IISWebSite.Bindings":                                        "[{\"protocol\":\"http\",\"port\":\"80\",\"host\":\"\",\"thumbprint\":null,\"certificateVariable\":null,\"requireSni\":false,\"enabled\":true}]",
							"Octopus.Action.IISWebSite.ApplicationPoolFrameworkVersion":                 applicationPoolFramework,
							"Octopus.Action.IISWebSite.ApplicationPoolIdentityType":                     applicationPoolIdentity,
							"Octopus.Action.IISWebSite.EnableAnonymousAuthentication":                   strconv.FormatBool(anonymousAuthentication),
							"Octopus.Action.IISWebSite.EnableBasicAuthentication":                       strconv.FormatBool(basicAuthentication),
							"Octopus.Action.IISWebSite.EnableWindowsAuthentication":                     strconv.FormatBool(windowsAuthentication),
							"Octopus.Action.IISWebSite.WebApplication.ApplicationPoolFrameworkVersion":  applicationPoolFramework,
							"Octopus.Action.IISWebSite.WebApplication.ApplicationPoolIdentityType":      applicationPoolIdentity,
							"Octopus.Action.Package.AutomaticallyRunConfigurationTransformationFiles":   strconv.FormatBool(configurationTransforms),
							"Octopus.Action.Package.AutomaticallyUpdateAppSettingsAndConnectionStrings": strconv.FormatBool(configurationVariables),
							"Octopus.Action.EnabledFeatures":                                            "Octopus.Features.IISWebSite,Octopus.Features.ConfigurationTransforms,Octopus.Features.ConfigurationVariables",
							"Octopus.Action.Package.FeedId":                                             feedID,
							"Octopus.Action.Package.DownloadOnTentacle":                                 "False",
							"Octopus.Action.IISWebSite.WebRootType":                                     "packageRoot",
							"Octopus.Action.IISWebSite.StartApplicationPool":                            "True",
							"Octopus.Action.IISWebSite.StartWebSite":                                    "True",
							"Octopus.Action.Package.PackageId":                                          packageID,
							"Octopus.Action.IISWebSite.WebSiteName":                                     websiteName,
							"Octopus.Action.IISWebSite.ApplicationPoolName":                             applicationPoolName,
						},
					},
				},
			}

			if jsonFileVariableReplacement != "" {
				deploymentStep.Actions[0].Properties["Octopus.Action.Package.JsonConfigurationVariablesTargets"] = jsonFileVariableReplacement
				deploymentStep.Actions[0].Properties["Octopus.Action.Package.JsonConfigurationVariablesEnabled"] = "True"

				deploymentStep.Actions[0].Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.JsonConfigurationVariables"
			}

			if variableSubstitutionInFiles != "" {
				deploymentStep.Actions[0].Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] = variableSubstitutionInFiles
				deploymentStep.Actions[0].Properties["Octopus.Action.SubstituteInFiles.Enabled"] = "True"

				deploymentStep.Actions[0].Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.SubstituteInFiles"
			}

			if targetRolesInterface, ok := localStep["target_roles"]; ok {
				var targetRoleSlice []string

				targetRoles := targetRolesInterface.([]interface{})

				for _, role := range targetRoles {
					targetRoleSlice = append(targetRoleSlice, role.(string))
				}

				deploymentStep.Properties = map[string]string{"Octopus.Action.TargetRoles": strings.Join(targetRoleSlice, ",")}
			}

			deploymentProcess.Steps = append(deploymentProcess.Steps, *deploymentStep)
		}
	}

	if v, ok := d.GetOk("deployment_step_inline_script"); ok {
		steps := v.([]interface{})
		for _, raw := range steps {

			localStep := raw.(map[string]interface{})

			scriptType := localStep["script_type"].(string)
			scriptBody := localStep["script_body"].(string)
			runOnServer := localStep["run_on_server"].(bool)
			stepCondition := localStep["step_condition"].(string)
			stepName := localStep["step_name"].(string)
			stepStartTrigger := localStep["step_start_trigger"].(string)

			deploymentStep := &octopusdeploy.DeploymentStep{
				Name:               stepName,
				PackageRequirement: "LetOctopusDecide",
				Condition:          octopusdeploy.DeploymentStepCondition(stepCondition),
				StartTrigger:       octopusdeploy.DeploymentStepStartTrigger(stepStartTrigger),
				Actions: []octopusdeploy.DeploymentAction{
					{
						Name:       stepName,
						ActionType: "Octopus.Script",
						Properties: map[string]string{
							"Octopus.Action.RunOnServer":                strconv.FormatBool(runOnServer),
							"Octopus.Action.Script.ScriptSource":        "Inline",
							"Octopus.Action.Package.DownloadOnTentacle": "False",
							"Octopus.Action.Script.ScriptBody":          scriptBody,
							"Octopus.Action.Script.Syntax":              scriptType,
						},
					},
				},
			}

			if targetRolesInterface, ok := localStep["target_roles"]; ok {
				var targetRoleSlice []string

				targetRoles := targetRolesInterface.([]interface{})

				for _, role := range targetRoles {
					targetRoleSlice = append(targetRoleSlice, role.(string))
				}

				deploymentStep.Properties = map[string]string{"Octopus.Action.TargetRoles": strings.Join(targetRoleSlice, ",")}
			}

			deploymentProcess.Steps = append(deploymentProcess.Steps, *deploymentStep)
		}
	}

	if v, ok := d.GetOk("deployment_step_kubernetes_helm"); ok {
		steps := v.([]interface{})
		for _, raw := range steps {

			localStep := raw.(map[string]interface{})

			resetValues := localStep["reset_values"].(bool)
			releaseName := localStep["release_name"].(string)
			namespace := localStep["namespace"].(string)
			yamlValues := localStep["yaml_values"].(string)
			tillerNamespace := localStep["tiller_namespace"].(string)
			feedID := localStep["feed_id"].(string)
			packageID := localStep["package_id"].(string)
			stepCondition := localStep["step_condition"].(string)
			stepName := localStep["step_name"].(string)
			stepStartTrigger := localStep["step_start_trigger"].(string)

			deploymentStep := &octopusdeploy.DeploymentStep{
				Name:               stepName,
				PackageRequirement: "LetOctopusDecide",
				Condition:          octopusdeploy.DeploymentStepCondition(stepCondition),
				StartTrigger:       octopusdeploy.DeploymentStepStartTrigger(stepStartTrigger),
				Actions: []octopusdeploy.DeploymentAction{
					{
						Name:       stepName,
						ActionType: "Octopus.HelmChartUpgrade",
						Properties: map[string]string{
							"Octopus.Action.Helm.ResetValues":           strconv.FormatBool(resetValues),
							"Octopus.Action.Helm.ReleaseName":           releaseName,
							"Octopus.Action.Helm.Namespace":             namespace,
							"Octopus.Action.Helm.YamlValues":            yamlValues,
							"Octopus.Action.Helm.TillerNamespace":       tillerNamespace,
							"Octopus.Action.Package.FeedId":             feedID,
							"Octopus.Action.Package.PackageId":          packageID,
							"Octopus.Action.Package.DownloadOnTentacle": "False",
						},
					},
				},
			}

			if targetRolesInterface, ok := localStep["target_roles"]; ok {
				var targetRoleSlice []string

				targetRoles := targetRolesInterface.([]interface{})

				for _, role := range targetRoles {
					targetRoleSlice = append(targetRoleSlice, role.(string))
				}

				deploymentStep.Properties = map[string]string{"Octopus.Action.TargetRoles": strings.Join(targetRoleSlice, ",")}
			}

			deploymentProcess.Steps = append(deploymentProcess.Steps, *deploymentStep)
		}
	}

	if v, ok := d.GetOk("deployment_step_kubernetes_yaml"); ok {
		steps := v.([]interface{})
		for _, raw := range steps {

			localStep := raw.(map[string]interface{})

			yamlValues := localStep["yaml_values"].(string)
			runOnServer := localStep["run_on_server"].(bool)
			stepCondition := localStep["step_condition"].(string)
			stepName := localStep["step_name"].(string)
			stepStartTrigger := localStep["step_start_trigger"].(string)

			deploymentStep := &octopusdeploy.DeploymentStep{
				Name:               stepName,
				PackageRequirement: "LetOctopusDecide",
				Condition:          octopusdeploy.DeploymentStepCondition(stepCondition),
				StartTrigger:       octopusdeploy.DeploymentStepStartTrigger(stepStartTrigger),
				Actions: []octopusdeploy.DeploymentAction{
					{
						Name:       stepName,
						ActionType: "Octopus.KubernetesDeployRawYaml",
						Properties: map[string]string{
							"Octopus.Action.RunOnServer":                             strconv.FormatBool(runOnServer),
							"Octopus.Action.Script.ScriptSource":                     "Inline",
							"Octopus.Action.KubernetesContainers.CustomResourceYaml": yamlValues,
						},
					},
				},
			}

			if targetRolesInterface, ok := localStep["target_roles"]; ok {
				var targetRoleSlice []string

				targetRoles := targetRolesInterface.([]interface{})

				for _, role := range targetRoles {
					targetRoleSlice = append(targetRoleSlice, role.(string))
				}

				deploymentStep.Properties = map[string]string{"Octopus.Action.TargetRoles": strings.Join(targetRoleSlice, ",")}
			}

			deploymentProcess.Steps = append(deploymentProcess.Steps, *deploymentStep)
		}
	}

	if v, ok := d.GetOk("deployment_step_package_script"); ok {
		steps := v.([]interface{})
		for _, raw := range steps {

			localStep := raw.(map[string]interface{})

			scriptFileName := localStep["script_file_name"].(string)
			scriptParameters := localStep["script_parameters"].(string)
			feedID := localStep["feed_id"].(string)
			packageID := localStep["package"].(string)
			jsonFileVariableReplacement := localStep["json_file_variable_replacement"].(string)
			variableSubstitutionInFiles := localStep["variable_substitution_in_files"].(string)
			configurationTransforms := localStep["configuration_transforms"].(bool)
			configurationVariables := localStep["configuration_variables"].(bool)
			stepCondition := localStep["step_condition"].(string)
			stepName := localStep["step_name"].(string)
			stepStartTrigger := localStep["step_start_trigger"].(string)
			runOnServer := localStep["run_on_server"].(bool)

			deploymentStep := &octopusdeploy.DeploymentStep{
				Name:               stepName,
				PackageRequirement: "LetOctopusDecide",
				Condition:          octopusdeploy.DeploymentStepCondition(stepCondition),
				StartTrigger:       octopusdeploy.DeploymentStepStartTrigger(stepStartTrigger),
				Actions: []octopusdeploy.DeploymentAction{
					{
						Name:       stepName,
						ActionType: "Octopus.Script",
						Properties: map[string]string{
							"Octopus.Action.RunOnServer":                strconv.FormatBool(runOnServer),
							"Octopus.Action.Script.ScriptSource":        "Package",
							"Octopus.Action.Package.DownloadOnTentacle": "False",
							"Octopus.Action.Package.FeedId":             feedID,
							"Octopus.Action.Package.PackageId":          packageID,
							"Octopus.Action.Script.ScriptFileName":      scriptFileName,
							"Octopus.Action.Script.ScriptParameters":    scriptParameters,
						},
					},
				},
			}

			if jsonFileVariableReplacement != "" {
				deploymentStep.Actions[0].Properties["Octopus.Action.Package.JsonConfigurationVariablesTargets"] = jsonFileVariableReplacement
				deploymentStep.Actions[0].Properties["Octopus.Action.Package.JsonConfigurationVariablesEnabled"] = "True"

				deploymentStep.Actions[0].Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.JsonConfigurationVariables"
			}

			if variableSubstitutionInFiles != "" {
				deploymentStep.Actions[0].Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] = variableSubstitutionInFiles
				deploymentStep.Actions[0].Properties["Octopus.Action.SubstituteInFiles.Enabled"] = "True"

				deploymentStep.Actions[0].Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.SubstituteInFiles"
			}

			if configurationTransforms {
				deploymentStep.Actions[0].Properties["Octopus.Action.Package.AutomaticallyRunConfigurationTransformationFiles"] = strconv.FormatBool(configurationTransforms)
				deploymentStep.Actions[0].Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.ConfigurationTransforms"
			}

			if configurationVariables {
				deploymentStep.Actions[0].Properties["Octopus.Action.Package.AutomaticallyUpdateAppSettingsAndConnectionStrings"] = strconv.FormatBool(configurationVariables)
				deploymentStep.Actions[0].Properties["Octopus.Action.EnabledFeatures"] += ",Octopus.Features.ConfigurationVariables"
			}

			if targetRolesInterface, ok := localStep["target_roles"]; ok {
				var targetRoleSlice []string

				targetRoles := targetRolesInterface.([]interface{})

				for _, role := range targetRoles {
					targetRoleSlice = append(targetRoleSlice, role.(string))
				}

				deploymentStep.Properties = map[string]string{"Octopus.Action.TargetRoles": strings.Join(targetRoleSlice, ",")}
			}

			deploymentProcess.Steps = append(deploymentProcess.Steps, *deploymentStep)
		}
	}

	if v, ok := d.GetOk("deployment_step_apply_terraform"); ok {
		steps := v.([]interface{})
		for _, raw := range steps {

			localStep := raw.(map[string]interface{})

			feedID := localStep["feed_id"].(string)
			packageID := localStep["package"].(string)
			stepCondition := localStep["step_condition"].(string)
			stepName := localStep["step_name"].(string)
			stepStartTrigger := localStep["step_start_trigger"].(string)
			runOnServer := localStep["run_on_server"].(bool)
			additionalInitParams := localStep["additional_init_params"].(string)

			deploymentStep := &octopusdeploy.DeploymentStep{
				Name:               stepName,
				PackageRequirement: "LetOctopusDecide",
				Condition:          octopusdeploy.DeploymentStepCondition(stepCondition),
				StartTrigger:       octopusdeploy.DeploymentStepStartTrigger(stepStartTrigger),
				Actions: []octopusdeploy.DeploymentAction{
					{
						Name:       stepName,
						ActionType: "Octopus.TerraformApply",
						Properties: map[string]string{
							"Octopus.Action.RunOnServer":                    strconv.FormatBool(runOnServer),
							"Octopus.Action.Script.ScriptSource":            "Package",
							"Octopus.Action.Package.DownloadOnTentacle":     "False",
							"Octopus.Action.Package.FeedId":                 feedID,
							"Octopus.Action.Package.PackageId":              packageID,
							"Octopus.Action.Aws.AssumeRole":                 "False",
							"Octopus.Action.AwsAccount.UseInstanceRole":     "False",
							"Octopus.Action.Terraform.AdditionalInitParams": additionalInitParams,
							"Octopus.Action.Terraform.AllowPluginDownloads": "True",
							"Octopus.Action.Terraform.ManagedAccount":       "None",
						},
					},
				},
			}

			if targetRolesInterface, ok := localStep["target_roles"]; ok {
				var targetRoleSlice []string

				targetRoles := targetRolesInterface.([]interface{})

				for _, role := range targetRoles {
					targetRoleSlice = append(targetRoleSlice, role.(string))
				}

				deploymentStep.Properties = map[string]string{"Octopus.Action.TargetRoles": strings.Join(targetRoleSlice, ",")}
			}

			if targetFilesInterface, ok := localStep["terraform_file_variable_replacement"]; ok {
				var targetFilesSlice []string

				targetFiles := targetFilesInterface.([]interface{})

				for _, file := range targetFiles {
					targetFilesSlice = append(targetFilesSlice, file.(string))
				}

				deploymentStep.Properties = map[string]string{"Octopus.Action.Terraform.FileSubstitution": strings.Join(targetFilesSlice, "\n")}
			}

			deploymentProcess.Steps = append(deploymentProcess.Steps, *deploymentStep)
		}
	}

	return deploymentProcess
}

func buildProjectResource(d *schema.ResourceData) *octopusdeploy.Project {
	name := d.Get("name").(string)
	lifecycleID := d.Get("lifecycle_id").(string)
	projectGroupID := d.Get("project_group_id").(string)

	project := octopusdeploy.NewProject(name, lifecycleID, projectGroupID)

	if attr, ok := d.GetOk("description"); ok {
		project.Description = attr.(string)
	}

	if attr, ok := d.GetOk("default_failure_mode"); ok {
		project.DefaultGuidedFailureMode = attr.(string)
	}

	if attr, ok := d.GetOk("skip_machine_behavior"); ok {
		project.ProjectConnectivityPolicy.SkipMachineBehavior = attr.(string)
	}

	if attr, ok := d.GetOk("allow_deployments_to_no_targets"); ok {
		project.ProjectConnectivityPolicy.AllowDeploymentsToNoTargets = attr.(bool)
	}

	if attr, ok := d.GetOk("tenanted_deployment_mode"); ok {
		project.TenantedDeploymentMode = attr.(string)
	}

	if attr, ok := d.GetOk("included_library_variable_sets"); ok {
		project.IncludedLibraryVariableSetIds = getSliceFromTerraformTypeList(attr)
	}

	if attr, ok := d.GetOk("discrete_channel_release"); ok {
		project.DiscreteChannelRelease = attr.(bool)
	}

	if attr, ok := d.GetOk("skip_package_steps_that_are_already_installed"); ok {
		project.DefaultToSkipIfAlreadyInstalled = attr.(bool)
	}

	return project
}

func updateDeploymentProcess(d *schema.ResourceData, client *octopusdeploy.Client, projectID string) error {
	deploymentProcess, err := client.DeploymentProcess.Get(projectID)

	if err != nil {
		return fmt.Errorf("error getting deployment process for project: %s", err.Error())
	}

	newDeploymentProcess := buildDeploymentProcess(d, deploymentProcess)
	// set the newly build deployment processes ID so it can be updated
	newDeploymentProcess.ID = deploymentProcess.ID

	updateDeploymentProcess, err := client.DeploymentProcess.Update(newDeploymentProcess)

	if err != nil {
		return fmt.Errorf("error creating deployment process for project: %s", err.Error())
	}

	d.Set("deployment_process_id", updateDeploymentProcess.ID)

	return nil
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newProject := buildProjectResource(d)

	createdProject, err := client.Project.Add(newProject)

	if err != nil {
		return fmt.Errorf("error creating project: %s", err.Error())
	}

	d.SetId(createdProject.ID)

	// set the deployment process
	errUpdatingDeploymentProcess := updateDeploymentProcess(d, client, createdProject.DeploymentProcessID)

	// deployment process is updated, not created, but log message makes more sense if it fails in a create step
	if errUpdatingDeploymentProcess != nil {
		return fmt.Errorf("error creating deploymentprocess: %s", errUpdatingDeploymentProcess.Error())
	}

	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	projectID := d.Id()

	project, err := client.Project.Get(projectID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading project id %s: %s", projectID, err.Error())
	}

	log.Printf("[DEBUG] project: %v", m)
	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("lifecycle_id", project.LifecycleID)
	d.Set("project_group_id", project.ProjectGroupID)
	d.Set("default_failure_mode", project.DefaultGuidedFailureMode)
	d.Set("skip_machine_behavior", project.ProjectConnectivityPolicy.SkipMachineBehavior)
	d.Set("allow_deployments_to_no_targets", project.ProjectConnectivityPolicy.AllowDeploymentsToNoTargets)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	project := buildProjectResource(d)
	project.ID = d.Id() // set project struct ID so octopus knows which project to update

	client := m.(*octopusdeploy.Client)

	project, err := client.Project.Update(project)

	if err != nil {
		return fmt.Errorf("error updating project id %s: %s", d.Id(), err.Error())
	}

	d.SetId(project.ID)

	// set the deployment process
	errUpdatingDeploymentProcess := updateDeploymentProcess(d, client, project.DeploymentProcessID)

	// deployment process is updated, not created, but log message makes more sense if it fails in a create step
	if errUpdatingDeploymentProcess != nil {
		return fmt.Errorf("error creating deploymentprocess: %s", errUpdatingDeploymentProcess.Error())
	}

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	projectID := d.Id()

	err := client.Project.Delete(projectID)

	if err != nil {
		return fmt.Errorf("error deleting project id %s: %s", projectID, err.Error())
	}

	d.SetId("")
	return nil
}
