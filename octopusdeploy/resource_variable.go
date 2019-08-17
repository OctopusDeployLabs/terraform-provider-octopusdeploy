package octopusdeploy

import (
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/encryption"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariableCreate,
		Read:   resourceVariableRead,
		Update: resourceVariableUpdate,
		Delete: resourceVariableDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVariableImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateValueFunc([]string{
					"String",
					"Sensitive",
					"Certificate",
					"AmazonWebServicesAccount",
					"AzureAccount",
				}),
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
				Type:     schema.TypeString,
				Optional: true,
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
						"label": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"required": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"pgp_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func resourceVariableRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	variableID := d.Id()
	projectID := d.Get("project_id").(string)
	isSensitive := d.Get("is_sensitive").(bool)
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
	if isSensitive {
		d.Set("value", nil)
	} else {
		d.Set("value", tfVar.Value)
	}
	d.Set("description", tfVar.Description)

	return nil
}

func buildVariableResource(d *schema.ResourceData) *octopusdeploy.Variable {
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
		varValue = d.Get("value").(string)
	}

	varScopeInterface := tfVariableScopetoODVariableScope(d)

	newVar := octopusdeploy.NewVariable(varName, varType, varValue, varDesc, varScopeInterface, varSensitive)

	varPrompt, ok := d.GetOk("prompt")
	if ok {
		tfPromptSettings := varPrompt.(*schema.Set)
		if len(tfPromptSettings.List()) == 1 {
			tfPromptList := tfPromptSettings.List()[0].(map[string]interface{})
			newPrompt := octopusdeploy.VariablePromptOptions{
				Description: tfPromptList["description"].(string),
				Label:       tfPromptList["label"].(string),
				Required:    tfPromptList["required"].(bool),
			}
			newVar.Prompt = &newPrompt
		}
	}

	return newVar
}

func encryptSensitiveValue(d *schema.ResourceData, ov *octopusdeploy.Variable) (isEncrypted bool, keyFingerprint, encryptedValue string, err error) {
	if isSensitive := d.Get("is_sensitive").(bool); !isSensitive {
		return isEncrypted, keyFingerprint, encryptedValue, nil
	}

	if v, ok := d.GetOk("pgp_key"); ok {
		pgpKey := strings.TrimSpace(v.(string))

		encryptionKey, err := encryption.RetrieveGPGKey(pgpKey)
		if err != nil {
			return isEncrypted, keyFingerprint, encryptedValue, fmt.Errorf("error retrieving PGP Key during Sensitive Variable (%s) creation: %s", ov.Name, err)
		}

		keyFingerprint, encryptedValue, err = encryption.EncryptValue(encryptionKey, ov.Value, "Sensitive Value")
		if err != nil {
			return isEncrypted, keyFingerprint, encryptedValue, fmt.Errorf("error encrypting value during Sensitive Variable (%s) creation: %s", ov.Name, err)
		}
		isEncrypted = true

	} else {
		return isEncrypted, keyFingerprint, encryptedValue, nil
	}

	return isEncrypted, keyFingerprint, encryptedValue, nil
}

func resourceVariableCreate(d *schema.ResourceData, m interface{}) error {
	octoMutex.Lock("atom-variable")
	defer octoMutex.Unlock("atom-variable")
	if err := validateVariable(d); err != nil {
		return err
	}

	client := m.(*octopusdeploy.Client)
	projID := d.Get("project_id").(string)

	newVariable := buildVariableResource(d)
	tfVar, err := client.Variable.AddSingle(projID, newVariable)

	if err != nil {
		return fmt.Errorf("error creating variable %s: %s", newVariable.Name, err.Error())
	}

	isEncrypted, fingerprint, encryptedValue, err := encryptSensitiveValue(d, newVariable)
	if err != nil {
		return fmt.Errorf("Error encrypting sensitive value: %s", err)
	}

	for _, v := range tfVar.Variables {
		if v.Name == newVariable.Name && v.Type == newVariable.Type && (v.IsSensitive || v.Value == newVariable.Value) && v.Description == newVariable.Description && v.IsSensitive == newVariable.IsSensitive {
			scopeMatches, _, err := client.Variable.MatchesScope(v.Scope, newVariable.Scope)
			if err != nil {
				return err
			}
			if scopeMatches {
				d.SetId(v.ID)
				if isEncrypted {
					d.Set("key_fingerprint", fingerprint)
					d.Set("encrypted_value", encryptedValue)
				}
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
	if err := validateVariable(d); err != nil {
		return err
	}

	tfVar := buildVariableResource(d)
	tfVar.ID = d.Id() // set project struct ID so octopus knows which project to update

	client := m.(*octopusdeploy.Client)
	projID := d.Get("project_id").(string)

	updatedVars, err := client.Variable.UpdateSingle(projID, tfVar)
	if err != nil {
		return fmt.Errorf("error updating variable id %s: %s", d.Id(), err.Error())
	}

	isEncrypted, fingerprint, encryptedValue, err := encryptSensitiveValue(d, tfVar)
	if err != nil {
		return fmt.Errorf("Error encrypting sensitive value: %s", err)
	}

	for _, v := range updatedVars.Variables {
		if v.Name == tfVar.Name && v.Type == tfVar.Type && (v.IsSensitive || v.Value == tfVar.Value) && v.Description == tfVar.Description && v.IsSensitive == tfVar.IsSensitive {
			scopeMatches, _, _ := client.Variable.MatchesScope(v.Scope, tfVar.Scope)
			if scopeMatches {
				d.SetId(v.ID)
				if isEncrypted {
					d.Set("key_fingerprint", fingerprint)
					d.Set("encrypted_value", encryptedValue)
				}
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

// Validating is done in its own function as we need to compare options once the entire
// schema has been parsed, which as far as I can tell we can't do in a normal validation
// function.
func validateVariable(d *schema.ResourceData) error {
	tfSensitive := d.Get("is_sensitive").(bool)
	tfType := d.Get("type").(string)

	if tfSensitive && tfType != "Sensitive" {
		return fmt.Errorf("when is_sensitive is set to true, type needs to be 'Sensitive'")
	}

	if !tfSensitive && tfType == "Sensitive" {
		return fmt.Errorf("when type is set to 'Sensitive', is_sensitive needs to be true")
	}

	return nil
}
