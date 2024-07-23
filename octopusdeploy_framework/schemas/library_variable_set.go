package schemas

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	types "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
)

type LibraryVariableSetResourceModel struct {
	Description   types.String `tfsdk:"description"`
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	SpaceID       types.String `tfsdk:"space_id"`
	Template      types.List   `tfsdk:"template"`
	TemplateIds   types.Map    `tfsdk:"template_ids"`
	VariableSetId types.String `tfsdk:"variable_set_id"`
}

func GetLibraryVariableSetDataSourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Attributes: getLibraryVariableSetDataSchema(),
	}
}

func getLibraryVariableSetDataSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"content_type": datasourceSchema.StringAttribute{
			Description: "A filter to search by content type.",
			Optional:    true,
		},
		"id":       util.GetIdDatasourceSchema(),
		"space_id": util.GetSpaceIdDatasourceSchema("library variable set"),
		"ids":      util.GetQueryIDsDatasourceSchema(),
		"library_variable_sets": datasourceSchema.ListAttribute{
			Computed:    true,
			Description: "A list of library variable sets that match the filter(s).",
			ElementType: types.ObjectType{AttrTypes: GetLibraryVariableSetObjectType()},
			Optional:    true,
		},
		"partial_name": util.GetQueryPartialNameDatasourceSchema(),
		"skip":         util.GetQuerySkipDatasourceSchema(),
		"take":         util.GetQueryTakeDatasourceSchema(),
	}
}

func GetLibraryVariableSetObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"description":     types.StringType,
		"id":              types.StringType,
		"name":            types.StringType,
		"space_id":        types.StringType,
		"template":        types.ListType{ElemType: types.ObjectType{AttrTypes: TemplateObjectType()}},
		"template_ids":    types.MapType{ElemType: types.StringType},
		"variable_set_id": types.StringType,
	}
}

func GetLibraryVariableSetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"description": GetDescriptionResourceSchema("library variable set"),
			"id":          GetIdResourceSchema(),
			"name":        GetNameResourceSchema(true),
			"space_id":    GetSpaceIdResourceSchema("library variable set"),
			"template": resourceSchema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.ObjectType{AttrTypes: TemplateObjectType()},
			},
			"template_ids": resourceSchema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"variable_set_id": resourceSchema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// fixTemplateIds uses the suggestion from https://github.com/hashicorp/terraform/issues/18863
// to ensure that the template_ids field has keys to match the list of template names.
func fixTemplateIds(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	templates := d.Get("template")
	templateIds := map[string]string{}
	if templates != nil {
		for _, t := range templates.([]interface{}) {
			template := t.(map[string]interface{})
			templateIds[template["name"].(string)] = template["id"].(string)
		}
	}
	if err := d.SetNew("template_ids", templateIds); err != nil {
		return err
	}

	return nil
}

func UpdateDataFromLibraryVariableSet(data *LibraryVariableSetResourceModel, spaceId string, libraryVariableSet *variables.LibraryVariableSet) {
	data.Description = types.StringValue(libraryVariableSet.Description)
	data.Name = types.StringValue(libraryVariableSet.Name)
	data.VariableSetId = types.StringValue(libraryVariableSet.VariableSetID)
	data.SpaceID = types.StringValue(spaceId)

	data.Template = FlattenTemplates(libraryVariableSet.Templates)
	data.TemplateIds = FlattenTemplateIds(libraryVariableSet.Templates)

	data.ID = types.StringValue(libraryVariableSet.GetID())
}

func CreateLibraryVariableSet(data *LibraryVariableSetResourceModel) *variables.LibraryVariableSet {
	libraryVariableSet := variables.NewLibraryVariableSet(data.Name.ValueString())
	libraryVariableSet.ID = data.ID.ValueString()
	libraryVariableSet.Description = data.Description.ValueString()
	libraryVariableSet.SpaceID = data.SpaceID.ValueString()

	if len(data.Template.Elements()) > 0 {
		for _, tfTemplate := range data.Template.Elements() {
			template := expandActionTemplateParameter(tfTemplate.(types.Object).Attributes())
			libraryVariableSet.Templates = append(libraryVariableSet.Templates, template)
		}
	}

	return libraryVariableSet
}

func FlattenTemplates(actionTemplateParameters []actiontemplates.ActionTemplateParameter) types.List {
	if actionTemplateParameters == nil {
		return types.ListNull(types.ObjectType{AttrTypes: TemplateObjectType()})
	}
	actionTemplateList := make([]attr.Value, 0, len(actionTemplateParameters))

	for _, actionTemplateParams := range actionTemplateParameters {
		attrs := map[string]attr.Value{
			"default_value":    types.StringValue(actionTemplateParams.DefaultValue.Value),
			"display_settings": flattenDisplaySettingsMap(actionTemplateParams.DisplaySettings),
			"help_text":        types.StringValue(actionTemplateParams.HelpText),
			"id":               types.StringValue(actionTemplateParams.ID),
			"label":            types.StringValue(actionTemplateParams.Label),
			"name":             types.StringValue(actionTemplateParams.Name),
		}
		actionTemplateList = append(actionTemplateList, types.ObjectValueMust(TemplateObjectType(), attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: TemplateObjectType()}, actionTemplateList)
}

func flattenDisplaySettingsMap(displaySettings map[string]string) types.Map {
	if len(displaySettings) == 0 {
		return types.MapNull(types.ObjectType{AttrTypes: TemplateObjectType()})
	}

	flattenedDisplaySettings := make(map[string]attr.Value, len(displaySettings))
	for key, displaySetting := range displaySettings {
		flattenedDisplaySettings[key] = types.StringValue(displaySetting)
	}

	displaySettingsMapValue, _ := types.MapValue(types.StringType, flattenedDisplaySettings)
	return displaySettingsMapValue
}

func FlattenTemplateIds(actionTemplateParameters []actiontemplates.ActionTemplateParameter) types.Map {
	if actionTemplateParameters == nil {
		return types.MapNull(types.ObjectType{AttrTypes: TemplateObjectType()})
	}

	templateIds := map[string]attr.Value{}
	for _, template := range actionTemplateParameters {
		templateIds[template.Name] = types.StringValue(template.GetID())
	}

	templateIdsValues, _ := types.MapValue(types.StringType, templateIds)
	return templateIdsValues
}
