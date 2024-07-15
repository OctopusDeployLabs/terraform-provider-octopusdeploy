package octopusdeploy_framework

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
)

type octopusDeployFrameworkProvider struct {
	Address types.String `tfsdk:"address"`
	ApiKey  types.String `tfsdk:"api_key"`
	SpaceID types.String `tfsdk:"space_id"`
}

var _ provider.Provider = (*octopusDeployFrameworkProvider)(nil)
var _ provider.ProviderWithMetaSchema = (*octopusDeployFrameworkProvider)(nil)
var ProviderTypeName = "octopusdeploy"

func NewOctopusDeployFrameworkProvider() *octopusDeployFrameworkProvider {
	return &octopusDeployFrameworkProvider{}
}

func (p *octopusDeployFrameworkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = ProviderTypeName
}

func (p *octopusDeployFrameworkProvider) MetaSchema(ctx context.Context, request provider.MetaSchemaRequest, response *provider.MetaSchemaResponse) {

}

func (p *octopusDeployFrameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var providerData octopusDeployFrameworkProvider
	resp.Diagnostics.Append(req.Config.Get(ctx, &providerData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := Config{}
	config.ApiKey = providerData.ApiKey.ValueString()
	if config.ApiKey == "" {
		config.ApiKey = os.Getenv("OCTOPUS_APIKEY")
	}
	config.Address = providerData.Address.ValueString()
	if config.Address == "" {
		config.Address = os.Getenv("OCTOPUS_URL")
	}
	config.SpaceID = providerData.SpaceID.ValueString()
	if err := config.GetClient(ctx); err != nil {
		resp.Diagnostics.AddError("failed to load client", err.Error())
	}
	if err := config.GetClient(ctx); err != nil {
		resp.Diagnostics.AddError("failed to load client", err.Error())
	}

	resp.DataSourceData = &config
	resp.ResourceData = &config
}

func (p *octopusDeployFrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectGroupsDataSource,
		NewSpaceDataSource,
		NewSpacesDataSource,
		NewLifecyclesDataSource,
	}
}

func (p *octopusDeployFrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	resources := []func() resource.Resource{
		NewProjectGroupResource,
	}

	resources = append(resources, getImportableResources()...)

	return resources
}

func getImportableResources() []func() resource.Resource {
	importableResources := []func() resource.ResourceWithImportState{
		NewLifecycleResource,
	}

	resources := make([]func() resource.Resource, 0, len(importableResources))

	for _, r := range importableResources {
		resourceFunc := r
		resources = append(resources, func() resource.Resource {
			return resourceFunc()
		})
	}
	return resources
}

func (p *octopusDeployFrameworkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Required:    true,
				Description: "The endpoint of the Octopus REST API",
			},
			"api_key": schema.StringAttribute{
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
