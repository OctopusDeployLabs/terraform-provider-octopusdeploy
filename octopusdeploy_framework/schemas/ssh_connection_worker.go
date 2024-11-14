package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SSHConnectionWorkerSchema struct{}

func (m SSHConnectionWorkerSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a SSH connection worker in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":                GetIdResourceSchema(),
			"space_id":          GetSpaceIdResourceSchema("Listening tentacle worker"),
			"name":              GetNameResourceSchema(true),
			"is_disabled":       GetOptionalBooleanResourceAttribute("When disabled, worker will not be included in any deployments", false),
			"machine_policy_id": GetRequiredStringResourceSchema("Select the machine policy"),
			"worker_pool_ids": resourceSchema.SetAttribute{
				ElementType: types.StringType,
				Description: "Select at least one worker pool for the worker",
				Required:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"account_id":  GetRequiredStringResourceSchema("Connection account"),
			"host":        GetRequiredStringResourceSchema("The hostname or IP address of the deployment target to connect to"),
			"port":        GetPortNumberResourceSchema(),
			"fingerprint": GetRequiredStringResourceSchema("The host fingerprint to be verified"),
			"proxy_id":    GetOptionalStringResourceSchema("Specify the connection type for the Tentacle: direct(when not set) or via a proxy server."),
			"dotnet_platform": resourceSchema.StringAttribute{
				Description: "NET Core platform of self-contained version of Calamari",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("linux-arm", "linux-arm64", "linux-x64", "osx-x64"),
				},
			},
		},
	}
}

func (m SSHConnectionWorkerSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

var _ EntitySchema = SSHConnectionWorkerSchema{}

type SSHConnectionWorkerResourceModel struct {
	Name            types.String `tfsdk:"name"`
	SpaceID         types.String `tfsdk:"space_id"`
	IsDisabled      types.Bool   `tfsdk:"is_disabled"`
	MachinePolicyID types.String `tfsdk:"machine_policy_id"`
	WorkerPoolIDs   types.Set    `tfsdk:"worker_pool_ids"`
	AccountId       types.String `tfsdk:"account_id"`
	Host            types.String `tfsdk:"host"`
	Port            types.Int64  `tfsdk:"port"`
	Fingerprint     types.String `tfsdk:"fingerprint"`
	ProxyID         types.String `tfsdk:"proxy_id"`
	DotnetPlatform  types.String `tfsdk:"dotnet_platform"`

	ResourceModel
}
