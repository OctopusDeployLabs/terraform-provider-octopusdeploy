package octopusdeploy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/require"
)

func TestExpandRunScriptAction(t *testing.T) {
	runScriptAction := expandRunScriptAction(nil)
	require.Nil(t, runScriptAction)

	runScriptAction = expandRunScriptAction(map[string]interface{}{})
	require.Nil(t, runScriptAction)

	runScriptAction = expandRunScriptAction(map[string]interface{}{
		"name": nil,
	})
	require.Nil(t, runScriptAction)

	runScriptAction = expandRunScriptAction(map[string]interface{}{
		"action_type": nil,
	})
	require.Nil(t, runScriptAction)

	runScriptAction = expandRunScriptAction(map[string]interface{}{
		"name": acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	})
	require.NotNil(t, runScriptAction)
}
