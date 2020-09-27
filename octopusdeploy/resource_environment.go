package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
	environmentID := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Environments.GetByID(environmentID)
	if err != nil {
		return fmt.Errorf(errorReadingEnvironment, environmentID, err.Error())
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

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
	apiClient := m.(*client.Client)

	newEnvironment := buildEnvironmentResource(d)
	env, err := apiClient.Environments.Add(newEnvironment)

	if err != nil {
		return fmt.Errorf("error creating environment %s: %s", newEnvironment.Name, err.Error())
	}

	d.SetId(env.ID)

	return nil
}

func resourceEnvironmentUpdate(d *schema.ResourceData, m interface{}) error {
	env := buildEnvironmentResource(d)
	env.ID = d.Id() // set project struct ID so octopus knows which project to update

	apiClient := m.(*client.Client)

	updatedEnv, err := apiClient.Environments.Update(env)

	if err != nil {
		return fmt.Errorf("error updating environment id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedEnv.ID)
	return nil
}

func resourceEnvironmentDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	environmentID := d.Id()

	err := apiClient.Environments.DeleteByID(environmentID)

	if err != nil {
		return fmt.Errorf("error deleting environment id %s: %s", environmentID, err.Error())
	}

	d.SetId(constEmptyString)
	return nil
}
