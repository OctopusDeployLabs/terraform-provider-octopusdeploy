package schemas

import (
	"context"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	types "github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
)

type LibraryVariableSetResourceModel struct {
	Description types.String `tfsdk:"description"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	SpaceID     types.String `tfsdk:"space_id"`
	//LibraryVariableSets types.List   `tfsdk:"library_variable_sets"`
	Template      types.List   `tfsdk:"template"`
	TemplateIds   types.Map    `tfsdk:"template_ids"`
	VariableSetId types.String `tfsdk:"variable_set_id"`
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

func FlattenLibraryVariableSet(libraryVariableSet *variables.LibraryVariableSet) map[string]attr.Value {
	if libraryVariableSet == nil {
		return nil
	}

	templateIds := map[string]attr.Value{}
	if libraryVariableSet.Templates != nil {
		for _, template := range libraryVariableSet.Templates {
			templateIds[template.Name] = types.StringValue(template.GetID())
		}
	}
	templateIdsValues, _ := types.MapValue(types.StringType, templateIds)

	templateValues, _ := types.ListValueFrom(
		context.Background(),
		types.ObjectType{AttrTypes: templateObjectType()},
		[]any{flattenActionTemplateParameters(libraryVariableSet.Templates)})

	libraryVariableSetMap := map[string]attr.Value{
		"description":     types.StringValue(libraryVariableSet.Description),
		"id":              types.StringValue(libraryVariableSet.GetID()),
		"name":            types.StringValue(libraryVariableSet.Name),
		"space_id":        types.StringValue(libraryVariableSet.SpaceID),
		"template":        templateValues,
		"variable_set_id": types.StringValue(libraryVariableSet.VariableSetID),
		"template_ids":    templateIdsValues,
	}

	return libraryVariableSetMap
}

func GetLibraryVariableSetDataSourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Attributes: getLibraryVariableSetDataSchema(),
	}
}

func getLibraryVariableSetDataSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		//dataSchema := getLibraryVariableSetSchema()
		//octopusdeploy.setDataSchema(&dataSchema)

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
		"template_ids":    types.MapType{ElemType: types.StringType},
		"variable_set_id": types.StringType,

		//	Blocks: map[string]resourceSchema.Block{
		//	"project_groups": resourceSchema.ListNestedBlock{
		//		Description: "A list of project groups that match the filter(s).",
		//		NestedObject: resourceSchema.NestedBlockObject{
		//			Attributes: getActionTemplateParameterSchema(),
		//		},
		//	},
		//},
	}
}

func GetLibraryVariableSetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"description": GetDescriptionResourceSchema("library variable set"),
			"id":          GetIdResourceSchema(),
			"name":        GetNameResourceSchema(true),
			"space_id":    GetSpaceIdResourceSchema("library variable set"),
			// This field is based on the suggestion at
			// https://discuss.hashicorp.com/t/custom-provider-how-to-reference-computed-attribute-of-typemap-list-set-defined-as-nested-block/22898/2
			"template": resourceSchema.ListAttribute{
				Optional:    true,
				ElementType: types.ObjectType{AttrTypes: templateObjectType()},
			},
			"template_ids": resourceSchema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    false,
			},
			"variable_set_id": resourceSchema.StringAttribute{
				Computed: true,
			},
		},
		//Blocks: map[string]resourceSchema.Block{
		//	"project_groups": resourceSchema.ListNestedBlock{
		//		Description: "A list of project groups that match the filter(s).",
		//		NestedObject: resourceSchema.NestedBlockObject{
		//			Attributes: getActionTemplateParameterSchema(),
		//		},
		//	},
		//},
	}
}

func UpdateDataFromLibraryVariableSet(data *LibraryVariableSetResourceModel, spaceId string, libraryVariableSet *variables.LibraryVariableSet) {
	data.Description = types.StringValue(libraryVariableSet.Description)
	data.Name = types.StringValue(libraryVariableSet.Name)
	data.VariableSetId = types.StringValue(libraryVariableSet.VariableSetID)
	data.Description = types.StringValue(libraryVariableSet.Description)
	data.SpaceID = types.StringValue(spaceId)

	if len(libraryVariableSet.Templates) > 0 {
		data.Template, _ = types.ListValueFrom(
			context.Background(),
			types.ObjectType{AttrTypes: templateObjectType()},
			[]any{flattenActionTemplateParameters(libraryVariableSet.Templates)},
		)
	}

	data.ID = types.StringValue(libraryVariableSet.GetID())
}
