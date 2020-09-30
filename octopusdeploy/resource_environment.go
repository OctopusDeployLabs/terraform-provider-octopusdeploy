package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constUseGuidedFailure: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			constAllowDynamicInfrastructure: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceEnvironmentRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Environments.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingEnvironment, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constEnvironment, m)

	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)
	d.Set(constUseGuidedFailure, resource.UseGuidedFailure)
	d.Set(constAllowDynamicInfrastructure, resource.AllowDynamicInfrastructure)

	return nil
}

func buildEnvironmentResource(d *schema.ResourceData) *model.Environment {
	envName := d.Get(constName).(string)

	var envDesc string
	var envGuided bool
	var envDynamic bool

	envDescInterface, ok := d.GetOk(constDescription)
	if ok {
		envDesc = envDescInterface.(string)
	}

	envGuidedInterface, ok := d.GetOk(constUseGuidedFailure)
	if ok {
		envGuided = envGuidedInterface.(bool)
	}

	allowDynamicInfrastructureInterface, ok := d.GetOk(constAllowDynamicInfrastructure)
	if ok {
		envDynamic = allowDynamicInfrastructureInterface.(bool)
	}

	var environment = model.NewEnvironment(envName, envDesc, envGuided)
	environment.AllowDynamicInfrastructure = envDynamic

	return environment
}

func resourceEnvironmentCreate(d *schema.ResourceData, m interface{}) error {
	environment := buildEnvironmentResource(d)

	apiClient := m.(*client.Client)
	resource, err := apiClient.Environments.Add(environment)
	if err != nil {
		return createResourceOperationError(errorCreatingEnvironment, environment.Name, err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func resourceEnvironmentUpdate(d *schema.ResourceData, m interface{}) error {
	environment := buildEnvironmentResource(d)
	environment.ID = d.Id() // set ID so Octopus API knows which environment to update

	apiClient := m.(*client.Client)
	resource, err := apiClient.Environments.Update(*environment)
	if err != nil {
		return createResourceOperationError(errorUpdatingEnvironment, d.Id(), err)
	}

	d.SetId(resource.ID)

	return nil
}

func resourceEnvironmentDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.Environments.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingEnvironment, id, err)
	}

	d.SetId(constEmptyString)
	return nil
}
