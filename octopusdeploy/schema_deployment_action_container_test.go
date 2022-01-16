package octopusdeploy

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/stretchr/testify/require"
)

func TestExpandDeploymentActionContainer(t *testing.T) {
	expected := octopusdeploy.DeploymentActionContainer{}
	actual := expandContainer(nil)
	require.Equal(t, expected, actual)

	var emptyInterface interface{}
	actual = expandContainer(emptyInterface)
	require.Equal(t, expected, actual)

	var emptyInterfaceArray []interface{}
	actual = expandContainer(emptyInterfaceArray)
	require.Equal(t, expected, actual)

	var testMap = make([]interface{}, 1)
	actual = expandContainer(testMap)
	require.Equal(t, expected, actual)

	testMap[0] = make(map[string]interface{}, 1)
	actual = expandContainer(testMap)
	require.Equal(t, expected, actual)

	testMap[0] = map[string]interface{}{
		"feed_id": "feeds-123",
		"image":   "image-123",
	}
	expected = octopusdeploy.DeploymentActionContainer{
		FeedID: "feeds-123",
		Image:  "image-123",
	}
	actual = expandContainer(testMap)

	require.Equal(t, expected, actual)
}
