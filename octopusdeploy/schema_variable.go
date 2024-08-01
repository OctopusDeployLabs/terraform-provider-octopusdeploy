package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getVariableDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Computed:    true,
			Description: "The identifier of the variable to find.",
			Type:        schema.TypeString,
		},
		"name": {
			Required:    true,
			Description: "The name of variable to find.",
			Type:        schema.TypeString,
		},
		"owner_id": {
			Required:    true,
			Description: "Owner ID for the variable to find.",
			Type:        schema.TypeString,
		},
		"scope": {
			Description: "As variable names can appear more than once under different scopes, a VariableScope must also be provided",
			Elem:        &schema.Resource{Schema: getVariableScopeSchema()},
			Required:    true,
			Type:        schema.TypeList,
			MaxItems:    1,
		},
		"space_id": getQuerySpaceID(),
		"description": {
			Description: "The description of this variable.",
			Computed:    true,
			Type:        schema.TypeString,
		},
		"is_editable": {
			Description: "Indicates whether or not this variable is considered editable.",
			Computed:    true,
			Type:        schema.TypeBool,
		},
		"is_sensitive": {
			Description: "Indicates whether or not this resource is considered sensitive and should be kept secret.",
			Computed:    true,
			Type:        schema.TypeBool,
		},
		"type": {
			Description: "The type of variable represented by this resource. Valid types are `AmazonWebServicesAccount`, `AzureAccount`, `GoogleCloudAccount`, `Certificate`, `Sensitive`, `String`, or `WorkerPool`.",
			Computed:    true,
			Type:        schema.TypeString,
		},
		"sensitive_value": {
			Computed:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"value": {
			Computed: true,
			Type:     schema.TypeString,
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
			Computed: true,
			Type:     schema.TypeList,
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
