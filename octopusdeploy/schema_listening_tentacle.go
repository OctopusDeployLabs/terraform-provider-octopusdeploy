package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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
