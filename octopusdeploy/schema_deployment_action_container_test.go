package octopusdeploy

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/stretchr/testify/require"
)

func TestExpandDeploymentActionContainer(t *testing.T) {
	actual := expandContainer(nil)
	require.Nil(t, actual)

	var emptyInterface interface{}
	actual = expandContainer(emptyInterface)
	require.Nil(t, actual)

	var emptyInterfaceArray []interface{}
	actual = expandContainer(emptyInterfaceArray)
	require.Nil(t, actual)

	var testMap = make([]interface{}, 1)
	actual = expandContainer(testMap)
	require.Nil(t, actual)

	testMap[0] = make(map[string]interface{}, 1)
	actual = expandContainer(testMap)
	require.Nil(t, actual)

	testMap[0] = map[string]interface{}{
		"feed_id": "feeds-123",
		"image":   "image-123",
	}
	expected := &deployments.DeploymentActionContainer{
		FeedID: "feeds-123",
		Image:  "image-123",
	}
	actual = expandContainer(testMap)

	require.Equal(t, expected, actual)
}
