package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandVariable(d *schema.ResourceData) *octopusdeploy.Variable {
	varName := d.Get("name").(string)
	varType := d.Get("type").(string)

	var varDesc, varValue string
	var varSensitive bool

	if varDescInterface, ok := d.GetOk("description"); ok {
		varDesc = varDescInterface.(string)
	}

	if varSensitiveInterface, ok := d.GetOk("is_sensitive"); ok {
		varSensitive = varSensitiveInterface.(bool)
	}

	if varSensitive {
		varValue = d.Get("sensitive_value").(string)
	} else {
		varValue = d.Get(constValue).(string)
	}

	varScopeInterface := tfVariableScopetoODVariableScope(d)

	newVar := octopusdeploy.NewVariable(varName, varType, varValue, varDesc, varScopeInterface, varSensitive)
	newVar.ID = d.Id()

	varPrompt, ok := d.GetOk(constPrompt)
	if ok {
		tfPromptSettings := varPrompt.(*schema.Set)
		if len(tfPromptSettings.List()) == 1 {
			tfPromptList := tfPromptSettings.List()[0].(map[string]interface{})
			newPrompt := octopusdeploy.VariablePromptOptions{
				Description: tfPromptList["description"].(string),
				Label:       tfPromptList["label"].(string),
				Required:    tfPromptList["is_required"].(bool),
			}
			newVar.Prompt = &newPrompt
		}
	}

	return newVar
}

func getVariableDataSchema() map[string]*schema.Schema {
	dataSchema := getVariableSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"id": getIDDataSchema(),
		"variables": {
			Computed:    true,
			Description: "A list of variables that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func getVariableSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema(),
		"encrypted_value": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"is_sensitive": getIsSensitiveSchema(),
		"key_fingerprint": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": getNameSchema(true),
		"pgp_key": {
			ForceNew:  true,
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"project_id": {
			Required: true,
			Type:     schema.TypeString,
		},
		"prompt": {
			Elem:     &schema.Resource{Schema: getVariablePromptOptionsSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"scope": getVariableScopeSchema(),
		"sensitive_value": {
			ConflictsWith: []string{"value"},
			Optional:      true,
			Sensitive:     true,
			Type:          schema.TypeString,
		},
		"type": getVariableTypeSchema(),
		"value": {
			ConflictsWith: []string{"sensitive_value"},
			Optional:      true,
			Type:          schema.TypeString,
		},
	}
}
