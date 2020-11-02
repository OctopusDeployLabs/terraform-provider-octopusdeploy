package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
					"EnvironmentDefault",
					"Off",
					"On",
				}, false)),
			},
			constSkipMachineBehavior: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "None",
				ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
					"SkipUnavailableMachines",
					"None",
				}, false)),
			},
			constAllowDeploymentsToNoTargets: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			constTenantedDeploymentMode: getTenantedDeploymentSchema(),
			constIncludedLibraryVariableSets: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constDiscreteChannelRelease: {
				Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			constSkipPackageStepsThatAreAlreadyInstalled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func buildProjectResource(d *schema.ResourceData) *octopusdeploy.Project {
	name := d.Get(constName).(string)
	lifecycleID := d.Get(constLifecycleID).(string)
	projectGroupID := d.Get(constProjectGroupID).(string)

	project := octopusdeploy.NewProject(name, lifecycleID, projectGroupID)

	if attr, ok := d.GetOk(constDescription); ok {
		project.Description = attr.(string)
	}

	if attr, ok := d.GetOk(constDefaultFailureMode); ok {
		project.DefaultGuidedFailureMode = attr.(string)
	}

	if attr, ok := d.GetOk(constSkipMachineBehavior); ok {
		project.ProjectConnectivityPolicy.SkipMachineBehavior = octopusdeploy.SkipMachineBehavior(attr.(string))
	}

	if attr, ok := d.GetOk(constAllowDeploymentsToNoTargets); ok {
		project.ProjectConnectivityPolicy.AllowDeploymentsToNoTargets = attr.(bool)
	}

	if attr, ok := d.GetOk(constTenantedDeploymentMode); ok {
		project.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(attr.(string))
	}

	if attr, ok := d.GetOk(constIncludedLibraryVariableSets); ok {
		project.IncludedLibraryVariableSetIDs = getSliceFromTerraformTypeList(attr)
	}

	if attr, ok := d.GetOk(constDiscreteChannelRelease); ok {
		project.DiscreteChannelRelease = attr.(bool)
	}

	if attr, ok := d.GetOk(constSkipPackageStepsThatAreAlreadyInstalled); ok {
		project.DefaultToSkipIfAlreadyInstalled = attr.(bool)
	}

	return project
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	project := buildProjectResource(d)

	client := m.(*octopusdeploy.Client)
	resource, err := client.Projects.Add(project)
	if err != nil {
		return createResourceOperationError(errorCreatingProject, project.ID, err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.Projects.GetByID(id)
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
	d.Set(constAllowDeploymentsToNoTargets, resource.ProjectConnectivityPolicy.AllowDeploymentsToNoTargets)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	project := buildProjectResource(d)
	project.ID = d.Id() // set ID so Octopus API knows which project to update

	client := m.(*octopusdeploy.Client)
	resource, err := client.Projects.Update(project)
	if err != nil {
		return createResourceOperationError(errorUpdatingProject, d.Id(), err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Projects.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingProject, d.Id(), err)
	}

	d.SetId(constEmptyString)

	return nil
}
