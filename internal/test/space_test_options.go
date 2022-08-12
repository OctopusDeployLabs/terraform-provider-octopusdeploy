package test

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

type SpaceTestOptions struct {
	Description string
	TestOptions[spaces.Space]
}

func NewSpaceTestOptions() *SpaceTestOptions {
	return &SpaceTestOptions{
		Description: acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		TestOptions: *NewTestOptions[spaces.Space]("space"),
	}
}
