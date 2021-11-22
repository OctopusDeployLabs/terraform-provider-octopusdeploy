package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandScriptModule(d *schema.ResourceData) *octopusdeploy.ScriptModule {
	name := d.Get("name").(string)

	scriptModule := octopusdeploy.NewScriptModule(name)
	scriptModule.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		scriptModule.Description = v.(string)
	}

	if v, ok := d.GetOk("script"); ok {
		scripts := v.(*schema.Set).List()
		for _, script := range scripts {
			rawScript := script.(map[string]interface{})

			if rawScript["body"] != nil {
				scriptModule.ScriptBody = rawScript["body"].(string)
			}

			if rawScript["syntax"] != nil {
				scriptModule.Syntax = rawScript["syntax"].(string)
			}
		}
	}

	if v, ok := d.GetOk("space_id"); ok {
		scriptModule.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("variable_set_id"); ok {
		scriptModule.VariableSetID = v.(string)
	}

	return scriptModule
}

func flattenScript(scriptModule *octopusdeploy.ScriptModule) []interface{} {
	if scriptModule == nil {
		return nil
	}

	flattenedScriptModules := make([]interface{}, 1)
	flattenedScriptModules[0] = map[string]interface{}{
		"body":   scriptModule.ScriptBody,
		"syntax": scriptModule.Syntax,
	}

	return flattenedScriptModules
}

func flattenScriptModule(scriptModule *octopusdeploy.ScriptModule) map[string]interface{} {
	if scriptModule == nil {
		return nil
	}

	return map[string]interface{}{
		"description":     scriptModule.Description,
		"id":              scriptModule.GetID(),
		"name":            scriptModule.Name,
		"script":          flattenScript(scriptModule),
		"space_id":        scriptModule.SpaceID,
		"variable_set_id": scriptModule.VariableSetID,
	}
}

func getScriptModuleDataSchema() map[string]*schema.Schema {
	dataSchema := getScriptModuleSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"id":  getDataSchemaID(),
		"ids": getQueryIDs(),
		"script_modules": {
			Computed:    true,
			Description: "A list of script modules that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
	}
}

func getScriptModuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema(),
		"id":          getIDSchema(),
		"name":        getNameSchema(true),
		"script": {
			Description: "The script associated with this script module.",
			Required:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"body": {
						Description: "The body of this script module.",
						Required:    true,
						Type:        schema.TypeString,
					},
					"syntax": {
						Description: "The syntax of the script. Valid types are `Bash`, `CSharp`, `FSharp`, `PowerShell`, or `Python`.",
						Required:    true,
						Type:        schema.TypeString,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
							"Bash",
							"CSharp",
							"FSharp",
							"PowerShell",
							"Python",
						}, false)),
					},
				},
			},
			MaxItems: 1,
			MinItems: 1,
			Type:     schema.TypeSet,
		},
		"space_id": getSpaceIDSchema(),
		"variable_set_id": {
			Computed:    true,
			Description: "The variable set ID for this script module.",
			Optional:    true,
			Type:        schema.TypeString,
		},
	}
}

func setScriptModule(ctx context.Context, d *schema.ResourceData, scriptModule *octopusdeploy.ScriptModule) error {
	d.Set("description", scriptModule.Description)
	d.Set("name", scriptModule.Name)

	if err := d.Set("script", flattenScript(scriptModule)); err != nil {
		return fmt.Errorf("error setting script: %s", err)
	}

	d.Set("space_id", scriptModule.SpaceID)
	d.Set("variable_set_id", scriptModule.VariableSetID)

	d.SetId(scriptModule.GetID())

	return nil
}
