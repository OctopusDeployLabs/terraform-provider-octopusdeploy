package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func expandActionTemplateParameters(actionTemplateParameters []interface{}) []actiontemplates.ActionTemplateParameter {
	if len(actionTemplateParameters) == 0 {
		return nil
	}

	expandedActionTemplateParameters := []actiontemplates.ActionTemplateParameter{}
	for _, actionTemplateParameter := range actionTemplateParameters {
		expandedActionTemplateParameters = append(expandedActionTemplateParameters, expandActionTemplateParameter(actionTemplateParameter.(types.Map).Elements()))
	}
	return expandedActionTemplateParameters
}

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

func mapTemplateNamesToIds(actionTemplateParameters []actiontemplates.ActionTemplateParameter) map[string]string {
	templateNameIds := map[string]string{}
	for _, actionTemplateParameter := range actionTemplateParameters {
		templateNameIds[actionTemplateParameter.Name] = actionTemplateParameter.ID
	}
	return templateNameIds
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

func getActionTemplateParameterSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"default_value": resourceSchema.StringAttribute{
			Description: "A default value for the parameter, if applicable. This can be a hard-coded value or a variable reference.",
			Optional:    true,
		},
		"display_settings": resourceSchema.MapAttribute{
			Description: "The display settings for the parameter.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"help_text": resourceSchema.StringAttribute{
			Description: "The help presented alongside the parameter input.",
			Optional:    true,
		},
		"id": util.GetIdResourceSchema(),
		"label": resourceSchema.StringAttribute{
			Description: "The label shown beside the parameter when presented in the deployment process. Example: `Server name`.",
			Optional:    true,
		},
		"name": util.GetNameResourceSchema(true),
	}
}

func flattenDisplaySettings(displaySettings map[string]attr.Value) map[string]string {
	flattenedDisplaySettings := make(map[string]string, len(displaySettings))
	for key, displaySetting := range displaySettings {
		flattenedDisplaySettings[key] = displaySetting.(types.String).ValueString()
	}
	return flattenedDisplaySettings
}
