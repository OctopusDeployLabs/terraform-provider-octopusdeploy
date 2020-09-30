package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		return nil
	}

	logResource(constMachine, m)
	d.SetId(name)

	return nil
}
