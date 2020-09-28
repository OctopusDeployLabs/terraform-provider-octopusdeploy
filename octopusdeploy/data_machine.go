package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataMachineReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constEndpointCommunicationStyle: {
				Type:     schema.TypeString,
				Computed: true,
			},
			constEndpointID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_proxyid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			constEndpointThumbprint: {
				Type:     schema.TypeString,
				Computed: true,
			},
			constEndpointURI: {
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

			constEnvironments: {
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
			constRoles: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			constStatus: {
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
	name := d.Get(constName).(string)

	apiClient := m.(*client.Client)
	resource, err := apiClient.Machines.GetByName(name)
	if err != nil {
		return createResourceOperationError(errorReadingMachine, name, err)
	}
	if resource == nil {
		// d.SetId(constEmptyString)
		return nil
	}

	logResource(constMachine, m)

	d.SetId(resource.ID)
	d.Set(constEndpointCommunicationStyle, resource.Endpoint.CommunicationStyle)
	d.Set(constEndpointID, resource.Endpoint.ID)
	d.Set("endpoint_proxyid", resource.Endpoint.ProxyID)
	d.Set("endpoint_tentacleversiondetails_upgradelocked", resource.Endpoint.TentacleVersionDetails.UpgradeLocked)
	d.Set("endpoint_tentacleversiondetails_upgraderequired", resource.Endpoint.TentacleVersionDetails.UpgradeRequired)
	d.Set("endpoint_tentacleversiondetails_upgradesuggested", resource.Endpoint.TentacleVersionDetails.UpgradeSuggested)
	d.Set("endpoint_tentacleversiondetails_version", resource.Endpoint.TentacleVersionDetails.Version)
	d.Set(constEndpointThumbprint, resource.Endpoint.Thumbprint)
	d.Set(constEndpointURI, resource.Endpoint.URI)
	d.Set(constEnvironments, resource.EnvironmentIDs)
	d.Set("haslatestcalamari", resource.HasLatestCalamari)
	d.Set("isdisabled", resource.IsDisabled)
	d.Set("isinprocess", resource.IsInProcess)
	d.Set("machinepolicy", resource.MachinePolicyID)
	d.Set(constRoles, resource.Roles)
	d.Set(constStatus, resource.Status)
	d.Set("statussummary", resource.StatusSummary)
	d.Set("tenanteddeploymentparticipation", resource.DeploymentMode)
	d.Set("tenantids", resource.TenantIDs)
	d.Set("tenanttags", resource.TenantTags)

	return nil
}
