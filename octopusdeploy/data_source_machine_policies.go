package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMachinePolicies() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing machine policies.",
		ReadContext: dataSourceMachinePoliciesRead,
		Schema:      getMachinePolicyDataSchema(),
	}
}

func dataSourceMachinePoliciesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := machinepolicies.MachinePoliciesQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	spaceID := d.Get("space_id").(string)
	client := m.(*client.Client)
	existingMachinePolicies, err := machinepolicies.Get(client, spaceID, query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedMachinePolicies := []interface{}{}
	for _, machinePolicy := range existingMachinePolicies.Items {
		flattenedMachinePolicies = append(flattenedMachinePolicies, flattenMachinePolicy(machinePolicy))
	}

	d.Set("machine_policies", flattenedMachinePolicies)
	d.SetId("Machine Policies " + time.Now().UTC().String())

	return nil
}
