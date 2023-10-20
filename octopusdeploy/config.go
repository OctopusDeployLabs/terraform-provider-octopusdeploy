package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Config holds Address and the APIKey of the Octopus Deploy server
type Config struct {
	Address string
	APIKey  string
	SpaceID string
}

// Client returns a new Octopus Deploy client
func (c *Config) Client() (*client.Client, diag.Diagnostics) {
	apiURL, err := url.Parse(c.Address)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	octopus, err := client.NewClient(nil, apiURL, c.APIKey, "")
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if len(c.SpaceID) > 0 {
		space, err := spaces.GetByID(octopus, c.SpaceID)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		octopus, err = client.NewClient(nil, apiURL, c.APIKey, space.GetID())
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return octopus, nil
}
