package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constLifecycleID: {
				Type:     schema.TypeString,
				Required: true,
			},
			constProjectGroupID: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constDefaultFailureMode: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "EnvironmentDefault",
				ValidateFunc: validateValueFunc([]string{
					"EnvironmentDefault",
					"Off",
					"On",
				}),
			},
			constSkipMachineBehavior: {
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

func buildProjectResource(d *schema.ResourceData) *model.Project {
	name := d.Get(constName).(string)
	lifecycleID := d.Get(constLifecycleID).(string)
	projectGroupID := d.Get(constProjectGroupID).(string)

	project := model.NewProject(name, lifecycleID, projectGroupID)

	if attr, ok := d.GetOk(constDescription); ok {
		project.Description = attr.(string)
	}

	if attr, ok := d.GetOk(constDefaultFailureMode); ok {
		project.DefaultGuidedFailureMode = attr.(string)
	}

	if attr, ok := d.GetOk(constSkipMachineBehavior); ok {
		project.ProjectConnectivityPolicy.SkipMachineBehavior = attr.(string)
	}

	if attr, ok := d.GetOk("allow_deployments_to_no_targets"); ok {
		project.ProjectConnectivityPolicy.AllowDeploymentsToNoTargets = attr.(bool)
	}

	if attr, ok := d.GetOk("tenanted_deployment_mode"); ok {
		project.TenantedDeploymentMode, _ = enum.ParseTenantedDeploymentMode(attr.(string))
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
	newProject := buildProjectResource(d)

	apiClient := m.(*client.Client)
	resource, err := apiClient.Projects.Add(newProject)
	if err != nil {
		return createResourceOperationError(errorCreatingProject, newProject.ID, err)
	}

	d.SetId(resource.ID)
	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Projects.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingProject, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constProject, m)

	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)
	d.Set(constLifecycleID, resource.LifecycleID)
	d.Set(constProjectGroupID, resource.ProjectGroupID)
	d.Set(constDefaultFailureMode, resource.DefaultGuidedFailureMode)
	d.Set(constSkipMachineBehavior, resource.ProjectConnectivityPolicy.SkipMachineBehavior)
	d.Set("allow_deployments_to_no_targets", resource.ProjectConnectivityPolicy.AllowDeploymentsToNoTargets)
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	project := buildProjectResource(d)
	project.ID = d.Id() // set project struct ID so octopus knows which project to update

	apiClient := m.(*client.Client)
	project, err := apiClient.Projects.Update(project)
	if err != nil {
		return createResourceOperationError(errorUpdatingProject, d.Id(), err)
	}

	d.SetId(project.ID)
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.Projects.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingProject, d.Id(), err)
	}

	d.SetId(constEmptyString)
	return nil
}
