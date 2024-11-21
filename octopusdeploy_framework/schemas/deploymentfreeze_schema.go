package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeploymentFreezeSchema struct{}

const DeploymentFreezeResourceName = "deployment_freeze"

func (d DeploymentFreezeSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":    GetIdResourceSchema(),
			"name":  GetNameResourceSchema(true),
			"start": GetDateTimeResourceSchema("start", true),
			"end":   GetDateTimeResourceSchema("end", true),
			"project_environment_scope": resourceSchema.MapAttribute{
				Description: "projects with environment scopes",
				Required:    true,
				ElementType: types.SetType{ElemType: types.StringType},
			},
		},
	}
}

func (d DeploymentFreezeSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

var _ EntitySchema = &DeploymentFreezeSchema{}
