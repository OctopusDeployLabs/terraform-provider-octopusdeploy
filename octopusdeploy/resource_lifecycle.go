package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceLifecycle() *schema.Resource {
	return &schema.Resource{
		Create: resourceLifecycleCreate,
		Read:   resourceLifecycleRead,
		Update: resourceLifecycleUpdate,
		Delete: resourceLifecycleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"release_retention_policy":  getRetentionPeriodSchema(),
			"tentacle_retention_policy": getRetentionPeriodSchema(),
			"phase":                     getPhasesSchema(),
		},
	}
}

func getRetentionPeriodSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"unit": {
					Type:        schema.TypeString,
					Description: "The unit of quantity_to_keep.",
					Optional:    true,
					Default:     (string)(octopusdeploy.RetentionUnit_Days),
					ValidateFunc: validateValueFunc([]string{
						(string)(octopusdeploy.RetentionUnit_Days),
						(string)(octopusdeploy.RetentionUnit_Items),
					}),
				},
				"quantity_to_keep": {
					Type:        schema.TypeInt,
					Description: "The number of days/releases to keep. If 0 all are kept.",
					Default:     0,
					Optional:    true,
				},
			},
		},
	}
}

func getPhasesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"minimum_environments_before_promotion": {
					Description: "The number of units required before a release can enter the next phase. If 0, all environments are required.",
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     0,
				},
				"is_optional_phase": {
					Description: "If false a release must be deployed to this phase before it can be deployed to the next phase.",
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
				},
				"automatic_deployment_targets": {
					Description: "Environment Ids in this phase that a release is automatically deployed to when it is eligible for this phase",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"optional_deployment_targets": {
					Description: "Environment Ids in this phase that a release can be deployed to, but is not automatically deployed to",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				//"release_retention_policy": getRetentionPeriodSchema(),
				//"tentacle_retention_policy": getRetentionPeriodSchema(),
			},
		},
	}
}

func resourceLifecycleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newLifecycle := buildLifecycleResource(d)

	createdLifecycle, err := client.Lifecycle.Add(newLifecycle)

	if err != nil {
		return fmt.Errorf("error creating project: %s", err.Error())
	}

	d.SetId(createdLifecycle.ID)

	return nil
}

func buildLifecycleResource(d *schema.ResourceData) *octopusdeploy.Lifecycle {
	name := d.Get("name").(string)

	lifecycle := octopusdeploy.NewLifecycle(name)

	if attr, ok := d.GetOk("description"); ok {
		lifecycle.Description = attr.(string)
	}

	releaseRetentionPolicy := getRetentionPeriod(d, "release_retention_policy")
	if releaseRetentionPolicy != nil {
		lifecycle.ReleaseRetentionPolicy = *releaseRetentionPolicy
	}

	tentacleRetentionPolicy := getRetentionPeriod(d, "tentacle_retention_policy")
	if tentacleRetentionPolicy != nil {
		lifecycle.TentacleRetentionPolicy = *tentacleRetentionPolicy
	}

	if attr, ok := d.GetOk("phase"); ok {
		tfPhases := attr.([]interface{})

		for _, tfPhase := range tfPhases {
			phase := buildPhaseResource(tfPhase.(map[string]interface{}))
			lifecycle.Phases = append(lifecycle.Phases, phase)
		}
	}

	return lifecycle
}

func getRetentionPeriod(d *schema.ResourceData, key string) *octopusdeploy.RetentionPeriod {
	attr, ok := d.GetOk(key)
	if ok {
		tfRetentionSettings := attr.(*schema.Set)
		if len(tfRetentionSettings.List()) == 1 {
			tfRetentionItem := tfRetentionSettings.List()[0].(map[string]interface{})
			retention := octopusdeploy.RetentionPeriod{
				Unit:           octopusdeploy.RetentionUnit(tfRetentionItem["unit"].(string)),
				QuantityToKeep: int32(tfRetentionItem["quantity_to_keep"].(int)),
			}
			return &retention
		}
	}

	return nil
}

func buildPhaseResource(tfPhase map[string]interface{}) octopusdeploy.Phase {
	phase := octopusdeploy.Phase{
		Name:                               tfPhase["name"].(string),
		MinimumEnvironmentsBeforePromotion: int32(tfPhase["minimum_environments_before_promotion"].(int)),
		IsOptionalPhase:                    tfPhase["is_optional_phase"].(bool),
		AutomaticDeploymentTargets:         getSliceFromTerraformTypeList(tfPhase["automatic_deployment_targets"]),
		OptionalDeploymentTargets:          getSliceFromTerraformTypeList(tfPhase["optional_deployment_targets"]),
	}

	if phase.AutomaticDeploymentTargets == nil {
		phase.AutomaticDeploymentTargets = []string{}
	}
	if phase.OptionalDeploymentTargets == nil {
		phase.OptionalDeploymentTargets = []string{}
	}

	return phase
}

func resourceLifecycleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	lifecycleID := d.Id()

	lifecycle, err := client.Lifecycle.Get(lifecycleID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading lifecycle id %s: %s", lifecycleID, err.Error())
	}

	log.Printf("[DEBUG] lifecycle: %v", m)
	d.Set("name", lifecycle.Name)
	d.Set("description", lifecycle.Description)

	return nil
}

func resourceLifecycleUpdate(d *schema.ResourceData, m interface{}) error {
	lifecycle := buildLifecycleResource(d)
	lifecycle.ID = d.Id() // set lifecycle struct ID so octopus knows which lifecycle to update

	client := m.(*octopusdeploy.Client)

	lifecycle, err := client.Lifecycle.Update(lifecycle)

	if err != nil {
		return fmt.Errorf("error updating lifecycle id %s: %s", d.Id(), err.Error())
	}

	d.SetId(lifecycle.ID)

	return nil
}

func resourceLifecycleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	lifecycleID := d.Id()

	err := client.Lifecycle.Delete(lifecycleID)

	if err != nil {
		return fmt.Errorf("error deleting lifecycle id %s: %s", lifecycleID, err.Error())
	}

	d.SetId("")
	return nil
}
