package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	resourceUserImporter := &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
	resourceUserSchema := map[string]*schema.Schema{
		constCanPasswordBeEdited: {
			Optional: true,
			Type:     schema.TypeBool,
		},
		constDisplayName: {
			Required: true,
			Type:     schema.TypeString,
		},
		constEmailAddress: {
			Optional: true,
			Type:     schema.TypeString,
		},
		constIdentities: {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		constIsActive: {
			Optional: true,
			Type:     schema.TypeBool,
		},
		constIsRequestor: {
			Optional: true,
			Type:     schema.TypeBool,
		},
		constIsService: {
			Optional: true,
			Type:     schema.TypeBool,
		},
		constUsername: {
			Required: true,
			Type:     schema.TypeString,
		},
		constPassword: {
			Optional: true,
			Type:     schema.TypeString,
		},
	}

	return &schema.Resource{
		CreateContext: resourceUserCreate,
		DeleteContext: resourceUserDelete,
		Importer:      resourceUserImporter,
		ReadContext:   resourceUserRead,
		Schema:        resourceUserSchema,
		UpdateContext: resourceUserUpdate,
	}
}

func buildUser(d *schema.ResourceData) *octopusdeploy.User {
	var username string
	if v, ok := d.GetOk(constUsername); ok {
		username = v.(string)
	}

	var displayName string
	if v, ok := d.GetOk(constDisplayName); ok {
		displayName = v.(string)
	}

	user := octopusdeploy.NewUser(username, displayName)

	if v, ok := d.GetOk(constCanPasswordBeEdited); ok {
		user.CanPasswordBeEdited = v.(bool)
	}

	if v, ok := d.GetOk(constEmailAddress); ok {
		user.EmailAddress = v.(string)
	}

	if v, ok := d.GetOk(constIsActive); ok {
		user.IsActive = v.(bool)
	}

	if v, ok := d.GetOk(constIsRequestor); ok {
		user.IsRequestor = v.(bool)
	}

	if v, ok := d.GetOk(constIsService); ok {
		user.IsService = v.(bool)
	}

	if v, ok := d.GetOk(constPassword); ok {
		user.Password = v.(string)
	}

	return user
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	user, err := client.Users.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(constCanPasswordBeEdited, user.CanPasswordBeEdited)
	d.Set(constDisplayName, user.DisplayName)
	d.Set(constEmailAddress, user.EmailAddress)
	d.Set(constIsActive, user.IsActive)
	d.Set(constIsRequestor, user.IsRequestor)
	d.Set(constIsService, user.IsService)
	d.Set(constPassword, user.Password)
	d.Set(constUsername, user.Username)
	d.SetId(user.GetID())

	return nil
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user := buildUser(d)

	client := m.(*octopusdeploy.Client)
	user, err := client.Users.Add(user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.GetID())

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user := buildUser(d)
	user.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.Users.Update(user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Users.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)

	return nil
}
