package octopusdeploy

import (
	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
)

func expandPollingTentacle(flattenedMap map[string]interface{}) *machines.PollingTentacleEndpoint {
	if octopusUrlString, ok := flattenedMap["octopus_url"]; ok {
		octopusURL, _ := url.Parse(octopusUrlString.(string))

		thumbprint := flattenedMap["thumbprint"].(string)
		endpoint := machines.NewPollingTentacleEndpoint(octopusURL, thumbprint)

		endpoint.CertificateSignatureAlgorithm = flattenedMap["certificate_signature_algorithm"].(string)
		endpoint.ID = flattenedMap["id"].(string)
		endpoint.TentacleVersionDetails = expandTentacleVersionDetails(flattenedMap["tentacle_version_details"])

		return endpoint
	}

	return nil
}
