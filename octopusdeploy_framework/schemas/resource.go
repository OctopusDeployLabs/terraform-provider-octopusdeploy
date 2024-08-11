package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IResourceModel interface {
	GetID() string
}

type ResourceModel struct {
	ID types.String `tfsdk:"id"`

	IResourceModel `tfsdk:"-"` // Ignore resource model interface in object conversion
}

func (r ResourceModel) GetID() string {
	return r.ID.ValueString()
}
