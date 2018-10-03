package octopusdeploy

import (
	"fmt"

	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariableCreate,
		Read:   resourceVariableRead,
		Update: resourceVariableUpdate,
		Delete: resourceVariableDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateValueFunc([]string{
					"String",
					"Sensitive",
					"Certificate",
					"AmazonWebServicesAccount",
				}),
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": schemaVariableScope,
			"is_sensitive": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceVariableRead(d *schema.ResourceData, m interface{}) error {
	octoMutex.Lock("atom-variable")
	defer octoMutex.Unlock("atom-variable")

	client := m.(*octopusdeploy.Client)

	variableID := d.Id()
	projectID := d.Get("project_id").(string)
	tfVar, err := client.Variable.GetByID(projectID, variableID)

	if err == octopusdeploy.ErrItemNotFound || tfVar == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Variable %s: %s", variableID, err.Error())
	}

	d.Set("name", tfVar.Name)
	d.Set("type", tfVar.Type)
	d.Set("value", tfVar.Value)
	d.Set("description", tfVar.Description)

	return nil
}

func buildVariableResource(d *schema.ResourceData) *octopusdeploy.Variable {
	varName := d.Get("name").(string)
	varType := d.Get("type").(string)
	varValue := d.Get("value").(string)

	var varDesc string
	var varSensitive bool

	varDescInterface, ok := d.GetOk("description")
	if ok {
		varDesc = varDescInterface.(string)
	}

	varSensitiveInterface, ok := d.GetOk("is_sensitive")
	if ok {
		varSensitive = varSensitiveInterface.(bool)
	}

	varScopeInterface := tfVariableScopetoODVariableScope(d)

	newVar := octopusdeploy.NewVariable(varName, varType, varValue, varDesc, varScopeInterface, varSensitive)
	newVar.Prompt.Required = false

	return newVar
}

func resourceVariableCreate(d *schema.ResourceData, m interface{}) error {
	octoMutex.Lock("atom-variable")
	defer octoMutex.Unlock("atom-variable")

	client := m.(*octopusdeploy.Client)
	projID := d.Get("project_id").(string)

	newVariable := buildVariableResource(d)
	tfVar, err := client.Variable.AddSingle(projID, newVariable)

	if err != nil {
		return fmt.Errorf("error creating variable %s: %s", newVariable.Name, err.Error())
	}

	for _, v := range tfVar.Variables {
		if v.Name == newVariable.Name && v.Type == newVariable.Type && v.Value == newVariable.Value && v.Description == newVariable.Description && v.IsSensitive == newVariable.IsSensitive {
			scopeMatches, _, err := client.Variable.MatchesScope(v.Scope, newVariable.Scope)
			if err != nil {
				return err
			}
			if scopeMatches {
				d.SetId(v.ID)
				return nil
			}
		}
	}

	d.SetId("")
	return fmt.Errorf("unable to locate variable in project %s", projID)
}

func resourceVariableUpdate(d *schema.ResourceData, m interface{}) error {
	octoMutex.Lock("atom-variable")
	defer octoMutex.Unlock("atom-variable")

	tfVar := buildVariableResource(d)
	tfVar.ID = d.Id() // set project struct ID so octopus knows which project to update

	client := m.(*octopusdeploy.Client)
	projID := d.Get("project_id").(string)

	updatedVars, err := client.Variable.UpdateSingle(projID, tfVar)

	if err != nil {
		return fmt.Errorf("error updating variable id %s: %s", d.Id(), err.Error())
	}

	for _, v := range updatedVars.Variables {
		if v.Name == tfVar.Name && v.Type == tfVar.Type && v.Value == tfVar.Value && v.Description == tfVar.Description && v.IsSensitive == tfVar.IsSensitive {
			scopeMatches, _, _ := client.Variable.MatchesScope(v.Scope, tfVar.Scope)
			if scopeMatches {
				d.SetId(v.ID)
				return nil
			}
		}
	}

	d.SetId("")
	return fmt.Errorf("unable to locate variable in project %s", projID)
}

func resourceVariableDelete(d *schema.ResourceData, m interface{}) error {
	octoMutex.Lock("atom-variable")
	defer octoMutex.Unlock("atom-variable")

	client := m.(*octopusdeploy.Client)
	projID := d.Get("project_id").(string)

	variableID := d.Id()

	_, err := client.Variable.DeleteSingle(projID, variableID)

	if err != nil {
		return fmt.Errorf("error deleting variable id %s: %s", variableID, err.Error())
	}

	d.SetId("")
	return nil
}
