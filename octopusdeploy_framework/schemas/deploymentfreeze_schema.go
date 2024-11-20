package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeploymentFreezeSchema struct{}

const DeploymentFreezeResourceName = "deployment_freeze"
const DeploymentFreezeDatasourceName = "deployment_freezes"

func (d DeploymentFreezeSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":    GetIdResourceSchema(),
			"name":  GetNameResourceSchema(true),
			"start": GetDateTimeResourceSchema("start"),
			"end":   GetDateTimeResourceSchema("end"),
			"ProjectEnvironmentScope": resourceSchema.MapAttribute{
				Description: "projects with environment scopes",
				Required:    true,
				ElementType: types.SetType{ElemType: types.StringType},
			},
		},
	}
}

func (d DeploymentFreezeSchema) GetDatasourceSchema() datasourceSchema.Schema {
	//TODO implement me
	panic("implement me")
}

var _ EntitySchema = &DeploymentFreezeSchema{}
