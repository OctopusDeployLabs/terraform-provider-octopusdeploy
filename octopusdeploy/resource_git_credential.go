package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGitCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitCredentialCreate,
		DeleteContext: resourceGitCredentialDelete,
		Description:   "This resource manages Git credentials in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceGitCredentialRead,
		Schema:        getGitCredentialSchema(),
		UpdateContext: resourceGitCredentialUpdate,
	}
}

func resourceGitCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resource := expandGitCredential(d)

	tflog.Info(ctx, fmt.Sprintf("creating Git credential, %s", resource.Name))

	client := m.(*client.Client)
	createdResource, err := client.GitCredentials.Add(resource)
	if err != nil {
		return diag.FromErr(err)
	}

	createdResource, err = client.GitCredentials.GetByID(createdResource.GetID())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGitCredential(ctx, d, createdResource); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdResource.GetID())

	tflog.Info(ctx, fmt.Sprintf("Git credential created (%s)", d.Id()))
	return nil
}

func resourceGitCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("deleting Git credential (%s)", d.Id()))

	client := m.(*client.Client)
	if err := client.GitCredentials.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	tflog.Info(ctx, "Git credential feed deleted")
	return nil
}

func resourceGitCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading Git credential (%s)", d.Id()))

	client := m.(*client.Client)
	resource, err := client.GitCredentials.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Git credential")
	}

	if err := setGitCredential(ctx, d, resource); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Git credential read (%s)", resource.GetID()))
	return nil
}

func resourceGitCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resource := expandGitCredential(d)

	tflog.Info(ctx, fmt.Sprintf("updating Git credential (%s)", resource.GetID()))

	client := m.(*client.Client)
	updatedResource, err := client.GitCredentials.Update(resource)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGitCredential(ctx, d, updatedResource); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Git credential updated (%s)", d.Id()))
	return nil
}
