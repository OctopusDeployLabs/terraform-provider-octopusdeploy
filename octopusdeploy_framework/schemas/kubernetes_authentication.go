package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const KubernetesAuthenticationDescription = "Kubernetes authentication"

type KubernetesAuthenticationSchema struct{}

type KubernetesAuthenticationModel struct {
	AccountId                 string `tfsdk:"account_id"`
	AdminLogin                string `tfsdk:"admin_login"`
	AssumedRoleArn            string `tfsdk:"assumed_role_arn"`
	AssumedRoleSession        string `tfsdk:"assumed_role_session"`
	AssumeRole                bool   `tfsdk:"assume_role"`
	AssumeRoleSessionDuration int    `tfsdk:"assume_role_session_duration"`
	AssumeRoleExternalId      string `tfsdk:"assume_role_external_id"`
	AuthenticationType        string `tfsdk:"authentication_type"`
	ClientCertificate         string `tfsdk:"client_certificate"`
	ClusterName               string `tfsdk:"cluster_name"`
	ClusterResourceGroup      string `tfsdk:"cluster_resource_group"`
	ImpersonateServiceAccount bool   `tfsdk:"impersonate_service_account"`
	Project                   string `tfsdk:"project"`
	Region                    string `tfsdk:"region"`
	ServiceAccountEmails      string `tfsdk:"service_account_emails"`
	UseInstanceRole           bool   `tfsdk:"use_instance_role"`
	UseVmServiceAccount       bool   `tfsdk:"use_vm_service_account"`
	Zone                      string `tfsdk:"zone"`
	TokenPath                 string `tfsdk:"token_path"`

	ResourceModel
}

func (d KubernetesAuthenticationSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(KubernetesAuthenticationDescription),
		Attributes:  GetKubernetesAuthenticationResourceSchema(),
	}
}

func GetKubernetesAuthenticationResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"account_id": resourceSchema.StringAttribute{
			Optional: true,
		},
		"admin_login": resourceSchema.StringAttribute{
			Optional: true,
		},
		"assumed_role_arn": resourceSchema.StringAttribute{
			Optional: true,
		},
		"assumed_role_session": resourceSchema.StringAttribute{
			Optional: true,
		},
		"assume_role": resourceSchema.BoolAttribute{
			Optional: true,
		},
		"assume_role_session_duration": resourceSchema.Int64Attribute{
			Optional: true,
		},
		"assume_role_external_id": resourceSchema.StringAttribute{
			Optional: true,
		},
		"authentication_type": resourceSchema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					"KubernetesAws",
					"KubernetesAzure",
					"KubernetesCertificate",
					"KubernetesGoogleCloud",
					"KubernetesStandard",
					"KubernetesPodService",
					"None",
				),
			},
		},
		"client_certificate": resourceSchema.StringAttribute{
			Optional: true,
		},
		"cluster_name": resourceSchema.StringAttribute{
			Optional: true,
		},
		"cluster_resource_group": resourceSchema.StringAttribute{
			Optional: true,
		},
		"impersonate_service_account": resourceSchema.BoolAttribute{
			Optional: true,
		},
		"project": resourceSchema.StringAttribute{
			Optional: true,
		},
		"region": resourceSchema.StringAttribute{
			Optional: true,
		},
		"service_account_emails": resourceSchema.StringAttribute{
			Optional: true,
		},
		"use_instance_role": resourceSchema.BoolAttribute{
			Optional: true,
		},
		"use_vm_service_account": resourceSchema.BoolAttribute{
			Optional: true,
		},
		"zone": resourceSchema.StringAttribute{
			Optional: true,
		},
		"token_path": resourceSchema.StringAttribute{
			Optional: true,
		},
	}
}
