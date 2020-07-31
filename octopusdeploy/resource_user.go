package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNewUserTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewUserCreate,
		//Read:   resourceNewUserTriggerRead,
		//Update: resourceNewUserTriggerUpdate,
		//Delete: resourceNewUserTriggerDelete,

		Schema: map[string]*schema.Schema{
			"UserName": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the new user.",
			},
			"DisplayName": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the new user.",
			},
		},
	}
}

func buildUserResource(d *schema.ResourceData) *octopusdeploy.User {
	envUserName := d.Get("UserName").(string)
	envDisplayName := d.Get("DisplayName").(string)

	var envDynamic bool

	allowDynamicInfrastructureInterface, ok := d.GetOk("allow_dynamic_infrastructure")
	if ok {
		envDynamic = allowDynamicInfrastructureInterface.(bool)
	}

	var User = octopusdeploy.NewUser(envUserName, envDesc, envGuided)
	User.AllowDynamicInfrastructure = envDynamic

	return User
}

func resourceNewUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newUser := buildUserResource(d)
	env, err := client.User.Add(newUser)

	if err != nil {
		return fmt.Errorf("error creating User %s: %s", newUser.UserName, err.Error())
	}

	d.SetId(env.ID)

	return nil
}
