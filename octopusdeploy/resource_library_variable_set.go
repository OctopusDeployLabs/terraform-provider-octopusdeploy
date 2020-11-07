package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLibraryVariableSet() *schema.Resource {
	resourceLibraryVariableSetImporter := &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
	resourceLibraryVariableSetSchema := map[string]*schema.Schema{
		constDescription: {
			Optional: true,
			Type:     schema.TypeString,
		},
		constName: {
			Required: true,
			Type:     schema.TypeString,
		},
		constSpaceID: {
			Computed: true,
			Type:     schema.TypeString,
		},
		constTemplate: {
			Optional: true,
			Elem: &schema.Resource{
				Schema: getTemplateSchema(),
			},
			Type: schema.TypeList,
		},
		constVariableSetID: {
			Computed: true,
			Type:     schema.TypeString,
		},
	}

	return &schema.Resource{
		CreateContext: resourceLibraryVariableSetCreate,
		DeleteContext: resourceLibraryVariableSetDelete,
		Importer:      resourceLibraryVariableSetImporter,
		ReadContext:   resourceLibraryVariableSetRead,
		Schema:        resourceLibraryVariableSetSchema,
		UpdateContext: resourceLibraryVariableSetUpdate,
	}
}

func getTemplateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		constDefaultValue: {
			Optional: true,
			Type:     schema.TypeString,
		},
		constDisplaySettings: {
			Optional: true,
			Type:     schema.TypeMap,
		},
		constHelpText: {
			Optional: true,
			Type:     schema.TypeString,
		},
		constID: {
			Computed: true,
			Type:     schema.TypeString,
		},
		constLabel: {
			Optional: true,
			Type:     schema.TypeString,
		},
		constName: {
			Required: true,
			Type:     schema.TypeString,
		},
	}
}

func resourceLibraryVariableSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	libraryVariableSet := buildLibraryVariableSetResource(d)

	client := m.(*octopusdeploy.Client)
	createdLibraryVariableSet, err := client.LibraryVariableSets.Add(libraryVariableSet)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenLibraryVariableSet(ctx, d, createdLibraryVariableSet)
	return nil
}

func buildLibraryVariableSetResource(d *schema.ResourceData) *octopusdeploy.LibraryVariableSet {
	var name string
	if v, ok := d.GetOk(constName); ok {
		name = v.(string)
	}

	libraryVariableSet := octopusdeploy.NewLibraryVariableSet(name)

	if v, ok := d.GetOk(constDescription); ok {
		libraryVariableSet.Description = v.(string)
	}

	if attr, ok := d.GetOk(constTemplate); ok {
		tfTemplates := attr.([]interface{})

		for _, tfTemplate := range tfTemplates {
			template := buildTemplateResource(tfTemplate.(map[string]interface{}))
			libraryVariableSet.Templates = append(libraryVariableSet.Templates, template)
		}
	}

	return libraryVariableSet
}

func buildTemplateResource(tfTemplate map[string]interface{}) *octopusdeploy.ActionTemplateParameter {
	actionTemplateParameter := octopusdeploy.NewActionTemplateParameter()

	propertyValue := octopusdeploy.PropertyValue(tfTemplate[constDefaultValue].(string))
	actionTemplateParameter.DefaultValue = &octopusdeploy.PropertyValueResource{
		PropertyValue: &propertyValue,
	}
	actionTemplateParameter.DisplaySettings = flattenDisplaySettings(tfTemplate[constDisplaySettings].(map[string]interface{}))
	actionTemplateParameter.HelpText = tfTemplate[constHelpText].(string)
	actionTemplateParameter.ID = tfTemplate[constID].(string)
	actionTemplateParameter.Label = tfTemplate[constLabel].(string)
	actionTemplateParameter.Name = tfTemplate[constName].(string)

	return actionTemplateParameter
}

func resourceLibraryVariableSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	libraryVariableSet, err := client.LibraryVariableSets.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	flattenLibraryVariableSet(ctx, d, libraryVariableSet)
	return nil
}

func resourceLibraryVariableSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	libraryVariableSet := buildLibraryVariableSetResource(d)
	libraryVariableSet.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	updatedLibraryVariableSet, err := client.LibraryVariableSets.Update(libraryVariableSet)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenLibraryVariableSet(ctx, d, updatedLibraryVariableSet)
	return nil
}

func resourceLibraryVariableSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.LibraryVariableSets.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)
	return nil
}

func flattenActionTemplateParameters(actionTemplateParameters []*octopusdeploy.ActionTemplateParameter) []interface{} {
	flattenedActionTemplateParameters := make([]interface{}, 0)
	for _, actionTemplateParameter := range actionTemplateParameters {
		a := make(map[string]interface{})
		a[constDefaultValue] = actionTemplateParameter.DefaultValue.PropertyValue
		a[constDisplaySettings] = actionTemplateParameter.DisplaySettings
		a[constHelpText] = actionTemplateParameter.HelpText
		a[constID] = actionTemplateParameter.ID
		a[constLabel] = actionTemplateParameter.Label
		a[constName] = actionTemplateParameter.Name
		flattenedActionTemplateParameters = append(flattenedActionTemplateParameters, a)
	}
	return flattenedActionTemplateParameters
}

func flattenDisplaySettings(displaySettings map[string]interface{}) map[string]string {
	flattenedDisplaySettings := make(map[string]string, len(displaySettings))
	for key, displaySetting := range displaySettings {
		flattenedDisplaySettings[key] = displaySetting.(string)
	}
	return flattenedDisplaySettings
}

func flattenLibraryVariableSet(ctx context.Context, d *schema.ResourceData, libraryVariableSet *octopusdeploy.LibraryVariableSet) {
	d.Set(constDescription, libraryVariableSet.Description)
	d.Set(constName, libraryVariableSet.Name)
	d.Set(constSpaceID, libraryVariableSet.SpaceID)
	d.Set(constVariableSetID, libraryVariableSet.VariableSetID)
	d.Set(constTemplate, flattenActionTemplateParameters(libraryVariableSet.Templates))

	d.SetId(libraryVariableSet.GetID())
}
