package octopusdeploy

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRunbook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRunbookCreate,
		DeleteContext: resourceRunbookDelete,
		Description:   "This resource manages runbooks in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceRunbookRead,
		Schema:        getRunbookSchema(),
		UpdateContext: resourceRunbookUpdate,
	}
}

func resourceRunbookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	runbook := expandRunbook(ctx, d)

	tflog.Info(ctx, fmt.Sprintf("creating runbook (%s)", runbook.Name))

	client := m.(*client.Client)
	createdRunbook, err := runbooks.Add(client, runbook)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setRunbook(ctx, d, createdRunbook); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdRunbook.GetID())

	tflog.Info(ctx, fmt.Sprintf("runbook created (%s)", d.Id()))
	return nil
}

func resourceRunbookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("deleting runbook (%s)", d.Id()))

	client := m.(*client.Client)
	if err := runbooks.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("runbook deleted (%s)", d.Id()))
	d.SetId("")
	return nil
}

func resourceRunbookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading runbook (%s)", d.Id()))

	client := m.(*client.Client)
	runbook, err := runbooks.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "runbook")
	}

	if err := setRunbook(ctx, d, runbook); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("runbook read (%s)", d.Id()))
	return nil
}

func resourceRunbookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("updating runbook (%s)", d.Id()))

	client := m.(*client.Client)
	runbook := expandRunbook(ctx, d)
	var updatedRunbook *runbooks.Runbook
	var err error

	runbookLinks, err := runbooks.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	runbook.Links = runbookLinks.Links

	updatedRunbook, err = runbooks.Update(client, runbook)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setRunbook(ctx, d, updatedRunbook); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("runbook updated (%s)", d.Id()))
	return nil
}
