package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	types "github.com/hashicorp/terraform-plugin-framework/types"
)

type LibraryVariableSetResourceModel struct {
	Description   types.String `tfsdk:"description"`
	Name          types.String `tfsdk:"name"`
	SpaceID       types.String `tfsdk:"space_id"`
	Template      types.List   `tfsdk:"template"`
	TemplateIds   types.Map    `tfsdk:"template_ids"`
	VariableSetId types.String `tfsdk:"variable_set_id"`

	ResourceModel
}

func GetLibraryVariableSetDataSourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing library variable sets.",
		Attributes: map[string]datasourceSchema.Attribute{
			"content_type": datasourceSchema.StringAttribute{
				Description: "A filter to search by content type.",
				Optional:    true,
			},
			"id":           GetIdDatasourceSchema(true),
			"space_id":     GetSpaceIdDatasourceSchema("library variable set", false),
			"ids":          GetQueryIDsDatasourceSchema(),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take":         GetQueryTakeDatasourceSchema(),
			"library_variable_sets": datasourceSchema.ListNestedAttribute{
				Computed: true,
				Optional: false,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: GetLibraryVariableSetObjectDatasourceSchema(),
				},
			},
		},
	}
}

func GetLibraryVariableSetObjectDatasourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"description": GetReadonlyDescriptionDatasourceSchema("library variable set"),
		"id":          GetIdDatasourceSchema(true),
		"name":        GetReadonlyNameDatasourceSchema(),
		"space_id":    GetSpaceIdDatasourceSchema("library variable set", true),
		"template_ids": datasourceSchema.MapAttribute{
			ElementType: types.StringType,
			Computed:    true,
		},
		"template": datasourceSchema.ListAttribute{
			Computed:    true,
			ElementType: types.ObjectType{AttrTypes: TemplateObjectType()},
		},
		"variable_set_id": datasourceSchema.StringAttribute{
			Computed: true,
		},
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
			"template": resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: GetActionTemplateParameterSchema(),
				},
			},
		},
	}
}

func MapFromLibraryVariableSet(data *LibraryVariableSetResourceModel, spaceId string, libraryVariableSet *variables.LibraryVariableSet) {
	data.Description = types.StringValue(libraryVariableSet.Description)
	data.Name = types.StringValue(libraryVariableSet.Name)
	data.VariableSetId = types.StringValue(libraryVariableSet.VariableSetID)
	data.SpaceID = types.StringValue(spaceId)

	data.Template = FlattenTemplates(libraryVariableSet.Templates)
	data.TemplateIds = FlattenTemplateIds(libraryVariableSet.Templates)

	data.ID = types.StringValue(libraryVariableSet.GetID())
}

func MapToLibraryVariableSet(data *LibraryVariableSetResourceModel) *variables.LibraryVariableSet {
	libraryVariableSet := variables.NewLibraryVariableSet(data.Name.ValueString())
	libraryVariableSet.ID = data.ID.ValueString()
	libraryVariableSet.Description = data.Description.ValueString()
	libraryVariableSet.SpaceID = data.SpaceID.ValueString()

	libraryVariableSet.Templates = ExpandActionTemplateParameters(data.Template)

	return libraryVariableSet
}

func FlattenTemplates(actionTemplateParameters []actiontemplates.ActionTemplateParameter) types.List {
	if len(actionTemplateParameters) == 0 {
		return types.ListValueMust(types.ObjectType{AttrTypes: TemplateObjectType()}, []attr.Value{})
	}
	actionTemplateList := make([]attr.Value, 0, len(actionTemplateParameters))

	for _, actionTemplateParams := range actionTemplateParameters {
		attrs := map[string]attr.Value{
			"default_value":    types.StringValue(actionTemplateParams.DefaultValue.Value),
			"display_settings": flattenDisplaySettingsMap(actionTemplateParams.DisplaySettings),
			"help_text":        util.Ternary(actionTemplateParams.HelpText != "", types.StringValue(actionTemplateParams.HelpText), types.StringValue("")),
			"id":               types.StringValue(actionTemplateParams.GetID()),
			"label":            util.Ternary(actionTemplateParams.Label != "", types.StringValue(actionTemplateParams.Label), types.StringNull()),
			"name":             types.StringValue(actionTemplateParams.Name),
		}
		actionTemplateList = append(actionTemplateList, types.ObjectValueMust(TemplateObjectType(), attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: TemplateObjectType()}, actionTemplateList)
}

func flattenDisplaySettingsMap(displaySettings map[string]string) types.Map {
	if len(displaySettings) == 0 {
		return types.MapNull(types.StringType)
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
		return types.MapNull(types.StringType)
	}

	templateIds := map[string]attr.Value{}
	for _, template := range actionTemplateParameters {
		templateIds[template.Name] = types.StringValue(template.ID)
	}

	templateIdsValues, _ := types.MapValue(types.StringType, templateIds)
	return templateIdsValues
}
