package octopusdeploy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/require"
)

func TestTenantConnectionGetID(t *testing.T) {
	tenantID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentID0 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentID1 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentID2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	tenantConnection := TenantConnection{
		TenantID:       tenantID,
		ProjectID:      projectID,
		EnvironmentIDs: []string{environmentID0, environmentID1, environmentID2},
	}

	expectedID := fmt.Sprintf(
		"%s:%s:%s+%s+%s",
		tenantID,
		projectID,
		environmentID0,
		environmentID1,
		environmentID2,
	)

	actualID := tenantConnection.GetID()
	require.True(t, reflect.DeepEqual(expectedID, actualID))
}

func TestExpandTenantConnectionFromTenant(t *testing.T) {
	tenantName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	tenantID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentID0 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentID1 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentID2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	tenant := tenants.NewTenant(tenantName)
	tenant.ID = tenantID
	tenant.ProjectEnvironments = map[string][]string{}
	tenant.ProjectEnvironments[projectID] = []string{environmentID0, environmentID1, environmentID2}

	expectedTenantConnection := TenantConnection{
		TenantID:       tenantID,
		ProjectID:      projectID,
		EnvironmentIDs: []string{environmentID0, environmentID1, environmentID2},
	}

	expectedTenantConnection.ID = expectedTenantConnection.GetID()

	actualTenantConnection, err := expandTenantConnectionFromTenant(tenant, projectID)
	require.NoError(t, err)
	require.Equal(t, expectedTenantConnection.ID, actualTenantConnection.ID)
	require.Equal(t, expectedTenantConnection.TenantID, actualTenantConnection.TenantID)
	require.Equal(t, expectedTenantConnection.ProjectID, actualTenantConnection.ProjectID)
	require.ElementsMatch(t, expectedTenantConnection.EnvironmentIDs, actualTenantConnection.EnvironmentIDs)
}

func TestExpandTenantConnectionFromID(t *testing.T) {
	tenantID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentID0 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentID1 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentID2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	id := fmt.Sprintf(
		"%s:%s:%s+%s+%s",
		tenantID,
		projectID,
		environmentID0,
		environmentID1,
		environmentID2,
	)

	expectedTenantConnection := TenantConnection{
		TenantID:       tenantID,
		ProjectID:      projectID,
		EnvironmentIDs: []string{environmentID0, environmentID1, environmentID2},
	}

	expectedTenantConnection.ID = id

	actualTenantConnection, err := expandTenantConnectionFromID(id)
	require.NoError(t, err)
	require.Equal(t, expectedTenantConnection.ID, actualTenantConnection.ID)
	require.Equal(t, expectedTenantConnection.TenantID, actualTenantConnection.TenantID)
	require.Equal(t, expectedTenantConnection.ProjectID, actualTenantConnection.ProjectID)
	require.ElementsMatch(t, expectedTenantConnection.EnvironmentIDs, actualTenantConnection.EnvironmentIDs)
}

func TestFlattenTenantConnection(t *testing.T) {
	tenantID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	projectID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	environmentIDs := []string{acctest.RandStringFromCharSet(20, acctest.CharSetAlpha), acctest.RandStringFromCharSet(20, acctest.CharSetAlpha), acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)}

	expectedExpanded := &TenantConnection{}
	expectedExpanded.TenantID = tenantID
	expectedExpanded.ProjectID = projectID
	expectedExpanded.EnvironmentIDs = environmentIDs
	expectedExpanded.ID = expectedExpanded.GetID()

	expectedFlattened := map[string]interface{}{
		"id":              fmt.Sprintf("%s:%s:%s", tenantID, projectID, strings.Join(environmentIDs, "+")),
		"tenant_id":       tenantID,
		"project_id":      projectID,
		"environment_ids": environmentIDs,
	}

	actualFlattened := flattenTenantConnection(expectedExpanded)
	require.True(t, reflect.DeepEqual(expectedFlattened, actualFlattened))
}
