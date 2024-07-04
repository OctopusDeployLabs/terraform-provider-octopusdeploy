package octopusdeploy

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

type octopusDeployProviderV6 struct {
	Address types.String `tfsdk:"address"`
	ApiKey  types.String `tfsdk:"api_key"`
	SpaceID types.String `tfsdk:"space_id"`
}

var _ provider.Provider = (*octopusDeployProviderV6)(nil)
var _ provider.ProviderWithMetaSchema = (*octopusDeployProviderV6)(nil)

func NewOctopusDeployProviderV6() *octopusDeployProviderV6 {
	return &octopusDeployProviderV6{}
}

func (p *octopusDeployProviderV6) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "octopus-deploy"
}

func (p *octopusDeployProviderV6) MetaSchema(ctx context.Context, request provider.MetaSchemaRequest, response *provider.MetaSchemaResponse) {

}

func (p *octopusDeployProviderV6) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var providerData octopusDeployProviderV6
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

func (p *octopusDeployProviderV6) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *octopusDeployProviderV6) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		project_group.NewProjectGroupResource,
	}
}

func (p *octopusDeployProviderV6) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
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
