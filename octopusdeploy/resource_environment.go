package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
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

	client := m.(*client.Client)
	createdEnvironment, err := environments.Add(client, environment)
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

	client := m.(*client.Client)
	if err := environments.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] environment deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading environment (%s)", d.Id())

	client := m.(*client.Client)
	environment, err := environments.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "environment")
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
	client := m.(*client.Client)
	updatedEnvironment, err := environments.Update(client, environment)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setEnvironment(ctx, d, updatedEnvironment); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] environment updated (%s)", d.Id())
	return nil
}
