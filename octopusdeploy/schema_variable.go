package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandVariable(d *schema.ResourceData) *variables.Variable {
	name := d.Get("name").(string)

	variable := variables.NewVariable(name)

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

	if v, ok := d.GetOk("space_id"); ok {
		variable.SpaceID = v.(string)
	}

	if variable.IsSensitive {
		variable.Type = "Sensitive"
		variable.Value = d.Get("sensitive_value").(string)
	} else {
		variable.Value = d.Get("value").(string)
	}

	if v, ok := d.GetOk("prompt"); ok {
		variable.Prompt = expandPromptedVariableSettings(v)
	}

	variable.ID = d.Id()

	return variable
}

func getVariableDataSchema() map[string]*schema.Schema {
	dataSchema := getVariableSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"id":       getDataSchemaID(),
		"space_id": getQuerySpaceID(),
		"ids":      getQueryIDs(),
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
		"description": getDescriptionSchema("variable"),
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
		"owner_id": {
			ConflictsWith: []string{"project_id"},
			Optional:      true,
			Type:          schema.TypeString,
		},
		"pgp_key": {
			ForceNew:  true,
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"project_id": {
			ConflictsWith: []string{"owner_id"},
			Deprecated:    "This attribute is deprecated; please use owner_id instead.",
			Optional:      true,
			Type:          schema.TypeString,
		},
		"prompt": {
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": getDescriptionSchema("variable prompt option"),
					"display_settings": {
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"control_type": {
									Description: "The type of control for rendering this prompted variable. Valid types are `SingleLineText`, `MultiLineText`, `Checkbox`, `Select`.",
									Required:    true,
									Type:        schema.TypeString,
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
										"Checkbox",
										"MultiLineText",
										"Select",
										"SingleLineText",
									}, false)),
								},
								"select_option": {
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"value": {
												Description: "The select value",
												Required:    true,
												Type:        schema.TypeString,
											},
											"display_name": {
												Description: "The display name for the select value",
												Required:    true,
												Type:        schema.TypeString,
											},
										},
									},
									Description: "If the `control_type` is `Select`, then this value defines an option.",
									Optional:    true,
									Type:        schema.TypeList,
								},
							},
						},
						MaxItems: 1,
						Optional: true,
						Type:     schema.TypeList,
					},
					"is_required": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"label": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
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
		"type":     getVariableTypeSchema(),
		"space_id": getSpaceIDSchema(),
		"value": {
			ConflictsWith: []string{"sensitive_value"},
			Optional:      true,
			Type:          schema.TypeString,
		},
	}
}

func setVariable(ctx context.Context, d *schema.ResourceData, variable *variables.Variable) error {
	if d == nil || variable == nil {
		return fmt.Errorf("error setting scope")
	}

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

	if err := d.Set("prompt", flattenPromptedVariableSettings(variable.Prompt)); err != nil {
		return fmt.Errorf("error setting prompted config: %s", err)
	}

	if err := d.Set("scope", flattenVariableScope(variable.Scope)); err != nil {
		return fmt.Errorf("error setting scope: %s", err)
	}

	d.SetId(variable.GetID())

	return nil
}
