package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMachinePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMachinePoliciesRead,
		Schema:      getMachinePolicyDataSchema(),
	}
}

func dataSourceMachinePoliciesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.MachinePoliciesQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	machinePolicies, err := client.MachinePolicies.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedMachinePolicies := []interface{}{}
	for _, machinePolicy := range machinePolicies.Items {
		flattenedMachinePolicies = append(flattenedMachinePolicies, flattenMachinePolicy(machinePolicy))
	}

	d.Set("machine_policies", flattenedMachinePolicies)
	d.SetId("Machine Policies " + time.Now().UTC().String())

	return nil
}
