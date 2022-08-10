package test

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
)

type ListeningTentacleDeploymentTargetTestOptions struct {
	TestOptions[machines.ListeningTentacleDeploymentTarget]
}

func NewListeningTentacleDeploymentTargetTestOptions() *ListeningTentacleDeploymentTargetTestOptions {
	return &ListeningTentacleDeploymentTargetTestOptions{
		TestOptions: *NewTestOptions[machines.ListeningTentacleDeploymentTarget]("listening_tentacle_deployment_target"),
	}
}

func ListeningTentacleDeploymentTargetConfiguration(options *ListeningTentacleDeploymentTargetTestOptions) string {
	configuration := fmt.Sprintf(`"resource "%s" "%s" {`, options.ResourceName, options.LocalName)

	if len(options.Resource.Name) > 0 {
		configuration += fmt.Sprintf(`name = "%s"`, options.Resource.Name)
	}

	configuration += "}"
	return configuration
}
