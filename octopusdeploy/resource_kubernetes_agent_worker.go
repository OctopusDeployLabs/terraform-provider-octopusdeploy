package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKubernetesAgentWorker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesAgentWorkerCreate,
		DeleteContext: resourceKubernetesAgentWorkerDelete,
		Description:   "This resource manages Kubernetes agent workers in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceKubernetesAgentWorkerRead,
		Schema:        getKubernetesAgentWorkerSchema(),
		UpdateContext: resourceKubernetesAgentWorkerUpdate,
	}
}

func resourceKubernetesAgentWorkerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	worker := expandKubernetesAgentWorker(d)
	client := m.(*client.Client)
	createdWorker, err := workers.Add(client, worker)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdWorker.GetID())
	return nil
}

func resourceKubernetesAgentWorkerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	Worker, err := workers.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "kubernetes tentacle worker")
	}

	flattenedKubernetesAgentWorker := flattenKubernetesAgentWorker(Worker)
	for key, value := range flattenedKubernetesAgentWorker {
		if key != "id" {
			err := d.Set(key, value)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func resourceKubernetesAgentWorkerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	if err := workers.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceKubernetesAgentWorkerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	worker := expandKubernetesAgentWorker(d)
	client := m.(*client.Client)

	worker.ID = d.Id()

	updatedWorker, err := workers.Update(client, worker)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(updatedWorker.GetID())

	return nil
}
