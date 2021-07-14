package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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

	client := m.(*octopusdeploy.Client)
	createdMachinePolicy, err := client.MachinePolicies.Add(machinePolicy)
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

	client := m.(*octopusdeploy.Client)
	if err := client.MachinePolicies.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] machine policy deleted")
	return nil
}

func resourceMachinePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading machine policy (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	machinePolicy, err := client.MachinePolicies.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] machine policy (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
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
	client := m.(*octopusdeploy.Client)
	updatedMachinePolicy, err := client.MachinePolicies.Update(machinePolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setMachinePolicy(ctx, d, updatedMachinePolicy); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] machine policy updated (%s)", d.Id())
	return nil
}
