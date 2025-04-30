package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const EndpointDescription = "Endpoint"

type EndpointSchema struct{}

type EndpointModel struct {
	AddClientCredentialSecret     types.String `tfsdk:"aad_client_credential_secret"`
	AadCredentialType             types.String `tfsdk:"aad_credential_type"`
	AadUserCredentialUsername     types.String `tfsdk:"aad_user_credential_username"`
	AccountId                     types.String `tfsdk:"account_id"`
	ApplicationsDirectory         types.String `tfsdk:"applications_directory"`
	Authentication                types.List   `tfsdk:"authentication"`
	Container                     types.List   `tfsdk:"container"`
	ContainerOptions              types.String `tfsdk:"container_options"`
	CommunicationStyle            types.String `tfsdk:"communication_style"`
	ConnectionEndpoint            types.String `tfsdk:"connection_endpoint"`
	ClusterCertificate            types.String `tfsdk:"cluster_certificate"`
	ClusterCertificatePath        types.String `tfsdk:"cluster_certificate_path"`
	CloudServiceName              types.String `tfsdk:"cloud_service_name"`
	CertificateSignatureAlgorithm types.String `tfsdk:"certificate_signature_algorithm"`
	CertificateStoreLocation      types.String `tfsdk:"certificate_store_location"`
	CertificateStoreName          types.String `tfsdk:"certificate_store_name"`
	ClientCertificateVariable     types.String `tfsdk:"client_certificate_variable"`
	DefaultWorkerPoolId           types.String `tfsdk:"default_worker_pool_id"`
	Host                          types.String `tfsdk:"host"`
	Namespace                     types.String `tfsdk:"namespace"`
	ProxyId                       types.String `tfsdk:"proxy_id"`
	Port                          types.Int64  `tfsdk:"port"`
	ResourceGroupName             types.String `tfsdk:"resource_group_name"`
	RunningInContainer            types.Bool   `tfsdk:"running_in_container"`
	SecurityMode                  types.String `tfsdk:"security_mode"`
	ServerCertificateThumbprint   types.String `tfsdk:"server_certificate_thumbprint"`
	SkipTlsVerification           types.Bool   `tfsdk:"skip_tls_verification"`
	Thumbprint                    types.String `tfsdk:"thumbprint"`
	Uri                           types.String `tfsdk:"uri"`
	WebAppName                    types.String `tfsdk:"web_app_name"`
	WebAppSlotName                types.String `tfsdk:"web_app_slot_name"`

	ResourceModel
}

func (d EndpointSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(EndpointDescription),
		Attributes:  GetEndpointResourceSchema(),
	}
}

func GetEndpointResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"aad_client_credential_secret": resourceSchema.StringAttribute{
			Optional: true,
		},
		"aad_credential_type": resourceSchema.StringAttribute{
			Optional: true,
		},
		"aad_user_credential_username": resourceSchema.StringAttribute{
			Optional: true,
		},
		"account_id": resourceSchema.StringAttribute{
			Optional: true,
		},
		"applications_directory": resourceSchema.StringAttribute{
			Optional: true,
		},
		"authentication": resourceSchema.SetNestedAttribute{
			Computed: true,
			NestedObject: resourceSchema.NestedAttributeObject{
				Attributes: GetKubernetesAuthenticationResourceSchema(),
			},
			Validators: []validator.Set{
				setvalidator.SizeAtMost(1),
			},
		},
		"certificate_signature_algorithm": resourceSchema.StringAttribute{
			Optional: true,
		},
		"certificate_store_location": resourceSchema.StringAttribute{
			Optional: true,
		},
		"certificate_store_name": resourceSchema.StringAttribute{
			Optional: true,
		},
		"client_certificate_variable": resourceSchema.StringAttribute{
			Optional: true,
		},
		"cluster_certificate": resourceSchema.StringAttribute{
			Optional: true,
		},
		"cluster_certificate_path": resourceSchema.StringAttribute{
			Optional: true,
		},
		"cluster_url": resourceSchema.StringAttribute{
			Optional: true,
		},
		"communication_style": resourceSchema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					"AzureCloudService",
					"AzureWebApp",
					"Ftp",
					"Kubernetes",
					"None",
					"OfflineDrop",
					"Ssh",
					"TentacleActive",
					"TentaclePassive",
				),
			},
		},
		"connection_endpoint": resourceSchema.StringAttribute{
			Optional: true,
		},
		"container": resourceSchema.ListNestedAttribute{
			Computed:     true,
			NestedObject: GetDeploymentActionContainerResourceSchema(),
			Optional:     true,
		},
		"container_options": resourceSchema.StringAttribute{
			Optional: true,
		},
		"default_worker_pool_id": resourceSchema.StringAttribute{
			Optional: true,
		},
		"destination": resourceSchema.ListNestedAttribute{
			Computed:     true,
			NestedObject: GetOfflinePackageDropDestinationResourceSchema(),
			Optional:     true,
		},
		"dot_net_core_platform": resourceSchema.StringAttribute{
			Optional: true,
		},
		"fingerprint": resourceSchema.StringAttribute{
			Optional: true,
		},
		"host": resourceSchema.StringAttribute{
			Optional: true,
		},
		"id": GetIdResourceSchema(),
		"proxy_id": resourceSchema.StringAttribute{
			Optional: true,
		},
		"port": resourceSchema.Int64Attribute{
			Optional: true,
		},
		"resource_group_name": resourceSchema.StringAttribute{
			Optional: true,
		},
		"running_in_container": resourceSchema.BoolAttribute{
			Optional: true,
		},
		"security_mode": resourceSchema.StringAttribute{
			Optional: true,
		},
		"server_certificate_thumbprint": resourceSchema.StringAttribute{
			Optional: true,
		},
		"skip_tls_verification": resourceSchema.BoolAttribute{
			Optional: true,
		},
		"slot": resourceSchema.StringAttribute{
			Optional: true,
		},
		"storage_account_name": resourceSchema.StringAttribute{
			Optional: true,
		},
		"swap_if_possible": resourceSchema.BoolAttribute{
			Optional: true,
		},
		"tentacle_version_details": resourceSchema.ListNestedAttribute{
			Computed:     true,
			NestedObject: GetTentacleVersionDetailsResourceSchema(),
			Optional:     true,
		},
		"thumbprint": resourceSchema.StringAttribute{
			Optional: true,
		},
		"working_directory": resourceSchema.StringAttribute{
			Optional: true,
		},
		"use_current_instance_count": resourceSchema.BoolAttribute{
			Optional: true,
		},
		"uri": resourceSchema.StringAttribute{
			Optional: true,
		},
		"web_app_name": resourceSchema.StringAttribute{
			Optional: true,
		},
		"web_app_slot_name": resourceSchema.StringAttribute{
			Optional: true,
		},
	}
}
