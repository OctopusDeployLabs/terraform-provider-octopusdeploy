package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScriptModule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScriptModuleCreate,
		DeleteContext: resourceScriptModuleDelete,
		Description:   "This resource manages script modules in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceScriptModuleRead,
		Schema:        getScriptModuleSchema(),
		UpdateContext: resourceScriptModuleUpdate,
	}
}

func resourceScriptModuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scriptModule := expandScriptModule(d)

	log.Printf("[INFO] creating script module: %#v", scriptModule)

	client := m.(*client.Client)
	createdScriptModule, err := client.ScriptModules.Add(scriptModule)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setScriptModule(ctx, d, createdScriptModule); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdScriptModule.GetID())

	log.Printf("[INFO] script module created (%s)", d.Id())
	return nil
}

func resourceScriptModuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting script module (%s)", d.Id())

	client := m.(*client.Client)
	err := client.ScriptModules.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] script module deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceScriptModuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading script module (%s)", d.Id())

	client := m.(*client.Client)
	scriptModule, err := client.ScriptModules.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "script module")
	}

	if err := setScriptModule(ctx, d, scriptModule); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] script module read (%s)", d.Id())
	return nil
}

func resourceScriptModuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating script module (%s)", d.Id())

	scriptModule := expandScriptModule(d)

	client := m.(*client.Client)
	updatedScriptModule, err := client.ScriptModules.Update(scriptModule)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setScriptModule(ctx, d, updatedScriptModule); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] script module updated (%s)", d.Id())
	return nil
}
