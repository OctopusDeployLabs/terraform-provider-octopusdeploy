package schemas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var VariableScopeFieldNames = struct {
	Actions      string
	Channels     string
	Environments string
	Machines     string
	Processes    string
	Roles        string
	TenantTags   string
}{
	Actions:      "actions",
	Channels:     "channels",
	Environments: "environments",
	Machines:     "machines",
	Processes:    "processes",
	Roles:        "roles",
	TenantTags:   "tenant_tags",
}

func getVariableScopeResourceSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				VariableScopeFieldNames.Actions:      getVariableScopeItemResourceSchema(VariableScopeFieldNames.Actions),
				VariableScopeFieldNames.Channels:     getVariableScopeItemResourceSchema(VariableScopeFieldNames.Channels),
				VariableScopeFieldNames.Environments: getVariableScopeItemResourceSchema(VariableScopeFieldNames.Environments),
				VariableScopeFieldNames.Machines:     getVariableScopeItemResourceSchema(VariableScopeFieldNames.Machines),
				VariableScopeFieldNames.Processes:    getVariableScopeItemResourceSchema(VariableScopeFieldNames.Processes),
				VariableScopeFieldNames.Roles:        getVariableScopeItemResourceSchema(VariableScopeFieldNames.Roles),
				VariableScopeFieldNames.TenantTags:   getVariableScopeItemResourceSchema(VariableScopeFieldNames.TenantTags),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func getVariableScopeItemResourceSchema(scopeDescription string) resourceSchema.ListAttribute {
	return resourceSchema.ListAttribute{
		Description: fmt.Sprintf("A list of %s that are scoped to this variable value.", scopeDescription),
		Optional:    true,
		ElementType: basetypes.StringType{},
	}
}
