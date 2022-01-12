package octopusdeploy

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Config holds Address and the APIKey of the Octopus Deploy server
type Config struct {
	Address   string
	APIKey    string
	SpaceID   string
	SpaceName string
}

// Client returns a new Octopus Deploy client
func (c *Config) Client() (*octopusdeploy.Client, diag.Diagnostics) {
	apiURL, err := url.Parse(c.Address)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client, err := octopusdeploy.NewClient(nil, apiURL, c.APIKey, "")
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if len(c.SpaceID) > 0 {
		log.Printf("[DEBUG] Config is locating space using ID '%s'", c.SpaceID)
		space, err := client.Spaces.GetByID(c.SpaceID)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		client, err = octopusdeploy.NewClient(nil, apiURL, c.APIKey, space.GetID())
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	if len(c.SpaceName) > 0 {
		log.Printf("[DEBUG] Config is locating space using name '%s'", c.SpaceName)
		spaces, err := client.Spaces.Get(octopusdeploy.SpacesQuery{
			PartialName: c.SpaceName,
		})
		if err != nil {
			return nil, diag.FromErr(err)
		}

		if spaces.TotalResults == 0 {
			return nil, diag.Errorf("Unable to locate space with name '%s', found no spaces", c.SpaceName)
		}
		if spaces.TotalResults > 1 {
			return nil, diag.Errorf("Unable to uniquely locate space with name '%s', found spaces %s", c.SpaceName, strings.Join(getQuotedSpaceNames(spaces.Items), ", "))
		}

		spaceID := spaces.Items[0].GetID()
		log.Printf("[DEBUG] Config located space using name '%s', which has ID '%s'", c.SpaceName, spaceID)

		client, err = octopusdeploy.NewClient(nil, apiURL, c.APIKey, spaceID)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return client, nil
}

func getQuotedSpaceNames(spaces []octopusdeploy.Space) []string {
	var newSlice []string
	for _, space := range spaces {
		newSlice = append(newSlice, fmt.Sprintf("'%s'", space.Name))
	}
	return newSlice
}
