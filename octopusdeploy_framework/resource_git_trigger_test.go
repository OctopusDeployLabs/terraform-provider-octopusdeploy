package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type gitTriggerSourcesTestData struct {
	deploymentActionSlug string
	gitDependencyName    string
	includeFilePaths     []string
	excludeFilePaths     []string
}

type gitTriggerTestData struct {
	name        string
	description string
	projectId   string
	spaceId     string
	channelId   string
	isDisabled  bool
	sources     []gitTriggerSourcesTestData
}

func TestAccOctopusDeployGitTrigger(t *testing.T) {
	projectId, spaceId, actionSlug, channelId := setupTestSpace(t)

	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_git_trigger." + localName
	createData := gitTriggerTestData{
		name:        acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		description: acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		projectId:   projectId,
		channelId:   channelId,
		spaceId:     spaceId,
		isDisabled:  false,
		sources: []gitTriggerSourcesTestData{
			{
				deploymentActionSlug: actionSlug,
				gitDependencyName:    "",
				includeFilePaths:     []string{acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)},
				excludeFilePaths:     []string{acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)},
			},
		},
	}
	//updateData := gitTriggerTestData{
	//	name:        createData.name + "-updated",
	//	description: createData.description + "-updated",
	//	projectId:   createData.projectId,
	//	channelId:   createData.channelId,
	//	spaceId:     createData.spaceId,
	//	isDisabled:  true,
	//	sources: []gitTriggerSourcesTestData{
	//		{
	//			deploymentActionSlug: createData.sources[0].deploymentActionSlug + "-updated",
	//			gitDependencyName:    createData.sources[0].gitDependencyName + "-updated",
	//			includeFilePaths:     []string{createData.sources[0].includeFilePaths[0] + "-updated"},
	//			excludeFilePaths:     []string{createData.sources[0].excludeFilePaths[0] + "-updated"},
	//		},
	//	},
	//}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             func(s *terraform.State) error { return testGitTriggerCheckDestroy(s) },
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testGitTriggerBasic(createData, localName),
				Check:  testAssertGitTriggerAttributes(createData, prefix),
			},
			//{
			//	Config: testGitTriggerBasic(updateData, localName),
			//	Check:  testAssertGitTriggerAttributes(updateData, prefix),
			//},
		},
	})
}

func setupTestSpace(t *testing.T) (string, string, string, string) {
	testFramework := test.OctopusContainerTest{}
	//newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "45a-projectwithgitdependency", []string{})
	err := testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "45a-projectwithgitdependency"), "Spaces-1", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	newSpaceId := "Spaces-1"

	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)

	query := projects.ProjectsQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	projectResources, err := client.Projects.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(projectResources.Items) == 0 {
		t.Fatalf("Space must have a project called \"Test\"")
	}
	project := projectResources.Items[0]

	channelResources, err := channels.GetAll(client, newSpaceId)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(channelResources) == 0 {
		t.Fatalf("Space must have a channel")
	}

	var projectChannel *channels.Channel

	for _, channel := range channelResources {
		if channel.ProjectID == project.ID {
			projectChannel = channel
			break
		}
	}

	if projectChannel == nil {
		t.Fatalf("No channel found for project ID: %s", project.ID)
	}

	deploymentProccessResource, err := deployments.GetDeploymentProcessByID(client, newSpaceId, project.DeploymentProcessID)

	if err != nil {
		t.Fatal(err.Error())
	}

	actionSlug := deploymentProccessResource.Steps[0].Actions[0].Slug

	return project.ID, newSpaceId, actionSlug, projectChannel.ID
}

func testGitTriggerBasic(data gitTriggerTestData, localName string) string {
	return fmt.Sprintf(`
		resource "octopusdeploy_git_trigger" "%s" {
			name        = "%s"
  			space_id    = "%s"
		    project_id  = "%s"
		  	channel_id  = "%s"
			is_disabled = "%t"
		  	sources 	= [%s]
		}
	`,
		localName,
		data.name,
		data.spaceId,
		data.projectId,
		data.channelId,
		data.isDisabled,
		convertGitTriggerSourcesToString(data.sources),
	)
}

func testAssertGitTriggerAttributes(expected gitTriggerTestData, prefix string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(prefix, "name", expected.name),
		resource.TestCheckResourceAttr(prefix, "space_id", expected.spaceId),
		resource.TestCheckResourceAttr(prefix, "project_id", expected.projectId),
		resource.TestCheckResourceAttr(prefix, "channel_id", expected.channelId),
		resource.TestCheckResourceAttr(prefix, "is_disabled", strconv.FormatBool(expected.isDisabled)),
		//resource.TestCheckResourceAttr(prefix, "sources.0.deploymentActionSlug", expected.sources[0].deploymentActionSlug),
	)
}

func testGitTriggerCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_git_trigger" {
			continue
		}

		projectTrigger, _ := octoClient.ProjectTriggers.GetByID(rs.Primary.ID)
		if projectTrigger == nil {
			return fmt.Errorf("git trigger (%s) NOT FOUND", rs.Primary.ID)

		}
		//if err == nil && projectTrigger != nil {
		//	return fmt.Errorf("git trigger (%s) still exists", rs.Primary.ID)
		//}
	}

	return nil
}

func convertGitTriggerSourcesToString(sources []gitTriggerSourcesTestData) string {
	var result string
	for _, source := range sources {
		result += fmt.Sprintf(`
		{
			deployment_action_slug = "%s"
			git_dependency_name    = "%s"
			include_file_paths     = [%s]
			exclude_file_paths     = [%s]
		}`,
			source.deploymentActionSlug,
			source.gitDependencyName,
			convertStringSliceToString(source.includeFilePaths),
			convertStringSliceToString(source.excludeFilePaths),
		)
	}
	return result
}

func convertGitTriggerSourceToString(source gitTriggerSourcesTestData) string {
	var result string

	result += fmt.Sprintf(`
		{
			deployment_action_slug = "%s"
			git_dependency_name    = "%s"
			include_file_paths     = [%s]
			exclude_file_paths     = [%s]
		}`,
		source.deploymentActionSlug,
		source.gitDependencyName,
		convertStringSliceToString(source.includeFilePaths),
		convertStringSliceToString(source.excludeFilePaths),
	)

	return result
}

func convertStringSliceToString(slice []string) string {
	return fmt.Sprintf(`"%s"`, strings.Join(slice, `", "`))
}
