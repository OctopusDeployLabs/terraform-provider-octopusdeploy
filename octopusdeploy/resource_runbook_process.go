package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
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
			Optional: true,
			Type:     schema.TypeString,
		},
		"project_id": {
			Description: "The project ID associated with this deployment process.",
			Optional:    true,
			Computed:    true,
			Type:        schema.TypeString,
		},
		"runbook_id": {
			Description: "The runbook ID associated with this deployment process.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"space_id": getSpaceIDSchema(),
		"step":     getDeploymentStepSchema(),
		"version": {
			Computed:    true,
			Description: "The version number of this deployment process.",
			Optional:    true,
			Type:        schema.TypeInt,
		},
	}
}

// resourceRunbookProcessCreate "creates" a new runbook deployment process. In reality every runbook has a deployment process
// already, so this function retrieves the existing process and updates it.
func resourceRunbookProcessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	runbookProcess := expandRunbookProcess(d, client)

	log.Printf("[INFO] creating runbook process: %#v", runbookProcess)

	runbook, err := client.Runbooks.GetByID(runbookProcess.RunbookID)
	if err != nil {
		return diag.FromErr(err)
	}

	var current *runbooks.RunbookProcess
	current, err = client.RunbookProcesses.GetByID(runbook.RunbookProcessID)

	runbookProcess.ID = current.ID
	runbookProcess.Links = current.Links
	runbookProcess.Version = current.Version
	runbookProcess.ProjectID = current.ProjectID
	runbookProcess.RunbookID = current.RunbookID

	createdRunbookProcess, err := client.RunbookProcesses.Update(runbookProcess)

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
	current, err := client.RunbookProcesses.GetByID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	runbookProcess := &runbooks.RunbookProcess{
		Version: current.Version,
	}
	runbookProcess.Links = current.Links
	runbookProcess.ID = d.Id()

	_, err = client.RunbookProcesses.Update(runbookProcess)

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
	runbookProcess, err := client.RunbookProcesses.GetByID(d.Id())

	if err != nil {
		return diag.FromErr(err)
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
	runbookProcess := expandRunbookProcess(d, client)
	current, err := client.RunbookProcesses.GetByID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	runbookProcess.Links = current.Links
	runbookProcess.Version = current.Version

	updatedRunbookProcess, err := client.RunbookProcesses.Update(runbookProcess)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := setRunbookProcess(ctx, d, updatedRunbookProcess); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] deployment process updated (%s)", d.Id())
	return nil
}
