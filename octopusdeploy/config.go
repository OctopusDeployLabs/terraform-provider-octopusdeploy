package octopusdeploy

import (
	"log"
	"net/http"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

// Config holds Address and the APIKey of the Octopus Deploy server
type Config struct {
	Address string
	APIKey  string
	Space   string
}

// Client returns a new Octopus Deploy client
func (c *Config) Client() (*octopusdeploy.Client, error) {
	client := octopusdeploy.NewClient(&(http.Client{}), c.Address, c.APIKey)

	if c.Space == "" {

		log.Printf("[INFO] Octopus Deploy Client configured against default space")

		return client, nil
	}

	log.Printf("[INFO] Octopus Deploy Client will be scoped to %s space", c.Space)

	space, err := client.Space.GetByName(c.Space)

	if err != nil {
		return nil, err
	}

	scopedClient := octopusdeploy.ForSpace(&(http.Client{}), c.Address, c.APIKey, space)

	log.Printf("[INFO] Octopus Deploy Client configured against %s space", c.Space)

	return scopedClient, nil
}
