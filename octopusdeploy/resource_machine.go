package octopusdeploy

import (
	"fmt"
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceMachineCreate,
		Read:   resourceMachineRead,
		Update: resourceMachineUpdate,
		Delete: resourceMachineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"communicationstyle": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validateValueFunc([]string{
								"None",
								"TentaclePassive",
								"TentacleActive",
								"Ssh",
								"OfflineDrop",
								"AzureWebApp",
								"Ftp",
								"AzureCloudService",
								"Kubernetes",
							}),
						},

						"proxyid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"thumbprint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Required: true,
						},
						"clustercertificate": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"clusterurl": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"skiptlsverification": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"defaultworkerpoolid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"authentication": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							MinItems: 0,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accountid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"clientcertificate": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"authenticationtype": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validateValueFunc([]string{
											"KubernetesCertificate",
											"KubernetesStandard",
										}),
									},
								},
							},
						},
					},
				},
			},
			"environments": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"haslatestcalamari": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"isdisabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"isinprocess": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"machinepolicy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"roles": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					MinItems: 1,
				},
				Required: true,
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
				Required: true,
				ValidateFunc: validateValueFunc([]string{
					"Untenanted",
					"TenantedOrUntenanted",
					"Tenanted",
				}),
			},
			"tenantids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"tenanttags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceMachineRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	machineID := d.Id()
	machine, err := client.Machine.Get(machineID)
	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error reading machine %s: %s", machineID, err.Error())
	}

	d.SetId(machine.ID)
	setMachineProperties(d, machine)

	return nil
}

func setMachineProperties(d *schema.ResourceData, m *octopusdeploy.Machine) {
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
}

func buildMachineResource(d *schema.ResourceData) *octopusdeploy.Machine {
	mName := d.Get("name").(string)
	mMachinepolicy := d.Get("machinepolicy").(string)
	mEnvironments := getSliceFromTerraformTypeList(d.Get("environments"))
	mRoles := getSliceFromTerraformTypeList(d.Get("roles"))
	mDisabled := d.Get("isdisabled").(bool)
	mTenantedDeploymentParticipation := d.Get("tenanteddeploymentparticipation").(string)
	mTenantIDs := getSliceFromTerraformTypeList(d.Get("tenantids"))
	mTenantTags := getSliceFromTerraformTypeList(d.Get("tenanttags"))

	//If we end up with a nil return, Octopus doesn't accept the API call. This ensure that we send
	//blank values rather than nil values.
	if mTenantIDs == nil {
		mTenantIDs = []string{}
	}
	if mTenantTags == nil {
		mTenantTags = []string{}
	}

	tfSchemaSetInterface, ok := d.GetOk("endpoint")
	if !ok {
		return nil
	}
	tfSchemaSet := tfSchemaSetInterface.(*schema.Set)
	if len(tfSchemaSet.List()) == 0 {
		return nil
	}
	//Get the first element in the list, which is a map of the interfaces
	tfSchemaList := tfSchemaSet.List()[0].(map[string]interface{})

	tfMachine := octopusdeploy.NewMachine(
		mName,
		mDisabled,
		mEnvironments,
		mRoles,
		mMachinepolicy,
		mTenantedDeploymentParticipation,
		mTenantIDs,
		mTenantTags,
	)

	tfMachine.URI = tfSchemaList["uri"].(string)
	tfMachine.Thumbprint = tfSchemaList["thumbprint"].(string)

	var proxyid *string
	if tfSchemaList["proxyid"] != nil {
		proxyString := tfSchemaList["proxyid"].(string)
		proxyid = &proxyString
	}

	tfMachine.Endpoint = &octopusdeploy.MachineEndpoint{
		URI:                 tfSchemaList["uri"].(string),
		Thumbprint:          tfSchemaList["thumbprint"].(string),
		CommunicationStyle:  tfSchemaList["communicationstyle"].(string),
		ProxyID:             proxyid,
		ClusterCertificate:  tfSchemaList["clustercertificate"].(string),
		ClusterURL:          tfSchemaList["clusterurl"].(string),
		Namespace:           tfSchemaList["namespace"].(string),
		SkipTLSVerification: strconv.FormatBool(tfSchemaList["skiptlsverification"].(bool)),
		DefaultWorkerPoolID: tfSchemaList["defaultworkerpoolid"].(string),
	}

	tfAuthenticationSchemaSetInterface, ok := tfSchemaList["authentication"]
	if ok {
		tfAuthenticationSchemaSet := tfAuthenticationSchemaSetInterface.(*schema.Set)
		if len(tfAuthenticationSchemaSet.List()) == 1 {
			//Get the first element in the list, which is a map of the interfaces
			tfAuthenticationSchemaList := tfAuthenticationSchemaSet.List()[0].(map[string]interface{})

			tfMachine.Endpoint.Authentication = &octopusdeploy.MachineEndpointAuthentication{
				AccountID:          tfAuthenticationSchemaList["accountid"].(string),
				ClientCertificate:  tfAuthenticationSchemaList["clientcertificate"].(string),
				AuthenticationType: tfAuthenticationSchemaList["authenticationtype"].(string),
			}
		}
	}

	return tfMachine
}

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)
	newMachine := buildMachineResource(d)
	newMachine.Status = "Unknown" //We don't want TF to attempt to update a machine just because its status has changed, so set it to Unknown on creation and let TF sort it out in the future.
	machine, err := client.Machine.Add(newMachine)
	if err != nil {
		return fmt.Errorf("error creating machine %s: %s", newMachine.Name, err.Error())
	}
	d.SetId(machine.ID)
	setMachineProperties(d, machine)
	return nil
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)
	machineID := d.Id()
	err := client.Machine.Delete(machineID)
	if err != nil {
		return fmt.Errorf("error deleting machine id %s: %s", machineID, err.Error())
	}
	d.SetId("")
	return nil
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {
	machine := buildMachineResource(d)
	machine.ID = d.Id() // set project struct ID so octopus knows which project to update
	client := m.(*octopusdeploy.Client)
	updatedMachine, err := client.Machine.Update(machine)
	if err != nil {
		return fmt.Errorf("error updating machine id %s: %s", d.Id(), err.Error())
	}
	d.SetId(updatedMachine.ID)
	setMachineProperties(d, machine)
	return nil
}
