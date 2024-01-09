package octopusdeploy

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var mutex = &sync.Mutex{}

func resourceVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVariableCreate,
		DeleteContext: resourceVariableDelete,
		Description:   "This resource manages variables in Octopus Deploy.",
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

func resourceVariableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	if err := validateVariable(d); err != nil {
		return diag.FromErr(err)
	}

	var spaceID string
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	projectID, projectOk := d.GetOk("project_id")
	ownerID, ownerOk := d.GetOk("owner_id")

	if !projectOk && !ownerOk {
		return diag.Errorf("one of project_id or owner_id must be configured")
	}

	var variableOwnerID string

	if projectOk {
		variableOwnerID = projectID.(string)
	} else {
		variableOwnerID = ownerID.(string)
	}

	variable := expandVariable(d)

	log.Printf("[INFO] creating variable: %#v", variable)

	client := m.(*client.Client)
	variableSet, err := variables.AddSingle(client, spaceID, variableOwnerID, variable)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range variableSet.Variables {
		if v.Name == variable.Name && v.Type == variable.Type && (v.IsSensitive || v.Value == variable.Value) && v.Description == variable.Description && v.IsSensitive == variable.IsSensitive {
			scopeMatches, _, err := matchesScope(v.Scope, &variable.Scope)
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
	return diag.Errorf("unable to locate variable for owner ID, %s", variableOwnerID)
}

func resourceVariableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading variable (%s)", d.Id())

	id := d.Id()

	var spaceID string
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	projectID, projectOk := d.GetOk("project_id")
	ownerID, ownerOk := d.GetOk("owner_id")

	if !projectOk && !ownerOk {
		return diag.Errorf("one of project_id or owner_id must be configured")
	}

	var variableOwnerID string

	if projectOk {
		variableOwnerID = projectID.(string)
	} else {
		variableOwnerID = ownerID.(string)
	}

	client := m.(*client.Client)
	variable, err := variables.GetByID(client, spaceID, variableOwnerID, id)
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "variable")
	}

	if err := setVariable(ctx, d, variable); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] variable read (%s)", d.Id())
	return nil
}

func resourceVariableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	log.Printf("[INFO] updating variable (%s)", d.Id())

	if err := validateVariable(d); err != nil {
		return diag.FromErr(err)
	}

	variable := expandVariable(d)

	var spaceID string
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	projectID, projectOk := d.GetOk("project_id")
	ownerID, ownerOk := d.GetOk("owner_id")

	if !projectOk && !ownerOk {
		return diag.Errorf("one of project_id or owner_id must be configured")
	}

	var variableOwnerID string

	if projectOk {
		variableOwnerID = projectID.(string)
	} else {
		variableOwnerID = ownerID.(string)
	}

	client := m.(*client.Client)
	variableSet, err := variables.UpdateSingle(client, spaceID, variableOwnerID, variable)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range variableSet.Variables {
		if v.Name == variable.Name && v.Type == variable.Type && (v.IsSensitive || v.Value == variable.Value) && v.Description == variable.Description && v.IsSensitive == variable.IsSensitive {
			scopeMatches, _, _ := matchesScope(v.Scope, &variable.Scope)
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
	return diag.Errorf("unable to locate variable for owner ID, %s", variableOwnerID)
}

func resourceVariableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	log.Printf("[INFO] deleting variable (%s)", d.Id())

	var spaceID string
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	projectID, projectOk := d.GetOk("project_id")
	ownerID, ownerOk := d.GetOk("owner_id")

	if !projectOk && !ownerOk {
		return diag.Errorf("one of project_id or owner_id must be configured")
	}

	var variableOwnerID string

	if projectOk {
		variableOwnerID = projectID.(string)
	} else {
		variableOwnerID = ownerID.(string)
	}

	client := m.(*client.Client)
	_, err := variables.DeleteSingle(client, spaceID, variableOwnerID, d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "variable")
	}

	log.Printf("[INFO] variable deleted (%s)", d.Id())
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

func matchesScope(variableScope variables.VariableScope, definedScope *variables.VariableScope) (bool, *variables.VariableScope, error) {
	if definedScope == nil {
		return true, &variables.VariableScope{}, nil
	}

	if definedScope.IsEmpty() {
		return true, &variables.VariableScope{}, nil
	}

	var matchedScopes variables.VariableScope
	var matched bool

	for _, e1 := range definedScope.Environments {
		for _, e2 := range variableScope.Environments {
			if e1 == e2 {
				matched = true
				matchedScopes.Environments = append(matchedScopes.Environments, e1)
			}
		}
	}

	for _, r1 := range definedScope.Roles {
		for _, r2 := range variableScope.Roles {
			if r1 == r2 {
				matched = true
				matchedScopes.Roles = append(matchedScopes.Roles, r1)
			}
		}
	}

	for _, m1 := range definedScope.Machines {
		for _, m2 := range variableScope.Machines {
			if m1 == m2 {
				matched = true
				matchedScopes.Machines = append(matchedScopes.Machines, m1)
			}
		}
	}

	for _, a1 := range definedScope.Actions {
		for _, a2 := range variableScope.Actions {
			if a1 == a2 {
				matched = true
				matchedScopes.Actions = append(matchedScopes.Actions, a1)
			}
		}
	}

	for _, c1 := range definedScope.Channels {
		for _, c2 := range variableScope.Channels {
			if c1 == c2 {
				matched = true
				matchedScopes.Channels = append(matchedScopes.Channels, c1)
			}
		}
	}

	for _, p1 := range definedScope.ProcessOwners {
		for _, p2 := range variableScope.ProcessOwners {
			if p1 == p2 {
				matched = true
				matchedScopes.ProcessOwners = append(matchedScopes.ProcessOwners, p1)
			}
		}
	}

	return matched, &matchedScopes, nil
}
