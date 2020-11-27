package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func flattenSSHConnection(endpoint *octopusdeploy.SSHEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	rawEndpoint := map[string]interface{}{
		"account_id":            endpoint.AccountID,
		"dot_net_core_platform": endpoint.DotNetCorePlatform,
		"fingerprint":           endpoint.Fingerprint,
		"host":                  endpoint.Host,
		"id":                    endpoint.GetID(),
		"proxy_id":              endpoint.ProxyID,
		"port":                  endpoint.Port,
	}

	if endpoint.URI != nil {
		rawEndpoint["uri"] = endpoint.URI.String()
	}

	return []interface{}{rawEndpoint}
}

func getSSHConnectionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"dot_net_core_platform": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"fingerprint": {
			Required: true,
			Type:     schema.TypeString,
		},
		"host": {
			Required: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"port": {
			Required: true,
			Type:     schema.TypeInt,
		},
		"proxy_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"uri": {
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}
