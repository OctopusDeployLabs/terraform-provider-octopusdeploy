package octopusdeployv6

import (
	"context"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeployv6/config"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeployv6/project_group"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OctopusDeployProviderV6 struct {
	Address types.String `tfsdk:"address"`
	ApiKey  types.String `tfsdk:"api_key"`
	SpaceID types.String `tfsdk:"space_id"`
}

var _ provider.Provider = (*OctopusDeployProviderV6)(nil)
var _ provider.ProviderWithMetaSchema = (*OctopusDeployProviderV6)(nil)

func NewOctopusDeployProviderV6() *OctopusDeployProviderV6 {
	return &OctopusDeployProviderV6{}
}

func (p *OctopusDeployProviderV6) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "octopus-deploy"
}

func (p *OctopusDeployProviderV6) MetaSchema(ctx context.Context, request provider.MetaSchemaRequest, response *provider.MetaSchemaResponse) {

}

func (p *OctopusDeployProviderV6) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var providerData OctopusDeployProviderV6
	resp.Diagnostics.Append(req.Config.Get(ctx, &providerData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := config.Config{}
	config.ApiKey = providerData.ApiKey.ValueString()
	config.Address = providerData.Address.ValueString()
	config.SpaceID = providerData.SpaceID.ValueString()
	if err := config.GetClient(ctx); err != nil {
		resp.Diagnostics.AddError("failed to load client", err.Error())
	}

	resp.DataSourceData = &config
	resp.ResourceData = &config
}

func (p *OctopusDeployProviderV6) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *OctopusDeployProviderV6) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		project_group.NewProjectGroupResource,
	}
}

func (p *OctopusDeployProviderV6) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Optional:    false,
				Required:    true,
				Description: "The endpoint of the Octopus REST API",
			},
			"api_key": schema.StringAttribute{
				Optional:    false,
				Required:    true,
				Description: "The API key to use with the Octopus REST API",
			},
			"space_id": schema.StringAttribute{
				Optional:    true,
				Description: "The space ID to target",
			},
		},
	}
}
