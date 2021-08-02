package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		DeleteContext: resourceEnvironmentDelete,
		Description:   "This resource manages environments in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceEnvironmentRead,
		Schema:        getEnvironmentSchema(),
		UpdateContext: resourceEnvironmentUpdate,
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	environment := expandEnvironment(d)

	log.Printf("[INFO] creating environment: %#v", environment)

	client := m.(*octopusdeploy.Client)
	createdEnvironment, err := client.Environments.Add(environment)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setEnvironment(ctx, d, createdEnvironment); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdEnvironment.GetID())

	log.Printf("[INFO] environment created (%s)", d.Id())
	return nil
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting environment (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Environments.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] environment deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading environment (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	environment, err := client.Environments.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] environment (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setEnvironment(ctx, d, environment); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] environment read (%s)", d.Id())
	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating environment (%s)", d.Id())

	environment := expandEnvironment(d)
	client := m.(*octopusdeploy.Client)
	updatedEnvironment, err := client.Environments.Update(environment)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setEnvironment(ctx, d, updatedEnvironment); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] environment updated (%s)", d.Id())
	return nil
}
