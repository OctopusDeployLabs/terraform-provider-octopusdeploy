package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IResourceModel interface {
	GetID() string
}

type ResourceModel struct {
	IResourceModel
	ID types.String `tfsdk:"id"`
}

func (r ResourceModel) GetID() string {
	return r.ID.ValueString()
}
