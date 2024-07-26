package schemas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	//"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

func GetIdDatasourceSchema() datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "The unique ID for this resource.",
		Computed:    true,
		Optional:    true,
	}
}

func GetSpaceIdDatasourceSchema(resourceDescription string) datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: "The space ID associated with this " + resourceDescription + ".",
		Computed:    true,
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
		Optional:    true,
	}
}

func GetIdResourceSchema() resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Description: "The unique ID for this resource.",
		Computed:    true,
		Optional:    true,
	}
}

func GetSpaceIdResourceSchema(resourceDescription string) resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Description: "The space ID associated with this " + resourceDescription + ".",
		Computed:    true,
		Optional:    true,
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

func GetSlugDatasourceSchema(resourceDescription string) datasourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: fmt.Sprintf("The unique slug of this %s", resourceDescription),
		Optional:    true,
		Computed:    true,
	}
}

func GetSlugResourceSchema(resourceDescription string) resourceSchema.Attribute {
	return datasourceSchema.StringAttribute{
		Description: fmt.Sprintf("The unique slug of this %s", resourceDescription),
		Optional:    true,
		Computed:    true,
	}
}

func GetIds(ids types.List) []string {
	var result = make([]string, 0, len(ids.Elements()))
	for _, id := range ids.Elements() {
		result = append(result, id.String())
	}
	return result
}

func GetBranchResourceSchema(resourceDescription string) resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Computed:    true,
		Description: fmt.Sprintf("The branch name associated with this %s (i.e. `main`). This value is optional and only applies to associated projects that are stored in version control.", resourceDescription),
		Optional:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func GetProjectIdResourceSchema(resourceDescription string) resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Description: fmt.Sprintf("The project ID associated with this %s.", resourceDescription),
		Required:    true,
	}
}

func getConditionExpressionResourceSchema() resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Computed:    true,
		Description: "The expression to evaluate to determine whether to run this step when 'condition' is 'Variable'",
		Optional:    true,
	}
}

func getPackageRequirementResourceSchema() resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Default:     stringdefault.StaticString("LetOctopusDecide"),
		Description: "Whether to run this step before or after package acquisition (if possible)",
		Optional:    true,
		Computed:    true,
		Validators: []validator.String{
			stringvalidator.OneOf(
				"AfterPackageAcquisition",
				"BeforePackageAcquisition",
				"LetOctopusDecide"),
		},
	}
}

func getPropertiesResourceSchema() resourceSchema.Attribute {
	return resourceSchema.MapAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Optional:    true,
	}
}

func getStartTriggerResourceSchema() resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Default:     stringdefault.StaticString("StartAfterPrevious"),
		Description: "Whether to run this step after the previous step ('StartAfterPrevious') or at the same time as the previous step ('StartWithPrevious')",
		Optional:    true,
		Computed:    true,
		Validators: []validator.String{
			stringvalidator.OneOf("StartAfterPrevious", "StartWithPrevious"),
		},
	}
}

func getTargetRolesResourceSchema() resourceSchema.Attribute {
	return resourceSchema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "The roles that this step run against, or runs on behalf of",
		Optional:    true,
	}
}

func getWindowSizeResourceSchema() resourceSchema.Attribute {
	return resourceSchema.StringAttribute{
		Description: "The maximum number of targets to deploy to simultaneously",
		Optional:    true,
	}
}

func GetNumber(val types.Int64) int {
	v := 0
	if !val.IsNull() {
		v = int(val.ValueInt64())
	}

	return v
}
