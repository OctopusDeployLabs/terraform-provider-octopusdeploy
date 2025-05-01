package schemas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

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

func GetQueryNameDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "A filter search by exact name",
		Optional:    true,
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

func GetUsernameDatasourceSchema(isRequired bool) datasourceSchema.Attribute {
	s := datasourceSchema.StringAttribute{
		Description: "The username associated with this resource.",
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

func GetValueDatasourceSchema(isRequired bool) datasourceSchema.Attribute {
	s := datasourceSchema.StringAttribute{
		Description: "The value of this resource.",
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

func GetIdResourceSchema() resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Description: "The unique ID for this resource.",
		Computed:    true,
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

func GetBooleanDatasourceAttribute(description string, isOptional bool) datasourceSchema.Attribute {
	return datasourceSchema.BoolAttribute{
		Description: description,
		Optional:    isOptional,
		Computed:    true,
	}
}

func GetSortOrderDatasourceSchema(resourceDescription string) datasourceSchema.Attribute {
	return datasourceSchema.Int64Attribute{
		Description: fmt.Sprintf("The order number to sort an %s", resourceDescription),
		Computed:    true,
	}
}

func GetSortOrderResourceSchema(resourceDescription string) resourceSchema.Attribute {
	return resourceSchema.Int64Attribute{
		Description: fmt.Sprintf("The order number to sort an %s.", resourceDescription),
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

func GetDownloadAttemptsResourceSchema() resourceSchema.Attribute {
	return resourceSchema.Int64Attribute{
		Default:     int64default.StaticInt64(5),
		Description: "The number of times a deployment should attempt to download a package from this feed before failing.",
		Optional:    true,
		Computed:    true,
	}
}

func GetUsernameResourceSchema(isRequired bool) resourceSchema.Attribute {
	s := resourceSchema.StringAttribute{
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

func GetDownloadRetryBackoffSecondsResourceSchema() resourceSchema.Attribute {
	return resourceSchema.Int64Attribute{
		Default:     int64default.StaticInt64(10),
		Description: "The number of seconds to apply as a linear back off between download attempts.",
		Optional:    true,
		Computed:    true,
	}
}

func GetPortNumberResourceSchema() resourceSchema.Attribute {
	return resourceSchema.Int64Attribute{
		Description: "The port number of the host to connect to (usually 22)",
		Required:    true,
	}
}

func GetValueResourceSchema(isRequired bool) resourceSchema.Attribute {
	s := resourceSchema.StringAttribute{
		Description: "The value of this resource.",
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

func GetOptionalBooleanResourceAttribute(description string, defaultValue bool) resourceSchema.Attribute {
	return resourceSchema.BoolAttribute{
		Default:     booldefault.StaticBool(defaultValue),
		Description: description,
		Optional:    true,
		Computed:    true,
	}
}

func GetRequiredBooleanResourceAttribute(description string) resourceSchema.Attribute {
	return resourceSchema.BoolAttribute{
		Description: description,
		Required:    true,
	}
}

func GetReadonlyBooleanResourceAttribute(description string) resourceSchema.Attribute {
	return resourceSchema.BoolAttribute{
		Description: description,
		Computed:    true,
	}
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

func GetOptionalStringResourceSchema(description string) resourceSchema.StringAttribute {
	return resourceSchema.StringAttribute{
		Optional:    true,
		Description: description,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}
}

func GetFeedUriResourceSchema() resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Required: true,
	}
}

func GetNumber(val types.Int64) int {
	v := 0
	if !val.IsNull() {
		v = int(val.ValueInt64())
	}

	return v
}

func GetSensitiveResourceSchema(description string, isRequired bool) resourceSchema.Attribute {
	s := resourceSchema.StringAttribute{
		Description: description,
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

func GetDateTimeResourceSchema(description string, isRequired bool) resourceSchema.Attribute {
	regex := "^((?:(\\d{4}-\\d{2}-\\d{2})T(\\d{2}:\\d{2}:\\d{2}(?:\\.\\d+)?))(?:Z|[\\+-]\\d{2}:\\d{2})?)"
	return resourceSchema.StringAttribute{
		Description: description,
		Required:    isRequired,
		Optional:    !isRequired,
		CustomType:  timetypes.RFC3339Type{},
		Validators: []validator.String{
			stringvalidator.RegexMatches(regexp.MustCompile(regex), fmt.Sprintf("must match RFC3339 format, %s", regex)),
		},
	}
}

func GetOidcSubjectKeysSchema(description string, isRequired bool) resourceSchema.Attribute {
	return resourceSchema.ListAttribute{
		Optional:    !isRequired,
		Required:    isRequired,
		Description: description,
		ElementType: types.StringType,
	}
}

func getCertificateDataFormatResourceSchema() resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Description: "Specifies the archive file format used for storing cryptography objects in the certificate. Valid formats are `Der`, `Pem`, `Pkcs12`, or `Unknown`.",
		Computed:    true,
		Optional:    true,
		Validators: []validator.String{
			stringvalidator.OneOf("Der", "Pem", "Pkcs12", "Unknown"),
		},
	}
}

func getEnvironmentsResourceSchema() resourceSchema.Attribute {
	return resourceSchema.ListAttribute{
		Description: "A list of environment IDs associated with this resource.",
		Computed:    true,
		Optional:    true,
		ElementType: types.StringType,
	}
}
