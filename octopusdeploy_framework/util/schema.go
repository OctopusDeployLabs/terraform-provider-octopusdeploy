package util

import (
	"fmt"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

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

func GetQueryNameDatasourceSchema() schema.Attribute {
	return schema.StringAttribute{
		Description: "A filter search by exact name",
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

func GetNameDatasourceWithMaxLengthSchema(isRequired bool, maxLength int) schema.Attribute {
	s := schema.StringAttribute{
		Description: fmt.Sprintf("The name of this resource, no more than %d characters long", maxLength),
		Validators: []validator.String{
			stringvalidator.LengthBetween(1, maxLength),
		},
	}

	if isRequired {
		s.Required = true
	} else {
		s.Optional = true
	}

	return s
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

func GetIdResourceSchema() resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Description: "The unique ID for this resource.",
		Computed:    true,
		Optional:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func GetSpaceIdResourceSchema(resourceDescription string) resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Description: "The space ID associated with this " + resourceDescription + ".",
		Computed:    true,
		Optional:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func GetNameResourceSchema(isRequired bool) schema.Attribute {
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

func GetDescriptionResourceSchema(resourceDescription string) schema.Attribute {
	return schema.StringAttribute{
		Description: "The description of this " + resourceDescription + ".",
		Optional:    true,
		Computed:    true,
	}
}

func GetSlugDatasourceSchema(resourceDescription string) schema.Attribute {
	return schema.StringAttribute{
		Description: fmt.Sprintf("The unique slug of this %s", resourceDescription),
		Optional:    true,
		Computed:    true,
	}
}

func GetSortOrderDataSourceSchema(resourceDescription string) schema.Attribute {
	return schema.Int64Attribute{
		Description: fmt.Sprintf("The order number to sort an %s", resourceDescription),
		Optional:    true,
		Computed:    true,
	}
}

func GetPasswordResourceSchema(isRequired bool) resourceSchema.Attribute {
	s := resourceSchema.StringAttribute{
		Description: "The password associated with this resource.",
		Sensitive:   true,
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

func GetUsernameResourceSchema(isRequired bool) resourceSchema.Attribute {
	s := &resourceSchema.StringAttribute{
		Description: "The username associated with this resource.",
		Sensitive:   true,
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

func GetIds(ids types.List) []string {
	var result = make([]string, 0, len(ids.Elements()))
	for _, id := range ids.Elements() {
		if str, ok := id.(types.String); ok {
			result = append(result, str.ValueString())
		}
	}
	return result
}

func GetNumber(val types.Int64) int {
	v := 0
	if !val.IsNull() {
		v = int(val.ValueInt64())
	}

	return v
}
