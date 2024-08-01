package schemas

import (
	"fmt"
	"strings"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var DisplaySettingsControlTypeNames = struct {
	CheckBox       string
	MultiLineText  string
	Select         string
	SingleLineText string
}{
	"CheckBox",
	"MultiLineText",
	"Select",
	"SingeLineText",
}

var displaySettingsControlTypes = []string{
	DisplaySettingsControlTypeNames.CheckBox,
	DisplaySettingsControlTypeNames.MultiLineText,
	DisplaySettingsControlTypeNames.Select,
	DisplaySettingsControlTypeNames.SingleLineText,
}

func getVariablePromptResourceSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				SchemaAttributeNames.Description:             util.GetDescriptionResourceSchema("variable prompt option"),
				VariableSchemaAttributeNames.DisplaySettings: getDisplaySettingsSchema(),
				VariableSchemaAttributeNames.IsRequired: resourceSchema.BoolAttribute{
					Optional: true,
				},
				VariableSchemaAttributeNames.Label: resourceSchema.StringAttribute{
					Optional: true,
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func getDisplaySettingsSchema() resourceSchema.ListNestedAttribute {
	return resourceSchema.ListNestedAttribute{
		Optional: true,
		NestedObject: resourceSchema.NestedAttributeObject{
			Attributes: map[string]resourceSchema.Attribute{
				VariableSchemaAttributeNames.ControlType: resourceSchema.StringAttribute{
					Description: fmt.Sprintf("The type of control for rendering this prompted variable. Valid types are %s", strings.Join(displaySettingsControlTypes, ", ")),
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							displaySettingsControlTypes...,
						),
					},
				},
				VariableSchemaAttributeNames.SelectOption: resourceSchema.ListNestedAttribute{
					Description: fmt.Sprintf("If the `%s` is `%s`, then this value defines an option.", VariableSchemaAttributeNames.ControlType, DisplaySettingsControlTypeNames.Select),
					Optional:    true,
					NestedObject: resourceSchema.NestedAttributeObject{
						Attributes: map[string]resourceSchema.Attribute{
							VariableSchemaAttributeNames.Value: resourceSchema.StringAttribute{
								Description: "The select value",
								Required:    true,
							},
							VariableSchemaAttributeNames.DisplayName: resourceSchema.StringAttribute{
								Description: "The display name for the select value",
								Required:    true,
							},
						},
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}
