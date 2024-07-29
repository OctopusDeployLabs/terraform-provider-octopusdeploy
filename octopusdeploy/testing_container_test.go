package octopusdeploy

import (
	"context"
	"flag"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"testing"
	"time"
)

var createSharedContainer = flag.Bool("createSharedContainer", true, "Set to true to run integration tests in containers")

// Note these could and should be removed and replaced by references to the IntegrationTestSuite.
// We haven't done this due to refactoring time.
var octoContainer *test.OctopusContainer
var octoClient *client.Client

type (
	IntegrationTestSuite struct {
		suite.Suite
		octoContainer      *test.OctopusContainer
		octoClient         *client.Client
		network            testcontainers.Network
		sqlServerContainer *test.MysqlContainer
		err                error
	}
)

func TestMySuite(t *testing.T) {
	suite.Run(t, &IntegrationTestSuite{})
}

// SetupSuite implements testify's SetupAllSuite interface and runs at the beginning of the test suite, when suite.Run() is called.
func (suite *IntegrationTestSuite) SetupSuite() {
	flag.Parse() // Parse the flags
	os.Setenv("TF_ACC", "1")
	if *createSharedContainer {
		testFramework := test.OctopusContainerTest{}
		suite.octoContainer, suite.octoClient, suite.sqlServerContainer, suite.network, suite.err = testFramework.ArrangeContainer()
		octoContainer = suite.octoContainer
		octoClient = suite.octoClient

		os.Setenv("OCTOPUS_URL", suite.octoContainer.URI)
		os.Setenv("OCTOPUS_APIKEY", test.ApiKey)
	} else {
		if os.Getenv("TF_ACC_LOCAL") != "" {
			var url = os.Getenv("OCTOPUS_URL")
			var apikey = os.Getenv("OCTOPUS_APIKEY")
			suite.octoClient, suite.err = octoclient.CreateClient(url, "", apikey)
			if suite.err != nil {
				log.Printf("Failed to create client: (%s)", suite.err.Error())
			}
		}
	}
}

// TearDownSuite implements testify's TearDownAllSuite interface and runs at the end of the suite.
func (suite *IntegrationTestSuite) TearDownSuite() {
	// Waiting for the container logs to clear.
	time.Sleep(10000 * time.Millisecond)
	testFramework := test.OctopusContainerTest{}
	ctx := context.Background()
	err := testFramework.CleanUp(ctx, suite.octoContainer, suite.sqlServerContainer, suite.network)

	if err != nil {
		log.Printf("Failed to clean up containers: (%s)", err.Error())
	}
}

// TearDownTest implements testify's TearDownTestSuite interface and runs after each suite of tests.
func (suite *IntegrationTestSuite) TearDownTest() {

}
