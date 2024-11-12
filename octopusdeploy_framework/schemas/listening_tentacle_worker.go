package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ListeningTentacleWorkerSchema struct{}

func (m ListeningTentacleWorkerSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a listening tentacle worker in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":                GetIdResourceSchema(),
			"name":              GetNameResourceSchema(true),
			"space_id":          GetSpaceIdResourceSchema("Listening tentacle worker"),
			"is_disabled":       GetOptionalBooleanResourceAttribute("When disabled, worker will not be included in any deployments", false),
			"machine_policy_id": GetRequiredStringResourceSchema("Select the machine policy"),
			"uri":               GetRequiredStringResourceSchema("The network address at which the Tentacle can be reached"),
			"thumbprint":        GetRequiredStringResourceSchema("The X509 certificate thumbprint that securely identifies the Tentacle"),
			"proxy_id":          GetOptionalStringResourceSchema("Specify the connection type for the Tentacle: direct(when not set) or via a proxy server."),
			"worker_pool_ids": resourceSchema.ListAttribute{
				ElementType: types.StringType,
				Description: "Select at least one worker pool for the worker",
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (m ListeningTentacleWorkerSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

var _ EntitySchema = ListeningTentacleWorkerSchema{}

type ListeningTentacleWorkerResourceModel struct {
	Name            types.String `tfsdk:"name"`
	SpaceID         types.String `tfsdk:"space_id"`
	IsDisabled      types.Bool   `tfsdk:"is_disabled"`
	WorkerPoolIDs   types.List   `tfsdk:"worker_pool_ids"`
	MachinePolicyID types.String `tfsdk:"machine_policy_id"`
	Uri             types.String `tfsdk:"uri"`
	Thumbprint      types.String `tfsdk:"thumbprint"`
	ProxyID         types.String `tfsdk:"proxy_id"`

	ResourceModel
}
