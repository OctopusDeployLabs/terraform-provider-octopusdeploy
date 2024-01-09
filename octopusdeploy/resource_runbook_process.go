package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbookprocess"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRunbookProcess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRunbookProcessCreate,
		DeleteContext: resourceRunbookProcessDelete,
		Description:   "This resource manages runbook processes in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceRunbookProcessRead,
		Schema:        getRunbookProcessSchema(),
		UpdateContext: resourceRunbookProcessUpdate,
	}
}

func getRunbookProcessSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": getIDSchema(),
		"last_snapshot_id": {
			Description: "Read only value containing the last snapshot ID.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"project_id": {
			Description: "The project ID associated with this runbook process.",
			Optional:    true,
			Computed:    true,
			Type:        schema.TypeString,
		},
		"runbook_id": {
			Description: "The runbook ID associated with this runbook process.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"space_id": getSpaceIDSchema(),
		"step":     getDeploymentStepSchema(),
		"version": {
			Computed:    true,
			Description: "The version number of this runbook process.",
			Optional:    true,
			Type:        schema.TypeInt,
		},
	}
}

// resourceRunbookProcessCreate "creates" a new runbook deployment process. In reality every runbook has a deployment process
// already, so this function retrieves the existing process and updates it.
func resourceRunbookProcessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	runbookProcess := expandRunbookProcess(ctx, d, client)

	log.Printf("[INFO] creating runbook process: %#v", runbookProcess)

	runbook, err := runbooks.GetByID(client, d.Get("space_id").(string), runbookProcess.RunbookID)
	if err != nil {
		return diag.FromErr(err)
	}

	var current *runbookprocess.RunbookProcess
	current, err = runbookprocess.GetByID(client, d.Get("space_id").(string), runbook.RunbookProcessID)

	runbookProcess.ID = current.ID
	runbookProcess.Links = current.Links
	runbookProcess.Version = current.Version

	createdRunbookProcess, err := runbookprocess.Update(client, runbookProcess)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := setRunbookProcess(ctx, d, createdRunbookProcess); err != nil {
		return diag.FromErr(err)
	}

	id := createdRunbookProcess.GetID()

	d.SetId(id)

	log.Printf("[INFO] deployment process created (%s)", d.Id())
	return nil
}

func resourceRunbookProcessDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting runbook process (%s)", d.Id())

	// "Deleting" a runbook process just means to clear it out
	client := m.(*client.Client)
	current, err := runbookprocess.GetByID(client, d.Get("space_id").(string), d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	runbookProcess := &runbookprocess.RunbookProcess{
		Version: current.Version,
	}
	runbookProcess.Links = current.Links
	runbookProcess.ID = d.Id()
	if v, ok := d.GetOk("space_id"); ok {
		runbookProcess.SpaceID = v.(string)
	}

	_, err = runbookprocess.Update(client, runbookProcess)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] deployment process deleted")
	return nil
}

func resourceRunbookProcessRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading runbook process (%s)", d.Id())

	client := m.(*client.Client)
	runbookProcess, err := runbookprocess.GetByID(client, d.Get("space_id").(string), d.Id())

	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "runbook_process")
	}

	if err := setRunbookProcess(ctx, d, runbookProcess); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] runbook process read (%s)", d.Id())
	return nil
}

func resourceRunbookProcessUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating runbook process (%s)", d.Id())

	client := m.(*client.Client)
	runbookProcess := expandRunbookProcess(ctx, d, client)
	current, err := runbookprocess.GetByID(client, runbookProcess.SpaceID, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	runbookProcess.Links = current.Links
	runbookProcess.Version = current.Version

	updatedRunbookProcess, err := runbookprocess.Update(client, runbookProcess)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := setRunbookProcess(ctx, d, updatedRunbookProcess); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] deployment process updated (%s)", d.Id())
	return nil
}
