package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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
