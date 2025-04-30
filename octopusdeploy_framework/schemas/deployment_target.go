package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const DeploymentTargetDescription = "Deployment target"

type DeploymentTargetSchema struct{}

type DeploymentTargetModel struct {
	Endpoint                        types.Object `tfsdk:"endpoint"`
	Environments                    types.List   `tfsdk:"environments"`
	HasLatestCalamari               types.Bool   `tfsdk:"has_latest_calamari"`
	HealthStatus                    types.String `tfsdk:"health_status"`
	IsDisabled                      types.Bool   `tfsdk:"is_disabled"`
	IsInProcess                     types.Bool   `tfsdk:"is_in_process"`
	MachinePolicyId                 types.String `tfsdk:"machine_policy_id"`
	Name                            types.String `tfsdk:"name"`
	OperatingSystem                 types.String `tfsdk:"operating_system"`
	Roles                           types.List   `tfsdk:"roles"`
	ShellName                       types.String `tfsdk:"shell_name"`
	ShellVersion                    types.String `tfsdk:"shell_version"`
	SpaceId                         types.String `tfsdk:"space_id"`
	Status                          types.String `tfsdk:"status"`
	StatusSummary                   types.String `tfsdk:"status_summary"`
	TenantedDeploymentParticipation types.String `tfsdk:"tenanted_deployment_participation"`
	Tenants                         types.List   `tfsdk:"tenants"`
	TenantTags                      types.List   `tfsdk:"tenant_tags"`
	Thumbprint                      types.String `tfsdk:"thumbprint"`
	Uri                             types.String `tfsdk:"uri"`

	ResourceModel
}

func (d DeploymentTargetSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(DeploymentTargetDescription),
		Attributes: map[string]resourceSchema.Attribute{
			"endpoint": resourceSchema.SingleNestedAttribute{
				Attributes: GetEndpointResourceSchema(),
			},
			"environments": resourceSchema.ListAttribute{
				Description: "A list of environment IDs associated with this resource.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"has_latest_calamari": resourceSchema.BoolAttribute{
				Computed: true,
			},
			"health_status": GetHealthStatusResourceSchema(),
			"id":            GetIdResourceSchema(),
			"is_disabled": resourceSchema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"is_in_process": resourceSchema.BoolAttribute{
				Computed: true,
			},
			"machine_policy_id": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"name": GetNameResourceSchema(true),
			"operating_system": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"roles": resourceSchema.ListAttribute{
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"shell_name": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"shell_version": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"space_id":       GetSpaceIdResourceSchema(DeploymentTargetDescription),
			"status":         GetStatusResourceSchema(),
			"status_summary": GetStatusSummaryResourceSchema(),
			"tenanted_deployment_participation": resourceSchema.StringAttribute{
				Description: "The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Untenanted",
						"TenantedOrUntenanted",
						"Tenanted",
					),
				},
			},
			"tenants": resourceSchema.ListAttribute{
				Description: "A list of tenant IDs associated with this certificate.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"tenant_tags": resourceSchema.ListAttribute{
				Description: "A list of tenant tags associated with this certificate.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"thumbprint": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"uri": resourceSchema.StringAttribute{
				Computed: true,
				Optional: true,
			},
		},
	}
}
