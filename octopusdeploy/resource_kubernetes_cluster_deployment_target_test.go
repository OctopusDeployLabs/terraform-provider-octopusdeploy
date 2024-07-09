package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	stdslices "slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKubernetesClusterDeploymentTargetBasic(t *testing.T) {
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountUsername := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	clusterURL := "https://example.com"
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	newClusterURL := "http://www.example.com"

	resource.Test(t, resource.TestCase{
		CheckDestroy: testDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterDeploymentTargetBasic(accountLocalName, accountName, accountUsername, environmentLocalName, environmentName, userRoleLocalName, userRoleName, localName, name, clusterURL),
			},
			{
				Config: testAccKubernetesClusterDeploymentTargetBasic(accountLocalName, accountName, accountUsername, environmentLocalName, environmentName, userRoleLocalName, userRoleName, localName, name, newClusterURL),
			},
		},
	})
}

func TestAccKubernetesClusterDeploymentTargetAws(t *testing.T) {
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountAccessKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountSecretKey := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	clusterName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	clusterURL := "https://example.com"
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterDeploymentTargetAws(
					accountLocalName,
					accountName,
					accountAccessKey,
					accountSecretKey,
					environmentLocalName,
					environmentName,
					userRoleLocalName,
					userRoleName,
					localName,
					name,
					clusterURL,
					clusterName),
			},
		},
	})
}

func TestAccKubernetesClusterDeploymentTargetGcp(t *testing.T) {
	accountLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	accountUsername := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleLocalName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	userRoleName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	clusterName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	clusterURL := "https://example.com"
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	project := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)
	region := acctest.RandStringFromCharSet(16, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testDeploymentTargetCheckDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterDeploymentTargetGcp(
					accountLocalName,
					accountName,
					accountUsername,
					environmentLocalName,
					environmentName,
					userRoleLocalName,
					userRoleName,
					localName,
					name,
					clusterURL,
					clusterName,
					project,
					region),
			},
		},
	})
}

func testAccKubernetesClusterDeploymentTargetBasic(accountLocalName string, accountName string, accountUsername string, environmentLocalName string, environmentName string, userRoleLocalName string, userRoleName string, localName string, name string, clusterURL string) string {
	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	usernamePasswordAccountID := "${octopusdeploy_username_password_account." + accountLocalName + ".id}"
	environmentID := "${octopusdeploy_environment." + environmentLocalName + ".id}"
	userRoleID := "${octopusdeploy_user_role." + userRoleLocalName + ".id}"

	return fmt.Sprintf(testUsernamePasswordMinimum(accountLocalName, accountName, accountUsername)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testUserRoleMinimum(userRoleLocalName, userRoleName)+"\n"+
		`resource "octopusdeploy_kubernetes_cluster_deployment_target" "%s" {
		   cluster_url  = "%s"
		   environments = ["%s"]
		   name         = "%s"
		   roles        = ["%s"]
		   tenanted_deployment_participation = "Untenanted"

		   authentication {
		     account_id = "%s"
		   }
	     }`, localName, clusterURL, environmentID, name, userRoleID, usernamePasswordAccountID)
}

func testAccKubernetesClusterDeploymentTargetGcp(
	accountLocalName string,
	accountName string,
	accountJSONKey string,
	environmentLocalName string,
	environmentName string,
	userRoleLocalName string,
	userRoleName string,
	localName string,
	name string,
	clusterURL string,
	clusterName string,
	project string,
	region string) string {
	gcpAccountID := "${octopusdeploy_gcp_account." + accountLocalName + ".id}"
	environmentID := "${octopusdeploy_environment." + environmentLocalName + ".id}"
	userRoleID := "${octopusdeploy_user_role." + userRoleLocalName + ".id}"

	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testGcpAccount(accountLocalName, accountName, accountJSONKey)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testUserRoleMinimum(userRoleLocalName, userRoleName)+"\n"+
		`resource "octopusdeploy_kubernetes_cluster_deployment_target" "%s" {
		   cluster_url  = "%s"
		   environments = ["%s"]
		   name         = "%s"
		   roles        = ["%s"]
		   tenanted_deployment_participation = "Untenanted"

		   gcp_account_authentication {
		     account_id = "%s"
			 cluster_name = "%s"
			 project = "%s"
			 region = "%s"
		   }
	     }`, localName, clusterURL, environmentID, name, userRoleID, gcpAccountID, clusterName, project, region)
}

func testAccKubernetesClusterDeploymentTargetAws(
	accountLocalName string,
	accountName string,
	accountAccessKey string,
	accountSecretKey string,
	environmentLocalName string,
	environmentName string,
	userRoleLocalName string,
	userRoleName string,
	localName string,
	name string,
	clusterURL string,
	clusterName string) string {
	awsAccountID := "${octopusdeploy_aws_account." + accountLocalName + ".id}"
	environmentID := "${octopusdeploy_environment." + environmentLocalName + ".id}"
	userRoleID := "${octopusdeploy_user_role." + userRoleLocalName + ".id}"

	allowDynamicInfrastructure := false
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 10)
	useGuidedFailure := false

	return fmt.Sprintf(testAwsAccount(accountLocalName, accountName, accountAccessKey, accountSecretKey)+"\n"+
		testAccEnvironment(environmentLocalName, environmentName, description, allowDynamicInfrastructure, sortOrder, useGuidedFailure)+"\n"+
		testUserRoleMinimum(userRoleLocalName, userRoleName)+"\n"+
		`resource "octopusdeploy_kubernetes_cluster_deployment_target" "%s" {
		   cluster_url  = "%s"
		   environments = ["%s"]
		   name         = "%s"
		   roles        = ["%s"]
		   tenanted_deployment_participation = "Untenanted"

		   aws_account_authentication {
		     account_id = "%s"
			 cluster_name = "%s"
		   }
	     }`, localName, clusterURL, environmentID, name, userRoleID, awsAccountID, clusterName)
}

// TestK8sTargetResource verifies that a k8s machine can be reimported with the correct settings
func TestK8sTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "29-k8starget", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "29a-k8stargetds"), newSpaceId, []string{})

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
		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "29a-k8stargetds"), "data_lookup")

		if err != nil {
			return err
		}

		if lookup != resource.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
		}

		return nil
	})
}

// TestK8sTargetResource verifies that a k8s machine can be reimported with the correct settings
func TestK8sTargetWithCertResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "47-k8stargetwithcert", []string{})

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
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "48-k8stargetpodauth", []string{})

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

func TestKubernetesDeploymentTargetData(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "55-kubernetesagentdeploymenttarget", []string{})

		if err != nil {
			return err
		}

		err = testFramework.TerraformInitAndApply(t, container, filepath.Join("../terraform", "55a-kubernetesagentdeploymenttargetds"), newSpaceId, []string{})

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := machines.MachinesQuery{
			DeploymentTargetTypes: []string{"KubernetesTentacle"},
			PartialName:           "minimum-agent",
			Skip:                  0,
			Take:                  1,
		}

		resources, err := machines.Get(client, newSpaceId, query)
		if err != nil {
			return err
		}

		var foundAgent = resources.Items[0]

		lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "55a-kubernetesagentdeploymenttargetds"), "data_lookup")
		if err != nil {
			return err
		}

		if lookup != foundAgent.ID {
			t.Fatal("The target lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + foundAgent.ID + "\".")
		}

		return nil
	})
}

func TestKubernetesDeploymentTargetResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(t, container, "../terraform", "55-kubernetesagentdeploymenttarget", []string{})

		if err != nil {
			return err
		}

		// Assert
		client, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)
		query := machines.MachinesQuery{
			DeploymentTargetTypes: []string{"KubernetesTentacle"},
			Skip:                  0,
			Take:                  3,
		}

		resources, err := machines.Get(client, newSpaceId, query)
		if err != nil {
			return err
		}

		if len(resources.Items) != 3 {
			t.Fatalf("Space must have three deployment targets with type KubernetesTentacle")
		}

		optionalAgentName := "optional-agent"
		optionalAgentIndex := stdslices.IndexFunc(resources.Items, func(t *machines.DeploymentTarget) bool { return t.Name == optionalAgentName })
		optionalAgentDeploymentTarget := resources.Items[optionalAgentIndex]
		optionalAgentEndpoint := optionalAgentDeploymentTarget.Endpoint.(*machines.KubernetesTentacleEndpoint)

		expectedDefaultNamespace := "kubernetes-namespace"
		if optionalAgentEndpoint.DefaultNamespace != expectedDefaultNamespace {
			t.Fatalf("Expected  \"%s\" to have a default namespace of \"%s\", instead has \"%s\"", optionalAgentName, expectedDefaultNamespace, optionalAgentEndpoint.DefaultNamespace)
		}

		if !optionalAgentDeploymentTarget.IsDisabled {
			t.Fatalf("Expected  \"%s\" to be disabled", optionalAgentName)
		}

		if !optionalAgentEndpoint.UpgradeLocked {
			t.Fatalf("Expected  \"%s\" to have upgrade locked", optionalAgentName)
		}

		tenantedAgentName := "tenanted-agent"
		tenantedAgentIndex := stdslices.IndexFunc(resources.Items, func(t *machines.DeploymentTarget) bool { return t.Name == tenantedAgentName })
		tenantedAgentDeploymentTarget := resources.Items[tenantedAgentIndex]

		if tenantedAgentDeploymentTarget.TenantedDeploymentMode != "Tenanted" {
			t.Fatalf("Expected \"%s\" to be tenanted, but it was \"%s\"", tenantedAgentName, tenantedAgentDeploymentTarget.TenantedDeploymentMode)
		}

		if len(tenantedAgentDeploymentTarget.TenantIDs) != 1 {
			t.Fatalf("Expected \"%s\" to have 1 tenant, but it has %d", tenantedAgentName, len(tenantedAgentDeploymentTarget.TenantIDs))
		}

		if len(tenantedAgentDeploymentTarget.TenantTags) != 2 {
			t.Fatalf("Expected \"%s\" to have 2 tenant tags, but it has %d", tenantedAgentName, len(tenantedAgentDeploymentTarget.TenantTags))
		}

		return nil
	})
}
