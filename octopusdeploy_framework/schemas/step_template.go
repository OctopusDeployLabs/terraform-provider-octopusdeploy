package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
				Optional: false,
				Computed: true,
			},
			"step_package_id": rs.StringAttribute{
				Required: true,
			},
			"action_type": rs.StringAttribute{
				Required: true,
			},
			"community_action_template_id": rs.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"packages": GetStepTemplatePackageSchema(),
			"properties": rs.MapAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"parameters": GetStepTemplateParameterSchema(),
		},
	}
}

func GetStepTemplateParameterSchema() rs.ListNestedAttribute {
	return rs.ListNestedAttribute{
		Required: true,
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
					Description: "The id for the attribute.",
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
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
		Required: true,
		NestedObject: rs.NestedAttributeObject{
			Attributes: map[string]rs.Attribute{
				"acquisition_location": rs.StringAttribute{
					Default:  stringdefault.StaticString("Server"),
					Computed: true,
				},
				"feed_id": rs.StringAttribute{
					Required: true,
				},
				"id":   GetIdResourceSchema(),
				"name": GetNameResourceSchema(true),
				"package_id": rs.StringAttribute{
					Optional: true,
					Required: false,
					Computed: true,
				},
				"properties": rs.SingleNestedAttribute{
					Required: true,
					Attributes: map[string]rs.Attribute{
						"extract": rs.StringAttribute{
							Default:  stringdefault.StaticString("True"),
							Computed: true,
						},
						"package_parameter_name": rs.StringAttribute{
							Required: true,
						},
						"purpose": rs.StringAttribute{
							Default:  stringdefault.StaticString(""),
							Computed: true,
						},
						"selection_mode": rs.StringAttribute{
							Required: true,
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
	Properties          types.Object `tfsdk:"Properties"`
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
