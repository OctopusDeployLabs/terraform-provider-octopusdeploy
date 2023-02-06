package main

/*
	To test the Octopus Terraform provider locally, save the following into a failed called ~/.terraformrc, replacing
	/var/home/yourname/Code/terraform-provider-octopusdeploy with the directory containing your clone
	of the git repo:

		provider_installation {
		  dev_overrides {
			"octopusdeploylabs/octopusdeploy" = "/var/home/yourname/Code/terraform-provider-octopusdeploy"
		  }

		  direct {}
		}

	Then build the provider executable with the command:

		go build -o terraform-provider-octopusdeploy main.go

	Terraform will then use the local executable rather than download the provider from the registry.
*/

import (
	"context"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/certificates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workerpools"
	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"k8s.io/utils/strings/slices"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const ApiKey = "API-ABCDEFGHIJKLMNOPQURTUVWXYZ12345"

type octopusContainer struct {
	testcontainers.Container
	URI string
}

type mysqlContainer struct {
	testcontainers.Container
	port string
	ip   string
}

type TestLogConsumer struct {
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	fmt.Println(string(l.Content))
}

func enableContainerLogging(container testcontainers.Container, ctx context.Context) error {
	// Display the container logs
	err := container.StartLogProducer(ctx)
	if err != nil {
		return err
	}
	g := TestLogConsumer{}
	container.FollowOutput(&g)
	return nil
}

// getReaperSkipped will return true if running in a podman environment
func getReaperSkipped() bool {
	if strings.Contains(os.Getenv("DOCKER_HOST"), "podman") {
		return true
	}

	return false
}

// getProvider returns the test containers provider
func getProvider() testcontainers.ProviderType {
	if strings.Contains(os.Getenv("DOCKER_HOST"), "podman") {
		return testcontainers.ProviderPodman
	}

	return testcontainers.ProviderDocker
}

// setupNetwork creates an internal network for Octopus and MS SQL
func setupNetwork(ctx context.Context) (testcontainers.Network, error) {
	return testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           "octopusterraformtests",
			CheckDuplicate: true,
			SkipReaper:     getReaperSkipped(),
		},
		ProviderType: getProvider(),
	})
}

// setupDatabase creates a MSSQL container
func setupDatabase(ctx context.Context) (*mysqlContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mcr.microsoft.com/mssql/server",
		ExposedPorts: []string{"1433/tcp"},
		Env: map[string]string{
			"ACCEPT_EULA": "Y",
			"SA_PASSWORD": "Password01!",
		},
		WaitingFor: wait.ForExec([]string{"/opt/mssql-tools/bin/sqlcmd", "-U", "sa", "-P", "Password01!", "-Q", "select 1"}).WithExitCodeMatcher(
			func(exitCode int) bool {
				return exitCode == 0
			}),
		SkipReaper: getReaperSkipped(),
		Networks: []string{
			"octopusterraformtests",
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "1433")
	if err != nil {
		return nil, err
	}

	return &mysqlContainer{
		Container: container,
		ip:        ip,
		port:      mappedPort.Port(),
	}, nil
}

// setupOctopus creates an Octopus container
func setupOctopus(ctx context.Context, connString string) (*octopusContainer, error) {
	if os.Getenv("LICENSE") == "" {
		return nil, errors.New("the LICENSE environment variable must be set to a base 64 encoded Octopus license key")
	}

	req := testcontainers.ContainerRequest{
		// Be aware that later versions of Octopus killed Github Actions.
		// I think maybe they used more memory? 2022.2 works fine though.
		Image:        "octopusdeploy/octopusdeploy:2022.2",
		ExposedPorts: []string{"8080/tcp"},
		Env: map[string]string{
			"ACCEPT_EULA":                   "Y",
			"DB_CONNECTION_STRING":          connString,
			"ADMIN_API_KEY":                 ApiKey,
			"DISABLE_DIND":                  "Y",
			"ADMIN_USERNAME":                "admin",
			"ADMIN_PASSWORD":                "Password01!",
			"OCTOPUS_SERVER_BASE64_LICENSE": os.Getenv("LICENSE"),
		},
		Privileged: false,
		WaitingFor: wait.ForLog("Listening for HTTP requests on").WithStartupTimeout(30 * time.Minute),
		SkipReaper: getReaperSkipped(),
		Networks: []string{
			"octopusterraformtests",
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Display the container logs
	enableContainerLogging(container, ctx)

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "8080")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &octopusContainer{Container: container, URI: uri}, nil
}

// arrangeTest is wrapper that initialises Octopus, runs a test, and cleans up the containers
func arrangeTest(t *testing.T, testFunc func(t *testing.T, container *octopusContainer) error) {
	err := retry.Do(
		func() error {

			if testing.Short() {
				t.Skip("skipping integration test")
			}

			ctx := context.Background()

			network, err := setupNetwork(ctx)
			if err != nil {
				return err
			}

			sqlServer, err := setupDatabase(ctx)
			if err != nil {
				return err
			}

			sqlIp, err := sqlServer.Container.ContainerIP(ctx)
			if err != nil {
				return err
			}

			t.Log("SQL Server IP: " + sqlIp)

			octopusContainer, err := setupOctopus(ctx, "Server="+sqlIp+",1433;Database=OctopusDeploy;User=sa;Password=Password01!")
			if err != nil {
				return err
			}

			// Clean up the container after the test is complete
			defer func() {
				// This fixes the "can not get logs from container which is dead or marked for removal" error
				// See https://github.com/testcontainers/testcontainers-go/issues/606
				octopusContainer.StopLogProducer()

				octoTerminateErr := octopusContainer.Terminate(ctx)
				sqlTerminateErr := sqlServer.Terminate(ctx)

				networkErr := network.Remove(ctx)

				if octoTerminateErr != nil || sqlTerminateErr != nil || networkErr != nil {
					t.Fatalf("failed to terminate container: %v %v", octoTerminateErr, sqlTerminateErr)
				}
			}()

			// give the server 5 minutes to start up
			success := false
			for start := time.Now(); ; {
				if time.Since(start) > 5*time.Minute {
					break
				}

				resp, err := http.Get(octopusContainer.URI + "/api")
				if err == nil && resp.StatusCode == http.StatusOK {
					success = true
					t.Log("Successfully contacted the Octopus API")
					break
				}

				time.Sleep(10 * time.Second)
			}

			if !success {
				t.Fatalf("Failed to access the Octopus API")
			}

			return testFunc(t, octopusContainer)
		},
		retry.Attempts(3),
	)

	if err != nil {
		t.Fatalf(err.Error())
	}
}

// initialiseOctopus uses Terraform to populate the test Octopus instance, making sure to clean up
// any files generated during previous Terraform executions to avoid conflicts and locking issues.
func initialiseOctopus(t *testing.T, container *octopusContainer, terraformDir string, spaceName string, initialiseVars []string, populateVars []string) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	t.Log("Working dir: " + path)

	// This test creates a new space and then populates the space.
	terraformProjectDirs := []string{}
	terraformProjectDirs = append(terraformProjectDirs, filepath.Join("test", "terraform", "1-singlespace"))
	terraformProjectDirs = append(terraformProjectDirs, filepath.Join(terraformDir))

	// First loop initialises the new space, second populates the space
	spaceId := "Spaces-1"
	for i, terraformProjectDir := range terraformProjectDirs {

		os.Remove(filepath.Join(terraformProjectDir, ".terraform.lock.hcl"))
		os.Remove(filepath.Join(terraformProjectDir, "terraform.tfstate"))

		args := []string{"init", "-no-color"}
		cmnd := exec.Command("terraform", args...)
		cmnd.Dir = terraformProjectDir
		out, err := cmnd.Output()

		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok {
				t.Log(string(exitError.Stderr))
			} else {
				t.Log(err.Error())
			}

			return err
		}

		t.Log(string(out))

		// when initialising the new space, we need to define a new space name as a variable
		vars := []string{}
		if i == 0 {
			vars = append(initialiseVars, "-var=octopus_space_name="+spaceName)
		} else {
			vars = populateVars
		}

		newArgs := append([]string{
			"apply",
			"-auto-approve",
			"-no-color",
			"-var=octopus_server=" + container.URI,
			"-var=octopus_apikey=" + ApiKey,
			"-var=octopus_space_id=" + spaceId,
		}, vars...)

		cmnd = exec.Command("terraform", newArgs...)
		cmnd.Dir = terraformProjectDir
		out, err = cmnd.Output()

		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok {
				t.Log(string(exitError.Stderr))
			} else {
				t.Log(err)
			}
			return err
		}

		t.Log(string(out))

		// get the ID of any new space created, which will be used in the subsequent Terraform executions
		spaceId, err = getOutputVariable(t, terraformProjectDir, "octopus_space_id")

		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok {
				t.Log(string(exitError.Stderr))
			} else {
				t.Log(err)
			}
			return err
		}
	}

	return nil
}

// getOutputVariable reads a Terraform output variable
func getOutputVariable(t *testing.T, terraformDir string, outputVar string) (string, error) {
	cmnd := exec.Command(
		"terraform",
		"output",
		"-raw",
		outputVar)
	cmnd.Dir = terraformDir
	out, err := cmnd.Output()

	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			t.Log(string(exitError.Stderr))
		} else {
			t.Log(err)
		}
		return "", err
	}

	return string(out), nil
}

// createClient creates a Octopus client to the given url
func createClient(uri string, spaceId string) (*client.Client, error) {
	url, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	return client.NewClient(nil, url, ApiKey, spaceId)
}

// act initialises Octopus and MSSQL
func act(t *testing.T, container *octopusContainer, terraformDir string, populateVars []string) (string, error) {
	t.Log("POPULATING TEST SPACE")

	spaceName := strings.ReplaceAll(fmt.Sprint(uuid.New()), "-", "")[:20]
	err := initialiseOctopus(t, container, terraformDir, spaceName, []string{}, populateVars)

	if err != nil {
		return "", err
	}

	return getOutputVariable(t, filepath.Join("test", "terraform", "1-singlespace"), "octopus_space_id")
}

// TestSpaceResource verifies that a space can be reimported with the correct settings
func TestSpaceResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/1-singlespace", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, "")
		query := spaces.SpacesQuery{
			IDs:  []string{newSpaceId},
			Skip: 0,
			Take: 1,
		}
		spaces, err := client.Spaces.Get(query)
		space := spaces.Items[0]

		if err != nil {
			return err
		}

		if space.Description != "My test space" {
			t.Fatalf("New space must have the name \"My test space\"")
		}

		if space.IsDefault {
			t.Fatalf("New space must not be the default one")
		}

		if space.TaskQueueStopped {
			t.Fatalf("New space must not have the task queue stopped")
		}

		if slices.Index(space.SpaceManagersTeams, "teams-administrators") == -1 {
			t.Fatalf("New space must have teams-administrators as a manager team")
		}

		return nil
	})
}

// TestProjectGroupResource verifies that a project group can be reimported with the correct settings
func TestProjectGroupResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/2-projectgroup", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := projectgroups.ProjectGroupsQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.ProjectGroups.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a project group called \"Test\"")
		}
		resource := resources.Items[0]

		if resource.Description != "Test Description" {
			t.Fatalf("The project group must be have a description of \"Test Description\"")
		}

		return nil
	})
}

// TestAwsAccountExport verifies that an AWS account can be reimported with the correct settings
func TestAwsAccountExport(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/3-awsaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := accounts.AccountsQuery{
			PartialName: "AWS Account",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Accounts.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an account called \"AWS Account\"")
		}
		resource := resources.Items[0].(*accounts.AmazonWebServicesAccount)

		if resource.AccessKey != "ABCDEFGHIJKLMNOPQRST" {
			t.Fatalf("The account must have an access key of \"ABCDEFGHIJKLMNOPQRST\"")
		}

		return nil
	})
}

// TestAzureAccountResource verifies that an Azure account can be reimported with the correct settings
func TestAzureAccountResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/4-azureaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := accounts.AccountsQuery{
			PartialName: "Azure",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Accounts.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an account called \"Azure\"")
		}
		resource := resources.Items[0].(*accounts.AzureServicePrincipalAccount)

		if fmt.Sprint(resource.SubscriptionID) != "95bf77d2-64b1-4ed2-9de1-b5451e3881f5" {
			t.Fatalf("The account must be have a client ID of \"95bf77d2-64b1-4ed2-9de1-b5451e3881f5\"")
		}

		if fmt.Sprint(resource.TenantID) != "18eb006b-c3c8-4a72-93cd-fe4b293f82ee" {
			t.Fatalf("The account must be have a client ID of \"18eb006b-c3c8-4a72-93cd-fe4b293f82ee\"")
		}

		if resource.Description != "Azure Account" {
			t.Fatalf("The account must be have a description of \"Azure Account\"")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatalf("The account must be have a tenanted deployment participation of \"Untenanted\"")
		}

		return nil
	})
}

// TestUsernamePasswordAccountResource verifies that a username/password account can be reimported with the correct settings
func TestUsernamePasswordAccountResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/5-userpassaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := accounts.AccountsQuery{
			PartialName: "GKE",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Accounts.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an account called \"GKE\"")
		}
		resource := resources.Items[0].(*accounts.UsernamePasswordAccount)

		if resource.Username != "admin" {
			t.Fatalf("The account must be have a username of \"admin\"")
		}

		if !resource.Password.HasValue {
			t.Fatalf("The account must be have a password")
		}

		if resource.Description != "A test account" {
			t.Fatalf("The account must be have a description of \"A test account\"")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatalf("The account must be have a tenanted deployment participation of \"Untenanted\"")
		}

		if len(resource.TenantTags) != 0 {
			t.Fatalf("The account must be have no tenant tags")
		}

		return nil
	})
}

// TestGcpAccountResource verifies that a GCP account can be reimported with the correct settings
func TestGcpAccountResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/6-gcpaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := accounts.AccountsQuery{
			PartialName: "Google",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Accounts.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an account called \"Google\"")
		}
		resource := resources.Items[0].(*accounts.GoogleCloudPlatformAccount)

		if !resource.JsonKey.HasValue {
			t.Fatalf("The account must be have a JSON key")
		}

		if resource.Description != "A test account" {
			t.Fatalf("The account must be have a description of \"A test account\"")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatalf("The account must be have a tenanted deployment participation of \"Untenanted\"")
		}

		if len(resource.TenantTags) != 0 {
			t.Fatalf("The account must be have no tenant tags")
		}

		return nil
	})
}

// TestSshAccountResource verifies that an SSH account can be reimported with the correct settings
func TestSshAccountResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/7-sshaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := accounts.AccountsQuery{
			PartialName: "SSH",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Accounts.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an account called \"SSH\"")
		}
		resource := resources.Items[0].(*accounts.SSHKeyAccount)

		if resource.AccountType != "SshKeyPair" {
			t.Fatal("The account must be have a type of \"SshKeyPair\"")
		}

		if resource.Username != "admin" {
			t.Fatal("The account must be have a username of \"admin\"")
		}

		if resource.Description != "A test account" {
			// This appears to be a bug in the provider where the description is not set
			t.Log("BUG: The account must be have a description of \"A test account\"")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The account must be have a tenanted deployment participation of \"Untenanted\"")
		}

		if len(resource.TenantTags) != 0 {
			t.Fatal("The account must be have no tenant tags")
		}

		return nil
	})
}

// TestAzureSubscriptionAccountResource verifies that an azure account can be reimported with the correct settings
func TestAzureSubscriptionAccountResource(t *testing.T) {
	// I could not figure out a combination of properties that made this resource work
	return

	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/8-azuresubscriptionaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := accounts.AccountsQuery{
			PartialName: "Subscription",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Accounts.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an account called \"Subscription\"")
		}
		resource := resources.Items[0].(*accounts.AzureSubscriptionAccount)

		if resource.AccountType != "AzureServicePrincipal" {
			t.Fatal("The account must be have a type of \"AzureServicePrincipal\"")
		}

		if resource.Description != "A test account" {
			t.Fatal("BUG: The account must be have a description of \"A test account\"")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The account must be have a tenanted deployment participation of \"Untenanted\"")
		}

		if len(resource.TenantTags) != 0 {
			t.Fatal("The account must be have no tenant tags")
		}

		return nil
	})
}

// TestTokenAccountResource verifies that a token account can be reimported with the correct settings
func TestTokenAccountResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/9-tokenaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := accounts.AccountsQuery{
			PartialName: "Token",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Accounts.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an account called \"Token\"")
		}
		resource := resources.Items[0].(*accounts.TokenAccount)

		if resource.AccountType != "Token" {
			t.Fatal("The account must be have a type of \"Token\"")
		}

		if !resource.Token.HasValue {
			t.Fatal("The account must be have a token")
		}

		if resource.Description != "A test account" {
			t.Fatal("The account must be have a description of \"A test account\"")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The account must be have a tenanted deployment participation of \"Untenanted\"")
		}

		if len(resource.TenantTags) != 0 {
			t.Fatal("The account must be have no tenant tags")
		}

		return nil
	})
}

// TestHelmFeedResource verifies that a helm feed can be reimported with the correct settings
func TestHelmFeedResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/10-helmfeed", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := feeds.FeedsQuery{
			PartialName: "Helm",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Feeds.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an feed called \"Helm\"")
		}
		resource := resources.Items[0].(*feeds.HelmFeed)

		if resource.FeedType != "Helm" {
			t.Fatal("The feed must have a type of \"Helm\"")
		}

		if resource.Username != "username" {
			t.Fatal("The feed must have a username of \"username\"")
		}

		if resource.FeedURI != "https://charts.helm.sh/stable/" {
			t.Fatal("The feed must be have a URI of \"https://charts.helm.sh/stable/\"")
		}

		foundExecutionTarget := false
		foundNotAcquired := false
		for _, o := range resource.PackageAcquisitionLocationOptions {
			if o == "ExecutionTarget" {
				foundExecutionTarget = true
			}

			if o == "NotAcquired" {
				foundNotAcquired = true
			}
		}

		if !(foundExecutionTarget && foundNotAcquired) {
			t.Fatal("The feed must be have a PackageAcquisitionLocationOptions including \"ExecutionTarget\" and \"NotAcquired\"")
		}

		return nil
	})
}

// TestDockerFeedResource verifies that a docker feed can be reimported with the correct settings
func TestDockerFeedResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/11-dockerfeed", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := feeds.FeedsQuery{
			PartialName: "Docker",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Feeds.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an feed called \"Docker\"")
		}
		resource := resources.Items[0].(*feeds.DockerContainerRegistry)

		if resource.FeedType != "Docker" {
			t.Fatal("The feed must have a type of \"Docker\"")
		}

		if resource.Username != "username" {
			t.Fatal("The feed must have a username of \"username\"")
		}

		if resource.APIVersion != "v1" {
			t.Fatal("The feed must be have a API version of \"v1\"")
		}

		if resource.FeedURI != "https://index.docker.io" {
			t.Fatal("The feed must be have a feed uri of \"https://index.docker.io\"")
		}

		foundExecutionTarget := false
		foundNotAcquired := false
		for _, o := range resource.PackageAcquisitionLocationOptions {
			if o == "ExecutionTarget" {
				foundExecutionTarget = true
			}

			if o == "NotAcquired" {
				foundNotAcquired = true
			}
		}

		if !(foundExecutionTarget && foundNotAcquired) {
			t.Fatal("The feed must be have a PackageAcquisitionLocationOptions including \"ExecutionTarget\" and \"NotAcquired\"")
		}

		return nil
	})
}

// TestEcrFeedResource verifies that a ecr feed can be reimported with the correct settings
func TestEcrFeedResource(t *testing.T) {
	if os.Getenv("ECR_ACCESS_KEY") == "" {
		t.Fatal("The ECR_ACCESS_KEY environment variable must be set a valid AWS access key")
	}

	if os.Getenv("ECR_SECRET_KEY") == "" {
		t.Fatal("The ECR_SECRET_KEY environment variable must be set a valid AWS secret key")
	}

	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act

		newSpaceId, err := act(t, container, "./test/terraform/12-ecrfeed", []string{
			"-var=feed_ecr_access_key=" + os.Getenv("ECR_ACCESS_KEY"),
			"-var=feed_ecr_secret_key=" + os.Getenv("ECR_SECRET_KEY"),
		})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := feeds.FeedsQuery{
			PartialName: "ECR",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Feeds.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an feed called \"ECR\"")
		}
		resource := resources.Items[0].(*feeds.AwsElasticContainerRegistry)

		if resource.FeedType != "AwsElasticContainerRegistry" {
			t.Fatal("The feed must have a type of \"AwsElasticContainerRegistry\" (was \"" + resource.FeedType + "\"")
		}

		if resource.AccessKey != os.Getenv("ECR_ACCESS_KEY") {
			t.Fatal("The feed must have a access key of \"" + os.Getenv("ECR_ACCESS_KEY") + "\" (was \"" + resource.AccessKey + "\"")
		}

		if resource.Region != "us-east-1" {
			t.Fatal("The feed must have a region of \"us-east-1\" (was \"" + resource.Region + "\"")
		}

		foundExecutionTarget := false
		foundNotAcquired := false
		for _, o := range resource.PackageAcquisitionLocationOptions {
			if o == "ExecutionTarget" {
				foundExecutionTarget = true
			}

			if o == "NotAcquired" {
				foundNotAcquired = true
			}
		}

		if !(foundExecutionTarget && foundNotAcquired) {
			t.Fatal("The feed must be have a PackageAcquisitionLocationOptions including \"ExecutionTarget\" and \"NotAcquired\"")
		}

		return nil
	})
}

// TestMavenFeedResource verifies that a maven feed can be reimported with the correct settings
func TestMavenFeedResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/13-mavenfeed", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := feeds.FeedsQuery{
			PartialName: "Maven",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Feeds.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an feed called \"Maven\"")
		}
		resource := resources.Items[0].(*feeds.MavenFeed)

		if resource.FeedType != "Maven" {
			t.Fatal("The feed must have a type of \"Maven\"")
		}

		if resource.Username != "username" {
			t.Fatal("The feed must have a username of \"username\"")
		}

		if resource.DownloadAttempts != 5 {
			t.Fatal("The feed must be have a downloads attempts set to \"5\"")
		}

		if resource.DownloadRetryBackoffSeconds != 10 {
			t.Fatal("The feed must be have a downloads retry backoff set to \"10\"")
		}

		if resource.FeedURI != "https://repo.maven.apache.org/maven2/" {
			t.Fatal("The feed must be have a feed uri of \"https://repo.maven.apache.org/maven2/\"")
		}

		foundExecutionTarget := false
		foundServer := false
		for _, o := range resource.PackageAcquisitionLocationOptions {
			if o == "ExecutionTarget" {
				foundExecutionTarget = true
			}

			if o == "Server" {
				foundServer = true
			}
		}

		if !(foundExecutionTarget && foundServer) {
			t.Fatal("The feed must be have a PackageAcquisitionLocationOptions including \"ExecutionTarget\" and \"Server\"")
		}

		return nil
	})
}

// TestNugetFeedResource verifies that a nuget feed can be reimported with the correct settings
func TestNugetFeedResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/14-nugetfeed", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := feeds.FeedsQuery{
			PartialName: "Nuget",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Feeds.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an feed called \"Nuget\"")
		}
		resource := resources.Items[0].(*feeds.NuGetFeed)

		if resource.FeedType != "NuGet" {
			t.Fatal("The feed must have a type of \"NuGet\"")
		}

		if !resource.EnhancedMode {
			t.Fatal("The feed must have enhanced mode set to true")
		}

		if resource.Username != "username" {
			t.Fatal("The feed must have a username of \"username\"")
		}

		if resource.DownloadAttempts != 5 {
			t.Fatal("The feed must be have a downloads attempts set to \"5\"")
		}

		if resource.DownloadRetryBackoffSeconds != 10 {
			t.Fatal("The feed must be have a downloads retry backoff set to \"10\"")
		}

		if resource.FeedURI != "https://index.docker.io" {
			t.Fatal("The feed must be have a feed uri of \"https://index.docker.io\"")
		}

		foundExecutionTarget := false
		foundServer := false
		for _, o := range resource.PackageAcquisitionLocationOptions {
			if o == "ExecutionTarget" {
				foundExecutionTarget = true
			}

			if o == "Server" {
				foundServer = true
			}
		}

		if !(foundExecutionTarget && foundServer) {
			t.Fatal("The feed must be have a PackageAcquisitionLocationOptions including \"ExecutionTarget\" and \"Server\"")
		}

		return nil
	})
}

// TestWorkerPoolResource verifies that a static worker pool can be reimported with the correct settings
func TestWorkerPoolResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/15-workerpool", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := workerpools.WorkerPoolsQuery{
			PartialName: "Docker",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.WorkerPools.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a worker pool called \"Docker\"")
		}
		resource := resources.Items[0].(*workerpools.StaticWorkerPool)

		if resource.WorkerPoolType != "StaticWorkerPool" {
			t.Fatal("The worker pool must be have a type of \"StaticWorkerPool\" (was \"" + resource.WorkerPoolType + "\"")
		}

		if resource.Description != "A test worker pool" {
			t.Fatal("The worker pool must be have a description of \"A test worker pool\" (was \"" + resource.Description + "\"")
		}

		if resource.SortOrder != 3 {
			t.Fatal("The worker pool must be have a sort order of \"3\" (was \"" + fmt.Sprint(resource.SortOrder) + "\"")
		}

		if resource.IsDefault {
			t.Fatal("The worker pool must be must not be the default")
		}

		return nil
	})
}

// TestEnvironmentResource verifies that an environment can be reimported with the correct settings
func TestEnvironmentResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/16-environment", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := environments.EnvironmentsQuery{
			PartialName: "Development",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Environments.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an environment called \"Development\"")
		}
		resource := resources.Items[0]

		if resource.Description != "A test environment" {
			t.Fatal("The environment must be have a description of \"A test environment\" (was \"" + resource.Description + "\"")
		}

		if !resource.AllowDynamicInfrastructure {
			t.Fatal("The environment must have dynamic infrastructure enabled.")
		}

		if resource.UseGuidedFailure {
			t.Fatal("The environment must not have guided failure enabled.")
		}

		return nil
	})
}

// TestLifecycleResource verifies that a lifecycle can be reimported with the correct settings
func TestLifecycleResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/17-lifecycle", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := lifecycles.Query{
			PartialName: "Simple",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Lifecycles.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an environment called \"Simple\"")
		}
		resource := resources.Items[0]

		if resource.Description != "A test lifecycle" {
			t.Fatal("The lifecycle must be have a description of \"A test lifecycle\" (was \"" + resource.Description + "\")")
		}

		if resource.TentacleRetentionPolicy.QuantityToKeep != 30 {
			t.Fatal("The lifecycle must be have a tentacle retention policy of \"30\" (was \"" + fmt.Sprint(resource.TentacleRetentionPolicy.QuantityToKeep) + "\")")
		}

		if resource.TentacleRetentionPolicy.ShouldKeepForever {
			t.Fatal("The lifecycle must be have a tentacle retention not set to keep forever")
		}

		if resource.TentacleRetentionPolicy.Unit != "Items" {
			t.Fatal("The lifecycle must be have a tentacle retention unit set to \"Items\" (was \"" + resource.TentacleRetentionPolicy.Unit + "\")")
		}

		if resource.ReleaseRetentionPolicy.QuantityToKeep != 1 {
			t.Fatal("The lifecycle must be have a release retention policy of \"1\" (was \"" + fmt.Sprint(resource.ReleaseRetentionPolicy.QuantityToKeep) + "\")")
		}

		if !resource.ReleaseRetentionPolicy.ShouldKeepForever {
			t.Log("BUG: The lifecycle must be have a release retention set to keep forever (known bug - the provider creates this field as false)")
		}

		if resource.ReleaseRetentionPolicy.Unit != "Days" {
			t.Fatal("The lifecycle must be have a release retention unit set to \"Days\" (was \"" + resource.ReleaseRetentionPolicy.Unit + "\")")
		}

		return nil
	})
}

// TestVariableSetResource verifies that a variable set can be reimported with the correct settings
func TestVariableSetResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/18-variableset", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := variables.LibraryVariablesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.LibraryVariableSets.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a library variable set called \"Test\"")
		}
		resource := resources.Items[0]

		if resource.Description != "Test variable set" {
			t.Fatal("The library variable set must be have a description of \"Test variable set\" (was \"" + resource.Description + "\")")
		}

		variableSet, err := client.Variables.GetAll(resource.ID)

		if len(variableSet.Variables) != 1 {
			t.Fatal("The library variable set must have one associated variable")
		}

		if variableSet.Variables[0].Name != "Test.Variable" {
			t.Fatal("The library variable set variable must have a name of \"Test.Variable\"")
		}

		if variableSet.Variables[0].Type != "String" {
			t.Fatal("The library variable set variable must have a type of \"String\"")
		}

		if variableSet.Variables[0].Description != "Test variable" {
			t.Fatal("The library variable set variable must have a description of \"Test variable\"")
		}

		if variableSet.Variables[0].Value != "test" {
			t.Fatal("The library variable set variable must have a value of \"test\"")
		}

		if variableSet.Variables[0].IsSensitive {
			t.Fatal("The library variable set variable must not be sensitive")
		}

		if !variableSet.Variables[0].IsEditable {
			t.Fatal("The library variable set variable must be editable")
		}

		return nil
	})
}

// TestProjectResource verifies that a project can be reimported with the correct settings
func TestProjectResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/19-project", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := projects.ProjectsQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Projects.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a project called \"Test\"")
		}
		resource := resources.Items[0]

		if resource.Description != "Test project" {
			t.Fatal("The project must be have a description of \"Test project\" (was \"" + resource.Description + "\")")
		}

		if resource.AutoCreateRelease {
			t.Fatal("The project must not have auto release create enabled")
		}

		if resource.DefaultGuidedFailureMode != "EnvironmentDefault" {
			t.Fatal("The project must be have a DefaultGuidedFailureMode of \"EnvironmentDefault\" (was \"" + resource.DefaultGuidedFailureMode + "\")")
		}

		if resource.DefaultToSkipIfAlreadyInstalled {
			t.Fatal("The project must not have DefaultToSkipIfAlreadyInstalled enabled")
		}

		if resource.IsDisabled {
			t.Fatal("The project must not have IsDisabled enabled")
		}

		if resource.IsVersionControlled {
			t.Fatal("The project must not have IsVersionControlled enabled")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The project must be have a TenantedDeploymentMode of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		if len(resource.IncludedLibraryVariableSets) != 0 {
			t.Fatal("The project must not have any library variable sets")
		}

		if resource.ConnectivityPolicy.AllowDeploymentsToNoTargets {
			t.Fatal("The project must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
		}

		if resource.ConnectivityPolicy.ExcludeUnhealthyTargets {
			t.Fatal("The project must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
		}

		if resource.ConnectivityPolicy.SkipMachineBehavior != "SkipUnavailableMachines" {
			t.Log("BUG: The project must be have a ConnectivityPolicy.SkipMachineBehavior of \"SkipUnavailableMachines\" (was \"" + resource.ConnectivityPolicy.SkipMachineBehavior + "\") - Known issue where the value returned by /api/Spaces-#/ProjectGroups/ProjectGroups-#/projects is different to /api/Spaces-/Projects")
		}

		return nil
	})
}

// TestProjectChannelResource verifies that a project channel can be reimported with the correct settings
func TestProjectChannelResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/20-channel", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := channels.Query{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Channels.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a channel called \"Test\"")
		}
		resource := resources.Items[0]

		if resource.Description != "Test channel" {
			t.Fatal("The channel must be have a description of \"Test channel\" (was \"" + resource.Description + "\")")
		}

		if !resource.IsDefault {
			t.Fatal("The channel must be be the default")
		}

		if len(resource.Rules) != 1 {
			t.Fatal("The channel must have one rule")
		}

		if resource.Rules[0].Tag != "^$" {
			t.Fatal("The channel rule must be have a tag of \"^$\" (was \"" + resource.Rules[0].Tag + "\")")
		}

		if resource.Rules[0].ActionPackages[0].DeploymentAction != "Test" {
			t.Fatal("The channel rule action step must be be set to \"Test\" (was \"" + resource.Rules[0].ActionPackages[0].DeploymentAction + "\")")
		}

		if resource.Rules[0].ActionPackages[0].PackageReference != "test" {
			t.Fatal("The channel rule action package must be be set to \"test\" (was \"" + resource.Rules[0].ActionPackages[0].PackageReference + "\")")
		}

		return nil
	})
}

// TestTagSetResource verifies that a tag set can be reimported with the correct settings
func TestTagSetResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/21-tagset", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := tagsets.TagSetsQuery{
			PartialName: "tag1",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.TagSets.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a tag set called \"tag1\"")
		}
		resource := resources.Items[0]

		if resource.Description != "Test tagset" {
			t.Fatal("The tag set must be have a description of \"Test tagset\" (was \"" + resource.Description + "\")")
		}

		if resource.SortOrder != 0 {
			t.Fatal("The tag set must be have a sort order of \"0\" (was \"" + fmt.Sprint(resource.SortOrder) + "\")")
		}

		tagAFound := false
		for _, u := range resource.Tags {
			if u.Name == "a" {
				tagAFound = true

				if u.Description != "tag a" {
					t.Fatal("The tag a must be have a description of \"tag a\" (was \"" + u.Description + "\")")
				}

				if u.Color != "#333333" {
					t.Fatal("The tag a must be have a color of \"#333333\" (was \"" + u.Color + "\")")
				}

				if u.SortOrder != 2 {
					t.Fatal("The tag a must be have a sort order of \"2\" (was \"" + fmt.Sprint(u.SortOrder) + "\")")
				}
			}
		}

		if !tagAFound {
			t.Fatal("Tag Set must have an tag called \"a\"")
		}

		return nil
	})
}

// TestGitCredentialsResource verifies that a git credential can be reimported with the correct settings
func TestGitCredentialsResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		_, err := act(t, container, "./test/terraform/22-gitcredentialtest", []string{})

		if err != nil {
			return err
		}

		// Octopus client does not expose git credentials, so we just test that they can be imported

		return nil
	})
}

// TestScriptModuleResource verifies that a script module set can be reimported with the correct settings
func TestScriptModuleResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/23-scriptmodule", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := variables.LibraryVariablesQuery{
			PartialName: "Test2",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.LibraryVariableSets.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a library variable set called \"Test2\"")
		}
		resource := resources.Items[0]

		if resource.Description != "Test script module" {
			t.Fatal("The library variable set must be have a description of \"Test script module\" (was \"" + resource.Description + "\")")
		}

		variables, err := client.Variables.GetAll(resource.ID)

		if len(variables.Variables) != 2 {
			t.Fatal("The library variable set must have two associated variables")
		}

		foundScript := false
		foundLanguage := false
		for _, u := range variables.Variables {
			if u.Name == "Octopus.Script.Module[Test2]" {
				foundScript = true

				if u.Type != "String" {
					t.Fatal("The library variable set variable must have a type of \"String\"")
				}

				if u.Value != "echo \"hi\"" {
					t.Fatal("The library variable set variable must have a value of \"\"echo \\\"hi\\\"\"\"")
				}

				if u.IsSensitive {
					t.Fatal("The library variable set variable must not be sensitive")
				}

				if !u.IsEditable {
					t.Fatal("The library variable set variable must be editable")
				}
			}

			if u.Name == "Octopus.Script.Module.Language[Test2]" {
				foundLanguage = true

				if u.Type != "String" {
					t.Fatal("The library variable set variable must have a type of \"String\"")
				}

				if u.Value != "PowerShell" {
					t.Fatal("The library variable set variable must have a value of \"PowerShell\"")
				}

				if u.IsSensitive {
					t.Fatal("The library variable set variable must not be sensitive")
				}

				if !u.IsEditable {
					t.Fatal("The library variable set variable must be editable")
				}
			}
		}

		if !foundLanguage || !foundScript {
			t.Fatal("Script module must create two variables for script and language")
		}

		return nil
	})
}

// TestTenantsResource verifies that a git credential can be reimported with the correct settings
func TestTenantsResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/24-tenants", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := tenants.TenantsQuery{
			PartialName: "Team A",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Tenants.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a tenant called \"Team A\"")
		}
		resource := resources.Items[0]

		if resource.Description != "Test tenant" {
			t.Fatal("The tenant must be have a description of \"tTest tenant\" (was \"" + resource.Description + "\")")
		}

		if len(resource.TenantTags) != 2 {
			t.Fatal("The tenant must have two tags")
		}

		if len(resource.ProjectEnvironments) != 1 {
			t.Fatal("The tenant must have one project environment")
		}

		for _, u := range resource.ProjectEnvironments {
			if len(u) != 3 {
				t.Fatal("The tenant must have be linked to three environments")
			}
		}

		return nil
	})
}

// TestCertificateResource verifies that a certificate can be reimported with the correct settings
func TestCertificateResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/25-certificates", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := certificates.CertificatesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Certificates.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a certificate called \"Test\"")
		}
		resource := resources.Items[0]

		if resource.Notes != "A test certificate" {
			t.Fatal("The tenant must be have a description of \"A test certificate\" (was \"" + resource.Notes + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The tenant must be have a tenant participation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		if resource.SubjectDistinguishedName != "CN=test.com" {
			t.Fatal("The tenant must be have a subject distinguished name of \"CN=test.com\" (was \"" + resource.SubjectDistinguishedName + "\")")
		}

		if len(resource.EnvironmentIDs) != 0 {
			t.Fatal("The tenant must have one project environment")
		}

		if len(resource.TenantTags) != 0 {
			t.Fatal("The tenant must have no tenant tags")
		}

		if len(resource.TenantIDs) != 0 {
			t.Fatal("The tenant must have no tenants")
		}

		return nil
	})
}

// TestTenantVariablesResource verifies that a tenant variables can be reimported with the correct settings
func TestTenantVariablesResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/26-tenant_variables", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		collection, err := client.TenantVariables.GetAll()
		if err != nil {
			return err
		}

		resourceName := "Test"
		found := false
		for _, tenantVariable := range collection {
			for _, project := range tenantVariable.ProjectVariables {
				if project.ProjectName == resourceName {
					for _, variables := range project.Variables {
						for _, value := range variables {
							// we expect one project variable to be defined
							found = true
							if value.Value != "my value" {
								t.Fatal("The tenant project variable must have a value of \"my value\" (was \"" + value.Value + "\")")
							}
						}
					}
				}
			}
		}

		if !found {
			t.Fatal("Space must have an tenant project variable for the project called \"" + resourceName + "\"")
		}

		return nil
	})
}

// TestMachinePolicyResource verifies that a machine policies can be reimported with the correct settings
func TestMachinePolicyResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/27-machinepolicy", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := machines.MachinePoliciesQuery{
			PartialName: "Testing",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.MachinePolicies.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine policy called \"Testing\"")
		}
		resource := resources.Items[0]

		if resource.Description != "test machine policy" {
			t.Fatal("The machine policy must have a description of \"test machine policy\" (was \"" + resource.Description + "\")")
		}

		if resource.ConnectionConnectTimeout.Minutes() != 1 {
			t.Fatal("The machine policy must have a ConnectionConnectTimeout of \"00:01:00\" (was \"" + fmt.Sprint(resource.ConnectionConnectTimeout) + "\")")
		}

		if resource.ConnectionRetryCountLimit != 5 {
			t.Fatal("The machine policy must have a ConnectionRetryCountLimit of \"5\" (was \"" + fmt.Sprint(resource.ConnectionRetryCountLimit) + "\")")
		}

		if resource.ConnectionRetrySleepInterval.Seconds() != 1 {
			t.Fatal("The machine policy must have a ConnectionRetrySleepInterval of \"00:00:01\" (was \"" + fmt.Sprint(resource.ConnectionRetrySleepInterval) + "\")")
		}

		if resource.ConnectionRetryTimeLimit.Minutes() != 5 {
			t.Fatal("The machine policy must have a ConnectionRetryTimeLimit of \"00:05:00\" (was \"" + fmt.Sprint(resource.ConnectionRetryTimeLimit) + "\")")
		}

		if resource.PollingRequestMaximumMessageProcessingTimeout.Minutes() != 10 {
			t.Fatal("The machine policy must have a PollingRequestMaximumMessageProcessingTimeout of \"00:10:00\" (was \"" + fmt.Sprint(resource.PollingRequestMaximumMessageProcessingTimeout) + "\")")
		}

		if resource.MachineCleanupPolicy.DeleteMachinesElapsedTimeSpan.Minutes() != 20 {
			t.Fatal("The machine policy must have a DeleteMachinesElapsedTimeSpan of \"00:20:00\" (was \"" + fmt.Sprint(resource.MachineCleanupPolicy.DeleteMachinesElapsedTimeSpan) + "\")")
		}

		if resource.MachineCleanupPolicy.DeleteMachinesBehavior != "DeleteUnavailableMachines" {
			t.Fatal("The machine policy must have a MachineCleanupPolicy.DeleteMachinesBehavior of \"DeleteUnavailableMachines\" (was \"" + resource.MachineCleanupPolicy.DeleteMachinesBehavior + "\")")
		}

		if resource.MachineConnectivityPolicy.MachineConnectivityBehavior != "ExpectedToBeOnline" {
			t.Fatal("The machine policy must have a MachineConnectivityPolicy.MachineConnectivityBehavior of \"ExpectedToBeOnline\" (was \"" + resource.MachineConnectivityPolicy.MachineConnectivityBehavior + "\")")
		}

		if resource.MachineHealthCheckPolicy.BashHealthCheckPolicy.RunType != "Inline" {
			t.Fatal("The machine policy must have a MachineHealthCheckPolicy.BashHealthCheckPolicy.RunType of \"Inline\" (was \"" + resource.MachineHealthCheckPolicy.BashHealthCheckPolicy.RunType + "\")")
		}

		if *resource.MachineHealthCheckPolicy.BashHealthCheckPolicy.ScriptBody != "" {
			t.Fatal("The machine policy must have a MachineHealthCheckPolicy.BashHealthCheckPolicy.ScriptBody of \"\" (was \"" + *resource.MachineHealthCheckPolicy.BashHealthCheckPolicy.ScriptBody + "\")")
		}

		if resource.MachineHealthCheckPolicy.PowerShellHealthCheckPolicy.RunType != "Inline" {
			t.Fatal("The machine policy must have a MachineHealthCheckPolicy.PowerShellHealthCheckPolicy.RunType of \"Inline\" (was \"" + resource.MachineHealthCheckPolicy.PowerShellHealthCheckPolicy.RunType + "\")")
		}

		if strings.HasPrefix(*resource.MachineHealthCheckPolicy.BashHealthCheckPolicy.ScriptBody, "$freeDiskSpaceThreshold") {
			t.Fatal("The machine policy must have a MachineHealthCheckPolicy.PowerShellHealthCheckPolicy.ScriptBody to start with \"$freeDiskSpaceThreshold\" (was \"" + *resource.MachineHealthCheckPolicy.PowerShellHealthCheckPolicy.ScriptBody + "\")")
		}

		if resource.MachineHealthCheckPolicy.HealthCheckCronTimezone != "UTC" {
			t.Fatal("The machine policy must have a MachineHealthCheckPolicy.HealthCheckCronTimezone of \"UTC\" (was \"" + resource.MachineHealthCheckPolicy.HealthCheckCronTimezone + "\")")
		}

		if resource.MachineHealthCheckPolicy.HealthCheckCron != "" {
			t.Fatal("The machine policy must have a MachineHealthCheckPolicy.HealthCheckCron of \"\" (was \"" + resource.MachineHealthCheckPolicy.HealthCheckCron + "\")")
		}

		if resource.MachineHealthCheckPolicy.HealthCheckType != "RunScript" {
			t.Fatal("The machine policy must have a MachineHealthCheckPolicy.HealthCheckType of \"RunScript\" (was \"" + resource.MachineHealthCheckPolicy.HealthCheckType + "\")")
		}

		if resource.MachineHealthCheckPolicy.HealthCheckInterval.Minutes() != 10 {
			t.Fatal("The machine policy must have a MachineHealthCheckPolicy.HealthCheckInterval of \"00:10:00\" (was \"" + fmt.Sprint(resource.MachineHealthCheckPolicy.HealthCheckInterval) + "\")")
		}

		if resource.MachineUpdatePolicy.CalamariUpdateBehavior != "UpdateOnDeployment" {
			t.Fatal("The machine policy must have a MachineUpdatePolicy.CalamariUpdateBehavior of \"UpdateOnDeployment\" (was \"" + resource.MachineUpdatePolicy.CalamariUpdateBehavior + "\")")
		}

		if resource.MachineUpdatePolicy.TentacleUpdateBehavior != "NeverUpdate" {
			t.Fatal("The machine policy must have a MachineUpdatePolicy.TentacleUpdateBehavior of \"NeverUpdate\" (was \"" + resource.MachineUpdatePolicy.CalamariUpdateBehavior + "\")")
		}

		return nil
	})
}

// TestProjectTriggerResource verifies that a project trigger can be reimported with the correct settings
func TestProjectTriggerResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/28-projecttrigger", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := projects.ProjectsQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Projects.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a project called \"Test\"")
		}
		resource := resources.Items[0]

		trigger, err := client.ProjectTriggers.GetByProjectID(resource.ID)

		if err != nil {
			return err
		}

		if trigger[0].Name != "test" {
			t.Fatal("The project must have a trigger called \"test\" (was \"" + trigger[0].Name + "\")")
		}

		if trigger[0].Filter.GetFilterType() != filters.MachineFilter {
			t.Fatal("The project trigger must have Filter.FilterType set to \"MachineFilter\" (was \"" + fmt.Sprint(trigger[0].Filter.GetFilterType()) + "\")")
		}

		return nil
	})
}

// TestK8sTargetResource verifies that a k8s machine can be reimported with the correct settings
func TestK8sTargetResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/29-k8starget", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := machines.MachinesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Test\"")
		}
		resource := resources.Items[0]

		if fmt.Sprint(resource.Endpoint.(*machines.KubernetesEndpoint).ClusterURL) != "https://cluster" {
			t.Fatal("The machine must have a Endpoint.ClusterUrl of \"https://cluster\" (was \"" + fmt.Sprint(resource.Endpoint.(*machines.KubernetesEndpoint).ClusterURL) + "\")")
		}

		return nil
	})
}

// TestSshTargetResource verifies that a ssh machine can be reimported with the correct settings
func TestSshTargetResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/30-sshtarget", []string{
			"-var=account_ec2_sydney=LS0tLS1CRUdJTiBFTkNSWVBURUQgUFJJVkFURSBLRVktLS0tLQpNSUlKbkRCT0Jna3Foa2lHOXcwQkJRMHdRVEFwQmdrcWhraUc5dzBCQlF3d0hBUUlwNEUxV1ZrejJEd0NBZ2dBCk1Bd0dDQ3FHU0liM0RRSUpCUUF3RkFZSUtvWklodmNOQXdjRUNIemFuVE1QbHA4ZkJJSUpTSncrdW5BL2ZaVFUKRGdrdWk2QnhOY0REUFg3UHZJZmNXU1dTc3V3YWRhYXdkVEdjY1JVd3pGNTNmRWJTUXJBYzJuWFkwUWVVcU1wcAo4QmdXUUthWlB3MEdqck5OQVJaTy9QYklxaU5ERFMybVRSekZidzREcFY5aDdlblZjL1ZPNlhJdzlxYVYzendlCnhEejdZSkJ2ckhmWHNmTmx1blErYTZGdlRUVkVyWkE1Ukp1dEZUVnhUWVR1Z3lvWWNXZzAzQWlsMDh3eDhyTHkKUkgvTjNjRlEzaEtLcVZuSHQvdnNZUUhTMnJJYkt0RTluelFPWDRxRDdVYXM3Z0c0L2ZkcmZQZjZFWTR1aGpBcApUeGZRTDUzcTBQZG85T09MZlRReFRxakVNaFpidjV1aEN5d0N2VHdWTmVLZ2MzN1pqdDNPSjI3NTB3U2t1TFZvCnllR0VaQmtML1VONjJjTUJuYlFsSTEzR2FXejBHZ0NJNGkwS3UvRmE4aHJZQTQwcHVFdkEwZFBYcVFGMDhYbFYKM1RJUEhGRWdBWlJpTmpJWmFyQW00THdnL1F4Z203OUR1SVM3VHh6RCtpN1pNSmsydjI1ck14Ly9MMXNNUFBtOQpWaXBwVnpLZmpqRmpwTDVjcVJucC9UdUZSVWpHaDZWMFBXVVk1eTVzYjJBWHpuSGZVd1lqeFNoUjBKWXpXejAwCjNHbklwNnlJa1UvL3dFVGJLcVliMjd0RjdETm1WMUxXQzl0ell1dm4yK2EwQkpnU0Jlc3c4WFJ1WWorQS92bVcKWk1YbkF2anZXR3RBUzA4d0ZOV3F3QUtMbzJYUHBXWGVMa3BZUHo1ZnY2QnJaNVNwYTg4UFhsa1VmOVF0VHRobwprZFlGOWVMdk5hTXpSSWJhbmRGWjdLcHUvN2I3L0tDWE9rMUhMOUxvdEpwY2tJdTAxWS81TnQwOHp5cEVQQ1RzClVGWG5DODNqK2tWMktndG5XcXlEL2k3Z1dwaHJSK0IrNE9tM3VZU1RuY042a2d6ZkV3WldpUVA3ZkpiNlYwTHoKc29yU09sK2g2WDRsMC9oRVdScktVQTBrOXpPZU9TQXhlbmpVUXFReWdUd0RqQTJWbTdSZXI2ZElDMVBwNmVETgpBVEJ0ME1NZjJJTytxbTJtK0VLd1FVSXY4ZXdpdEpab016MFBaOHB6WEM0ZFMyRTErZzZmbnE2UGJ5WWRISDJnCmVraXk4Y2duVVJmdHJFaVoyMUxpMWdpdTJaeVM5QUc0Z1ZuT0E1Y05oSzZtRDJUaGl5UUl2M09yUDA0aDFTNlEKQUdGeGJONEhZK0tCYnVITTYwRG1PQXR5c3o4QkJheHFwWjlXQkVhV01ubFB6eEI2SnFqTGJrZ1BkQ2wycytUWAphcWx0UDd6QkpaenVTeVNQc2tQR1NBREUvaEF4eDJFM1RQeWNhQlhQRVFUM2VkZmNsM09nYXRmeHBSYXJLV09PCnFHM2lteW42ZzJiNjhWTlBDSnBTYTNKZ1Axb0NNVlBpa2RCSEdSVUV3N2dXTlJVOFpXRVJuS292M2c0MnQ4dkEKU2Z0a3VMdkhoUnlPQW91SUVsNjJIems0WC9CeVVOQ2J3MW50RzFQeHpSaERaV2dPaVhPNi94WFByRlpKa3BtcQpZUUE5dW83OVdKZy9zSWxucFJCdFlUbUh4eU9mNk12R2svdXlkZExkcmZ6MHB6QUVmWm11YTVocWh5M2Y4YlNJCmpxMlJwUHE3eHJ1Y2djbFAwTWFjdHkrbm9wa0N4M0lNRUE4NE9MQ3dxZjVtemtwY0U1M3hGaU1hcXZTK0dHZmkKZlZnUGpXTXRzMFhjdEtCV2tUbVFFN3MxSE5EV0g1dlpJaDY2WTZncXR0cjU2VGdtcHRLWHBVdUJ1MEdERFBQbwp1aGI4TnVRRjZwNHNoM1dDbXlzTU9uSW5jaXRxZWE4NTFEMmloK2lIY3VqcnJidkVYZGtjMnlxUHBtK3Q3SXBvCm1zWkxVemdXRlZpNWY3KzZiZU56dGJ3T2tmYmdlQVAyaklHTzdtR1pKWWM0L1d1eXBqeVRKNlBQVC9IMUc3K3QKUTh5R3FDV3BzNFdQM2srR3hrbW90cnFROFcxa0J1RDJxTEdmSTdMMGZUVE9lWk0vQUZ1VDJVSkcxKzQ2czJVVwp2RlF2VUJmZ0dTWlh3c1VUeGJRTlZNaTJib1BCRkNxbUY2VmJTcmw2YVgrSm1NNVhySUlqUUhGUFZWVGxzeUtpClVDUC9PQTJOWlREdW9IcC9EM0s1Qjh5MlIyUTlqZlJ0RkcwL0dnMktCbCtObzdTbXlPcWlsUlNkZ1VJb0p5QkcKRGovZXJ4ZkZNMlc3WTVsNGZ2ZlNpdU1OZmlUTVdkY3cxSStnVkpGMC9mTHRpYkNoUlg0OTlIRWlXUHZkTGFKMwppcDJEYU9ReS9QZG5zK3hvaWlMNWtHV25BVUVwanNjWno0YU5DZFowOXRUb1FhK2RZd3g1R1ovNUtmbnVpTURnClBrWjNXalFpOVlZRWFXbVIvQ2JmMjAyRXdoNjdIZzVqWE5kb0RNendXT0V4RFNkVFFYZVdzUUI0LzNzcjE2S2MKeitGN2xhOXhHVEVhTDllQitwcjY5L2JjekJLMGVkNXUxYUgxcXR3cjcrMmliNmZDdlMyblRGQTM1ZG50YXZlUwp4VUJVZ0NzRzVhTTl4b2pIQ0o4RzRFMm9iRUEwUDg2SFlqZEJJSXF5U0txZWtQYmFybW4xR1JrdUVlbU5hTVdyCkM2bWZqUXR5V2ZMWnlSbUlhL1dkSVgzYXhqZHhYa3kydm4yNVV6MXZRNklrNnRJcktPYUJnRUY1cmYwY014dTUKN1BYeTk0dnc1QjE0Vlcra2JqQnkyY3hIajJhWnJEaE53UnVQNlpIckg5MHZuN2NmYjYwU0twRWxxdmZwdlN0VQpvQnVXQlFEUUE3bHpZajhhT3BHend3LzlYTjI5MGJrUnd4elVZRTBxOVl4bS9VSHJTNUlyRWtKSml2SUlEb3hICjF4VTVLd2ErbERvWDJNcERrZlBQVE9XSjVqZG8wbXNsN0dBTmc1WGhERnBpb2hFMEdSS2lGVytYcjBsYkJKU2oKUkxibytrbzhncXU2WHB0OWU4U0Y5OEJ4bFpEcFBVMG5PcGRrTmxwTVpKYVlpaUUzRjRFRG9DcE56bmxpY2JrcApjZ2FrcGVrbS9YS21RSlJxWElXci8wM29SdUVFTXBxZzlRbjdWRG8zR0FiUTlnNUR5U1Bid0xvT25xQ0V3WGFJCkF6alFzWU4rc3VRd2FqZHFUcEthZ1FCbWRaMmdNZDBTMTV1Ukt6c2wxOHgzK1JabmRiNWoxNjNuV0NkMlQ5VDgKald3NURISDgvVUFkSGZoOHh0RTJ6bWRHbEg5T3I5U2hIMzViMWgxVm8rU2pNMzRPeWpwVjB3TmNVL1psOTBUdAp1WnJwYnBwTXZCZUVmRzZTczVXVGhySm9LaGl0RkNwWlVqaDZvdnk3Mzd6ditKaUc4aDRBNG1GTmRPSUtBd0I0Cmp2Nms3V3poUVlEa2Q0ZXRoajNndVJCTGZQNThNVEJKaWhZemVINkUzclhjSGE5b0xnREgzczd4bU8yVEtUY24Kd3VIM3AvdC9WWFN3UGJ0QXBXUXdTRFNKSnA5WkF4S0Q1eVdmd3lTU2ZQVGtwM2c1b2NmKzBhSk1Kc2FkU3lwNQpNR1Vic1oxd1hTN2RXMDhOYXZ2WmpmbElNUm8wUFZDbkRVcFp1bjJuekhTRGJDSjB1M0ZYd1lFQzFFejlJUnN0ClJFbDdpdTZQRlVMSldSU0V0SzBKY1lLS0ltNXhQWHIvbTdPc2duMUNJL0F0cTkrWEFjODk1MGVxeTRwTFVQYkYKZkhFOFhVYWFzUU82MDJTeGpnOTZZaWJ3ZnFyTDF2Vjd1MitUYzJleUZ1N3oxUGRPZDQyWko5M2wvM3lOUW92egora0JuQVdObzZ3WnNKSitHNDZDODNYRVBLM0h1bGw1dFg2UDU4NUQ1b3o5U1oyZGlTd1FyVFN1THVSL0JCQUpVCmd1K2FITkJGRmVtUXNEL2QxMllud1h3d3FkZXVaMDVmQlFiWUREdldOM3daUjJJeHZpd1E0bjZjZWl3OUZ4QmcKbWlzMFBGY2NZOWl0SnJrYXlWQVVZUFZ3Sm5XSmZEK2pQNjJ3UWZJWmhhbFQrZDJpUzVQaDEwdWlMNHEvY1JuYgo1c1Mvc2o0Tm5QYmpxc1ZmZWlKTEh3PT0KLS0tLS1FTkQgRU5DUllQVEVEIFBSSVZBVEUgS0VZLS0tLS0K",
			"-var=account_ec2_sydney_cert=whatever",
		})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := machines.MachinesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Test\"")
		}
		resource := resources.Items[0]

		if resource.Endpoint.(*machines.SSHEndpoint).Host != "3.25.215.87" {
			t.Fatal("The machine must have a Endpoint.Host of \"3.25.215.87\" (was \"" + resource.Endpoint.(*machines.SSHEndpoint).Host + "\")")
		}

		if resource.Endpoint.(*machines.SSHEndpoint).DotNetCorePlatform != "linux-x64" {
			t.Fatal("The machine must have a Endpoint.DotNetCorePlatform of \"linux-x64\" (was \"" + resource.Endpoint.(*machines.SSHEndpoint).DotNetCorePlatform + "\")")
		}

		return nil
	})
}

// TestListeningTargetResource verifies that a listening machine can be reimported with the correct settings
func TestListeningTargetResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/31-listeningtarget", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := machines.MachinesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Test\"")
		}
		resource := resources.Items[0]

		if resource.URI != "https://tentacle/" {
			t.Fatal("The machine must have a Uri of \"https://tentacle/\" (was \"" + resource.URI + "\")")
		}

		if resource.Thumbprint != "55E05FD1B0F76E60F6DA103988056CE695685FD1" {
			t.Fatal("The machine must have a Thumbprint of \"55E05FD1B0F76E60F6DA103988056CE695685FD1\" (was \"" + resource.Thumbprint + "\")")
		}

		if len(resource.Roles) != 1 {
			t.Fatal("The machine must have 1 role")
		}

		if resource.Roles[0] != "vm" {
			t.Fatal("The machine must have a role of \"vm\" (was \"" + resource.Roles[0] + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		return nil
	})
}

// TestPollingTargetResource verifies that a polling machine can be reimported with the correct settings
func TestPollingTargetResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/32-pollingtarget", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := machines.MachinesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Test\"")
		}
		resource := resources.Items[0]

		if resource.Endpoint.(*machines.PollingTentacleEndpoint).URI.Host != "abcdefghijklmnopqrst" {
			t.Fatal("The machine must have a Uri of \"poll://abcdefghijklmnopqrst/\" (was \"" + resource.Endpoint.(*machines.PollingTentacleEndpoint).URI.Host + "\")")
		}

		if resource.Thumbprint != "1854A302E5D9EAC1CAA3DA1F5249F82C28BB2B86" {
			t.Fatal("The machine must have a Thumbprint of \"1854A302E5D9EAC1CAA3DA1F5249F82C28BB2B86\" (was \"" + resource.Thumbprint + "\")")
		}

		if len(resource.Roles) != 1 {
			t.Fatal("The machine must have 1 role")
		}

		if resource.Roles[0] != "vm" {
			t.Fatal("The machine must have a role of \"vm\" (was \"" + resource.Roles[0] + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		return nil
	})
}

// TestCloudRegionTargetResource verifies that a cloud region can be reimported with the correct settings
func TestCloudRegionTargetResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/33-cloudregiontarget", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := machines.MachinesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Test\"")
		}
		resource := resources.Items[0]

		if len(resource.Roles) != 1 {
			t.Fatal("The machine must have 1 role")
		}

		if resource.Roles[0] != "cloud" {
			t.Fatal("The machine must have a role of \"cloud\" (was \"" + resource.Roles[0] + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		return nil
	})
}

// TestOfflineDropTargetResource verifies that an offline drop can be reimported with the correct settings
func TestOfflineDropTargetResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/34-offlinedroptarget", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := machines.MachinesQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Test\"")
		}
		resource := resources.Items[0]

		if len(resource.Roles) != 1 {
			t.Fatal("The machine must have 1 role")
		}

		if resource.Roles[0] != "offline" {
			t.Fatal("The machine must have a role of \"offline\" (was \"" + resource.Roles[0] + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		if resource.Endpoint.(*machines.OfflinePackageDropEndpoint).ApplicationsDirectory != "c:\\temp" {
			t.Fatal("The machine must have a Endpoint.ApplicationsDirectory of \"c:\\temp\" (was \"" + resource.Endpoint.(*machines.OfflinePackageDropEndpoint).ApplicationsDirectory + "\")")
		}

		if resource.Endpoint.(*machines.OfflinePackageDropEndpoint).WorkingDirectory != "c:\\temp" {
			t.Fatal("The machine must have a Endpoint.OctopusWorkingDirectory of \"c:\\temp\" (was \"" + resource.Endpoint.(*machines.OfflinePackageDropEndpoint).WorkingDirectory + "\")")
		}

		return nil
	})
}

// TestAzureCloudServiceTargetResource verifies that a azure cloud service target can be reimported with the correct settings
func TestAzureCloudServiceTargetResource(t *testing.T) {
	// I could not figure out a combination of properties that made the octopusdeploy_azure_subscription_account resource work
	return

	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/35-azurecloudservicetarget", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := machines.MachinesQuery{
			PartialName: "Azure",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Azure\"")
		}
		resource := resources.Items[0]

		if len(resource.Roles) != 1 {
			t.Fatal("The machine must have 1 role")
		}

		if resource.Roles[0] != "cloud" {
			t.Fatal("The machine must have a role of \"cloud\" (was \"" + resource.Roles[0] + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		if resource.Endpoint.(*machines.AzureCloudServiceEndpoint).CloudServiceName != "servicename" {
			t.Fatal("The machine must have a Endpoint.CloudServiceName of \"c:\\temp\" (was \"" + resource.Endpoint.(*machines.AzureCloudServiceEndpoint).CloudServiceName + "\")")
		}

		if resource.Endpoint.(*machines.AzureCloudServiceEndpoint).StorageAccountName != "accountname" {
			t.Fatal("The machine must have a Endpoint.StorageAccountName of \"accountname\" (was \"" + resource.Endpoint.(*machines.AzureCloudServiceEndpoint).StorageAccountName + "\")")
		}

		if !resource.Endpoint.(*machines.AzureCloudServiceEndpoint).UseCurrentInstanceCount {
			t.Fatal("The machine must have Endpoint.UseCurrentInstanceCount set")
		}

		return nil
	})
}

// TestAzureServiceFabricTargetResource verifies that a service fabric target can be reimported with the correct settings
func TestAzureServiceFabricTargetResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/36-servicefabrictarget", []string{
			"-var=target_service_fabric=whatever",
		})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := machines.MachinesQuery{
			PartialName: "Service Fabric",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Service Fabric\"")
		}
		resource := resources.Items[0]

		if len(resource.Roles) != 1 {
			t.Fatal("The machine must have 1 role")
		}

		if resource.Roles[0] != "cloud" {
			t.Fatal("The machine must have a role of \"cloud\" (was \"" + resource.Roles[0] + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		if resource.Endpoint.(*machines.AzureServiceFabricEndpoint).ConnectionEndpoint != "http://endpoint" {
			t.Fatal("The machine must have a Endpoint.ConnectionEndpoint of \"http://endpoint\" (was \"" + resource.Endpoint.(*machines.AzureServiceFabricEndpoint).ConnectionEndpoint + "\")")
		}

		if resource.Endpoint.(*machines.AzureServiceFabricEndpoint).AadCredentialType != "UserCredential" {
			t.Fatal("The machine must have a Endpoint.AadCredentialType of \"UserCredential\" (was \"" + resource.Endpoint.(*machines.AzureServiceFabricEndpoint).AadCredentialType + "\")")
		}

		if resource.Endpoint.(*machines.AzureServiceFabricEndpoint).AadUserCredentialUsername != "username" {
			t.Fatal("The machine must have a Endpoint.AadUserCredentialUsername of \"username\" (was \"" + resource.Endpoint.(*machines.AzureServiceFabricEndpoint).AadUserCredentialUsername + "\")")
		}

		return nil
	})
}

// TestAzureWebAppTargetResource verifies that a web app target can be reimported with the correct settings
func TestAzureWebAppTargetResource(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/37-webapptarget", []string{
			"-var=account_sales_account=whatever",
		})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := machines.MachinesQuery{
			PartialName: "Web App",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Machines.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a machine called \"Web App\"")
		}
		resource := resources.Items[0]

		if len(resource.Roles) != 1 {
			t.Fatal("The machine must have 1 role")
		}

		if resource.Roles[0] != "cloud" {
			t.Fatal("The machine must have a role of \"cloud\" (was \"" + resource.Roles[0] + "\")")
		}

		if resource.TenantedDeploymentMode != "Untenanted" {
			t.Fatal("The machine must have a TenantedDeploymentParticipation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
		}

		if resource.Endpoint.(*machines.AzureWebAppEndpoint).ResourceGroupName != "mattc-webapp" {
			t.Fatal("The machine must have a Endpoint.ResourceGroupName of \"mattc-webapp\" (was \"" + resource.Endpoint.(*machines.AzureWebAppEndpoint).ResourceGroupName + "\")")
		}

		if resource.Endpoint.(*machines.AzureWebAppEndpoint).WebAppName != "mattc-webapp" {
			t.Fatal("The machine must have a Endpoint.WebAppName of \"mattc-webapp\" (was \"" + resource.Endpoint.(*machines.AzureWebAppEndpoint).WebAppName + "\")")
		}

		if resource.Endpoint.(*machines.AzureWebAppEndpoint).WebAppSlotName != "slot1" {
			t.Fatal("The machine must have a Endpoint.WebAppSlotName of \"slot1\" (was \"" + resource.Endpoint.(*machines.AzureWebAppEndpoint).WebAppSlotName + "\")")
		}

		return nil
	})
}

// TestProjectWithGitUsernameExport verifies that a project can be reimported with the correct git settings
func TestProjectWithGitUsernameExport(t *testing.T) {
	if os.Getenv("GIT_CREDENTIAL") == "" {
		t.Fatal("The GIT_CREDENTIAL environment variable must be set")
	}

	if os.Getenv("GIT_USERNAME") == "" {
		t.Fatal("The GIT_USERNAME environment variable must be set")
	}

	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		_, err := act(t, container, "./test/terraform/39-projectgitusername", []string{
			"-var=project_git_password=" + os.Getenv("GIT_CREDENTIAL"),
			"-var=project_git_username=" + os.Getenv("GIT_USERNAME"),
		})

		if err != nil {
			return err
		}

		// The client does not expose git credentials, so just test the import worked ok

		return nil
	})
}

// TestProjectWithDollarSignsExport verifies that a project can be reimported with terraform string interpolation
func TestProjectWithDollarSignsExport(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/40-escapedollar", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := projects.ProjectsQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Projects.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a project called \"Test\"")
		}

		return nil
	})
}

// TestProjectTerraformInlineScriptExport verifies that a project can be reimported with a terraform inline template step.
// See https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/issues/478
func TestProjectTerraformInlineScriptExport(t *testing.T) {
	arrangeTest(t, func(t *testing.T, container *octopusContainer) error {
		// Act
		newSpaceId, err := act(t, container, "./test/terraform/41-terraforminlinescript", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := createClient(container.URI, newSpaceId)
		query := projects.ProjectsQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Projects.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a project called \"Test\"")
		}
		resource := resources.Items[0]

		deploymentProcess, err := client.DeploymentProcesses.GetByID(resource.DeploymentProcessID)

		if deploymentProcess.Steps[0].Actions[0].Properties["Octopus.Action.Terraform.Template"].Value != "#test" {
			t.Fatalf("The inline Terraform template must be set to \"#test\"")
		}

		return nil
	})
}
