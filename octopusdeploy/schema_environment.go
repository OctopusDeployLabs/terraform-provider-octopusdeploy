package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandEnvironment(d *schema.ResourceData) *octopusdeploy.Environment {
	name := d.Get("name").(string)

	environment := octopusdeploy.NewEnvironment(name)
	environment.ID = d.Id()

	if v, ok := d.GetOk("allow_dynamic_infrastructure"); ok {
		environment.AllowDynamicInfrastructure = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		environment.Description = v.(string)
	}

	if v, ok := d.GetOk("sort_order"); ok {
		environment.SortOrder = v.(int)
	}

	if v, ok := d.GetOk("use_guided_failure"); ok {
		environment.UseGuidedFailure = v.(bool)
	}

	return environment
}

func getEnvironmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_dynamic_infrastructure": &schema.Schema{
			Optional: true,
			Type:     schema.TypeBool,
		},
		"description": &schema.Schema{
			Optional: true,
			Type:     schema.TypeString,
		},
		"name": &schema.Schema{
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"sort_order": &schema.Schema{
			Computed: true,
			Type:     schema.TypeInt,
		},
		"use_guided_failure": &schema.Schema{
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}

func setEnvironment(ctx context.Context, d *schema.ResourceData, environment *octopusdeploy.Environment) {
	d.Set("allow_dynamic_infrastructure", environment.AllowDynamicInfrastructure)
	d.Set("description", environment.Description)
	d.Set("name", environment.Name)
	d.Set("sort_order", environment.SortOrder)
	d.Set("use_guided_failure", environment.UseGuidedFailure)

	d.SetId(environment.GetID())
}
