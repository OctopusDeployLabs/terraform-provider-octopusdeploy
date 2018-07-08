package octopusdeploy

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"deployment_process_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"project_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_failure_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "EnvironmentDefault",
				ValidateFunc: validateValueFunc([]string{
					"EnvironmentDefault",
					"Off",
					"On",
				}),
			},
			"skip_machine_behavior": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "None",
				ValidateFunc: validateValueFunc([]string{
					"SkipUnavailableMachines",
					"None",
				}),
			},
			"deployment_step_windows_service": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_transforms": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"configuration_variables": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"executable_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"feed_id": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "feeds-builtin",
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
						"step_condition": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validateValueFunc([]string{
								"success",
								"failure",
								"always",
								"variable",
							}),
							Default: "success",
						},
						"step_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"step_start_trigger": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "StartAfterPrevious",
							ValidateFunc: validateValueFunc([]string{
								"startafterprevious",
								"startwithprevious",
							}),
						},
						"target_roles": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
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
			serviceAccount := localStep["service_account"].(string)
			serviceName := localStep["service_name"].(string)
			serviceStartMode := localStep["service_start_mode"].(string)
			stepName := localStep["step_name"].(string)
			stepCondition := localStep["step_condition"].(string)
			stepStartTrigger := localStep["step_start_trigger"].(string)

			deploymentStep := &octopusdeploy.DeploymentStep{
				Name:               stepName,
				PackageRequirement: "LetOctopusDecide",
				Condition:          stepCondition,
				StartTrigger:       stepStartTrigger,
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
							"Octopus.Action.Package.DownloadOnTentacle":                                 "False",
							"Octopus.Action.WindowsService.ServiceName":                                 serviceName,
							"Octopus.Action.WindowsService.ExecutablePath":                              executablePath,
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
