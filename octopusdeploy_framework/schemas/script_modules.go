package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
)

type ScriptModuleResourceModel struct {
	Description   types.String `tfsdk:"description"`
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	SpaceID       types.String `tfsdk:"space_id"`
	VariableSetId types.String `tfsdk:"variable_set_id"`
	Script        types.List   `tfsdk:"script"`
}

func ScriptModuleObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"description":     types.StringType,
		"id":              types.StringType,
		"name":            types.StringType,
		"space_id":        types.StringType,
		"variable_set_id": types.StringType,
		"script": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"body":   types.StringType,
			"syntax": types.StringType,
		}}},
	}
}

func GetScriptModuleSchemaBlock() map[string]resourceSchema.Block {
	return map[string]resourceSchema.Block{
		"script": resourceSchema.ListNestedBlock{
			Description: "The script associated with this script module.",
			NestedObject: resourceSchema.NestedBlockObject{
				Attributes: map[string]resourceSchema.Attribute{
					"body": resourceSchema.StringAttribute{
						Description: "The body of this script module.",
						Required:    true,
					},
					"syntax": resourceSchema.StringAttribute{
						Description: "The syntax of the script. Valid types are `Bash`, `CSharp`, `FSharp`, `PowerShell`, or `Python`.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOfCaseInsensitive(
								"Bash",
								"CSharp",
								"FSharp",
								"PowerShell",
								"Python"),
						},
					},
				},
			},
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
				listvalidator.SizeAtLeast(1),
			},
		},
	}
}

func GetScriptModuleResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"description": GetDescriptionResourceSchema("script module"),
		"id":          GetIdResourceSchema(),
		"name":        GetNameResourceSchema(true),
		"space_id":    GetSpaceIdResourceSchema("Script Module"),
		"variable_set_id": resourceSchema.StringAttribute{
			Computed:    true,
			Description: "The variable set ID for this script module.",
			Optional:    true,
		},
	}
}

func MapFromScriptModuleToState(data *ScriptModuleResourceModel) *variables.ScriptModule {
	name := data.Name.ValueString()
	scriptModule := variables.NewScriptModule(name)
	scriptModule.ID = data.ID.ValueString()
	scriptModule.Description = data.Description.ValueString()
	// We enforce on the schema a single required script
	scriptDetails := data.Script.Elements()[0].(types.Object).Attributes()
	scriptModule.Syntax = scriptDetails["syntax"].(types.String).ValueString()
	scriptModule.ScriptBody = scriptDetails["body"].(types.String).ValueString()
	scriptModule.SpaceID = data.SpaceID.ValueString()
	scriptModule.VariableSetID = data.VariableSetId.ValueString()

	return scriptModule
}

func flattenScript(scriptModule *variables.ScriptModule) []attr.Value {
	return []attr.Value{
		types.ObjectValueMust(map[string]attr.Type{
			"body":   types.StringType,
			"syntax": types.StringType,
		}, map[string]attr.Value{
			"body":   types.StringValue(scriptModule.ScriptBody),
			"syntax": types.StringValue(scriptModule.Syntax),
		}),
	}
}

func MapToScriptModuleFromState(data *ScriptModuleResourceModel, scriptModule *variables.ScriptModule) {
	data.Description = types.StringValue(scriptModule.Description)
	data.Name = types.StringValue(scriptModule.Name)
	data.SpaceID = types.StringValue(scriptModule.SpaceID)
	data.VariableSetId = types.StringValue(scriptModule.VariableSetID)
	data.ID = types.StringValue(scriptModule.ID)

	flattenScript(scriptModule)

	var script, _ = types.ListValue(types.StringType, flattenScript(scriptModule))
	data.Script = script
}
