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
}

// Client returns a new Octopus Deploy client
func (c *Config) Client() *octopusdeploy.Client {
	httpClient := http.Client{}
	client := octopusdeploy.NewClient(&httpClient, c.Address, c.APIKey)
	log.Printf("[INFO] Octopus Deploy Client configured ")

	return client
}
