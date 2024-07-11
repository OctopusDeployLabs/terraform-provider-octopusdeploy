package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/url"
)

type Config struct {
	Address string
	ApiKey  string
	SpaceID string
	Client  *client.Client
}

func (c *Config) GetClient(ctx context.Context) error {
	tflog.Debug(ctx, "GetClient")
	apiURL, err := url.Parse(c.Address)
	if err != nil {
		return err
	}

	octopus, err := client.NewClient(nil, apiURL, c.ApiKey, "")
	if err != nil {
		return err
	}

	if len(c.SpaceID) > 0 {
		space, err := spaces.GetByID(octopus, c.SpaceID)
		if err != nil {
			return err
		}

		octopus, err = client.NewClient(nil, apiURL, c.ApiKey, space.GetID())
		if err != nil {
			return err
		}
	}

	c.Client = octopus

	createdClient := octopus != nil
	tflog.Debug(ctx, fmt.Sprintf("GetClient completed: %t", createdClient))
	return nil
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
