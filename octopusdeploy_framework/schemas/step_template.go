package schemas

import (
	"fmt"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

const (
	StepTemplateResourceDescription = "step_template"
)

type StepTemplateSchema struct{}

var _ EntitySchema = StepTemplateSchema{}

func (s StepTemplateSchema) GetDatasourceSchema() ds.Schema {
	return ds.Schema{}
}

func (s StepTemplateSchema) GetResourceSchema() rs.Schema {
	return rs.Schema{
		Description: util.GetResourceSchemaDescription(StepTemplateResourceDescription),
		Attributes: map[string]rs.Attribute{
			"id":          GetIdResourceSchema(),
			"name":        GetNameResourceSchema(true),
			"description": GetDescriptionResourceSchema(EnvironmentResourceDescription),
			"space_id":    GetSpaceIdResourceSchema(EnvironmentResourceDescription),
			"version": rs.Int32Attribute{
				Description: "The version of the step template",
				Optional:    false,
				Computed:    true,
			},
			"step_package_id": rs.StringAttribute{
				Description: "The ID of the step package",
				Required:    true,
			},
			"action_type": rs.StringAttribute{
				Description: "The action type of the step template",
				Required:    true,
			},
			"community_action_template_id": rs.StringAttribute{
				Description: "The ID of the community action template",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"packages": GetStepTemplatePackageSchema(),
			"properties": rs.MapAttribute{
				Description: "Properties for the step template",
				Required:    true,
				ElementType: types.StringType,
			},
			"parameters": GetStepTemplateParameterSchema(),
		},
	}
}

func GetStepTemplateParameterSchema() rs.ListNestedAttribute {
	return rs.ListNestedAttribute{
		Description: "List of parameters that can be used in Step Template.",
		Required:    true,
		NestedObject: rs.NestedAttributeObject{
			Attributes: map[string]rs.Attribute{
				"default_value": rs.StringAttribute{
					Description: "A default value for the parameter, if applicable. This can be a hard-coded value or a variable reference.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"display_settings": rs.MapAttribute{
					Description: "The display settings for the parameter.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"help_text": rs.StringAttribute{
					Description: "The help presented alongside the parameter input.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"id": rs.StringAttribute{
					Description: "The id for the property.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"), fmt.Sprintf("must be a valid UUID, unique within this list. Here is one you could use: %s.\nExpect uuid", uuid.New())),
					},
				},
				"label": rs.StringAttribute{
					Description: "The label shown beside the parameter when presented in the deployment process. Example: `Server name`.",
					Optional:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"name": rs.StringAttribute{
					Description: "The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods. Example: `ServerName`",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
			},
		},
	}
}

func GetStepTemplatePackageSchema() rs.ListNestedAttribute {
	return rs.ListNestedAttribute{
		Description: "Package information for the step template",
		Required:    true,
		NestedObject: rs.NestedAttributeObject{
			Attributes: map[string]rs.Attribute{
				"acquisition_location": rs.StringAttribute{
					Description: "Acquisition location for the package.",
					Default:     stringdefault.StaticString("Server"),
					Optional:    true,
					Computed:    true,
				},
				"feed_id": rs.StringAttribute{
					Description: "ID of the feed.",
					Required:    true,
				},
				"id":   GetIdResourceSchema(),
				"name": GetNameResourceSchema(true),
				"package_id": rs.StringAttribute{
					Description: "The ID of the package to use.",
					Optional:    true,
					Required:    false,
					Computed:    true,
				},
				"properties": rs.SingleNestedAttribute{
					Description: "Properties for the package.",
					Required:    true,
					Attributes: map[string]rs.Attribute{
						"extract": rs.StringAttribute{
							Description: "If the package should extract.",
							Default:     stringdefault.StaticString("True"),
							Optional:    true,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("^(True|Fasle)$"), "Extract must be True or False"),
							},
						},
						"package_parameter_name": rs.StringAttribute{
							Description: "The name of the package parameter",
							Default:     stringdefault.StaticString(""),
							Optional:    true,
							Computed:    true,
						},
						"purpose": rs.StringAttribute{
							Description: "The purpose of this property.",
							Default:     stringdefault.StaticString(""),
							Optional:    true,
							Required:    false,
							Computed:    true,
						},
						"selection_mode": rs.StringAttribute{
							Description: "The selection mode.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

type StepTemplateTypeResourceModel struct {
	ActionType                types.String `tfsdk:"action_type"`
	SpaceID                   types.String `tfsdk:"space_id"`
	CommunityActionTemplateId types.String `tfsdk:"community_action_template_id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	Packages                  types.List   `tfsdk:"packages"`
	Parameters                types.List   `tfsdk:"parameters"`
	Properties                types.Map    `tfsdk:"properties"`
	StepPackageId             types.String `tfsdk:"step_package_id"`
	Version                   types.Int32  `tfsdk:"version"`

	ResourceModel
}

type StepTemplatePackageType struct {
	ID                  types.String `tfsdk:"id"`
	AcquisitionLocation types.String `tfsdk:"acquisition_location"`
	Name                types.String `tfsdk:"name"`
	FeedID              types.String `tfsdk:"feed_id"`
	PackageID           types.String `tfsdk:"package_id"`
	Properties          types.Object `tfsdk:"properties"`
}

func GetStepTemplatePackageTypeAttributes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                   types.StringType,
		"acquisition_location": types.StringType,
		"name":                 types.StringType,
		"feed_id":              types.StringType,
		"package_id":           types.StringType,
		"properties":           types.ObjectType{AttrTypes: GetStepTemplatePackagePropertiesTypeAttributes()},
	}
}

func GetStepTemplatePackagePropertiesTypeAttributes() map[string]attr.Type {
	return map[string]attr.Type{
		"extract":                types.StringType,
		"package_parameter_name": types.StringType,
		"purpose":                types.StringType,
		"selection_mode":         types.StringType,
	}
}

type StepTemplateParameterType struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Label           types.String `tfsdk:"label"`
	HelpText        types.String `tfsdk:"help_text"`
	DisplaySettings types.Map    `tfsdk:"display_settings"`
	DefaultValue    types.String `tfsdk:"default_value"`
}

func GetStepTemplateParameterTypeAttributes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"name":             types.StringType,
		"label":            types.StringType,
		"help_text":        types.StringType,
		"display_settings": types.MapType{ElemType: types.StringType},
		"default_value":    types.StringType,
	}
}
