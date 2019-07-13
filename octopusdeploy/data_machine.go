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
	machine, err := client.Machine.GetByName(machineName)
	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading machine with name %s: %s", machineName, err.Error())
	}

	d.SetId(machine.ID)
	d.Set("endpoint_communicationstyle", machine.Endpoint.CommunicationStyle)
	d.Set("endpoint_id", machine.Endpoint.ID)
	d.Set("endpoint_proxyid", machine.Endpoint.ProxyID)
	d.Set("endpoint_tentacleversiondetails_upgradelocked", machine.Endpoint.TentacleVersionDetails.UpgradeLocked)
	d.Set("endpoint_tentacleversiondetails_upgraderequired", machine.Endpoint.TentacleVersionDetails.UpgradeRequired)
	d.Set("endpoint_tentacleversiondetails_upgradesuggested", machine.Endpoint.TentacleVersionDetails.UpgradeSuggested)
	d.Set("endpoint_tentacleversiondetails_version", machine.Endpoint.TentacleVersionDetails.Version)
	d.Set("endpoint_thumbprint", machine.Endpoint.Thumbprint)
	d.Set("endpoint_uri", machine.Endpoint.URI)
	d.Set("environments", machine.EnvironmentIDs)
	d.Set("haslatestcalamari", machine.HasLatestCalamari)
	d.Set("isdisabled", machine.IsDisabled)
	d.Set("isinprocess", machine.IsInProcess)
	d.Set("machinepolicy", machine.MachinePolicyID)
	d.Set("roles", machine.Roles)
	d.Set("status", machine.Status)
	d.Set("statussummary", machine.StatusSummary)
	d.Set("tenanteddeploymentparticipation", machine.TenantedDeploymentParticipation)
	d.Set("tenantids", machine.TenantIDs)
	d.Set("tenanttags", machine.TenantTags)
	//d.Set("thumbprint", machine.Thumbprint)
	//d.Set("uri", machine.URI)

	return nil
}
