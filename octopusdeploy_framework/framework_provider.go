package octopusdeploy_framework

import (
	"context"
	"os"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type octopusDeployFrameworkProvider struct {
	Address types.String `tfsdk:"address"`
	ApiKey  types.String `tfsdk:"api_key"`
	SpaceID types.String `tfsdk:"space_id"`
}

var _ provider.Provider = (*octopusDeployFrameworkProvider)(nil)
var _ provider.ProviderWithMetaSchema = (*octopusDeployFrameworkProvider)(nil)
var _ provider.ProviderWithFunctions

func NewOctopusDeployFrameworkProvider() *octopusDeployFrameworkProvider {
	return &octopusDeployFrameworkProvider{}
}

func (p *octopusDeployFrameworkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = util.GetProviderName()
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

	resp.DataSourceData = &config
	resp.ResourceData = &config
}

func (p *octopusDeployFrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectGroupsDataSource,
		NewSpaceDataSource,
		NewSpacesDataSource,
		NewLifecyclesDataSource,
		NewEnvironmentsDataSource,
		NewGitCredentialsDataSource,
		NewFeedsDataSource,
		NewLibraryVariableSetDataSource,
		NewVariablesDataSource,
		NewProjectsDataSource,
		NewScriptModuleDataSource,
	}
}

func (p *octopusDeployFrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSpaceResource,
		NewProjectGroupResource,
		NewMavenFeedResource,
		NewLifecycleResource,
		NewEnvironmentResource,
		NewGitCredentialResource,
		NewHelmFeedResource,
		NewArtifactoryGenericFeedResource,
		NewGitHubRepositoryFeedResource,
		NewAwsElasticContainerRegistryFeedResource,
		NewNugetFeedResource,
		NewTenantProjectResource,
		NewTenantProjectVariableResource,
		NewTenantCommonVariableResource,
		NewLibraryVariableSetFeedResource,
		NewVariableResource,
		NewProjectResource,
		NewScriptModuleResource,
	}
}

func (p *octopusDeployFrameworkProvider) Schema(_ context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Optional:    true,
				Description: "The endpoint of the Octopus REST API",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Description: "The API key to use with the Octopus REST API",
			},
			"space_id": schema.StringAttribute{
				Optional:    true,
				Description: "The space ID to target",
			},
		},
	}
}
