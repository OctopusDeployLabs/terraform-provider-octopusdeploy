package octopusdeploy

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestFlattenAndExpandTenant(t *testing.T) {
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	tenant := octopusdeploy.NewTenant(name)
	tenant.ID = "Tenants-123"
	tenant.Description = description
	tenant.ClonedFromTenantID = "Tenants-321"
	tenant.ProjectEnvironments["Projects-123"] = []string{"Environments-123"}

	flattenedTenant := flattenTenant(tenant)

	t.Log(flattenedTenant)
}
