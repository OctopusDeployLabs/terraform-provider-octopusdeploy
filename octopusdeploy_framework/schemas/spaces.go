package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type SpacesSchema struct{}

var _ EntitySchema = SpacesSchema{}

func (s SpacesSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{}
}

func (s SpacesSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Attributes: map[string]datasourceSchema.Attribute{
			// request
			"ids":          GetQueryIDsDatasourceSchema(),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take":         GetQueryTakeDatasourceSchema(),

			// response
			"id": GetIdDatasourceSchema(true),
			"spaces": datasourceSchema.ListNestedAttribute{
				Computed: true,
				Optional: false,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: SpaceSchema{}.GetDatasourceSchema().Attributes,
				},
			},
		},
	}
}
