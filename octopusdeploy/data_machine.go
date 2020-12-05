package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataMachine() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "use a machine-specific data resource instead",
		Read:               dataMachineReadByName,
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"endpoint_communicationstyle": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"endpoint_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"endpoint_proxyid": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"endpoint_thumbprint": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"endpoint_uri": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"endpoint_tentacleversiondetails_upgradelocked": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"endpoint_tentacleversiondetails_upgraderequired": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"endpoint_tentacleversiondetails_upgradesuggested": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"endpoint_tentacleversiondetails_version": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"environments": {
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Type:     schema.TypeList,
			},
			"haslatestcalamari": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"isdisabled": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"isinprocess": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"machinepolicy": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"roles": {
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Type:     schema.TypeList,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"statussummary": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"tenanteddeploymentparticipation": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"tenantids": {
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Type:     schema.TypeList,
			},
			"tenanttags": {
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Type:     schema.TypeList,
			},
		},
	}
}

func dataMachineReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	machineName := d.Get("name").(string)
	machines, err := client.Machines.GetByName(machineName)
	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading machine with name %s: %s", machineName, err.Error())
	}

	for _, machine := range machines {
		if machine.Name == machineName {
			endpointResource, err := octopusdeploy.ToEndpointResource(machine.Endpoint)
			if err != nil {
				return err
			}

			d.Set("endpoint_communicationstyle", machine.Endpoint.GetCommunicationStyle())
			d.Set("endpoint_id", endpointResource.GetID())
			d.Set("endpoint_proxyid", endpointResource.ProxyID)
			d.Set("endpoint_tentacleversiondetails_upgradelocked", endpointResource.TentacleVersionDetails.UpgradeLocked)
			d.Set("endpoint_tentacleversiondetails_upgraderequired", endpointResource.TentacleVersionDetails.UpgradeRequired)
			d.Set("endpoint_tentacleversiondetails_upgradesuggested", endpointResource.TentacleVersionDetails.UpgradeSuggested)
			d.Set("endpoint_tentacleversiondetails_version", endpointResource.TentacleVersionDetails.Version)
			d.Set("endpoint_thumbprint", endpointResource.Thumbprint)
			d.Set("endpoint_uri", endpointResource.URI)
			d.Set("environments", machine.EnvironmentIDs)
			d.Set("haslatestcalamari", machine.HasLatestCalamari)
			d.Set("isdisabled", machine.IsDisabled)
			d.Set("isinprocess", machine.IsInProcess)
			d.Set("machinepolicy", machine.MachinePolicyID)
			d.Set("roles", machine.Roles)
			d.Set("status", machine.Status)
			d.Set("statussummary", machine.StatusSummary)
			d.Set("tenanteddeploymentparticipation", machine.TenantedDeploymentMode)
			d.Set("tenantids", machine.TenantIDs)
			d.Set("tenanttags", machine.TenantTags)
			d.Set("thumbprint", machine.Thumbprint)
			d.Set("uri", machine.URI)
			d.SetId(machine.GetID())

			return nil
		}
	}

	return nil
}
