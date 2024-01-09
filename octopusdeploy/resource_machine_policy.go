package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMachinePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMachinePolicyCreate,
		DeleteContext: resourceMachinePolicyDelete,
		Description:   "This resource manages machine policies in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceMachinePolicyRead,
		Schema:        getMachinePolicySchema(),
		UpdateContext: resourceMachinePolicyUpdate,
	}
}

func resourceMachinePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	machinePolicy := expandMachinePolicy(d)

	log.Printf("[INFO] creating machine policy: %#v", machinePolicy)

	client := m.(*client.Client)
	createdMachinePolicy, err := machinepolicies.Add(client, machinePolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setMachinePolicy(ctx, d, createdMachinePolicy); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdMachinePolicy.GetID())

	log.Printf("[INFO] machine policy created (%s)", d.Id())
	return nil
}

func resourceMachinePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting machine policy (%s)", d.Id())

	spaceID := d.Get("space_id").(string)
	client := m.(*client.Client)
	if err := machinepolicies.DeleteByID(client, spaceID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] machine policy deleted")
	return nil
}

func resourceMachinePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading machine policy (%s)", d.Id())

	spaceID := d.Get("space_id").(string)
	client := m.(*client.Client)
	machinePolicy, err := machinepolicies.GetByID(client, spaceID, d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "machine policy")
	}

	if err := setMachinePolicy(ctx, d, machinePolicy); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] machine policy read (%s)", d.Id())
	return nil
}

func resourceMachinePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating machine policy (%s)", d.Id())

	machinePolicy := expandMachinePolicy(d)
	client := m.(*client.Client)
	updatedMachinePolicy, err := machinepolicies.Update(client, machinePolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setMachinePolicy(ctx, d, updatedMachinePolicy); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] machine policy updated (%s)", d.Id())
	return nil
}
