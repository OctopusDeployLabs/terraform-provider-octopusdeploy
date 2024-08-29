package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type EntitySchema interface {
	GetResourceSchema() resourceSchema.Schema
	GetDatasourceSchema() datasourceSchema.Schema
}
