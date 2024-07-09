package octopusdeployv6

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"testing"
)

var octoContainer *test.OctopusContainer
var octoClient *client.Client
var network testcontainers.Network
var sqlServerContainer *test.MysqlContainer
var err error

func TestMain(m *testing.M) {
	testFramework := test.OctopusContainerTest{}
	octoContainer, octoClient, sqlServerContainer, network, err = testFramework.ArrangeContainer(m)
	os.Setenv("OCTOPUS_URL", octoContainer.URI)
	os.Setenv("OCTOPUS_APIKEY", test.ApiKey)

	code := m.Run()
	ctx := context.Background()
	err := testFramework.CleanUp(ctx, octoContainer, sqlServerContainer, network)

	if err != nil {
		log.Printf("Failed to clean up containers: (%s)", err.Error())
		panic(m)
	}

	log.Printf("Exit code: (%d)", code)
	os.Exit(code)
}
