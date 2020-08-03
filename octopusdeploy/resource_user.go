package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		//Read:   resourceUserTriggerRead,
		//Update: resourceUserTriggerUpdate,
		//Delete: resourcesUserTriggerDelete,

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

	UserName := d.Get("UserName").(string)
	DisplayName := d.Get("DisplayName").(string)

	user := octopusdeploy.User(UserName, DisplayName)

	if attr, ok := d.GetOk("description"); ok {
		user.Description = attr.(string)
	}

	return user
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newUser := buildUserResource(d)

	newUserCreated, err := client.User.Add(newUser)

	if err != nil {
		return fmt.Errorf("error creating User %s: %s", err.Error())
	}

	d.SetId(newUserCreated.ID)

	return nil
}
