package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			constCanPasswordBeEdited: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			constDisplayName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constEmailAddress: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constIdentities: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			constIsActive: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			constIsRequestor: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			constIsService: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			constUsername: {
				Type:     schema.TypeString,
				Required: true,
			},
			constPassword: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	user, err := client.Users.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(constUsername, user.Username)

	return nil
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	username := d.Get(constUsername).(string)
	displayName := d.Get(constDisplayName).(string)

	user := octopusdeploy.NewUser(username, displayName)

	client := m.(*octopusdeploy.Client)
	user, err := client.Users.Add(user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.GetID())

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	username := d.Get(constUsername).(string)
	displayName := d.Get(constDisplayName).(string)
	user := octopusdeploy.NewUser(username, displayName)
	user.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.Users.Update(*user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userID := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Users.DeleteByID(userID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)

	return nil
}
