package octopusdeploy_framework

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/testcontainers/testcontainers-go"
)

var createSharedContainer = flag.Bool("createSharedContainer", false, "Set to true to run integration tests in containers")

var octoContainer *test.OctopusContainer
var octoClient *client.Client
var network testcontainers.Network
var sqlServerContainer *test.MysqlContainer
var err error
var customTestEnvironment = map[string]string{}

func TestMain(m *testing.M) {
	flag.Parse() // Parse the flags
	os.Setenv("TF_ACC", "1")
	if *createSharedContainer {

		testFramework := test.OctopusContainerTest{
			CustomEnvironment: map[string]string{
				"OCTOPUS__FeatureToggles__DeploymentFreezeByTenantFeatureToggle": "true",
			},
		}

		octoContainer, octoClient, sqlServerContainer, network, err = testFramework.ArrangeContainer()
		if err != nil {
			log.Printf("Failed to arrange containers: (%s)", err.Error())
		}

		os.Setenv("OCTOPUS_URL", octoContainer.URI)
		os.Setenv("OCTOPUS_APIKEY", test.ApiKey)
		code := m.Run()
		ctx := context.Background()

		// Waiting for the container logs to clear.
		time.Sleep(10000 * time.Millisecond)
		err := testFramework.CleanUp(ctx, octoContainer, sqlServerContainer, network)

		if err != nil {
			log.Printf("Failed to clean up containers: (%s)", err.Error())
			panic(m)
		}

		log.Printf("Exit code: (%d)", code)
		os.Exit(code)
	} else {
		if os.Getenv("TF_ACC_LOCAL") != "" {
			var url = os.Getenv("OCTOPUS_URL")
			var apikey = os.Getenv("OCTOPUS_APIKEY")
			octoClient, err = octoclient.CreateClient(url, "", apikey)
			if err != nil {
				log.Printf("Failed to create client: (%s)", err.Error())
				panic(m)
			}
			octoContainer = &test.OctopusContainer{
				Container: nil,
				URI:       url,
			}
		}
		code := m.Run()
		os.Exit(code)
	}
}
