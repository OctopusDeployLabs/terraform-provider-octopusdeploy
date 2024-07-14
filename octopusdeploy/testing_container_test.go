package octopusdeploy

import (
	"context"
	"flag"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"testing"
	"time"
)

var createSharedContainer = flag.Bool("createSharedContainer", false, "Set to true to run integration tests in containers")

var octoContainer *test.OctopusContainer
var octoClient *client.Client
var network testcontainers.Network
var sqlServerContainer *test.MysqlContainer
var err error

func TestMain(m *testing.M) {
	flag.Parse() // Parse the flags

	if *createSharedContainer {

		testFramework := test.OctopusContainerTest{}
		octoContainer, octoClient, sqlServerContainer, network, err = testFramework.ArrangeContainer(m)
		os.Setenv("OCTOPUS_URL", octoContainer.URI)
		os.Setenv("OCTOPUS_APIKEY", test.ApiKey)
		os.Setenv("TF_ACC", "1")

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
		code := m.Run()
		os.Exit(code)
	}
}
