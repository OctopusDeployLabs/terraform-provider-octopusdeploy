package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandListeningTentacle(flattenedMap map[string]interface{}) *octopusdeploy.ListeningTentacleEndpoint {
	tentacleURL, _ := url.Parse(flattenedMap["tentacle_url"].(string))
	thumbprint := flattenedMap["thumbprint"].(string)

	endpoint := octopusdeploy.NewListeningTentacleEndpoint(tentacleURL, thumbprint)
	endpoint.CertificateSignatureAlgorithm = flattenedMap["certificate_signature_algorithm"].(string)
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.ProxyID = flattenedMap["proxy_id"].(string)
	endpoint.TentacleVersionDetails = expandTentacleVersionDetails(flattenedMap["tentacle_version_details"])

	return endpoint
}

func flattenListeningTentacle(endpoint *octopusdeploy.ListeningTentacleEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	rawEndpoint := map[string]interface{}{
		"certificate_signature_algorithm": endpoint.CertificateSignatureAlgorithm,
		"id":                              endpoint.GetID(),
		"proxy_id":                        endpoint.ProxyID,
		"tentacle_version_details":        flattenTentacleVersionDetails(endpoint.TentacleVersionDetails),
		"thumbprint":                      endpoint.Thumbprint,
	}

	if endpoint.URI != nil {
		rawEndpoint["tentacle_url"] = endpoint.URI.String()
	}

	return []interface{}{rawEndpoint}
}

func getListeningTentacleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"certificate_signature_algorithm": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"proxy_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"tentacle_version_details": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getTentacleVersionDetailsSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tentacle_url": {
			Required: true,
			Type:     schema.TypeString,
		},
		"thumbprint": {
			Required: true,
			Type:     schema.TypeString,
		},
	}
}
