package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
)

func expandPollingTentacle(flattenedMap map[string]interface{}) *machines.PollingTentacleEndpoint {
	octopusURL, _ := url.Parse(flattenedMap["octopus_url"].(string))
	thumbprint := flattenedMap["thumbprint"].(string)

	endpoint := machines.NewPollingTentacleEndpoint(octopusURL, thumbprint)
	endpoint.CertificateSignatureAlgorithm = flattenedMap["certificate_signature_algorithm"].(string)
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.TentacleVersionDetails = expandTentacleVersionDetails(flattenedMap["tentacle_version_details"])

	return endpoint
}
