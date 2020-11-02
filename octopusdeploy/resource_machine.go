package octopusdeploy

import (
	"log"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
							ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
								"None",
								"TentaclePassive",
								"TentacleActive",
								"Ssh",
								"OfflineDrop",
								"AzureWebApp",
								"Ftp",
								"AzureCloudService",
								"Kubernetes",
							}, false)),
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
										ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
											"KubernetesCertificate",
											"KubernetesStandard",
										}, false)),
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
				ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
					"Untenanted",
					"TenantedOrUntenanted",
					"Tenanted",
				}, false)),
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

	client := m.(*octopusdeploy.Client)
	resource, err := client.Machines.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingMachine, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constMachine, m)

	d.SetId(resource.GetID())
	setMachineProperties(d, resource)

	return nil
}

func setMachineProperties(d *schema.ResourceData, m *octopusdeploy.DeploymentTarget) {
	d.Set(constEnvironments, m.EnvironmentIDs)
	d.Set("haslatestcalamari", m.HasLatestCalamari)
	d.Set("isdisabled", m.IsDisabled)
	d.Set("isinprocess", m.IsInProcess)
	d.Set("machinepolicy", m.MachinePolicyID)
	d.Set(constRoles, m.Roles)
	d.Set(constStatus, m.Status)
	d.Set("statussummary", m.StatusSummary)
	d.Set("tenanteddeploymentparticipation", m.TenantedDeploymentMode)
	d.Set("tenantids", m.TenantIDs)
	d.Set("tenanttags", m.TenantTags)
}

func buildMachineResource(d *schema.ResourceData) *octopusdeploy.DeploymentTarget {
	name := d.Get(constName).(string)
	machinePolicy := d.Get("machinepolicy").(string)
	environments := getSliceFromTerraformTypeList(d.Get(constEnvironments))
	roles := getSliceFromTerraformTypeList(d.Get(constRoles))
	isDisabled := d.Get("isdisabled").(bool)
	deploymentMode := octopusdeploy.TenantedDeploymentMode(d.Get("tenanteddeploymentparticipation").(string))
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

	var proxyID string
	if tfSchemaList["proxyid"] != nil {
		proxyString := tfSchemaList["proxyid"].(string)
		proxyID = proxyString
	}

	communicationStyle := octopusdeploy.CommunicationStyle(tfSchemaList["communicationstyle"].(string))

	var endpoint octopusdeploy.IEndpoint
	switch communicationStyle {
	case "AzureCloudService":
		azureCloudServiceEndpoint := octopusdeploy.NewAzureCloudServiceEndpoint()
		azureCloudServiceEndpoint.DefaultWorkerPoolID = tfSchemaList["defaultworkerpoolid"].(string)
		endpoint = azureCloudServiceEndpoint
	case "AzureServiceFabricCluster":
		endpoint = octopusdeploy.NewServiceFabricEndpoint()
	case "AzureWebApp":
		endpoint = octopusdeploy.NewAzureWebAppEndpoint()
	case "Kubernetes":
		clusterURL := d.Get("clusterURL").(url.URL)
		kubernetesEndpoint := octopusdeploy.NewKubernetesEndpoint(clusterURL)
		kubernetesEndpoint.ClusterCertificate = tfSchemaList["clustercertificate"].(string)
		kubernetesEndpoint.ClusterURL, _ = url.Parse(tfSchemaList["clusterurl"].(string))
		kubernetesEndpoint.Namespace = tfSchemaList[constNamespace].(string)
		kubernetesEndpoint.ProxyID = proxyID
		kubernetesEndpoint.SkipTLSVerification = tfSchemaList["skiptlsverification"].(bool)
		endpoint = kubernetesEndpoint
	case "None":
		endpoint = octopusdeploy.NewCloudRegionEndpoint()
	case "OfflineDrop":
		endpoint = octopusdeploy.NewOfflineDropEndpoint()
	case "Ssh":
		host := d.Get("host").(string)
		port := d.Get("port").(int)
		fingerprint := d.Get("fingerprint").(string)
		sshEndpoint := octopusdeploy.NewSSHEndpoint(host, port, fingerprint)
		sshEndpoint.ProxyID = proxyID
		endpoint = sshEndpoint
	case "TentacleActive":
		uri, _ := url.Parse(tfSchemaList[constURI].(string))
		thumbprint := tfSchemaList[constThumbprint].(string)
		endpoint = octopusdeploy.NewPollingTentacleEndpoint(uri, thumbprint)
	case "TentaclePassive":
		uri, _ := url.Parse(tfSchemaList[constURI].(string))
		thumbprint := tfSchemaList[constThumbprint].(string)
		endpoint = octopusdeploy.NewListeningTentacleEndpoint(uri, thumbprint)
	}

	machine := octopusdeploy.NewDeploymentTarget(name, endpoint, environments, roles)
	machine.TenantedDeploymentMode = deploymentMode
	machine.IsDisabled = isDisabled
	machine.MachinePolicyID = machinePolicy
	machine.TenantIDs = tenantIDs
	machine.TenantTags = tenantTags
	machine.Thumbprint = tfSchemaList[constThumbprint].(string)
	machine.URI = tfSchemaList[constURI].(string)

	tfAuthenticationSchemaSetInterface, ok := tfSchemaList[constAuthentication]
	if ok {
		tfAuthenticationSchemaSet := tfAuthenticationSchemaSetInterface.(*schema.Set)
		if len(tfAuthenticationSchemaSet.List()) == 1 {
			// Get the first element in the list, which is a map of the interfaces
			// tfAuthenticationSchemaList := tfAuthenticationSchemaSet.List()[0].(map[string]interface{})

			// machine.Endpoint.Authentication = &octopusdeploy.MachineEndpointAuthentication{
			// 	AccountID:          tfAuthenticationSchemaList["accountid"].(string),
			// 	ClientCertificate:  tfAuthenticationSchemaList["clientcertificate"].(string),
			// 	AuthenticationType: tfAuthenticationSchemaList["authenticationtype"].(string),
			// }
		}
	}

	return machine
}

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {
	machine := buildMachineResource(d)
	machine.Status = "Unknown" // We don't want TF to attempt to update a machine just because its status has changed, so set it to Unknown on creation and let TF sort it out in the future.

	client := m.(*octopusdeploy.Client)
	resource, err := client.Machines.Add(machine)
	if err != nil {
		return createResourceOperationError(errorCreatingMachine, machine.Name, err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	setMachineProperties(d, resource)

	return nil
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingMachine, id, err)
	}

	d.SetId(constEmptyString)

	return nil
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {
	machine := buildMachineResource(d)
	machine.ID = d.Id() // set ID so Octopus API knows which machine to update

	client := m.(*octopusdeploy.Client)
	updatedMachine, err := client.Machines.Update(machine)
	if err != nil {
		return createResourceOperationError(errorUpdatingMachine, d.Id(), err)
	}

	d.SetId(updatedMachine.ID)
	setMachineProperties(d, machine)

	return nil
}
