package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"
	"reflect"
	"testing"
)

func TestExpandChannelGitResourceRules_WithValidData_ReturnsExpectedResult(t *testing.T) {
	input := map[string]interface{}{
		"id": "rule-1",
		"git_dependency_actions": []interface{}{
			map[string]interface{}{
				"deployment_action_slug": "deploy-action-1",
				"git_dependency_name":    "",
			},
		},
		"rules": []string{"rule1", "rule2"},
	}

	expected := channels.ChannelGitResourceRule{
		Id: "rule-1",
		GitDependencyActions: []gitdependencies.DeploymentActionGitDependency{
			{
				DeploymentActionSlug: "deploy-action-1",
				GitDependencyName:    "",
			},
		},
		Rules: []string{"rule1", "rule2"},
	}

	actual := expandChannelGitResourceRules(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}

func TestExpandChannelGitResourceRules_WithEmptyData_ReturnsEmptyStruct(t *testing.T) {
	input := map[string]interface{}{}

	expected := channels.ChannelGitResourceRule{}

	actual := expandChannelGitResourceRules(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}

func TestFlattenChannelGitResourceRules_WithValidData_ReturnsExpectedMap(t *testing.T) {
	input := []channels.ChannelGitResourceRule{
		{
			Id: "rule-1",
			GitDependencyActions: []gitdependencies.DeploymentActionGitDependency{
				{
					DeploymentActionSlug: "deploy-action-1",
					GitDependencyName:    "ref-1",
				},
			},
			Rules: []string{"rule1", "rule2"},
		},
	}

	expected := []map[string]interface{}{
		{
			"id": "rule-1",
			"git_dependency_actions": []interface{}{
				map[string]interface{}{
					"deployment_action_slug": "deploy-action-1",
					"git_dependency_name":    "ref-1",
				},
			},
			"rules": []string{"rule1", "rule2"},
		},
	}

	actual := flattenChannelGitResourceRules(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}

func TestFlattenChannelGitResourceRules_WithEmptyData_ReturnsEmptySlice(t *testing.T) {
	input := []channels.ChannelGitResourceRule{}

	expected := []map[string]interface{}{}

	actual := flattenChannelGitResourceRules(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %+v, got %+v", expected, actual)
	}
}
