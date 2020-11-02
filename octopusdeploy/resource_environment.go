package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		DeleteContext: resourceEnvironmentDelete,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			constAllowDynamicInfrastructure: {
				Default:  false,
				Optional: true,
				Type:     schema.TypeBool,
			},
			constDescription: {
				Optional: true,
				Type:     schema.TypeString,
			},
			constName: {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			constSortOrder: {
				Default:  -1,
				Optional: true,
				Type:     schema.TypeInt,
			},
			constUseGuidedFailure: {
				Default:  false,
				Optional: true,
				Type:     schema.TypeBool,
			},
		},
	}
}

func buildEnvironmentResource(d *schema.ResourceData) *octopusdeploy.Environment {
	name := d.Get(constName).(string)
	environment := octopusdeploy.NewEnvironment(name)

	if v, ok := d.GetOk(constAllowDynamicInfrastructure); ok {
		environment.AllowDynamicInfrastructure = v.(bool)
	}

	if v, ok := d.GetOk(constDescription); ok {
		environment.Description = v.(string)
	}

	if v, ok := d.GetOk(constSortOrder); ok {
		environment.SortOrder = v.(int)
	}

	if v, ok := d.GetOk(constUseGuidedFailure); ok {
		environment.UseGuidedFailure = v.(bool)
	}

	return environment
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	environment, err := client.Environments.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	logResource(constEnvironment, m)

	d.Set(constAllowDynamicInfrastructure, environment.AllowDynamicInfrastructure)
	d.Set(constDescription, environment.Description)
	d.Set(constName, environment.Name)
	d.Set(constSortOrder, environment.SortOrder)
	d.Set(constUseGuidedFailure, environment.UseGuidedFailure)

	return nil
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	environment := buildEnvironmentResource(d)

	client := m.(*octopusdeploy.Client)
	createdEnvironment, err := client.Environments.Add(environment)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdEnvironment.GetID())

	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	environment := buildEnvironmentResource(d)
	environment.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	updatedEnvironment, err := client.Environments.Update(environment)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(updatedEnvironment.GetID())

	return nil
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Environments.DeleteByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)
	return nil
}
