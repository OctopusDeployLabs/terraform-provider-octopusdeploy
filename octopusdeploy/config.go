package octopusdeploy

import (
	"log"

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
	apiClient, err := client.NewClient(nil, c.Address, c.APIKey, c.Space)

	if err != nil {
		log.Println(err)
	}

	log.Printf("[INFO] Octopus Deploy Client Ready")

	return apiClient, err
}
