package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		//Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the new user.",
			},
			"displayname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the new user.",
			},
		},
	}
}

func buildUserResource(d *schema.ResourceData) *octopusdeploy.User {
	Username := d.Get("username").(string)
	Displayname := d.Get("displayname").(string)

	user := octopusdeploy.NewUser(Username, Displayname)

	if attr, ok := d.GetOk("displayname"); ok {
		user.Username = attr.(string)
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

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	userID := d.Id()

	user, err := client.User.Get(userID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading user id %s: %s", user.ID, err.Error())
	}

	log.Printf("[DEBUG] user: %v", m)
	d.Set("name", user.Username)
	d.Set("displayname", user.Displayname)
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	user := buildUserResource(d)
	user.ID = d.Id()

	client := m.(*octopusdeploy.Client)

	updatedUser, err := client.User.Update(user)

	if err != nil {
		return fmt.Errorf("error updating user id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedUser.ID)
	return nil
}
