package model

import (
	"log"
	"net/http"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
)

// Config holds Address and the APIKey of the Octopus Deploy server
type Config struct {
	Address string
	APIKey  string
	Space   string
}

// Client returns a new Octopus Deploy client
func (c *Config) Client() (*client.Client, error) {
	apiClient, err := client.NewClient(&(http.Client{}), c.Address, c.APIKey)

	if c.Space == "" {

		log.Printf("[INFO] Octopus Deploy Client configured against default space")

		return apiClient, nil
	}

	log.Printf("[INFO] Octopus Deploy Client will be scoped to %s space", c.Space)

	space, err := apiClient.Spaces.GetByName(c.Space)

	if err != nil {
		return nil, err
	}

	scopedClient, err := client.ForSpace(&(http.Client{}), c.Address, c.APIKey, space)

	log.Printf("[INFO] Octopus Deploy Client configured against %s space", c.Space)

	return scopedClient, nil
}
