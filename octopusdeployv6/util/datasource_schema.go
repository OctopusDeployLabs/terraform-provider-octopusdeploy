package util

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetQueryIDsDatasourceSchema() schema.Attribute {
	return schema.ListAttribute{
		Description: "A filter to search by a list of IDs.",
		ElementType: types.StringType,
		Optional:    true,
	}
}

func GetQueryPartialNameDatasourceSchema() schema.Attribute {
	return schema.StringAttribute{
		Description: "A filter to search by a partial name.",
		Optional:    true,
	}
}

func GetQuerySkipDatasourceSchema() schema.Attribute {
	return schema.Int64Attribute{
		Description: "A filter to specify the number of items to skip in the response.",
		Optional:    true,
	}
}

func GetQueryTakeDatasourceSchema() schema.Attribute {
	return schema.Int64Attribute{
		Description: "A filter to specify the number of items to take (or return) in the response.",
		Optional:    true,
	}
}

func GetIdDatasourceSchema() schema.Attribute {
	return schema.StringAttribute{
		Description: "The unique ID for this resource.",
		Computed:    true,
		Optional:    true,
	}
}

func GetSpaceIdDatasourceSchema(resourceDescription string) schema.Attribute {
	return schema.StringAttribute{
		Description: "The space ID associated with this " + resourceDescription + ".",
		Computed:    true,
		Optional:    true,
	}
}

func GetNameDatasourceSchema(isRequired bool) schema.Attribute {
	s := schema.StringAttribute{
		Description: "The name of this resource.",
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}

	if isRequired {
		s.Required = true
	} else {
		s.Optional = true
	}

	return s
}

func GetDescriptionDatasourceSchema(resourceDescription string) schema.Attribute {
	return schema.StringAttribute{
		Description: "The description of this " + resourceDescription + ".",
		Optional:    true,
		Computed:    true,
	}
}
