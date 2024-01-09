package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/scriptmodules"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceScriptModules() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing script modules.",
		ReadContext: dataSourceScriptModulesRead,
		Schema:      getScriptModuleDataSchema(),
	}
}

func dataSourceScriptModulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := variables.LibraryVariablesQuery{
		ContentType: "ScriptModule",
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	spaceID := d.Get("space_id").(string)

	client := m.(*client.Client)
	existingScriptModules, err := scriptmodules.Get(client, spaceID, query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedScriptModules := []interface{}{}
	for _, scriptModule := range existingScriptModules.Items {
		flattenedScriptModules = append(flattenedScriptModules, flattenScriptModule(scriptModule))
	}

	d.Set("script_modules", flattenedScriptModules)
	d.SetId("Script Modules " + time.Now().UTC().String())

	return nil
}
