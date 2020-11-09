package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserReadByName,
		Schema: map[string]*schema.Schema{
			"username": {
				Required:  true,
				Sensitive: true,
				Type:      schema.TypeString,
			},
			"can_password_be_edited": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"display_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"email_address": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"identity": {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"claim": {
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"is_identifying_claim": {
										Computed: true,
										Type:     schema.TypeBool,
									},
									"value": {
										Computed: true,
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
				Computed: true,
				Type:     schema.TypeBool,
			},
			"is_requestor": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"is_service": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"modified_by": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"modified_on": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"password": {
				Computed:  true,
				Sensitive: true,
				Type:      schema.TypeString,
			},
		},
	}
}

func dataSourceUserReadByName(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*octopusdeploy.Client)
	username := d.Get("username").(string)
	query := octopusdeploy.UsersQuery{
		Filter: username,
		Take:   1,
	}

	users, err := client.Users.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}
	if users == nil || len(users.Items) == 0 {
		return diag.Errorf("unable to retrieve user (filter: %s)", username)
	}

	// NOTE: two or more users can have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, user := range users.Items {
		if user.Username == username {
			flattenUser(ctx, d, user)
			return nil
		}
	}

	return nil
}
