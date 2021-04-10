package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandSSHConnection(flattenedMap map[string]interface{}) *octopusdeploy.SSHEndpoint {
	host := flattenedMap["host"].(string)
	port := flattenedMap["port"].(int)
	fingerprint := flattenedMap["fingerprint"].(string)

	endpoint := octopusdeploy.NewSSHEndpoint(host, port, fingerprint)
	endpoint.AccountID = flattenedMap["account_id"].(string)
	endpoint.DotNetCorePlatform = flattenedMap["dot_net_core_platform"].(string)
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.ProxyID = flattenedMap["proxy_id"].(string)

	endpointURI, _ := url.Parse(flattenedMap["uri"].(string))
	endpoint.URI = endpointURI

	return endpoint
}
