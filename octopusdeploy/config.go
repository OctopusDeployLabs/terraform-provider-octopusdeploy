package octopusdeploy

import (
	"fmt"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Config holds Address and the APIKey of the Octopus Deploy server
type Config struct {
	Address     string
	APIKey      string
	AccessToken string
	SpaceID     string
}

// Client returns a new Octopus Deploy client
func (c *Config) Client() (*client.Client, diag.Diagnostics) {
	octopus, err := getClientForDefaultSpace(c)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if len(c.SpaceID) > 0 {
		space, err := spaces.GetByID(octopus, c.SpaceID)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		octopus, err = getClientForSpace(c, space.GetID())
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return octopus, nil
}

func getClientForDefaultSpace(c *Config) (*client.Client, error) {
	return getClientForSpace(c, "")
}

func getClientForSpace(c *Config, spaceID string) (*client.Client, error) {
	apiURL, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	credential, err := getApiCredential(c)
	if err != nil {
		return nil, err
	}

	return client.NewClientWithCredentials(nil, apiURL, credential, spaceID, "TerraformProvider")
}

func getApiCredential(c *Config) (client.ICredential, error) {
	if c.APIKey != "" {
		credential, err := client.NewApiKey(c.APIKey)
		if err != nil {
			return nil, err
		}

		return credential, nil
	}

	if c.AccessToken != "" {
		credential, err := client.NewAccessToken(c.AccessToken)
		if err != nil {
			return nil, err
		}

		return credential, nil
	}

	return nil, fmt.Errorf("either an APIKey or an AccessToken is required to connect to the Octopus Server instance")
}
