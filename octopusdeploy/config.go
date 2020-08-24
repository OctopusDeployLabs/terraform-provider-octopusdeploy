package octopusdeploy

import (
	"log"
	"net/http"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/OctopusDeploy/Go-Swagger-Client-Gen-Terraform"
)

// Config holds Address and the APIKey of the Octopus Deploy server
type Config struct {
	Address string
	APIKey  string
	Space   string
}

// Client returns a new Octopus Deploy client
func (c *Config) Client() (*Client, error) {
	
	transportConfig := TransportConfig{}
	client := NewHTTPClientWithConfig(nil, transportConfig.WithHost(c.Address))
	authInfo := httptransport.APIKeyAuth("X-Octopus-ApiKey", "header", c.APIKey)

	return client

	if c.Space == "" {
		log.Printf("[INFO] Octopus Deploy Client configured against default space")
		return client, nil
	}

	log.Printf("[INFO] Octopus Deploy Client will be scoped to %s space", c.Space)

	newGetSpaceByIDParams := spaces.NewGetSpaceByIDParams()
	newGetSpaceByIDParams.ID = c.Space
	getSpaceByIDOK, err := client.Spaces.GetSpaceByID(newGetSpaceByIDParams, authInfo)

	// space, err := client.Space.GetByName(c.Space)

	if err != nil {
		return nil, err
	}

	// default base path --> [address]/api
	// new base bath     --> [address]/api/[space-id]
	
	// transportConfig.BasePath = "/" + space
	// scopedClient := octopusdeploy.NewHTTPClientWithConfig(nil, &transportConfig)

	// scopedClient := octopusdeploy.ForSpace(&(http.Client{}), c.Address, c.APIKey, space)

	log.Printf("[INFO] Octopus Deploy Client configured against %s space", c.Space)

	return scopedClient, nil
}
