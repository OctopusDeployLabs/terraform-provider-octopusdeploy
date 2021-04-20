package octopusdeploy

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVariableCreate,
		DeleteContext: resourceVariableDelete,
		Importer:      &schema.ResourceImporter{State: resourceVariableImport},
		ReadContext:   resourceVariableRead,
		Schema:        getVariableSchema(),
		UpdateContext: resourceVariableUpdate,
	}
}

func resourceVariableImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[INFO] importing variable (%s)", d.Id())

	importStrings := strings.Split(d.Id(), ":")
	if len(importStrings) != 2 {
		return nil, fmt.Errorf("octopusdeploy_variable import must be in the form of OwnerID:VariableID (e.g. Projects-62:0906031f-68ba-4a15-afaa-657c1564e07b")
	}

	d.Set("owner_id", importStrings[0])
	d.SetId(importStrings[1])

	return []*schema.ResourceData{d}, nil
}

func resourceVariableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading variable (%s)", d.Id())

	id := d.Id()
	ownerID := d.Get("owner_id").(string)

	log.Printf(`ID: %s OwnerID: %s`, id, ownerID)

	client := m.(*octopusdeploy.Client)
	variable, err := client.Variables.GetByID(ownerID, id)
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] variable (%s) not found; deleting from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := setVariable(ctx, d, variable); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] variable read (%s)", d.Id())
	return nil
}

func resourceVariableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := validateVariable(d); err != nil {
		return diag.FromErr(err)
	}

	ownerID := d.Get("owner_id").(string)
	variable := expandVariable(d)

	log.Printf("[INFO] creating variable: %#v", variable)

	client := m.(*octopusdeploy.Client)
	tfVar, err := client.Variables.AddSingle(ownerID, variable)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range tfVar.Variables {
		if v.Name == variable.Name && v.Type == variable.Type && (v.IsSensitive || v.Value == variable.Value) && v.Description == variable.Description && v.IsSensitive == variable.IsSensitive {
			scopeMatches, _, err := client.Variables.MatchesScope(v.Scope, &variable.Scope)
			if err != nil {
				return diag.FromErr(err)
			}
			if scopeMatches {
				d.SetId(v.ID)
				log.Printf("[INFO] variable created (%s)", d.Id())
				return nil
			}
		}
	}

	d.SetId("")
	return diag.Errorf("unable to locate variable in project %s", ownerID)
}

func resourceVariableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating variable (%s)", d.Id())

	if err := validateVariable(d); err != nil {
		return diag.FromErr(err)
	}

	tfVar := expandVariable(d)
	ownerID := d.Get("owner_id").(string)

	client := m.(*octopusdeploy.Client)
	updatedVars, err := client.Variables.UpdateSingle(ownerID, tfVar)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range updatedVars.Variables {
		if v.Name == tfVar.Name && v.Type == tfVar.Type && (v.IsSensitive || v.Value == tfVar.Value) && v.Description == tfVar.Description && v.IsSensitive == tfVar.IsSensitive {
			scopeMatches, _, _ := client.Variables.MatchesScope(v.Scope, &tfVar.Scope)
			if scopeMatches {
				if err := setVariable(ctx, d, v); err != nil {
					return diag.FromErr(err)
				}
				log.Printf("[INFO] variable updated (%s)", d.Id())
				return nil
			}
		}
	}

	d.SetId("")
	return diag.Errorf("unable to locate variable in project %s", ownerID)
}

func resourceVariableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting variable (%s)", d.Id())

	ownerID := d.Get("owner_id").(string)

	client := m.(*octopusdeploy.Client)
	_, err := client.Variables.DeleteSingle(ownerID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] variable deleted")
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
