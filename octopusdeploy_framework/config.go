package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/url"
)

type Config struct {
	Address          string
	ApiKey           string
	AccessToken      string
	SpaceID          string
	Client           *client.Client
	TerraformVersion string
	OctopusVersion   string
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
