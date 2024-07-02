package pluginprovider

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Client returns a new Octopus Deploy client
func (c *Config) Client() (*client.Client, diag.Diagnostics) {
	apiURL, err := url.Parse(c.Address)
	var diags diag.Diagnostics
	if err != nil {
		diags.AddError("Invalid Address", "Could not parse the Address URL: "+err.Error())
		return nil, diags
	}

	octopus, err := client.NewClient(nil, apiURL, c.APIKey, "")
	if err != nil {
		diags.AddError("Client Error", "Could not create the Octopus Deploy client: "+err.Error())
		return nil, diags
	}

	if len(c.SpaceID) > 0 {
		space, err := spaces.GetByID(octopus, c.SpaceID)
		if err != nil {
			diags.AddError("Space ID Error", "Could not get space by ID: "+err.Error())
			return nil, diags
		}

		octopus, err = client.NewClient(nil, apiURL, c.APIKey, space.GetID())
		if err != nil {
			diags.AddError("Client Error", "Could not create the Octopus Deploy client with space ID: "+err.Error())
			return nil, diags
		}
	}

	return octopus, nil
}
