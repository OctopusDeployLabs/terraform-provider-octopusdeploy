package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Type:     schema.TypeList,
			},
			"skip": {
				Default:  0,
				Type:     schema.TypeInt,
				Optional: true,
			},
			"take": {
				Default:  1,
				Type:     schema.TypeInt,
				Optional: true,
			},
			"users": {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"can_password_be_edited": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"display_name": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"email_address": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"id": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"identity": {
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"provider": {
										Optional: true,
										Type:     schema.TypeString,
									},
									"claim": {
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Required: true,
													Type:     schema.TypeString,
												},
												"is_identifying_claim": {
													Required: true,
													Type:     schema.TypeBool,
												},
												"value": {
													Required: true,
													Type:     schema.TypeString,
												},
											},
										},
										Type: schema.TypeSet,
									},
								},
							},
							Type: schema.TypeSet,
						},
						"is_active": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"is_service": {
							Optional: true,
							Type:     schema.TypeBool,
						},
						"username": {
							Optional: true,
							Type:     schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	query := octopusdeploy.UsersQuery{
		Filter: d.Get("filter").(string),
		IDs:    expandArray(d.Get("ids").([]interface{})),
		Skip:   d.Get("skip").(int),
		Take:   d.Get("take").(int),
	}

	client := meta.(*octopusdeploy.Client)
	users, err := client.Users.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedUsers := []interface{}{}
	for _, user := range users.Items {
		flattenedUser := map[string]interface{}{
			"can_password_be_edited": user.CanPasswordBeEdited,
			"display_name":           user.DisplayName,
			"email_address":          user.EmailAddress,
			"id":                     user.GetID(),
			"identity":               flattenIdentities(user.Identities),
			"is_active":              user.IsActive,
			"is_service":             user.IsService,
			"username":               user.Username,
		}
		flattenedUsers = append(flattenedUsers, flattenedUser)
	}

	d.Set("users", flattenedUsers)
	d.SetId("Users " + time.Now().UTC().String())

	return nil
}
