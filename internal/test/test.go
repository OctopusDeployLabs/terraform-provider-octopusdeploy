package test

import "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

type ProjectGroupTestOptions struct {
	Project ProjectTestOptions
	TestOptions
}

func NewProjectGroupTestOptions() *ProjectGroupTestOptions {
	return &ProjectGroupTestOptions{
		Project:     *NewProjectTestOptions(),
		TestOptions: *NewTestOptions(),
	}
}

type ProjectCreateTestOptions struct {
	Lifecycle    *TestOptions
	ProjectGroup *ProjectGroupTestOptions
}

type ProjectTestOptions struct {
	Channel     *TestOptions
	Description string
	TestOptions
}

func NewProjectTestOptions() *ProjectTestOptions {
	return &ProjectTestOptions{
		Channel:     NewTestOptions(),
		Description: acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		TestOptions: *NewTestOptions(),
	}
}

type TestOptions struct {
	LocalName string
	Name      string
}

func NewTestOptions() *TestOptions {
	return &TestOptions{
		LocalName: acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		Name:      acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	}
}

func (d *DeploymentProcessTestOptions) ProjectCreateTestOptions() *ProjectCreateTestOptions {
	return &ProjectCreateTestOptions{
		Lifecycle:    d.Lifecycle,
		ProjectGroup: d.ProjectGroup,
	}
}

type DeploymentProcessTestOptions struct {
	ActionName   string
	ActionType   string
	Lifecycle    *TestOptions
	PackageID    string
	PackageName  string
	ProjectGroup *ProjectGroupTestOptions
	StepName     string
	Space        *TestOptions
	TestOptions
}

func NewDeploymentProcessTestOptions() *DeploymentProcessTestOptions {
	return &DeploymentProcessTestOptions{
		ActionName:   acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		ActionType:   "Octopus.DeployTentaclePackage",
		Lifecycle:    NewTestOptions(),
		PackageID:    "Octopus.Cli",
		PackageName:  acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		ProjectGroup: NewProjectGroupTestOptions(),
		StepName:     acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		Space:        NewTestOptions(),
		TestOptions:  *NewTestOptions(),
	}
}
