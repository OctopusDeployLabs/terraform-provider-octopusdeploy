package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
)

func expandListeningTentacle(flattenedMap map[string]interface{}) *machines.ListeningTentacleEndpoint {
	tentacleURL, _ := url.Parse(flattenedMap["tentacle_url"].(string))
	thumbprint := flattenedMap["thumbprint"].(string)

	endpoint := machines.NewListeningTentacleEndpoint(tentacleURL, thumbprint)
	endpoint.CertificateSignatureAlgorithm = flattenedMap["certificate_signature_algorithm"].(string)
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.ProxyID = flattenedMap["proxy_id"].(string)
	endpoint.TentacleVersionDetails = expandTentacleVersionDetails(flattenedMap["tentacle_version_details"])

	return endpoint
}
