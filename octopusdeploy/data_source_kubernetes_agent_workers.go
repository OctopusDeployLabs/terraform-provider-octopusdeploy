package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func dataSourceKubernetesAgentWorkers() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing kubernetes agent deployment targets.",
		ReadContext: dataSourceKubernetesAgentDeploymentTargetsRead,
		Schema:      getKubernetesAgentDeploymentTargetDataSchema(),
	}
}

func dataSourceKubernetesAgentWorkersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := machines.WorkersQuery{
		CommunicationStyles: []string{"KubernetesTentacle"},
		HealthStatuses:      expandArray(d.Get("health_statuses").([]interface{})),
		IDs:                 expandArray(d.Get("ids").([]interface{})),
		IsDisabled:          d.Get("is_disabled").(bool),
		Name:                d.Get("name").(string),
		PartialName:         d.Get("partial_name").(string),
		ShellNames:          expandArray(d.Get("shell_names").([]interface{})),
		Skip:                d.Get("skip").(int),
		Take:                d.Get("take").(int),
		Thumbprint:          d.Get("thumbprint").(string),
	}

	client := m.(*client.Client)
	existingWorkers, err := client.Workers.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedKubernetesAgents := []interface{}{}
	for _, worker := range existingWorkers.Items {
		flattenedKubernetesAgents = append(flattenedKubernetesAgents, flattenKubernetesAgentWorker(worker))
	}

	err = d.Set("kubernetes_agent_workers", flattenedKubernetesAgents)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("KubernetesAgentDeploymentWorkers " + time.Now().UTC().String())
	return nil
}
