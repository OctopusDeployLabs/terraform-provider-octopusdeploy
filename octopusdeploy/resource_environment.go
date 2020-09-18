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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_guided_failure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_dynamic_infrastructure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceEnvironmentRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	environmentID := d.Id()
	env, err := apiClient.Environments.Get(environmentID)

	if err == client.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading environment %s: %s", environmentID, err.Error())
	}

	d.Set("name", env.Name)
	d.Set("description", env.Description)
	d.Set("use_guided_failure", env.UseGuidedFailure)
	d.Set("allow_dynamic_infrastructure", env.AllowDynamicInfrastructure)

	return nil
}

func buildEnvironmentResource(d *schema.ResourceData) *model.Environment {
	envName := d.Get("name").(string)

	var envDesc string
	var envGuided bool
	var envDynamic bool

	envDescInterface, ok := d.GetOk("description")
	if ok {
		envDesc = envDescInterface.(string)
	}

	envGuidedInterface, ok := d.GetOk("use_guided_failure")
	if ok {
		envGuided = envGuidedInterface.(bool)
	}

	allowDynamicInfrastructureInterface, ok := d.GetOk("allow_dynamic_infrastructure")
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

	err := apiClient.Environments.Delete(environmentID)

	if err != nil {
		return fmt.Errorf("error deleting environment id %s: %s", environmentID, err.Error())
	}

	d.SetId("")
	return nil
}
