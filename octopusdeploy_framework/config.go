package octopusdeploy_framework

import (
	"context"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Start of OctoAI patch

type HeaderRoundTripper struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

func (h *HeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}
	return h.Transport.RoundTrip(req)
}

func getHttpClient(octopusUrl string) (*http.Client, error) {
	if isDirectlyAccessibleOctopusInstance(octopusUrl) {
		return createHttpClient(octopusUrl)
	}

	return nil, nil
}

// isDirectlyAccessibleOctopusInstance determines if the host should be contacted directly
func isDirectlyAccessibleOctopusInstance(octopusUrl string) bool {
	serviceEnabled, found := os.LookupEnv("REDIRECTION_SERVICE_ENABLED")

	if !found || serviceEnabled != "true" {
		return true
	}

	parsedUrl, err := url.Parse(octopusUrl)

	// Contact the server directly if the URL is invalid
	if err != nil {
		return true
	}

	return strings.HasSuffix(parsedUrl.Hostname(), ".octopus.app") ||
		strings.HasSuffix(parsedUrl.Hostname(), ".testoctopus.com") ||
		parsedUrl.Hostname() == "localhost" ||
		parsedUrl.Hostname() == "127.0.0.1"
}

func createHttpClient(octopusUrl string) (*http.Client, error) {

	serviceApiKey, found := os.LookupEnv("REDIRECTION_SERVICE_API_KEY")

	if !found {
		return nil, errors.New("REDIRECTION_SERVICE_API_KEY is required")
	}

	parsedUrl, err := url.Parse(octopusUrl)

	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X_REDIRECTION_UPSTREAM_HOST":   parsedUrl.Hostname(),
		"X_REDIRECTION_SERVICE_API_KEY": serviceApiKey,
	}

	return &http.Client{
		Transport: &HeaderRoundTripper{
			Transport: http.DefaultTransport,
			Headers:   headers,
		},
	}, nil
}

// End of OctoAI patch

type Config struct {
	Address     string
	ApiKey      string
	AccessToken string
	SpaceID     string
	Client      *client.Client
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

	// OctoAI patch
	httpClient, err := getHttpClient(c.Address)
	if err != nil {
		return nil, err
	}

	return client.NewClientWithCredentials(httpClient, apiURL, credential, spaceID, "TerraformProvider")
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
