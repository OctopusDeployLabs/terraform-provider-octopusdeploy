package test

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

type LifecycleTestOptions struct {
	TestOptions[lifecycles.Lifecycle]
}

func NewLifecycleTestOptions() *LifecycleTestOptions {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	options := &LifecycleTestOptions{
		TestOptions: *NewTestOptions[lifecycles.Lifecycle]("lifecycle"),
	}
	options.Resource = lifecycles.NewLifecycle(name)
	options.Resource.Description = description

	return options
}
func LifecycleConfiguration(options *LifecycleTestOptions) string {
	configuration := fmt.Sprintf(`resource "%s" "%s" {`, options.ResourceName, options.LocalName)
	configuration += "\n"
	configuration += fmt.Sprintf(`name = "%s"`, options.Resource.Name)
	configuration += "\n"

	if len(options.Resource.Description) > 0 {
		configuration += fmt.Sprintf(`description = "%s"`, options.Resource.Description)
		configuration += "\n"
	}

	// Add release retention policy
	configuration += `
    release_retention_policy {
        unit             = "Days"
        quantity_to_keep = 30
    }
    `

	// Add tentacle retention policy
	configuration += `
    tentacle_retention_policy {
        unit             = "Days"
        quantity_to_keep = 30
    }
    `

	configuration += "}"
	return configuration
}
