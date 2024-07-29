package octopusdeploy

import (
	"reflect"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/require"
)

func TestFlattenTenant(t *testing.T) {
	clonedFromTenantID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	id := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectEnvironments := map[string][]string{}
	spaceID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantTags := []string{acctest.RandStringFromCharSet(20, acctest.CharSetAlpha), acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)}

	expectedExpanded := tenants.NewTenant(name)
	expectedExpanded.ClonedFromTenantID = clonedFromTenantID
	expectedExpanded.Description = description
	expectedExpanded.ID = id
	expectedExpanded.ProjectEnvironments = projectEnvironments
	expectedExpanded.SpaceID = spaceID
	expectedExpanded.TenantTags = tenantTags

	expectedFlattened := map[string]interface{}{
		"cloned_from_tenant_id": clonedFromTenantID,
		"description":           description,
		"id":                    id,
		"name":                  name,
		"space_id":              spaceID,
		"tenant_tags":           tenantTags,
	}

	actualFlattened := flattenTenant(expectedExpanded)
	require.True(t, reflect.DeepEqual(expectedFlattened, actualFlattened))
}
