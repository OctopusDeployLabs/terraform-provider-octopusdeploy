package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/configuration"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"go/version"
	"net/url"
)

type Config struct {
	Address        string
	ApiKey         string
	AccessToken    string
	SpaceID        string
	Client         *client.Client
	OctopusVersion string
	FeatureToggles map[string]bool
}

func (c *Config) SetOctopus(ctx context.Context) diag.Diagnostics {
	tflog.Debug(ctx, "SetOctopus")

	diags := diag.Diagnostics{}

	if clientError := c.GetClient(ctx); clientError != nil {
		diags.AddError("failed to load client", clientError.Error())
		return diags
	}

	if versionError := c.SetOctopusVersion(ctx); versionError != nil {
		diags.AddError("failed to load Octopus Server version", versionError.Error())
		return diags
	}

	if featuresError := c.SetFeatureToggles(ctx); featuresError != nil {
		diags.AddError("failed to load feature toggles", featuresError.Error())
		return diags
	}

	tflog.Debug(ctx, "SetOctopus completed")
	return diags
}

func (c *Config) GetClient(ctx context.Context) error {
	tflog.Debug(ctx, "GetClient")

	octopus, err := getClientForDefaultSpace(c, ctx)
	if err != nil {
		return err
	}

	if len(c.SpaceID) > 0 {
		space, err := spaces.GetByID(octopus, c.SpaceID)
		if err != nil {
			return err
		}

		octopus, err = getClientForSpace(c, ctx, space.GetID())
		if err != nil {
			return err
		}
	}

	c.Client = octopus

	createdClient := octopus != nil
	tflog.Debug(ctx, fmt.Sprintf("GetClient completed: %t", createdClient))
	return nil
}

func (c *Config) SetFeatureToggles(ctx context.Context) error {
	tflog.Debug(ctx, "SetFeatureToggles")

	response, err := configuration.Get(c.Client, &configuration.FeatureToggleConfigurationQuery{})
	if err != nil {
		return err
	}

	features := make(map[string]bool, len(response.FeatureToggles))
	for _, feature := range response.FeatureToggles {
		features[feature.Name] = feature.IsEnabled
	}

	c.FeatureToggles = features

	tflog.Debug(ctx, fmt.Sprintf("SetFeatureToggles completed with %d features", len(c.FeatureToggles)))
	return nil
}

func (c *Config) SetOctopusVersion(ctx context.Context) error {
	tflog.Debug(ctx, "SetOctopusVersion")

	root, err := client.GetServerRoot(c.Client)
	if err != nil {
		return err
	}

	c.OctopusVersion = root.Version
	tflog.Debug(ctx, fmt.Sprintf("SetOctopusVersion completed with %s", c.OctopusVersion))

	return nil
}

func getClientForDefaultSpace(c *Config, ctx context.Context) (*client.Client, error) {
	return getClientForSpace(c, ctx, "")
}

func getClientForSpace(c *Config, ctx context.Context, spaceID string) (*client.Client, error) {
	apiURL, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	credential, err := getApiCredential(c, ctx)
	if err != nil {
		return nil, err
	}

	return client.NewClientWithCredentials(nil, apiURL, credential, spaceID, "TerraformProvider")
}

func getApiCredential(c *Config, ctx context.Context) (client.ICredential, error) {
	tflog.Debug(ctx, "GetClient: Trying the following auth methods in order of priority - APIKey, AccessToken")

	if c.ApiKey != "" {
		tflog.Debug(ctx, "GetClient: Attempting to authenticate with API Key")
		credential, err := client.NewApiKey(c.ApiKey)
		if err != nil {
			return nil, err
		}

		return credential, nil
	} else {
		tflog.Debug(ctx, "GetClient: No API Key found")
	}

	if c.AccessToken != "" {
		tflog.Debug(ctx, "GetClient: Attempting to authenticate with Access Token")
		credential, err := client.NewAccessToken(c.AccessToken)
		if err != nil {
			return nil, err
		}

		return credential, nil
	} else {
		tflog.Debug(ctx, "GetClient: No Access Token found")
	}

	return nil, fmt.Errorf("either an APIKey or an AccessToken is required to connect to the Octopus Server instance")
}

func DataSourceConfiguration(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *Config {
	if req.ProviderData == nil {
		return nil
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return config
}

func ResourceConfiguration(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *Config {
	if req.ProviderData == nil {
		return nil
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return config
}

// FeatureToggleEnabled Reports whether feature toggle enabled on connected Octopus Server instance.
//
// Returns true for enabled toggle and false for disabled or non-existent feature toggle
func (c *Config) FeatureToggleEnabled(toggle string) bool {
	if enabled, ok := c.FeatureToggles[toggle]; ok {
		return enabled
	}

	return false
}

// EnsureResourceCompatibilityByFeature Reports whether resource is compatible with current instance of Octopus Server by .
// Returns diagnostics with error when resource is incompatible and empty diagnostics for compatible resources
func (c *Config) EnsureResourceCompatibilityByFeature(resourceName string, toggle string) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if c.FeatureToggleEnabled(toggle) {
		return diags
	}

	summary := fmt.Sprintf("The '%s' resource is not supported by the connected Octopus Deploy instance", resourceName)
	detail := fmt.Sprintf("This resource requires feature toggle '%s' to be enabled.", toggle)
	diags.AddError(summary, detail)

	return diags
}

// EnsureResourceCompatibilityByVersion Reports whether resource is compatible with current version of Octopus Server.
// Returns diagnostics with error when resource is incompatible and empty diagnostics for compatible resources
//
// Example: '2025.1' - first version where resource can be used
func (c *Config) EnsureResourceCompatibilityByVersion(resourceName string, version string) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if c.IsVersionSameOrGreaterThan(version) {
		return diags
	}

	summary := fmt.Sprintf("The '%s' resource is not supported by the current Octopus Deploy server version", resourceName)
	detail := fmt.Sprintf("This resource requires Octopus Deploy server version %s or later. The connected server is running version %s, which is incompatible with this resource.", version, c.OctopusVersion)
	diags.AddError(summary, detail)

	return diags
}

func (c *Config) IsVersionSameOrGreaterThan(minVersion string) bool {
	if c.OctopusVersion == "0.0.0-local" {
		return true // Always true for local instance
	}

	diff := version.Compare(fmt.Sprintf("go%s", c.OctopusVersion), fmt.Sprintf("go%s", minVersion))

	return diff == 1 || diff == 0
}
