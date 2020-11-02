package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Config holds Address and the APIKey of the Octopus Deploy server
type Config struct {
	Address string
	APIKey  string
	Space   string
}

// Client returns a new Octopus Deploy client
func (c *Config) Client() (*octopusdeploy.Client, diag.Diagnostics) {
	apiURL, err := url.Parse(c.Address)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client, err := octopusdeploy.NewClient(nil, apiURL, c.APIKey, c.Space)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
