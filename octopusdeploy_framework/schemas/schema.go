package schemas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	//"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SchemaAttributeNames = struct {
	ID          string
	Name        string
	Description string
	SpaceID     string
}{
	ID:          "id",
	Name:        "name",
	Description: "description",
	SpaceID:     "space_id",
}

func GetQueryIDsDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.ListAttribute{
		Description: "A filter to search by a list of IDs.",
		ElementType: types.StringType,
		Optional:    true,
	}
}

func GetQueryPartialNameDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "A filter to search by a partial name.",
		Optional:    true,
	}
}

func GetQuerySkipDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.Int64Attribute{
		Description: "A filter to specify the number of items to skip in the response.",
		Optional:    true,
	}
}

func GetQueryTakeDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.Int64Attribute{
		Description: "A filter to specify the number of items to take (or return) in the response.",
		Optional:    true,
	}
}

func GetReadonlyNameDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "The name of this resource.",
		Computed:    true,
	}
}

func GetIdDatasourceSchema(isReadOnly bool) datasourceSchema.Attribute {
	s := datasourceSchema.StringAttribute{
		Description: "The unique ID for this resource.",
	}

	if isReadOnly {
		s.Computed = true
	} else {
		s.Computed = true
		s.Optional = true
	}

	return s
}

func GetSpaceIdDatasourceSchema(resourceDescription string, isReadOnly bool) datasourceSchema.Attribute {
	s := datasourceSchema.StringAttribute{
		Description: "The space ID associated with this " + resourceDescription + ".",
	}

	if isReadOnly {
		s.Computed = true
	} else {
		s.Computed = true
		s.Optional = true
	}

	return s
}

func GetReadonlyDescriptionDatasourceSchema(resourceDescription string) datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "The description of this " + resourceDescription + ".",
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

func GetNameResourceSchema(isRequired bool) resourceSchema.Attribute {
	s := resourceSchema.StringAttribute{
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

func GetDescriptionResourceSchema(resourceDescription string) resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Description: "The description of this " + resourceDescription + ".",
		Optional:    true,
		Computed:    true,
		Default:     stringdefault.StaticString(""),
	}
}

func GetSlugDatasourceSchema(resourceDescription string, isReadOnly bool) datasourceSchema.Attribute {
	s := datasourceSchema.StringAttribute{
		Description: fmt.Sprintf("The unique slug of this %s", resourceDescription),
	}

	if isReadOnly {
		s.Computed = true
	} else {
		s.Optional = true
		s.Computed = true
	}

	return s
}

func GetSlugResourceSchema(resourceDescription string) resourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: fmt.Sprintf("The unique slug of this %s", resourceDescription),
		Optional:    true,
		Computed:    true,
	}
}

func GetBooleanDatasourceAttribute(description string, isOptional bool) datasourceSchema.Attribute {
	return datasourceSchema.BoolAttribute{
		Description: description,
		Optional:    isOptional,
		Computed:    true,
	}
}

func GetBooleanResourceAttribute(description string, defaultValue bool, isOptional bool) resourceSchema.Attribute {
	return resourceSchema.BoolAttribute{
		Default:     booldefault.StaticBool(defaultValue),
		Description: description,
		Optional:    isOptional,
		Computed:    true,
	}
}

func GetIds(ids types.List) []string {
	var result = make([]string, 0, len(ids.Elements()))
	for _, id := range ids.Elements() {
		strVal, ok := id.(types.String)

		if !ok || strVal.IsNull() || strVal.IsUnknown() {
			continue
		}
		result = append(result, strVal.ValueString())
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
