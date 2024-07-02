package pluginprovider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
)

var _ provider.Provider = &OctopusDeployProvider{}

type Config struct {
	Address string
	APIKey  string
	SpaceID string
}

type OctopusDeployProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func GetDefaultFromEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (p *OctopusDeployProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config Config

	// Retrieve provider data from the request
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize the client
	client, diags := config.Client()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the client for the provider
	resp.DataSourceData = client
	resp.ResourceData = client

}

// OctopusDeployProviderModel describes the provider data model.
type OctopusDeployProviderModel struct {
	Address types.String `tfsdk:"address"`
	APIKey  types.String `tfsdk:"api_key"`
	SpaceID types.String `tfsdk:"email"`
}

func (p *OctopusDeployProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "octopusdeploy"
	resp.Version = p.version
}

func (p *OctopusDeployProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Description: "The endpoint of the Octopus REST API",
				Optional:    true,  // Use Optional with a default value
				Required:    false, // Required should be false since it is Optional
				Validators:  []validator.String{},
			},
			"api_key": schema.StringAttribute{
				Description: "The endpoint of the Octopus REST API",
				Required:    true, // Required should be false since it is Optional
				Validators:  []validator.String{},
			},
			"space_id": schema.StringAttribute{
				Description: "The endpoint of the Octopus REST API",
				Optional:    true, // Use Optional with a default value
				Validators:  []validator.String{},
			},
		},
	}
}

func (p *OctopusDeployProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// NewExampleResource,
	}
}

func (p *OctopusDeployProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// NewExampleDataSource,
	}
}

func Provider() provider.Provider {
	return &OctopusDeployProvider{}
}
