package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMachinePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMachinePolicyCreate,
		DeleteContext: resourceMachinePolicyDelete,
		Importer:      getImporter(),
		ReadContext:   resourceMachinePolicyRead,
		Schema:        getMachinePolicySchema(),
		UpdateContext: resourceMachinePolicyUpdate,
	}
}

func resourceMachinePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	machinePolicy := expandMachinePolicy(d)

	client := m.(*octopusdeploy.Client)
	createdMachinePolicy, err := client.MachinePolicies.Add(machinePolicy)
	if createdMachinePolicy != nil && err == nil {
		d.SetId(createdMachinePolicy.ID)
		setMachinePolicy(ctx, d, createdMachinePolicy)
		return nil
	}

	return diag.FromErr(err)
}

func resourceMachinePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	machinePolicy, err := client.MachinePolicies.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	setMachinePolicy(ctx, d, machinePolicy)
	return nil
}

func resourceMachinePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	machinePolicy := expandMachinePolicy(d)

	client := m.(*octopusdeploy.Client)
	updatedMachinePolicy, err := client.MachinePolicies.Update(machinePolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	setMachinePolicy(ctx, d, updatedMachinePolicy)
	return nil
}

func resourceMachinePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.MachinePolicies.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
