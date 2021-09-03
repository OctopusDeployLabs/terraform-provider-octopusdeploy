package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGoogleCloudPlatformAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGoogleCloudPlatformAccountCreate,
		DeleteContext: resourceGoogleCloudPlatformAccountDelete,
		Description:   "This resource manages GCP accounts in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceGoogleCloudPlatformAccountRead,
		Schema:        getGoogleCloudPlatformAccountSchema(),
		UpdateContext: resourceGoogleCloudPlatformAccountUpdate,
	}
}

func resourceGoogleCloudPlatformAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandGoogleCloudPlatformAccount(d)

	log.Printf("[INFO] creating GCP account: %#v", account)

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGoogleCloudPlatformAccount(ctx, d, createdAccount.(*octopusdeploy.GoogleCloudPlatformAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] GCP account created (%s)", d.Id())
	return nil
}

func resourceGoogleCloudPlatformAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting GCP account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] GCP account deleted")
	return nil
}

func resourceGoogleCloudPlatformAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading GCP account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] GCP account (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	amazonWebServicesAccount := accountResource.(*octopusdeploy.GoogleCloudPlatformAccount)
	if err := setGoogleCloudPlatformAccount(ctx, d, amazonWebServicesAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] GCP account read: %#v", amazonWebServicesAccount)
	return nil
}

func resourceGoogleCloudPlatformAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandGoogleCloudPlatformAccount(d)

	log.Printf("[INFO] updating GCP account: %#v", account)

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGoogleCloudPlatformAccount(ctx, d, updatedAccount.(*octopusdeploy.GoogleCloudPlatformAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] GCP account updated (%s)", d.Id())
	return nil
}
