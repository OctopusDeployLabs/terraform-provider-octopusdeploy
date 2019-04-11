package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataMachineReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint_communicationstyle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_proxyid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_tentacleversiondetails_upgradelocked": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_tentacleversiondetails_upgraderequired": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_tentacleversiondetails_upgradesuggested": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_tentacleversiondetails_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"environments": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"haslatestcalamari": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"isdisabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"isinprocess": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"machinepolicy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"statussummary": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenanteddeploymentparticipation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenantids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"tenanttags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataMachineReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	machineName := d.Get("name").(string)
	machines, err := client.Machine.GetAll()
	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading machine with name %s: %s", machineName, err.Error())
	}

	for _, m := range *machines {
		if m.Name == machineName {
			d.SetId(m.ID)
			d.Set("endpoint_communicationstyle", m.Endpoint.CommunicationStyle)
			d.Set("endpoint_id", m.Endpoint.ID)
			d.Set("endpoint_proxyid", m.Endpoint.ProxyID)
			d.Set("endpoint_tentacleversiondetails_upgradelocked", m.Endpoint.TentacleVersionDetails.UpgradeLocked)
			d.Set("endpoint_tentacleversiondetails_upgraderequired", m.Endpoint.TentacleVersionDetails.UpgradeRequired)
			d.Set("endpoint_tentacleversiondetails_upgradesuggested", m.Endpoint.TentacleVersionDetails.UpgradeSuggested)
			d.Set("endpoint_tentacleversiondetails_version", m.Endpoint.TentacleVersionDetails.Version)
			d.Set("endpoint_thumbprint", m.Endpoint.Thumbprint)
			d.Set("endpoint_uri", m.Endpoint.URI)
			d.Set("environments", m.EnvironmentIDs)
			d.Set("haslatestcalamari", m.HasLatestCalamari)
			d.Set("isdisabled", m.IsDisabled)
			d.Set("isinprocess", m.IsInProcess)
			d.Set("machinepolicy", m.MachinePolicyID)
			d.Set("roles", m.Roles)
			d.Set("status", m.Status)
			d.Set("statussummary", m.StatusSummary)
			d.Set("tenanteddeploymentparticipation", m.TenantedDeploymentParticipation)
			d.Set("tenantids", m.TenantIDs)
			d.Set("tenanttags", m.TenantTags)
			//d.Set("thumbprint", m.Thumbprint)
			//d.Set("uri", m.URI)

		}
	}

	return nil
}
