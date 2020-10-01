package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

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

	if attr, ok := d.GetOk(constAllowDeploymentsToNoTargets); ok {
		project.ProjectConnectivityPolicy.AllowDeploymentsToNoTargets = attr.(bool)
	}

	if attr, ok := d.GetOk(constTenantedDeploymentMode); ok {
		project.TenantedDeploymentMode, _ = enum.ParseTenantedDeploymentMode(attr.(string))
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

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	project := buildProjectResource(d)
	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Projects.Add(project)
	if err != nil {
		// return createResourceOperationError(errorCreatingProject, project.ID, err)
		return diag.FromErr(err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Projects.GetByID(id)
	if err != nil {
		// return createResourceOperationError(errorReadingProject, id, err)
		return diag.FromErr(err)
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

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	project := buildProjectResource(d)
	project.ID = d.Id() // set ID so Octopus API knows which project to update

	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Projects.Update(*project)
	if err != nil {
		// return createResourceOperationError(errorUpdatingProject, d.Id(), err)
		return diag.FromErr(err)
	}

	d.SetId(resource.ID)

	return nil
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	diagValidate()

	apiClient := m.(*client.Client)
	err := apiClient.Projects.DeleteByID(id)
	if err != nil {
		// return createResourceOperationError(errorDeletingProject, d.Id(), err)
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)

	return nil
}
