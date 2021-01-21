package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandLibraryVariableSet(d *schema.ResourceData) *octopusdeploy.LibraryVariableSet {
	name := d.Get("name").(string)

	libraryVariableSet := octopusdeploy.NewLibraryVariableSet(name)
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

func getLibraryVariableSetDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": getNameSchema(true),
	}
}

func getLibraryVariableSetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema(),
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

func setLibraryVariableSet(ctx context.Context, d *schema.ResourceData, libraryVariableSet *octopusdeploy.LibraryVariableSet) error {
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
