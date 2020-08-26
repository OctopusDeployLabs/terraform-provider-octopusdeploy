package octopusdeploy

import (
	"fmt"
	"log"

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
		},
	}
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
		project.TenantedDeploymentMode, _ = octopusdeploy.ParseTenantedDeploymentMode(attr.(string))
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

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newProject := buildProjectResource(d)

	createdProject, err := client.Project.Add(newProject)

	if err != nil {
		return fmt.Errorf("error creating project: %s", err.Error())
	}

	d.SetId(createdProject.ID)

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
