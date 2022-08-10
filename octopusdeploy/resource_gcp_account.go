package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
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

	client := m.(*client.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGoogleCloudPlatformAccount(ctx, d, createdAccount.(*accounts.GoogleCloudPlatformAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] GCP account created (%s)", d.Id())
	return nil
}

func resourceGoogleCloudPlatformAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting GCP account (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] GCP account deleted")
	return nil
}

func resourceGoogleCloudPlatformAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading GCP account (%s)", d.Id())

	client := m.(*client.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "GCP account")
	}

	amazonWebServicesAccount := accountResource.(*accounts.GoogleCloudPlatformAccount)
	if err := setGoogleCloudPlatformAccount(ctx, d, amazonWebServicesAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] GCP account read: %#v", amazonWebServicesAccount)
	return nil
}

func resourceGoogleCloudPlatformAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandGoogleCloudPlatformAccount(d)

	log.Printf("[INFO] updating GCP account: %#v", account)

	client := m.(*client.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGoogleCloudPlatformAccount(ctx, d, updatedAccount.(*accounts.GoogleCloudPlatformAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] GCP account updated (%s)", d.Id())
	return nil
}
