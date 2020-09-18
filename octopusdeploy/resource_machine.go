package octopusdeploy

import (
	"fmt"
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	apiClient := m.(*client.Client)

	machineID := d.Id()
	machine, err := apiClient.Machines.Get(machineID)
	if err == client.ErrItemNotFound {
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

func setMachineProperties(d *schema.ResourceData, m *model.Machine) {
	d.Set("environments", m.EnvironmentIDs)
	d.Set("haslatestcalamari", m.HasLatestCalamari)
	d.Set("isdisabled", m.IsDisabled)
	d.Set("isinprocess", m.IsInProcess)
	d.Set("machinepolicy", m.MachinePolicyID)
	d.Set("roles", m.Roles)
	d.Set("status", m.Status)
	d.Set("statussummary", m.StatusSummary)
	d.Set("tenanteddeploymentparticipation", m.DeploymentMode)
	d.Set("tenantids", m.TenantIDs)
	d.Set("tenanttags", m.TenantTags)
}

func buildMachineResource(d *schema.ResourceData) *model.Machine {
	mName := d.Get("name").(string)
	mMachinepolicy := d.Get("machinepolicy").(string)
	mEnvironments := getSliceFromTerraformTypeList(d.Get("environments"))
	mRoles := getSliceFromTerraformTypeList(d.Get("roles"))
	mDisabled := d.Get("isdisabled").(bool)
	mTenantedDeploymentParticipation, _ := d.Get("tenanteddeploymentparticipation").(string)
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

	tfMachine := &model.Machine{
		Name:            mName,
		IsDisabled:      mDisabled,
		EnvironmentIDs:  mEnvironments,
		Roles:           mRoles,
		MachinePolicyID: mMachinepolicy,
		DeploymentMode:  mTenantedDeploymentParticipation,
		TenantIDs:       mTenantIDs,
		TenantTags:      mTenantTags,
	}

	tfMachine.URI = tfSchemaList["uri"].(string)
	tfMachine.Thumbprint = tfSchemaList["thumbprint"].(string)

	var proxyid *string
	if tfSchemaList["proxyid"] != nil {
		proxyString := tfSchemaList["proxyid"].(string)
		proxyid = &proxyString
	}

	tfMachine.Endpoint = &model.MachineEndpoint{
		URI:                 tfSchemaList["uri"].(string),
		Thumbprint:          tfSchemaList["thumbprint"].(string),
		CommunicationStyle:  tfSchemaList["communicationstyle"].(string),
		ProxyID:             proxyid,
		DefaultWorkerPoolID: tfSchemaList["defaultworkerpoolid"].(string),
	}

	tfMachine.Endpoint.ClusterCertificate = tfSchemaList["clustercertificate"].(string)
	tfMachine.Endpoint.ClusterURL = tfSchemaList["clusterurl"].(string)
	tfMachine.Endpoint.Namespace = tfSchemaList["namespace"].(string)
	tfMachine.Endpoint.SkipTLSVerification = strconv.FormatBool(tfSchemaList["skiptlsverification"].(bool))

	tfAuthenticationSchemaSetInterface, ok := tfSchemaList["authentication"]
	if ok {
		tfAuthenticationSchemaSet := tfAuthenticationSchemaSetInterface.(*schema.Set)
		if len(tfAuthenticationSchemaSet.List()) == 1 {
			//Get the first element in the list, which is a map of the interfaces
			tfAuthenticationSchemaList := tfAuthenticationSchemaSet.List()[0].(map[string]interface{})

			tfMachine.Endpoint.Authentication = &model.MachineEndpointAuthentication{
				AccountID:          tfAuthenticationSchemaList["accountid"].(string),
				ClientCertificate:  tfAuthenticationSchemaList["clientcertificate"].(string),
				AuthenticationType: tfAuthenticationSchemaList["authenticationtype"].(string),
			}
		}
	}

	return tfMachine
}

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	newMachine := buildMachineResource(d)
	newMachine.Status = "Unknown" //We don't want TF to attempt to update a machine just because its status has changed, so set it to Unknown on creation and let TF sort it out in the future.
	machine, err := apiClient.Machines.Add(newMachine)
	if err != nil {
		return fmt.Errorf("error creating machine %s: %s", newMachine.Name, err.Error())
	}
	d.SetId(machine.ID)
	setMachineProperties(d, machine)
	return nil
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	machineID := d.Id()
	err := apiClient.Machines.Delete(machineID)
	if err != nil {
		return fmt.Errorf("error deleting machine id %s: %s", machineID, err.Error())
	}
	d.SetId("")
	return nil
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {
	machine := buildMachineResource(d)
	machine.ID = d.Id() // set project struct ID so octopus knows which project to update
	apiClient := m.(*client.Client)
	updatedMachine, err := apiClient.Machines.Update(machine)
	if err != nil {
		return fmt.Errorf("error updating machine id %s: %s", d.Id(), err.Error())
	}
	d.SetId(updatedMachine.ID)
	setMachineProperties(d, machine)
	return nil
}
