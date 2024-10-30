package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MachineProxyResourceModel struct {
	SpaceID  types.String `tfsdk:"space_id"`
	Name     types.String `tfsdk:"name"`
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Port     types.Int32  `tfsdk:"port"`
	Password types.String `tfsdk:"password"`
	ResourceModel
}

type MachineProxyDataSourceModel struct {
	ID          types.String           `tfsdk:"id"`
	SpaceID     types.String           `tfsdk:"space_id"`
	IDs         types.List             `tfsdk:"ids"`
	PartialName types.String           `tfsdk:"partial_name"`
	Skip        types.Int64            `tfsdk:"skip"`
	Take        types.Int64            `tfsdk:"take"`
	Proxies     []ProxyDatasourceModel `tfsdk:"machine_proxies"`
}

type ProxyDatasourceModel struct {
	ID       types.String `tfsdk:"id"`
	SpaceID  types.String `tfsdk:"space_id"`
	Name     types.String `tfsdk:"name"`
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Port     types.Int32  `tfsdk:"port"`
}

const (
	MachineProxyResourceName   = "machine_proxy"
	MachineProxyDataSourceName = "machine_proxies"
)

type MachineProxySchema struct{}

var _ EntitySchema = MachineProxySchema{}

func (p MachineProxySchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages machine proxies in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":       GetIdResourceSchema(),
			"space_id": GetSpaceIdResourceSchema(MachineProxyResourceName),
			"name":     GetNameResourceSchema(true),
			"host": util.ResourceString().
				Required().
				Description("DNS hostname of the proxy server").
				Build(),
			"username": util.ResourceString().
				Required().
				Description("Username of the proxy server").
				Build(),
			"password": util.ResourceString().
				Required().
				Description("Password of the proxy server").
				Sensitive().
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"port": util.ResourceInt32().
				Optional().
				Computed().
				Default(80).
				Description("The port number for the proxy server.").
				Build(),
		},
	}
}

func (p MachineProxySchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing Octopus Deploy machine proxies.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":              util.DataSourceString().Computed().Description("An auto-generated identifier that includes the timestamp when this data source was last modified.").Build(),
			"space_id":        util.DataSourceString().Optional().Description("A Space ID to filter by. Will revert what is specified on the provider if not set").Build(),
			"ids":             GetQueryIDsDatasourceSchema(),
			"partial_name":    GetQueryPartialNameDatasourceSchema(),
			"skip":            GetQuerySkipDatasourceSchema(),
			"take":            GetQueryTakeDatasourceSchema(),
			"machine_proxies": getMachineProxiesDataSourceAttribute(),
		},
	}
}

func getMachineProxiesDataSourceAttribute() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "A list of machine proxies that match the filter(s).",
		Computed:    true,
		Optional:    false,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"id": util.DataSourceString().Computed().Build(),
				"space_id": util.DataSourceString().
					Computed().
					Description("The space ID associated with this machine proxy.").
					Build(),
				"name": util.DataSourceString().Computed().Build(),
				"host": util.DataSourceString().
					Computed().
					Description("DNS hostname of the proxy server").
					Build(),
				"username": util.DataSourceString().
					Computed().
					Description("Username for the proxy server").
					Build(),
				"port": util.DataSourceInt64().
					Computed().
					Description("The port number for the proxy server.").
					Build(),
			},
		},
	}
}
