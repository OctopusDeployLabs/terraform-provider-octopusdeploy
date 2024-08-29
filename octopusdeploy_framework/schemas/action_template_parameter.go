package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ActionTemplateParameterSchema struct{}

var _ EntitySchema = ActionTemplateParameterSchema{}

func expandActionTemplateParameter(tfTemplate map[string]attr.Value) actiontemplates.ActionTemplateParameter {
	actionTemplateParameter := actiontemplates.NewActionTemplateParameter()

	propertyValue := core.NewPropertyValue(tfTemplate["default_value"].(types.String).ValueString(), false)
	actionTemplateParameter.DefaultValue = &propertyValue

	actionTemplateParameter.DisplaySettings = flattenDisplaySettings(tfTemplate["display_settings"].(types.Map).Elements())
	actionTemplateParameter.HelpText = tfTemplate["help_text"].(types.String).ValueString()
	actionTemplateParameter.ID = tfTemplate["id"].(types.String).ValueString()
	actionTemplateParameter.Label = tfTemplate["label"].(types.String).ValueString()
	actionTemplateParameter.Name = tfTemplate["name"].(types.String).ValueString()

	return *actionTemplateParameter
}

func ExpandActionTemplateParameters(actionTemplateParameters types.List) []actiontemplates.ActionTemplateParameter {
	if len(actionTemplateParameters.Elements()) == 0 {
		return []actiontemplates.ActionTemplateParameter{}
	}

	expandedActionTemplateParameters := []actiontemplates.ActionTemplateParameter{}
	for _, actionTemplateParameter := range actionTemplateParameters.Elements() {
		expandedActionTemplateParameters = append(expandedActionTemplateParameters, expandActionTemplateParameter(actionTemplateParameter.(types.Object).Attributes()))
	}
	return expandedActionTemplateParameters
}

// Note these are used in projects and were copied there during the Library Variable Set migration. This is currently unused.
func flattenActionTemplateParameters(actionTemplateParameters []actiontemplates.ActionTemplateParameter) []interface{} {
	flattenedActionTemplateParameters := make([]interface{}, 0)
	for _, actionTemplateParameter := range actionTemplateParameters {
		a := make(map[string]interface{})
		a["default_value"] = actionTemplateParameter.DefaultValue.Value
		a["display_settings"] = actionTemplateParameter.DisplaySettings
		a["help_text"] = actionTemplateParameter.HelpText
		a["id"] = actionTemplateParameter.ID
		a["label"] = actionTemplateParameter.Label
		a["name"] = actionTemplateParameter.Name
		flattenedActionTemplateParameters = append(flattenedActionTemplateParameters, a)
	}
	return flattenedActionTemplateParameters
}

func TemplateObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"default_value":    types.StringType,
		"display_settings": types.MapType{ElemType: types.StringType},
		"help_text":        types.StringType,
		"id":               types.StringType,
		"label":            types.StringType,
		"name":             types.StringType,
	}
}

func (a ActionTemplateParameterSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

func (a ActionTemplateParameterSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"description": GetDescriptionResourceSchema("library variable set"),
			"id":          GetIdResourceSchema(),
			"name":        GetNameResourceSchema(true),
			"space_id":    GetSpaceIdResourceSchema("library variable set"),
			"template_ids": resourceSchema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"variable_set_id": resourceSchema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Description: "This resource manages library variable sets in Octopus Deploy.",
		Blocks: map[string]resourceSchema.Block{
			"template": GetActionTemplateParameterSchema(),
		},
	}
}

func GetActionTemplateParameterSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"default_value": resourceSchema.StringAttribute{
					Description: "A default value for the parameter, if applicable. This can be a hard-coded value or a variable reference.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"display_settings": resourceSchema.MapAttribute{
					Description: "The display settings for the parameter.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"help_text": resourceSchema.StringAttribute{
					Description: "The help presented alongside the parameter input.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"id": GetIdResourceSchema(),
				"label": resourceSchema.StringAttribute{
					Description: "The label shown beside the parameter when presented in the deployment process. Example: `Server name`.",
					Optional:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"name": resourceSchema.StringAttribute{
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

func flattenDisplaySettings(displaySettings map[string]attr.Value) map[string]string {
	flattenedDisplaySettings := make(map[string]string, len(displaySettings))
	for key, displaySetting := range displaySettings {
		flattenedDisplaySettings[key] = displaySetting.(types.String).ValueString()
	}
	return flattenedDisplaySettings
}
