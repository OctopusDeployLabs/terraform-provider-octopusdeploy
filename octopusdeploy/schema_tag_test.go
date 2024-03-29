package octopusdeploy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestExpandTag(t *testing.T) {
	canonicalTagName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	color := "#FF0000"
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	id := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	sortOrder := acctest.RandIntRange(0, 1000)

	resourceDataMap := map[string]interface{}{
		"canonical_tag_name": canonicalTagName,
		"color":              color,
		"description":        description,
		"id":                 id,
		"name":               name,
		"sort_order":         sortOrder,
	}

	d := schema.TestResourceDataRaw(t, getTagSchema(), resourceDataMap)
	tag := expandTag(d)

	require.Equal(t, tag.CanonicalTagName, canonicalTagName)
	require.Equal(t, tag.Color, color)
	require.Equal(t, tag.Description, description)
	require.Equal(t, tag.Name, name)
	require.Equal(t, tag.SortOrder, sortOrder)
}
