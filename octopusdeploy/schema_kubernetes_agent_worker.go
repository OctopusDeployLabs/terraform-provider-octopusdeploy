package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"net/url"
)

func expandKubernetesAgentWorker(kubernetesAgent *schema.ResourceData) *machines.Worker {
	uri, _ := url.Parse(kubernetesAgent.Get("uri").(string))
	thumbprint := kubernetesAgent.Get("thumbprint").(string)

	communicationsMode := kubernetesAgent.Get("communication_mode").(string)
	upgradeLocked := kubernetesAgent.Get("upgrade_locked").(bool)
	var endpoint machines.IEndpoint = machines.NewKubernetesTentacleEndpoint(uri, thumbprint, upgradeLocked, communicationsMode, "")

	name := kubernetesAgent.Get("name").(string)
	Worker := machines.NewWorker(name, endpoint)

	Worker.IsDisabled = kubernetesAgent.Get("is_disabled").(bool)
	Worker.Thumbprint = thumbprint
	Worker.WorkerPoolIDs = getSliceFromTerraformTypeList(kubernetesAgent.Get("worker_pool_ids"))

	Worker.SpaceID = kubernetesAgent.Get("space_id").(string)

	return Worker
}

func flattenKubernetesAgentWorker(Worker *machines.Worker) map[string]interface{} {
	if Worker == nil {
		return nil
	}

	if Worker.Endpoint.GetCommunicationStyle() != "KubernetesTentacle" {
		return nil
	}

	endpoint := Worker.Endpoint.(*machines.KubernetesTentacleEndpoint)

	flattenedWorker := map[string]interface{}{}
	flattenedWorker["id"] = Worker.GetID()
	flattenedWorker["space_id"] = Worker.SpaceID
	flattenedWorker["name"] = Worker.Name
	flattenedWorker["machine_policy_id"] = Worker.MachinePolicyID
	flattenedWorker["is_disabled"] = Worker.IsDisabled

	flattenedWorker["thumbprint"] = endpoint.TentacleEndpointConfiguration.Thumbprint
	flattenedWorker["uri"] = endpoint.TentacleEndpointConfiguration.URI.String()
	flattenedWorker["communication_mode"] = endpoint.TentacleEndpointConfiguration.CommunicationMode
	flattenedWorker["worker_pool_ids"] = Worker.WorkerPoolIDs

	if endpoint.KubernetesAgentDetails != nil {
		flattenedWorker["agent_version"] = endpoint.KubernetesAgentDetails.AgentVersion
		flattenedWorker["agent_tentacle_version"] = endpoint.KubernetesAgentDetails.TentacleVersion
		flattenedWorker["agent_upgrade_status"] = endpoint.KubernetesAgentDetails.UpgradeStatus
		flattenedWorker["agent_helm_release_name"] = endpoint.KubernetesAgentDetails.HelmReleaseName
		flattenedWorker["agent_kubernetes_namespace"] = endpoint.KubernetesAgentDetails.KubernetesNamespace
	}

	return flattenedWorker
}

func getKubernetesAgentWorkerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id":       getIDSchema(),
		"space_id": getSpaceIDSchema(),
		"name":     getNameSchema(true),
		"communication_mode": {
			Optional:         true,
			Description:      "The communication mode used by the Kubernetes agent to communicate with Octopus Server. Currently, the only supported value is 'Polling'.",
			Type:             schema.TypeString,
			Default:          "Polling",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Polling"}, false)),
		},
		"machine_policy_id": {
			Description: "Optional ID of the machine policy that the Kubernetes agent will use. If not provided the default machine policy will be used.",
			Computed:    true,
			Optional:    true,
			Type:        schema.TypeString,
		},
		"thumbprint": {
			Description: "The thumbprint of the Kubernetes agent's certificate used by server to verify the identity of the agent. This is the same thumbprint that was used when installing the agent.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"uri": {
			Description: "The URI of the Kubernetes agent's used by the server to queue messages. This is the same subscription uri that was used when installing the agent.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"upgrade_locked": {
			Description: "If enabled the Kubernetes agent will not automatically upgrade and will stay on the currently installed version, even if the associated machine policy is configured to automatically upgrade.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"is_disabled": {
			Description: "Whether the Kubernetes agent is disabled. If the agent is disabled, it will not be included in any deployments.",
			Optional:    true,
			Default:     false,
			Type:        schema.TypeBool,
		},
		"worker_pool_ids": {
			Description: "A list of worker pool Ids specifying the pools in which this worker belongs",
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			Required:    true,
			Type:        schema.TypeList,
		},

		// Read-only Values
		"agent_version": {
			Description: "Current Helm chart version of the agent.",
			Computed:    true,
			Type:        schema.TypeString,
		},
		"agent_tentacle_version": {
			Description: "Current Tentacle version of the agent",
			Computed:    true,
			Type:        schema.TypeString,
		},
		"agent_upgrade_status": {
			Description: "Current upgrade availability status of the agent. One of 'NoUpgrades', 'UpgradeAvailable', 'UpgradeSuggested', 'UpgradeRequired'",
			Computed:    true,
			Type:        schema.TypeString,
		},
		"agent_helm_release_name": {
			Description: "Name of the Helm release that the agent belongs to.",
			Computed:    true,
			Type:        schema.TypeString,
		},
		"agent_kubernetes_namespace": {
			Description: "Name of the Kubernetes namespace where the agent is installed.",
			Computed:    true,
			Type:        schema.TypeString,
		},
	}
}

func getKubernetesAgentWorkerDataSchema() map[string]*schema.Schema {
	dataSchema := getKubernetesAgentWorkerSchema()
	setDataSchema(&dataSchema)

	WorkerDataSchema := getWorkerDataSchema()
	WorkerDataSchema["kubernetes_agent_workers"] = &schema.Schema{
		Computed:    true,
		Description: "A list of kubernetes agent workers that match the filter(s).",
		Elem:        &schema.Resource{Schema: dataSchema},
		Optional:    false,
		Type:        schema.TypeList,
	}

	delete(WorkerDataSchema, "communication_styles")
	delete(WorkerDataSchema, "workers")
	WorkerDataSchema["id"] = getDataSchemaID()
	return WorkerDataSchema
}
