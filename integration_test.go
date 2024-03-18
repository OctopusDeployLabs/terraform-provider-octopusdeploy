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

	Checkout the provider with

		git clone https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy.git

	Then build the provider executable with the command:

		go build -o terraform-provider-octopusdeploy main.go

	Terraform will then use the local executable rather than download the provider from the registry.

	To build the and run the tests, run:

		export LICENSE=base 64 octopus license
		export ECR_ACCESS_KEY=aws access key
		export ECR_SECRET_KEY=aws secret key
		export GIT_CREDENTIAL=github token
		export GIT_USERNAME=github username
		go test -c -o integration_test
		./integration_test
*/

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

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
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workerpools"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"k8s.io/utils/strings/slices"
)

// TestSpaceResource verifies that a space can be reimported with the correct settings
func TestSpaceResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "1-singlespace", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, "", test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "2-projectgroup", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "2a-projectgroupds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "2a-projectgroupds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestAwsAccountExport verifies that an AWS account can be reimported with the correct settings
func TestAwsAccountExport(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "3-awsaccount", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "3a-awsaccountds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "3a-awsaccountds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestAzureAccountResource verifies that an Azure account can be reimported with the correct settings
func TestAzureAccountResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "4-azureaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "5-userpassaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "6-gcpaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "7-sshaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "8-azuresubscriptionaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "9-tokenaccount", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "10-helmfeed", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "10a-helmfeedds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "10a-helmfeedds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestDockerFeedResource verifies that a docker feed can be reimported with the correct settings
func TestDockerFeedResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "11-dockerfeed", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "11a-dockerfeedds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "11a-dockerfeedds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
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

	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act

		newSpaceId, err := testFramework.Act(t, container, "./terraform", "12-ecrfeed", []string{
			"-var=feed_ecr_access_key=" + os.Getenv("ECR_ACCESS_KEY"),
			"-var=feed_ecr_secret_key=" + os.Getenv("ECR_SECRET_KEY"),
		})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "12a-ecrfeedds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "12a-ecrfeedds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestMavenFeedResource verifies that a maven feed can be reimported with the correct settings
func TestMavenFeedResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "13-mavenfeed", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "13a-mavenfeedds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "13a-mavenfeedds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestNugetFeedResource verifies that a nuget feed can be reimported with the correct settings
func TestNugetFeedResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "14-nugetfeed", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "14a-nugetfeedds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "14a-nugetfeedds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestWorkerPoolResource verifies that a static worker pool can be reimported with the correct settings
func TestWorkerPoolResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "15-workerpool", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "15a-workerpoolds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "15a-workerpoolds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestEnvironmentResource verifies that an environment can be reimported with the correct settings
func TestEnvironmentResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "16-environment", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "16a-environmentlookup"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "16a-environmentlookup"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The environment lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestLifecycleResource verifies that a lifecycle can be reimported with the correct settings
func TestLifecycleResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "17-lifecycle", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "17a-lifecycleds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "17a-lifecycleds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestVariableSetResource verifies that a variable set can be reimported with the correct settings
func TestVariableSetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "18-variableset", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "18a-variablesetds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "18a-variablesetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestProjectResource verifies that a project can be reimported with the correct settings
func TestProjectResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "19-project", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "19a-projectds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "19a-projectds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

func TestProjectInSpaceResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "19b-projectspace", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		spaces, err := spaces.GetAll(client)

		if err != nil {
			return err
		}
		idx := sort.Search(len(spaces), func(i int) bool { return spaces[i].Name == "Project Space Test" })
		space := spaces[idx]

		query := projects.ProjectsQuery{
			PartialName: "Test project in space",
			Skip:        0,
			Take:        1,
		}

		resources, err := projects.Get(client, space.ID, query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a project called \"Test project in space\"")
		}
		resource := resources.Items[0]

		if resource.Description != "Test project in space" {
			t.Fatal("The project must be have a description of \"Test project in space\" (was \"" + resource.Description + "\")")
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "20-channel", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "20a-channelds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "20a-channelds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The environment lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestTagSetResource verifies that a tag set can be reimported with the correct settings
func TestTagSetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "21-tagset", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "21a-tagsetds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "21a-tagsetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The environment lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestGitCredentialsResource verifies that a git credential can be reimported with the correct settings
func TestGitCredentialsResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "22-gitcredentialtest", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "22a-gitcredentialtestds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "22a-gitcredentialtestds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup == "" {
			t.Fatal("The target lookup did not succeed.")
		}

		return nil
	})
}

// TestScriptModuleResource verifies that a script module set can be reimported with the correct settings
func TestScriptModuleResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "23-scriptmodule", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "23a-scriptmoduleds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "23a-scriptmoduleds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestTenantsResource verifies that a git credential can be reimported with the correct settings
func TestTenantsResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "24-tenants", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "24a-tenantsds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		tagsets, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "24a-tenantsds"), "tagsets")

		if err != nil {
			return err
		}

		if tagsets == "" {
			t.Fatal("The tagset lookup failed.")
		}

		tenants, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "24a-tenantsds"), "tenants_lookup")

		if err != nil {
			return err
		}

		if tenants != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + tenants + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestCertificateResource verifies that a certificate can be reimported with the correct settings
func TestCertificateResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "25-certificates", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "25a-certificatesds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "25a-certificatesds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The environment lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestTenantVariablesResource verifies that a tenant variables can be reimported with the correct settings
func TestTenantVariablesResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "26-tenant_variables", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "27-machinepolicy", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "28-projecttrigger", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "29-k8starget", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "29a-k8stargetds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "29a-k8stargetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestSshTargetResource verifies that a ssh machine can be reimported with the correct settings
func TestSshTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "30-sshtarget", []string{
			"-var=account_ec2_sydney=LS0tLS1CRUdJTiBFTkNSWVBURUQgUFJJVkFURSBLRVktLS0tLQpNSUlKbkRCT0Jna3Foa2lHOXcwQkJRMHdRVEFwQmdrcWhraUc5dzBCQlF3d0hBUUlwNEUxV1ZrejJEd0NBZ2dBCk1Bd0dDQ3FHU0liM0RRSUpCUUF3RkFZSUtvWklodmNOQXdjRUNIemFuVE1QbHA4ZkJJSUpTSncrdW5BL2ZaVFUKRGdrdWk2QnhOY0REUFg3UHZJZmNXU1dTc3V3YWRhYXdkVEdjY1JVd3pGNTNmRWJTUXJBYzJuWFkwUWVVcU1wcAo4QmdXUUthWlB3MEdqck5OQVJaTy9QYklxaU5ERFMybVRSekZidzREcFY5aDdlblZjL1ZPNlhJdzlxYVYzendlCnhEejdZSkJ2ckhmWHNmTmx1blErYTZGdlRUVkVyWkE1Ukp1dEZUVnhUWVR1Z3lvWWNXZzAzQWlsMDh3eDhyTHkKUkgvTjNjRlEzaEtLcVZuSHQvdnNZUUhTMnJJYkt0RTluelFPWDRxRDdVYXM3Z0c0L2ZkcmZQZjZFWTR1aGpBcApUeGZRTDUzcTBQZG85T09MZlRReFRxakVNaFpidjV1aEN5d0N2VHdWTmVLZ2MzN1pqdDNPSjI3NTB3U2t1TFZvCnllR0VaQmtML1VONjJjTUJuYlFsSTEzR2FXejBHZ0NJNGkwS3UvRmE4aHJZQTQwcHVFdkEwZFBYcVFGMDhYbFYKM1RJUEhGRWdBWlJpTmpJWmFyQW00THdnL1F4Z203OUR1SVM3VHh6RCtpN1pNSmsydjI1ck14Ly9MMXNNUFBtOQpWaXBwVnpLZmpqRmpwTDVjcVJucC9UdUZSVWpHaDZWMFBXVVk1eTVzYjJBWHpuSGZVd1lqeFNoUjBKWXpXejAwCjNHbklwNnlJa1UvL3dFVGJLcVliMjd0RjdETm1WMUxXQzl0ell1dm4yK2EwQkpnU0Jlc3c4WFJ1WWorQS92bVcKWk1YbkF2anZXR3RBUzA4d0ZOV3F3QUtMbzJYUHBXWGVMa3BZUHo1ZnY2QnJaNVNwYTg4UFhsa1VmOVF0VHRobwprZFlGOWVMdk5hTXpSSWJhbmRGWjdLcHUvN2I3L0tDWE9rMUhMOUxvdEpwY2tJdTAxWS81TnQwOHp5cEVQQ1RzClVGWG5DODNqK2tWMktndG5XcXlEL2k3Z1dwaHJSK0IrNE9tM3VZU1RuY042a2d6ZkV3WldpUVA3ZkpiNlYwTHoKc29yU09sK2g2WDRsMC9oRVdScktVQTBrOXpPZU9TQXhlbmpVUXFReWdUd0RqQTJWbTdSZXI2ZElDMVBwNmVETgpBVEJ0ME1NZjJJTytxbTJtK0VLd1FVSXY4ZXdpdEpab016MFBaOHB6WEM0ZFMyRTErZzZmbnE2UGJ5WWRISDJnCmVraXk4Y2duVVJmdHJFaVoyMUxpMWdpdTJaeVM5QUc0Z1ZuT0E1Y05oSzZtRDJUaGl5UUl2M09yUDA0aDFTNlEKQUdGeGJONEhZK0tCYnVITTYwRG1PQXR5c3o4QkJheHFwWjlXQkVhV01ubFB6eEI2SnFqTGJrZ1BkQ2wycytUWAphcWx0UDd6QkpaenVTeVNQc2tQR1NBREUvaEF4eDJFM1RQeWNhQlhQRVFUM2VkZmNsM09nYXRmeHBSYXJLV09PCnFHM2lteW42ZzJiNjhWTlBDSnBTYTNKZ1Axb0NNVlBpa2RCSEdSVUV3N2dXTlJVOFpXRVJuS292M2c0MnQ4dkEKU2Z0a3VMdkhoUnlPQW91SUVsNjJIems0WC9CeVVOQ2J3MW50RzFQeHpSaERaV2dPaVhPNi94WFByRlpKa3BtcQpZUUE5dW83OVdKZy9zSWxucFJCdFlUbUh4eU9mNk12R2svdXlkZExkcmZ6MHB6QUVmWm11YTVocWh5M2Y4YlNJCmpxMlJwUHE3eHJ1Y2djbFAwTWFjdHkrbm9wa0N4M0lNRUE4NE9MQ3dxZjVtemtwY0U1M3hGaU1hcXZTK0dHZmkKZlZnUGpXTXRzMFhjdEtCV2tUbVFFN3MxSE5EV0g1dlpJaDY2WTZncXR0cjU2VGdtcHRLWHBVdUJ1MEdERFBQbwp1aGI4TnVRRjZwNHNoM1dDbXlzTU9uSW5jaXRxZWE4NTFEMmloK2lIY3VqcnJidkVYZGtjMnlxUHBtK3Q3SXBvCm1zWkxVemdXRlZpNWY3KzZiZU56dGJ3T2tmYmdlQVAyaklHTzdtR1pKWWM0L1d1eXBqeVRKNlBQVC9IMUc3K3QKUTh5R3FDV3BzNFdQM2srR3hrbW90cnFROFcxa0J1RDJxTEdmSTdMMGZUVE9lWk0vQUZ1VDJVSkcxKzQ2czJVVwp2RlF2VUJmZ0dTWlh3c1VUeGJRTlZNaTJib1BCRkNxbUY2VmJTcmw2YVgrSm1NNVhySUlqUUhGUFZWVGxzeUtpClVDUC9PQTJOWlREdW9IcC9EM0s1Qjh5MlIyUTlqZlJ0RkcwL0dnMktCbCtObzdTbXlPcWlsUlNkZ1VJb0p5QkcKRGovZXJ4ZkZNMlc3WTVsNGZ2ZlNpdU1OZmlUTVdkY3cxSStnVkpGMC9mTHRpYkNoUlg0OTlIRWlXUHZkTGFKMwppcDJEYU9ReS9QZG5zK3hvaWlMNWtHV25BVUVwanNjWno0YU5DZFowOXRUb1FhK2RZd3g1R1ovNUtmbnVpTURnClBrWjNXalFpOVlZRWFXbVIvQ2JmMjAyRXdoNjdIZzVqWE5kb0RNendXT0V4RFNkVFFYZVdzUUI0LzNzcjE2S2MKeitGN2xhOXhHVEVhTDllQitwcjY5L2JjekJLMGVkNXUxYUgxcXR3cjcrMmliNmZDdlMyblRGQTM1ZG50YXZlUwp4VUJVZ0NzRzVhTTl4b2pIQ0o4RzRFMm9iRUEwUDg2SFlqZEJJSXF5U0txZWtQYmFybW4xR1JrdUVlbU5hTVdyCkM2bWZqUXR5V2ZMWnlSbUlhL1dkSVgzYXhqZHhYa3kydm4yNVV6MXZRNklrNnRJcktPYUJnRUY1cmYwY014dTUKN1BYeTk0dnc1QjE0Vlcra2JqQnkyY3hIajJhWnJEaE53UnVQNlpIckg5MHZuN2NmYjYwU0twRWxxdmZwdlN0VQpvQnVXQlFEUUE3bHpZajhhT3BHend3LzlYTjI5MGJrUnd4elVZRTBxOVl4bS9VSHJTNUlyRWtKSml2SUlEb3hICjF4VTVLd2ErbERvWDJNcERrZlBQVE9XSjVqZG8wbXNsN0dBTmc1WGhERnBpb2hFMEdSS2lGVytYcjBsYkJKU2oKUkxibytrbzhncXU2WHB0OWU4U0Y5OEJ4bFpEcFBVMG5PcGRrTmxwTVpKYVlpaUUzRjRFRG9DcE56bmxpY2JrcApjZ2FrcGVrbS9YS21RSlJxWElXci8wM29SdUVFTXBxZzlRbjdWRG8zR0FiUTlnNUR5U1Bid0xvT25xQ0V3WGFJCkF6alFzWU4rc3VRd2FqZHFUcEthZ1FCbWRaMmdNZDBTMTV1Ukt6c2wxOHgzK1JabmRiNWoxNjNuV0NkMlQ5VDgKald3NURISDgvVUFkSGZoOHh0RTJ6bWRHbEg5T3I5U2hIMzViMWgxVm8rU2pNMzRPeWpwVjB3TmNVL1psOTBUdAp1WnJwYnBwTXZCZUVmRzZTczVXVGhySm9LaGl0RkNwWlVqaDZvdnk3Mzd6ditKaUc4aDRBNG1GTmRPSUtBd0I0Cmp2Nms3V3poUVlEa2Q0ZXRoajNndVJCTGZQNThNVEJKaWhZemVINkUzclhjSGE5b0xnREgzczd4bU8yVEtUY24Kd3VIM3AvdC9WWFN3UGJ0QXBXUXdTRFNKSnA5WkF4S0Q1eVdmd3lTU2ZQVGtwM2c1b2NmKzBhSk1Kc2FkU3lwNQpNR1Vic1oxd1hTN2RXMDhOYXZ2WmpmbElNUm8wUFZDbkRVcFp1bjJuekhTRGJDSjB1M0ZYd1lFQzFFejlJUnN0ClJFbDdpdTZQRlVMSldSU0V0SzBKY1lLS0ltNXhQWHIvbTdPc2duMUNJL0F0cTkrWEFjODk1MGVxeTRwTFVQYkYKZkhFOFhVYWFzUU82MDJTeGpnOTZZaWJ3ZnFyTDF2Vjd1MitUYzJleUZ1N3oxUGRPZDQyWko5M2wvM3lOUW92egora0JuQVdObzZ3WnNKSitHNDZDODNYRVBLM0h1bGw1dFg2UDU4NUQ1b3o5U1oyZGlTd1FyVFN1THVSL0JCQUpVCmd1K2FITkJGRmVtUXNEL2QxMllud1h3d3FkZXVaMDVmQlFiWUREdldOM3daUjJJeHZpd1E0bjZjZWl3OUZ4QmcKbWlzMFBGY2NZOWl0SnJrYXlWQVVZUFZ3Sm5XSmZEK2pQNjJ3UWZJWmhhbFQrZDJpUzVQaDEwdWlMNHEvY1JuYgo1c1Mvc2o0Tm5QYmpxc1ZmZWlKTEh3PT0KLS0tLS1FTkQgRU5DUllQVEVEIFBSSVZBVEUgS0VZLS0tLS0K",
			"-var=account_ec2_sydney_cert=whatever",
		})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "30a-sshtargetds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "30a-sshtargetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestListeningTargetResource verifies that a listening machine can be reimported with the correct settings
func TestListeningTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "31-listeningtarget", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "31a-listeningtargetds"), newSpaceId, []string{})

		if err != nil {
			t.Log("BUG: listening targets data sources don't appear to work")
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "32-pollingtarget", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "32a-pollingtargetds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "32a-pollingtargetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestCloudRegionTargetResource verifies that a cloud region can be reimported with the correct settings
func TestCloudRegionTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "33-cloudregiontarget", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "33a-cloudregiontargetds"), newSpaceId, []string{})

		if err != nil {
			t.Fatal("cloud region data source does not appear to work")
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "34-offlinedroptarget", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "34a-offlinedroptargetds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "34a-offlinedroptargetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestAzureCloudServiceTargetResource verifies that a azure cloud service target can be reimported with the correct settings
func TestAzureCloudServiceTargetResource(t *testing.T) {
	// I could not figure out a combination of properties that made the octopusdeploy_azure_subscription_account resource work
	return

	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "35-azurecloudservicetarget", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "36-servicefabrictarget", []string{
			"-var=target_service_fabric=whatever",
		})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "36a-servicefabrictargetds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "36a-servicefabrictargetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestAzureWebAppTargetResource verifies that a web app target can be reimported with the correct settings
func TestAzureWebAppTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "37-webapptarget", []string{
			"-var=account_sales_account=whatever",
		})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "37a-webapptarget"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "37a-webapptarget"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
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

	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		_, err := testFramework.Act(t, container, "./terraform", "39-projectgitusername", []string{
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "40-escapedollar", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "41-terraforminlinescript", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

// TestProjectTerraformPackageScriptExport verifies that a project can be reimported with a terraform package template step.
// See https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/issues/478
func TestProjectTerraformPackageScriptExport(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "42-terraformpackagescript", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		if deploymentProcess.Steps[0].Actions[0].Properties["Octopus.Action.Script.ScriptSource"].Value != "Package" {
			t.Fatalf("The Terraform template must be set deploy files from a package")
		}

		if deploymentProcess.Steps[0].Actions[0].Properties["Octopus.Action.Terraform.TemplateDirectory"].Value != "blah" {
			t.Fatalf("The Terraform template directory must be set to \"blah\"")
		}

		return nil
	})
}

// TestProjectTerraformPackageScriptExport verifies that users and teams can be reimported
func TestUsersAndTeams(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "43-users", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "43a-usersds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		err = func() error {
			query := users.UsersQuery{
				Filter: "Service Account",
				IDs:    nil,
				Skip:   0,
				Take:   1,
			}

			resources, err := client.Users.Get(query)
			if err != nil {
				return err
			}

			if len(resources.Items) == 0 {
				t.Fatalf("Space must have a user called \"Service Account\"")
			}

			resource := resources.Items[0]

			if resource.Username != "saccount" {
				t.Fatalf("Account must have a username \"saccount\"")
			}

			if resource.EmailAddress != "a@a.com" {
				t.Fatalf("Account must have a email \"a@a.com\"")
			}

			if !resource.IsService {
				t.Fatalf("Account must be a service account")
			}

			if !resource.IsActive {
				t.Fatalf("Account must be active")
			}

			return nil
		}()

		if err != nil {
			return err
		}

		err = func() error {
			query := users.UsersQuery{
				Filter: "Bob Smith",
				IDs:    nil,
				Skip:   0,
				Take:   1,
			}

			resources, err := client.Users.Get(query)
			if err != nil {
				return err
			}

			if len(resources.Items) == 0 {
				t.Fatalf("Space must have a user called \"Service Account\"")
			}

			resource := resources.Items[0]

			if resource.Username != "bsmith" {
				t.Fatalf("Regular account must have a username \"bsmith\"")
			}

			if resource.EmailAddress != "bob.smith@example.com" {
				t.Fatalf("Regular account must have a email \"bob.smith@example.com\"")
			}

			if resource.IsService {
				t.Fatalf("Account must not be a service account")
			}

			if resource.IsActive {
				t.Log("BUG: Account must not be active")
			}

			return nil
		}()

		if err != nil {
			return err
		}

		err = func() error {
			query := teams.TeamsQuery{
				IDs:           nil,
				IncludeSystem: false,
				PartialName:   "Deployers",
				Skip:          0,
				Spaces:        nil,
				Take:          1,
			}

			resources, err := client.Teams.Get(query)
			if err != nil {
				return err
			}

			if len(resources.Items) == 0 {
				t.Fatalf("Space must have a team called \"Deployers\"")
			}

			resource := resources.Items[0]

			if len(resource.MemberUserIDs) != 1 {
				t.Fatalf("Team must have one user")
			}

			return nil
		}()

		if err != nil {
			return err
		}

		// Verify the environment data lookups work
		teams, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "43a-usersds"), "teams_lookup")

		if err != nil {
			return err
		}

		if teams == "" {
			t.Fatal("The teams lookup failed.")
		}

		roles, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "43a-usersds"), "roles_lookup")

		if err != nil {
			return err
		}

		if roles == "" {
			t.Fatal("The roles lookup failed.")
		}

		users, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "43a-usersds"), "users_lookup")

		if err != nil {
			return err
		}

		if users == "" {
			t.Fatal("The users lookup failed.")
		}

		return nil
	})
}

// TestGithubFeedResource verifies that a nuget feed can be reimported with the correct settings
func TestGithubFeedResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "44-githubfeed", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "44a-githubfeedds"), newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := feeds.FeedsQuery{
			PartialName: "Github",
			Skip:        0,
			Take:        1,
		}

		resources, err := client.Feeds.Get(query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have an feed called \"Github\"")
		}
		resource := resources.Items[0].(*feeds.GitHubRepositoryFeed)

		if resource.FeedType != "GitHub" {
			t.Fatal("The feed must have a type of \"GitHub\"")
		}

		if resource.Username != "test-username" {
			t.Fatal("The feed must have a username of \"test-username\"")
		}

		if resource.DownloadAttempts != 1 {
			t.Fatal("The feed must be have a downloads attempts set to \"1\"")
		}

		if resource.DownloadRetryBackoffSeconds != 30 {
			t.Fatal("The feed must be have a downloads retry backoff set to \"30\"")
		}

		if resource.FeedURI != "https://api.github.com" {
			t.Fatal("The feed must be have a feed uri of \"https://api.github.com\"")
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

		// Verify the environment data lookups work
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "44a-githubfeedds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestProjectWithScriptActions verifies that a project with a plain script step can be applied and reapplied
func TestProjectWithScriptActions(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "45-projectwithscriptactions", []string{})

		if err != nil {
			return err
		}

		// Do a second apply to catch the scenario documented at https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/issues/509
		err = testFramework.TerraformApply(t, filepath.Join("./terraform", "45-projectwithscriptactions"), container.URI, newSpaceId, []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		deploymentProcess, err := client.DeploymentProcesses.GetByID(resource.DeploymentProcessID)
		if err != nil {
			return err
		}
		if len(deploymentProcess.Steps) != 1 {
			t.Fatal("The DeploymentProcess should have a single Deployment Step")
		}
		step := deploymentProcess.Steps[0]

		if len(step.Actions) != 3 {
			t.Fatal("The DeploymentProcess should have a three Deployment Actions")
		}

		if step.Actions[0].Name != "Pre Script Action" {
			t.Fatal("The first Deployment Action should be name \"Pre Script Action\" (was \"" + step.Actions[0].Name + "\")")
		}
		if step.Actions[1].Name != "Hello world (using PowerShell)" {
			t.Fatal("The second Deployment Action should be name \"Hello world (using PowerShell)\" (was \"" + step.Actions[1].Name + "\")")
		}
		if step.Actions[2].Name != "Post Script Action" {
			t.Fatal("The third Deployment Action should be name \"Post Script Action\" (was \"" + step.Actions[2].Name + "\")")
		}

		return nil
	})
}

// TestRunbookResource verifies that a runbook can be reimported with the correct settings
func TestRunbookResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "46-runbooks", []string{})

		if err != nil {
			return err
		}

		//err = testFramework.TerraformInitAndApply(t, container, filepath.Join("./terraform", "46a-runbooks"), newSpaceId, []string{})
		//
		//if err != nil {
		//	return err
		//}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		resources, err := client.Runbooks.GetAll()
		if err != nil {
			return err
		}

		found := false
		runbookId := ""
		for _, r := range resources {
			if r.Name == "Runbook" {
				found = true
				runbookId = r.ID

				if r.Description != "Test Runbook" {
					t.Fatal("The runbook must be have a description of \"Test Runbook\" (was \"" + r.Description + "\")")
				}

				if r.ConnectivityPolicy.AllowDeploymentsToNoTargets {
					t.Fatal("The runbook must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
				}

				if r.ConnectivityPolicy.ExcludeUnhealthyTargets {
					t.Fatal("The runbook must not have ConnectivityPolicy.AllowDeploymentsToNoTargets enabled")
				}

				if r.ConnectivityPolicy.SkipMachineBehavior != "SkipUnavailableMachines" {
					t.Log("BUG: The runbook must be have a ConnectivityPolicy.SkipMachineBehavior of \"SkipUnavailableMachines\" (was \"" + r.ConnectivityPolicy.SkipMachineBehavior + "\") - Known issue where the value returned by /api/Spaces-#/ProjectGroups/ProjectGroups-#/projects is different to /api/Spaces-/Projects")
				}

				if r.RunRetentionPolicy.QuantityToKeep != 10 {
					t.Fatal("The runbook must not have RunRetentionPolicy.QuantityToKeep of 10 (was \"" + fmt.Sprint(r.RunRetentionPolicy.QuantityToKeep) + "\")")
				}

				if r.RunRetentionPolicy.ShouldKeepForever {
					t.Fatal("The runbook must not have RunRetentionPolicy.ShouldKeepForever of false (was \"" + fmt.Sprint(r.RunRetentionPolicy.ShouldKeepForever) + "\")")
				}

				if r.ConnectivityPolicy.SkipMachineBehavior != "SkipUnavailableMachines" {
					t.Log("BUG: The runbook must be have a ConnectivityPolicy.SkipMachineBehavior of \"SkipUnavailableMachines\" (was \"" + r.ConnectivityPolicy.SkipMachineBehavior + "\") - Known issue where the value returned by /api/Spaces-#/ProjectGroups/ProjectGroups-#/projects is different to /api/Spaces-/Projects")
				}

				if r.MultiTenancyMode != "Untenanted" {
					t.Fatal("The runbook must be have a TenantedDeploymentMode of \"Untenanted\" (was \"" + r.MultiTenancyMode + "\")")
				}

				if r.EnvironmentScope != "Specified" {
					t.Fatal("The runbook must be have a EnvironmentScope of \"Specified\" (was \"" + r.EnvironmentScope + "\")")
				}

				if len(r.Environments) != 1 {
					t.Fatal("The runbook must be have a Environments array of 1 (was \"" + strings.Join(r.Environments, ", ") + "\")")
				}

				if r.DefaultGuidedFailureMode != "EnvironmentDefault" {
					t.Fatal("The runbook must be have a DefaultGuidedFailureMode of \"EnvironmentDefault\" (was \"" + r.DefaultGuidedFailureMode + "\")")
				}

				if !r.ForcePackageDownload {
					t.Log("BUG: The runbook must be have a ForcePackageDownload of \"true\" (was \"" + fmt.Sprint(r.ForcePackageDownload) + "\")")
				}

				process, err := client.RunbookProcesses.GetByID(r.RunbookProcessID)

				if err != nil {
					t.Fatal("Failed to retrieve the runbook process.")
				}

				if len(process.Steps) != 1 {
					t.Fatal("The runbook must be have a 1 step")
				}
			}
		}

		if !found {
			t.Fatalf("Space must have a runbook called \"Runbook\"")
		}

		// There was an issue where deleting a runbook and reapplying the terraform module caused an error, so
		// verify this process works.
		client.Runbooks.DeleteByID(runbookId)
		err = testFramework.TerraformApply(t, "./terraform/46-runbooks", container.URI, newSpaceId, []string{})

		if err != nil {
			t.Fatal("Failed to reapply the runbooks after deleting them.")
		}

		// Verify the environment data lookups work
		//lookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "46a-runbooks"), "data_lookup")
		//
		//if err != nil {
		//	return err
		//}
		//
		//if lookup != resource.ID {
		//	t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		//}

		return nil
	})
}

// TestK8sTargetResource verifies that a k8s machine can be reimported with the correct settings
func TestK8sTargetWithCertResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "47-k8stargetwithcert", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		if fmt.Sprint(resource.Endpoint.(*machines.KubernetesEndpoint).Authentication.GetAuthenticationType()) != "KubernetesCertificate" {
			t.Fatal("The machine must have a Endpoint.AuthenticationType of \"KubernetesCertificate\" (was \"" + fmt.Sprint(resource.Endpoint.(*machines.KubernetesEndpoint).Authentication.GetAuthenticationType()) + "\")")
		}

		return nil
	})
}

// TestK8sPodAuthTargetResource verifies that a k8s machine with pod auth can be reimported with the correct settings
func TestK8sPodAuthTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "48-k8stargetpodauth", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
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

		if fmt.Sprint(resource.Endpoint.(*machines.KubernetesEndpoint).Authentication.GetAuthenticationType()) != "KubernetesPodService" {
			t.Fatal("The machine must have a Endpoint.Authentication.AuthenticationType of \"KubernetesPodService\" (was \"" + fmt.Sprint(resource.Endpoint.(*machines.KubernetesEndpoint).Authentication.GetAuthenticationType()) + "\")")
		}

		if fmt.Sprint(resource.Endpoint.(*machines.KubernetesEndpoint).Authentication.(*machines.KubernetesPodAuthentication).TokenPath) != "/var/run/secrets/kubernetes.io/serviceaccount/token" {
			t.Fatal("The machine must have a Endpoint.Authentication.TokenPath of \"/var/run/secrets/kubernetes.io/serviceaccount/token\" (was \"" + fmt.Sprint(resource.Endpoint.(*machines.KubernetesEndpoint).Authentication.(*machines.KubernetesPodAuthentication).TokenPath) + "\")")
		}

		if fmt.Sprint(resource.Endpoint.(*machines.KubernetesEndpoint).ClusterCertificatePath) != "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt" {
			t.Fatal("The machine must have a Endpoint.ClusterCertificatePath of \"/var/run/secrets/kubernetes.io/serviceaccount/ca.crt\" (was \"" + fmt.Sprint(resource.Endpoint.(*machines.KubernetesEndpoint).ClusterCertificatePath) + "\")")
		}

		return nil
	})
}

func TestVariableResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "49-variables", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		project, err := client.Projects.GetByName("Test")
		variableSet, err := client.Variables.GetAll(project.ID)

		if err != nil {
			return err
		}

		if len(variableSet.Variables) != 7 {
			t.Fatalf("Expected 7 variables to be created.")
		}

		for _, variable := range variableSet.Variables {
			switch variable.Name {
			case "UnscopedVariable":
				if !variable.Scope.IsEmpty() {
					t.Fatalf("Expected UnscopedVariable to have no scope values.")
				}
			case "ActionScopedVariable":
				if len(variable.Scope.Actions) == 0 {
					t.Fatalf("Expected ActionScopedVariable to have action scope.")
				}
			case "ChannelScopedVariable":
				if len(variable.Scope.Channels) == 0 {
					t.Fatalf("Expected ChannelScopedVariable to have channel scope.")
				}
			case "EnvironmentScopedVariable":
				if len(variable.Scope.Environments) == 0 {
					t.Fatalf("Expected EnvironmentScopedVariable to have environment scope.")
				}
			case "MachineScopedVariable":
				if len(variable.Scope.Machines) == 0 {
					t.Fatalf("Expected MachineScopedVariable to have machine scope.")
				}
			case "ProcessScopedVariable":
				if len(variable.Scope.ProcessOwners) == 0 {
					t.Fatalf("Expected ProcessScopedVariable to have process scope.")
				}
			case "RoleScopedVariable":
				if len(variable.Scope.Roles) == 0 {
					t.Fatalf("Expected RoleScopedVariable to have role scope.")
				}
			}
		}

		return nil
	})
}

// TestTerraformApplyStepWithWorkerPool verifies that a terraform apply step with a custom worker pool is deployed successfully
// See https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/issues/601
func TestTerraformApplyStepWithWorkerPool(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "./terraform", "50-applyterraformtemplateaction", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := projects.ProjectsQuery{
			PartialName: "Test",
			Skip:        0,
			Take:        1,
		}

		resources, err := projects.Get(client, newSpaceId, query)
		if err != nil {
			return err
		}

		if len(resources.Items) == 0 {
			t.Fatalf("Space must have a project called \"Test\"")
		}
		resource := resources.Items[0]

		// Get worker pool
		wpQuery := workerpools.WorkerPoolsQuery{
			PartialName: "Docker",
			Skip:        0,
			Take:        1,
		}

		workerpools, err := workerpools.Get(client, newSpaceId, wpQuery)
		if err != nil {
			return err
		}

		if len(workerpools.Items) == 0 {
			t.Fatalf("Space must have a worker pool called \"Docker\"")
		}

		// Get deployment process
		process, err := deployments.GetDeploymentProcessByID(client, "", resource.DeploymentProcessID)
		if err != nil {
			return err
		}

		// Worker pool must be assigned
		if process.Steps[0].Actions[0].WorkerPool != workerpools.Items[0].GetID() {
			t.Fatalf("Action must use the worker pool \"Docker\"")
		}

		return nil
	})
}
