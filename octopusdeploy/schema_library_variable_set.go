package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandLibraryVariableSet(d *schema.ResourceData) *variables.LibraryVariableSet {
	name := d.Get("name").(string)

	libraryVariableSet := variables.NewLibraryVariableSet(name)
	libraryVariableSet.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		libraryVariableSet.Description = v.(string)
	}

	if attr, ok := d.GetOk("template"); ok {
		tfTemplates := attr.([]interface{})

		for _, tfTemplate := range tfTemplates {
			template := expandActionTemplateParameter(tfTemplate.(map[string]interface{}))
			libraryVariableSet.Templates = append(libraryVariableSet.Templates, template)
		}
	}

	return libraryVariableSet
}

func flattenLibraryVariableSet(libraryVariableSet *variables.LibraryVariableSet) map[string]interface{} {
	if libraryVariableSet == nil {
		return nil
	}

	return map[string]interface{}{
		"description":     libraryVariableSet.Description,
		"id":              libraryVariableSet.GetID(),
		"name":            libraryVariableSet.Name,
		"space_id":        libraryVariableSet.SpaceID,
		"template":        flattenActionTemplateParameters(libraryVariableSet.Templates),
		"variable_set_id": libraryVariableSet.VariableSetID,
	}
}

func getLibraryVariableSetDataSchema() map[string]*schema.Schema {
	dataSchema := getLibraryVariableSetSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"content_type": getQueryContentType(),
		"id":           getDataSchemaID(),
		"ids":          getQueryIDs(),
		"library_variable_sets": {
			Computed:    true,
			Description: "A list of library variable sets that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
	}
}

func getLibraryVariableSetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema("library variable set"),
		"id":          getIDSchema(),
		"name":        getNameSchema(true),
		"space_id":    getSpaceIDSchema(),
		"template": {
			Optional: true,
			Elem:     &schema.Resource{Schema: getActionTemplateParameterSchema()},
			Type:     schema.TypeList,
		},
		"variable_set_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}

func setLibraryVariableSet(ctx context.Context, d *schema.ResourceData, libraryVariableSet *variables.LibraryVariableSet) error {
	d.Set("description", libraryVariableSet.Description)
	d.Set("name", libraryVariableSet.Name)
	d.Set("space_id", libraryVariableSet.SpaceID)
	d.Set("variable_set_id", libraryVariableSet.VariableSetID)

	if err := d.Set("template", flattenActionTemplateParameters(libraryVariableSet.Templates)); err != nil {
		return fmt.Errorf("error setting template: %s", err)
	}

	d.SetId(libraryVariableSet.GetID())

	return nil
}
