package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandSSHConnection(d *schema.ResourceData) *octopusdeploy.SSHEndpoint {
	host := d.Get("host").(string)
	port := d.Get("port").(int)
	fingerprint := d.Get("fingerprint").(string)

	endpoint := octopusdeploy.NewSSHEndpoint(host, port, fingerprint)
	endpoint.ID = d.Id()

	if v, ok := d.GetOk("account_id"); ok {
		endpoint.AccountID = v.(string)
	}

	if v, ok := d.GetOk("dot_net_core_platform"); ok {
		endpoint.DotNetCorePlatform = v.(string)
	}

	if v, ok := d.GetOk("proxy_id"); ok {
		endpoint.ProxyID = v.(string)
	}

	if v, ok := d.GetOk("uri"); ok {
		endpointURI, _ := url.Parse(v.(string))
		endpoint.URI = endpointURI
	}

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
