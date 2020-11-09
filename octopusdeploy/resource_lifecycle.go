package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceLifecycle() *schema.Resource {
	resourceLifecycleImporter := &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
	resourceLifecycleSchema := map[string]*schema.Schema{
		"description": &schema.Schema{
			Optional: true,
			Type:     schema.TypeString,
		},
		constName: &schema.Schema{
			Required: true,
			Type:     schema.TypeString,
		},
		constPhase: {
			Elem:     phaseSchema(),
			Optional: true,
			Type:     schema.TypeList,
		},
		constReleaseRetentionPolicy: getRetentionPeriodSchema(),
		constSpaceID: {
			Type:     schema.TypeString,
			Computed: true,
		},
		constTentacleRetentionPolicy: getRetentionPeriodSchema(),
	}

	return &schema.Resource{
		CreateContext: resourceLifecycleCreate,
		DeleteContext: resourceLifecycleDelete,
		Importer:      resourceLifecycleImporter,
		ReadContext:   resourceLifecycleRead,
		Schema:        resourceLifecycleSchema,
		UpdateContext: resourceLifecycleUpdate,
	}
}

func getRetentionPeriodSchema() *schema.Schema {
	return &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				constQuantityToKeep: {
					Default:     30,
					Description: "The number of days/releases to keep. If 0 all are kept.",
					Optional:    true,
					Type:        schema.TypeInt,
				},
				constShouldKeepForever: {
					Default:  false,
					Optional: true,
					Type:     schema.TypeBool,
				},
				constUnit: {
					Default:     octopusdeploy.RetentionUnitDays,
					Description: "The unit of quantity_to_keep.",
					Optional:    true,
					Type:        schema.TypeString,
					ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
						octopusdeploy.RetentionUnitDays,
						octopusdeploy.RetentionUnitItems,
					}, false)),
				},
			},
		},
		Optional: true,
		Type:     schema.TypeList,
	}
}

func phaseSchema() *schema.Resource {
	phaseSchema := map[string]*schema.Schema{
		constAutomaticDeploymentTargets: {
			Description: "Environment IDs in this phase that a release is automatically deployed to when it is eligible for this phase",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		constID: {
			Type:     schema.TypeString,
			Computed: true,
		},
		constIsOptionalPhase: {
			Default:     false,
			Description: "If false a release must be deployed to this phase before it can be deployed to the next phase.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		constMinimumEnvironmentsBeforePromotion: {
			Default:     0,
			Description: "The number of units required before a release can enter the next phase. If 0, all environments are required.",
			Optional:    true,
			Type:        schema.TypeInt,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		constOptionalDeploymentTargets: {
			Description: "Environment IDs in this phase that a release can be deployed to, but is not automatically deployed to",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		constReleaseRetentionPolicy:  getRetentionPeriodSchema(),
		constTentacleRetentionPolicy: getRetentionPeriodSchema(),
	}

	return &schema.Resource{
		Schema: phaseSchema,
	}
}

func resourceLifecycleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lifecycle := buildLifecycleResource(d)

	client := m.(*octopusdeploy.Client)
	createdLifecycle, err := client.Lifecycles.Add(lifecycle)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenLifecycle(ctx, d, createdLifecycle)
	return nil
}

func buildLifecycleResource(d *schema.ResourceData) *octopusdeploy.Lifecycle {
	var name string
	if v, ok := d.GetOk(constName); ok {
		name = v.(string)
	}

	lifecycle := octopusdeploy.NewLifecycle(name)

	if v, ok := d.GetOk("description"); ok {
		lifecycle.Description = v.(string)
	}

	if v, ok := d.GetOk(constSpaceID); ok {
		lifecycle.SpaceID = v.(string)
	}

	releaseRetentionPolicy := getRetentionPeriod(d, constReleaseRetentionPolicy)
	if releaseRetentionPolicy != nil {
		lifecycle.ReleaseRetentionPolicy = *releaseRetentionPolicy
	}

	tentacleRetentionPolicy := getRetentionPeriod(d, constTentacleRetentionPolicy)
	if tentacleRetentionPolicy != nil {
		lifecycle.TentacleRetentionPolicy = *tentacleRetentionPolicy
	}

	if attr, ok := d.GetOk(constPhases); ok {
		tfPhases := attr.([]interface{})

		for _, tfPhase := range tfPhases {
			phase := buildPhaseResource(tfPhase.(map[string]interface{}))
			lifecycle.Phases = append(lifecycle.Phases, phase)
		}
	}

	return lifecycle
}

func getRetentionPeriod(d *schema.ResourceData, key string) *octopusdeploy.RetentionPeriod {
	v, ok := d.GetOk(key)
	if ok {
		retentionPeriod := v.([]interface{})
		if len(retentionPeriod) == 1 {
			tfRetentionItem := retentionPeriod[0].(map[string]interface{})
			retention := octopusdeploy.RetentionPeriod{
				QuantityToKeep:    int32(tfRetentionItem[constQuantityToKeep].(int)),
				ShouldKeepForever: tfRetentionItem[constShouldKeepForever].(bool),
				Unit:              tfRetentionItem[constUnit].(string),
			}
			return &retention
		}
	}

	return nil
}

func buildPhaseResource(tfPhase map[string]interface{}) octopusdeploy.Phase {
	phase := octopusdeploy.Phase{
		AutomaticDeploymentTargets:         getSliceFromTerraformTypeList(tfPhase[constAutomaticDeploymentTargets]),
		IsOptionalPhase:                    tfPhase[constIsOptionalPhase].(bool),
		MinimumEnvironmentsBeforePromotion: int32(tfPhase[constMinimumEnvironmentsBeforePromotion].(int)),
		Name:                               tfPhase[constName].(string),
		OptionalDeploymentTargets:          getSliceFromTerraformTypeList(tfPhase[constOptionalDeploymentTargets]),
	}

	if phase.AutomaticDeploymentTargets == nil {
		phase.AutomaticDeploymentTargets = []string{}
	}
	if phase.OptionalDeploymentTargets == nil {
		phase.OptionalDeploymentTargets = []string{}
	}

	return phase
}

func resourceLifecycleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	lifecycle, err := client.Lifecycles.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	flattenLifecycle(ctx, d, lifecycle)
	return nil
}

func resourceLifecycleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lifecycle := buildLifecycleResource(d)
	lifecycle.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	updatedLifecycle, err := client.Lifecycles.Update(lifecycle)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenLifecycle(ctx, d, updatedLifecycle)
	return nil
}

func resourceLifecycleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Lifecycles.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func flattenStringArray(values []string) []interface{} {
	s := make([]interface{}, len(values))
	for i, v := range values {
		s[i] = v
	}
	return s
}

func flattenPhases(phases []octopusdeploy.Phase) []interface{} {
	flattenedPhases := make([]interface{}, 0)
	for _, phase := range phases {
		p := make(map[string]interface{})
		p[constAutomaticDeploymentTargets] = flattenStringArray(phase.AutomaticDeploymentTargets)
		p["id"] = phase.ID
		p[constIsOptionalPhase] = phase.IsOptionalPhase
		p[constMinimumEnvironmentsBeforePromotion] = int(phase.MinimumEnvironmentsBeforePromotion)
		p[constName] = phase.Name
		p[constOptionalDeploymentTargets] = flattenStringArray(phase.OptionalDeploymentTargets)
		if phase.ReleaseRetentionPolicy != nil {
			p[constReleaseRetentionPolicy] = flattenRetentionPeriod(*phase.ReleaseRetentionPolicy)
		}
		if phase.TentacleRetentionPolicy != nil {
			p[constTentacleRetentionPolicy] = flattenRetentionPeriod(*phase.TentacleRetentionPolicy)
		}
		flattenedPhases = append(flattenedPhases, p)
	}
	return flattenedPhases
}

func flattenRetentionPeriod(r octopusdeploy.RetentionPeriod) []interface{} {
	retentionPeriod := make(map[string]interface{})
	retentionPeriod["unit"] = r.Unit
	retentionPeriod["quantity_to_keep"] = int(r.QuantityToKeep)
	retentionPeriod["should_keep_forever"] = r.ShouldKeepForever
	return []interface{}{retentionPeriod}
}

func flattenLifecycle(ctx context.Context, d *schema.ResourceData, lifecycle *octopusdeploy.Lifecycle) {
	d.Set("description", lifecycle.Description)
	d.Set("name", lifecycle.Name)
	d.Set("phase", flattenPhases(lifecycle.Phases))
	d.Set("space_id", lifecycle.SpaceID)
	d.Set(constReleaseRetentionPolicy, flattenRetentionPeriod(lifecycle.ReleaseRetentionPolicy))
	d.Set(constTentacleRetentionPolicy, flattenRetentionPeriod(lifecycle.TentacleRetentionPolicy))

	d.SetId(lifecycle.GetID())
}
