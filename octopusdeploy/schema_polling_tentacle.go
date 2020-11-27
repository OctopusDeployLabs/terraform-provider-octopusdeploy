package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandPollingTentacle(flattenedMap map[string]interface{}) *octopusdeploy.PollingTentacleEndpoint {
	octopusURL, _ := url.Parse(flattenedMap["octopus_url"].(string))
	thumbprint := flattenedMap["thumbprint"].(string)

	endpoint := octopusdeploy.NewPollingTentacleEndpoint(octopusURL, thumbprint)
	endpoint.CertificateSignatureAlgorithm = flattenedMap["certificate_signature_algorithm"].(string)
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.TentacleVersionDetails = expandTentacleVersionDetails(flattenedMap["tentacle_version_details"])

	return endpoint
}

func flattenPollingTentacle(endpoint *octopusdeploy.PollingTentacleEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	rawEndpoint := map[string]interface{}{
		"certificate_signature_algorithm": endpoint.CertificateSignatureAlgorithm,
		"id":                              endpoint.GetID(),
		"tentacle_version_details":        flattenTentacleVersionDetails(endpoint.TentacleVersionDetails),
		"thumbprint":                      endpoint.Thumbprint,
	}

	if endpoint.URI != nil {
		rawEndpoint["octopus_url"] = endpoint.URI.String()
	}

	return []interface{}{rawEndpoint}
}

func getPollingTentacleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"certificate_signature_algorithm": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"octopus_url": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"tentacle_version_details": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getTentacleVersionDetailsSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"thumbprint": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
