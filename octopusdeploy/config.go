package octopusdeploy

import (
	"net/http"
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
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

	// This is intentional on the feature branch - Todo: remove when merging to main branch
	proxyStr := "http://172.21.224.1:8866"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return nil, nil
	}

	tr := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	httpClient := http.Client{Transport: tr}

	octopus, err := client.NewClient(&httpClient, apiURL, c.APIKey, "")
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if len(c.SpaceID) > 0 {
		space, err := octopus.Spaces.GetByID(c.SpaceID)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		octopus, err = client.NewClient(&httpClient, apiURL, c.APIKey, space.GetID())
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return octopus, nil
}
