package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IResourceModel interface {
	GetID() string
}

type ResourceModel struct {
	ID types.String `tfsdk:"id"`

	IResourceModel `tfsdk:"-"`
}

func (r ResourceModel) GetID() string {
	return r.ID.ValueString()
}
