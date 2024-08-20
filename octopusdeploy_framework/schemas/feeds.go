package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type FeedsSchema struct{}

var _ EntitySchema = FeedsSchema{}

func (f FeedsSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{}
}

func (f FeedsSchema) GetDatasourceSchemaAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"feed_type": datasourceSchema.StringAttribute{
			Description: "A filter to search by feed type. Valid feed types are `AwsElasticContainerRegistry`, `BuiltIn`, `Docker`, `GitHub`, `Helm`, `Maven`, `NuGet`, or `OctopusProject`.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					"AwsElasticContainerRegistry",
					"BuiltIn",
					"Docker",
					"GitHub",
					"Helm",
					"Maven",
					"NuGet",
					"OctopusProject"),
			},
		},
		"ids":          GetQueryIDsDatasourceSchema(),
		"name":         GetNameDatasourceSchema(false),
		"partial_name": GetQueryPartialNameDatasourceSchema(),
		"skip":         GetQuerySkipDatasourceSchema(),
		"take":         GetQueryTakeDatasourceSchema(),
		"space_id":     GetSpaceIdDatasourceSchema("feeds", false),

		// response
		"id": GetIdDatasourceSchema(true),
	}
}
