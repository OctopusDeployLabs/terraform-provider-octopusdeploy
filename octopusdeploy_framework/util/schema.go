package util

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetQueryIDsDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.ListAttribute{
		Description: "A filter to search by a list of IDs.",
		ElementType: types.StringType,
		Optional:    true,
	}
}

func GetQuerySpaceIDDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "A Space ID to filter by. Will revert what is specified on the provider if not set.",
		Optional:    true,
	}
}

func GetQueryNameDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "A filter search by exact name",
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

func GetNameDatasourceWithMaxLengthSchema(isRequired bool, maxLength int) datasourceSchema.Attribute {
	s := datasourceSchema.StringAttribute{
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

func GetNameDatasourceSchema(isRequired bool) datasourceSchema.Attribute {
	s := datasourceSchema.StringAttribute{
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

func GetDescriptionDatasourceSchema(resourceDescription string) datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "The description of this " + resourceDescription + ".",
		Computed:    true,
	}
}

func GetQueryDatasourceTags() datasourceSchema.Attribute {
	return datasourceSchema.ListAttribute{
		Description: "A filter to search by a list of tags.",
		ElementType: types.StringType,
		Optional:    true,
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
	return resourceSchema.StringAttribute{
		Description: fmt.Sprintf("The unique slug of this %s", resourceDescription),
		Optional:    true,
		Computed:    true,
	}
}

func GetSortOrderDataSourceSchema(resourceDescription string) resourceSchema.Attribute {
	return resourceSchema.Int64Attribute{
		Description: fmt.Sprintf("The order number to sort an %s", resourceDescription),
		Computed:    true,
	}
}

func GetSortOrderResourceSchema(resourceDescription string) resourceSchema.Attribute {
	return resourceSchema.Int64Attribute{
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

func GetPasswordDataSourceSchema(isRequired bool) datasourceSchema.Attribute {
	s := datasourceSchema.StringAttribute{
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

func GetRequiredStringResourceSchema(description string) resourceSchema.StringAttribute {
	return resourceSchema.StringAttribute{
		Required:    true,
		Description: description,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}
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

func GetDownloadAttemptsResourceSchema() resourceSchema.Attribute {
	return resourceSchema.Int64Attribute{
		Default:     int64default.StaticInt64(5),
		Description: "The number of times a deployment should attempt to download a package from this feed before failing.",
		Optional:    true,
		Computed:    true,
	}
}

func GetDownloadRetryBackoffSecondsResourceSchema() resourceSchema.Attribute {
	return resourceSchema.Int64Attribute{
		Default:     int64default.StaticInt64(10),
		Description: "The number of seconds to apply as a linear back off between download attempts.",
		Optional:    true,
		Computed:    true,
	}
}

func GetPackageAcquisitionLocationOptionsResourceSchema() resourceSchema.Attribute {
	return resourceSchema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
		Optional:    true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
	}
}

func GetFeedUriResourceSchema() resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Required: true,
	}
}
