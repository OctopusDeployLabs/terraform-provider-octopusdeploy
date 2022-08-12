package test

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

type DeploymentProcessTestOptions struct {
	ActionName   string
	ActionType   string
	Lifecycle    *LifecycleTestOptions
	PackageID    string
	PackageName  string
	ProjectGroup *ProjectGroupTestOptions
	Project      *ProjectTestOptions
	StepName     string
	Space        *SpaceTestOptions
	TestOptions[deployments.DeploymentProcess]
}

func NewDeploymentProcessTestOptions() *DeploymentProcessTestOptions {
	lifecycleTestOptions := NewLifecycleTestOptions()
	projectGroupTestOptions := NewProjectGroupTestOptions()
	projectTestOptions := NewProjectTestOptions(lifecycleTestOptions, projectGroupTestOptions)

	return &DeploymentProcessTestOptions{
		ActionName:   acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		ActionType:   "Octopus.DeployTentaclePackage",
		Lifecycle:    lifecycleTestOptions,
		PackageID:    "Octopus.Cli",
		PackageName:  acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		ProjectGroup: projectGroupTestOptions,
		Project:      projectTestOptions,
		StepName:     acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		Space:        NewSpaceTestOptions(),
		TestOptions:  *NewTestOptions[deployments.DeploymentProcess]("deployment_process"),
	}
}

func (d *DeploymentProcessTestOptions) ProjectCreateTestOptions() *ProjectCreateTestOptions {
	return &ProjectCreateTestOptions{
		Lifecycle:    d.Lifecycle,
		ProjectGroup: d.ProjectGroup,
	}
}
