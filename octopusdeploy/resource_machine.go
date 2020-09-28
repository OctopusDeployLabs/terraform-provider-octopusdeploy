package octopusdeploy

import (
	"log"
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/enum"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceMachineCreate,
		Read:   resourceMachineRead,
		Update: resourceMachineUpdate,
		Delete: resourceMachineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constEndpoint: {
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
						constThumbprint: {
							Type:     schema.TypeString,
							Required: true,
						},
						constURI: {
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
						constNamespace: {
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
						constAuthentication: {
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
			constEnvironments: {
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
			constRoles: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					MinItems: 1,
				},
				Required: true,
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
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Machines.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingMachine, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constMachine, m)

	d.SetId(resource.ID)
	setMachineProperties(d, resource)

	return nil
}

func setMachineProperties(d *schema.ResourceData, m *model.Machine) {
	d.Set(constEnvironments, m.EnvironmentIDs)
	d.Set("haslatestcalamari", m.HasLatestCalamari)
	d.Set("isdisabled", m.IsDisabled)
	d.Set("isinprocess", m.IsInProcess)
	d.Set("machinepolicy", m.MachinePolicyID)
	d.Set(constRoles, m.Roles)
	d.Set(constStatus, m.Status)
	d.Set("statussummary", m.StatusSummary)
	d.Set("tenanteddeploymentparticipation", m.DeploymentMode)
	d.Set("tenantids", m.TenantIDs)
	d.Set("tenanttags", m.TenantTags)
}

func buildMachineResource(d *schema.ResourceData) *model.Machine {
	name := d.Get(constName).(string)
	machinePolicy := d.Get("machinepolicy").(string)
	environments := getSliceFromTerraformTypeList(d.Get(constEnvironments))
	roles := getSliceFromTerraformTypeList(d.Get(constRoles))
	isDisabled := d.Get("isdisabled").(bool)
	deploymentMode, _ := d.Get("tenanteddeploymentparticipation").(string)
	tenantIDs := getSliceFromTerraformTypeList(d.Get("tenantids"))
	tenantTags := getSliceFromTerraformTypeList(d.Get("tenanttags"))

	// If we end up with a nil return, Octopus doesn't accept the API call. This ensure that we send
	// blank values rather than nil values.
	if tenantIDs == nil {
		tenantIDs = []string{}
	}
	if tenantTags == nil {
		tenantTags = []string{}
	}

	tfSchemaSetInterface, ok := d.GetOk(constEndpoint)
	if !ok {
		return nil
	}
	tfSchemaSet := tfSchemaSetInterface.(*schema.Set)
	if len(tfSchemaSet.List()) == 0 {
		return nil
	}
	// Get the first element in the list, which is a map of the interfaces
	tfSchemaList := tfSchemaSet.List()[0].(map[string]interface{})

	machine, err := model.NewMachine(name, isDisabled, environments, roles, machinePolicy, deploymentMode, tenantIDs, tenantTags)
	if err != nil {
		return nil
	}

	machine.URI = tfSchemaList[constURI].(string)
	machine.Thumbprint = tfSchemaList[constThumbprint].(string)

	var proxyID string
	if tfSchemaList["proxyid"] != nil {
		proxyString := tfSchemaList["proxyid"].(string)
		proxyID = proxyString
	}

	communicationStyle, err := enum.ParseCommunicationStyle(tfSchemaList["communicationstyle"].(string))
	if err != nil {
		return nil
	}

	endpoint, err := model.NewMachineEndpoint(
		tfSchemaList[constURI].(string),
		tfSchemaList[constThumbprint].(string),
		communicationStyle,
		proxyID,
		tfSchemaList["defaultworkerpoolid"].(string),
	)
	if err != nil {
		return nil
	}

	machine.Endpoint = endpoint
	machine.Endpoint.ClusterCertificate = tfSchemaList["clustercertificate"].(string)
	machine.Endpoint.ClusterURL = tfSchemaList["clusterurl"].(string)
	machine.Endpoint.Namespace = tfSchemaList[constNamespace].(string)
	machine.Endpoint.SkipTLSVerification = strconv.FormatBool(tfSchemaList["skiptlsverification"].(bool))

	tfAuthenticationSchemaSetInterface, ok := tfSchemaList[constAuthentication]
	if ok {
		tfAuthenticationSchemaSet := tfAuthenticationSchemaSetInterface.(*schema.Set)
		if len(tfAuthenticationSchemaSet.List()) == 1 {
			// Get the first element in the list, which is a map of the interfaces
			tfAuthenticationSchemaList := tfAuthenticationSchemaSet.List()[0].(map[string]interface{})

			machine.Endpoint.Authentication = &model.MachineEndpointAuthentication{
				AccountID:          tfAuthenticationSchemaList["accountid"].(string),
				ClientCertificate:  tfAuthenticationSchemaList["clientcertificate"].(string),
				AuthenticationType: tfAuthenticationSchemaList["authenticationtype"].(string),
			}
		}
	}

	return machine
}

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {
	machine := buildMachineResource(d)
	machine.Status = "Unknown" // We don't want TF to attempt to update a machine just because its status has changed, so set it to Unknown on creation and let TF sort it out in the future.

	apiClient := m.(*client.Client)
	resource, err := apiClient.Machines.Add(machine)
	if err != nil {
		return createResourceOperationError(errorCreatingMachine, machine.Name, err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	setMachineProperties(d, resource)

	return nil
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.Machines.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingMachine, id, err)
	}

	d.SetId(constEmptyString)

	return nil
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {
	machine := buildMachineResource(d)
	machine.ID = d.Id() // set ID so Octopus API knows which machine to update

	apiClient := m.(*client.Client)
	updatedMachine, err := apiClient.Machines.Update(machine)
	if err != nil {
		return createResourceOperationError(errorUpdatingMachine, d.Id(), err)
	}

	d.SetId(updatedMachine.ID)
	setMachineProperties(d, machine)

	return nil
}
