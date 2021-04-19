package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandVariable(d *schema.ResourceData) octopusdeploy.Variable {
	name := d.Get("name").(string)

	variable := octopusdeploy.NewVariable(name)

	if v, ok := d.GetOk("description"); ok {
		variable.Description = v.(string)
	}

	if v, ok := d.GetOk("is_editable"); ok {
		variable.IsEditable = v.(bool)
	}

	if v, ok := d.GetOk("is_sensitive"); ok {
		variable.IsSensitive = v.(bool)
	}

	if v, ok := d.GetOk("type"); ok {
		variable.Type = v.(string)
	}

	if v, ok := d.GetOk("scope"); ok {
		variable.Scope = expandVariableScope(v)
	}

	if variable.IsSensitive {
		variable.Type = "Sensitive"
		variable.Value = d.Get("sensitive_value").(string)
	} else {
		variable.Value = d.Get("value").(string)
	}

	variable.ID = d.Id()

	varPrompt, ok := d.GetOk("prompt")
	if ok {
		tfPromptSettings := varPrompt.(*schema.Set)
		if len(tfPromptSettings.List()) == 1 {
			tfPromptList := tfPromptSettings.List()[0].(map[string]interface{})
			newPrompt := octopusdeploy.VariablePromptOptions{
				Description: tfPromptList["description"].(string),
				Label:       tfPromptList["label"].(string),
				Required:    tfPromptList["is_required"].(bool),
			}
			variable.Prompt = &newPrompt
		}
	}

	return *variable
}

func getVariableDataSchema() map[string]*schema.Schema {
	dataSchema := getVariableSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"id": getDataSchemaID(),
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
		"is_editable": {
			Default:     true,
			Description: "Indicates whether or not this variable is considered editable.",
			Optional:    true,
			Type:        schema.TypeBool,
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
		"scope": {
			Elem:     &schema.Resource{Schema: getVariableScopeSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
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

func setVariable(ctx context.Context, d *schema.ResourceData, variable octopusdeploy.Variable) error {
	d.Set("description", variable.Description)
	d.Set("is_editable", variable.IsEditable)
	d.Set("is_sensitive", variable.IsSensitive)
	d.Set("name", variable.Name)
	d.Set("type", variable.Type)

	if variable.IsSensitive {
		d.Set("value", nil)
	} else {
		d.Set("value", variable.Value)
	}

	if err := d.Set("scope", flattenVariableScope(variable.Scope)); err != nil {
		return fmt.Errorf("error setting scope: %s", err)
	}

	d.SetId(variable.GetID())

	return nil
}
