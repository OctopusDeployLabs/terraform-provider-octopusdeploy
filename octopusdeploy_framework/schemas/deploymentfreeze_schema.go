package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type DeploymentFreezeSchema struct{}

const DeploymentFreezeResourceName = "deployment_freeze"

func (d DeploymentFreezeSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":    GetIdResourceSchema(),
			"name":  GetNameResourceSchema(true),
			"start": GetDateTimeResourceSchema("The start time of the freeze, must be RFC3339 format", true),
			"end":   GetDateTimeResourceSchema("The end time of the freeze, must be RFC3339 format", true),
		},
	}
}

func (d DeploymentFreezeSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

var _ EntitySchema = &DeploymentFreezeSchema{}
