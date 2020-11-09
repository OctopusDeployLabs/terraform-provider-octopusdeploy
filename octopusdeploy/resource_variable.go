package octopusdeploy

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var mutex = &sync.Mutex{}

func resourceVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVariableCreate,
		ReadContext:   resourceVariableRead,
		UpdateContext: resourceVariableUpdate,
		DeleteContext: resourceVariableDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVariableImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"type": {
				Required: true,
				Type:     schema.TypeString,
				ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
					"String",
					"Sensitive",
					"Certificate",
					"AmazonWebServicesAccount",
					"AzureAccount",
				}, false)),
			},
			"value": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"sensitive_value"},
			},
			"sensitive_value": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"value"},
				Sensitive:     true,
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"scope": schemaVariableScope,
			"is_sensitive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"prompt": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"label": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_required": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"pgp_key": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"key_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVariableImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importStrings := strings.Split(d.Id(), ":")
	if len(importStrings) != 2 {
		return nil, fmt.Errorf("octopusdeploy_variable import must be in the form of ProjectID:VariableID (e.g. Projects-62:0906031f-68ba-4a15-afaa-657c1564e07b")
	}

	d.Set("project_id", importStrings[0])
	d.SetId(importStrings[1])

	return []*schema.ResourceData{d}, nil
}

func resourceVariableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	projectID := d.Get("project_id").(string)

	client := m.(*octopusdeploy.Client)
	resource, err := client.Variables.GetByID(projectID, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if resource == nil {
		d.SetId("")
		return nil
	}

	logResource(constVariable, m)

	d.Set("name", resource.Name)
	d.Set("type", resource.Type)

	isSensitive := d.Get(constIsSensitive).(bool)
	if isSensitive {
		d.Set(constValue, nil)
	} else {
		d.Set(constValue, resource.Value)
	}

	d.Set("description", resource.Description)

	return nil
}

func buildVariableResource(d *schema.ResourceData) *octopusdeploy.Variable {
	varName := d.Get(constName).(string)
	varType := d.Get(constType).(string)

	var varDesc, varValue string
	var varSensitive bool

	if varDescInterface, ok := d.GetOk("description"); ok {
		varDesc = varDescInterface.(string)
	}

	if varSensitiveInterface, ok := d.GetOk(constIsSensitive); ok {
		varSensitive = varSensitiveInterface.(bool)
	}

	if varSensitive {
		varValue = d.Get(constSensitiveValue).(string)
	} else {
		varValue = d.Get(constValue).(string)
	}

	varScopeInterface := tfVariableScopetoODVariableScope(d)

	newVar := octopusdeploy.NewVariable(varName, varType, varValue, varDesc, varScopeInterface, varSensitive)

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

func resourceVariableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()
	if err := validateVariable(d); err != nil {
		return diag.FromErr(err)
	}

	projID := d.Get("project_id").(string)
	newVariable := buildVariableResource(d)

	client := m.(*octopusdeploy.Client)
	tfVar, err := client.Variables.AddSingle(projID, newVariable)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range tfVar.Variables {
		if v.Name == newVariable.Name && v.Type == newVariable.Type && (v.IsSensitive || v.Value == newVariable.Value) && v.Description == newVariable.Description && v.IsSensitive == newVariable.IsSensitive {
			scopeMatches, _, err := client.Variables.MatchesScope(v.Scope, newVariable.Scope)
			if err != nil {
				return diag.FromErr(err)
			}
			if scopeMatches {
				d.SetId(v.ID)
				return nil
			}
		}
	}

	d.SetId("")
	return diag.Errorf("unable to locate variable in project %s", projID)
}

func resourceVariableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	if err := validateVariable(d); err != nil {
		return diag.FromErr(err)
	}

	tfVar := buildVariableResource(d)
	tfVar.ID = d.Id() // set ID so Octopus API knows which variable to update
	projID := d.Get("project_id").(string)

	client := m.(*octopusdeploy.Client)
	updatedVars, err := client.Variables.UpdateSingle(projID, tfVar)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range updatedVars.Variables {
		if v.Name == tfVar.Name && v.Type == tfVar.Type && (v.IsSensitive || v.Value == tfVar.Value) && v.Description == tfVar.Description && v.IsSensitive == tfVar.IsSensitive {
			scopeMatches, _, _ := client.Variables.MatchesScope(v.Scope, tfVar.Scope)
			if scopeMatches {
				d.SetId(v.ID)
				return nil
			}
		}
	}

	d.SetId("")
	return diag.Errorf("unable to locate variable in project %s", projID)
}

func resourceVariableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	projID := d.Get("project_id").(string)
	variableID := d.Id()

	client := m.(*octopusdeploy.Client)
	_, err := client.Variables.DeleteSingle(projID, variableID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// Validating is done in its own function as we need to compare options once the entire
// schema has been parsed, which as far as I can tell we can't do in a normal validation
// function.
func validateVariable(d *schema.ResourceData) error {
	tfSensitive := d.Get(constIsSensitive).(bool)
	tfType := d.Get(constType).(string)

	if tfSensitive && tfType != "Sensitive" {
		return fmt.Errorf("when is_sensitive is set to true, type needs to be 'Sensitive'")
	}

	if !tfSensitive && tfType == "Sensitive" {
		return fmt.Errorf("when type is set to 'Sensitive', is_sensitive needs to be true")
	}

	return nil
}
