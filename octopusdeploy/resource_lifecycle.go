package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLifecycle() *schema.Resource {
	return &schema.Resource{
		Create: resourceLifecycleCreate,
		Read:   resourceLifecycleRead,
		Update: resourceLifecycleUpdate,
		Delete: resourceLifecycleDelete,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constReleaseRetentionPolicy:  getRetentionPeriodSchema(),
			constTentacleRetentionPolicy: getRetentionPeriodSchema(),
			constPhase:                   getPhasesSchema(),
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
				constUnit: {
					Type:        schema.TypeString,
					Description: "The unit of quantity_to_keep.",
					Optional:    true,
					Default:     (string)(model.RetentionUnitDays),
					ValidateFunc: validateValueFunc([]string{
						(string)(model.RetentionUnitDays),
						(string)(model.RetentionUnitItems),
					}),
				},
				constQuantityToKeep: {
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
				constName: {
					Type:     schema.TypeString,
					Required: true,
				},
				constMinimumEnvironmentsBeforePromotion: {
					Description: "The number of units required before a release can enter the next phase. If 0, all environments are required.",
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     0,
				},
				constIsOptionalPhase: {
					Description: "If false a release must be deployed to this phase before it can be deployed to the next phase.",
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
				},
				constAutomaticDeploymentTargets: {
					Description: "Environment Ids in this phase that a release is automatically deployed to when it is eligible for this phase",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				constOptionalDeploymentTargets: {
					Description: "Environment Ids in this phase that a release can be deployed to, but is not automatically deployed to",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				// constReleaseRetentionPolicy: getRetentionPeriodSchema(),
				// constTentacleRetentionPolicy: getRetentionPeriodSchema(),
			},
		},
	}
}

func resourceLifecycleCreate(d *schema.ResourceData, m interface{}) error {
	lifecycle, err := buildLifecycleResource(d)
	if err != nil {
		return err
	}

	apiClient := m.(*client.Client)
	resource, err := apiClient.Lifecycles.Add(lifecycle)
	if err != nil {
		return createResourceOperationError(errorCreatingLifecycle, lifecycle.Name, err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func buildLifecycleResource(d *schema.ResourceData) (*model.Lifecycle, error) {
	name := d.Get(constName).(string)

	lifecycle, err := model.NewLifecycle(name)
	if err != nil {
		return nil, err
	}

	if attr, ok := d.GetOk(constDescription); ok {
		lifecycle.Description = attr.(string)
	}

	releaseRetentionPolicy := getRetentionPeriod(d, constReleaseRetentionPolicy)
	if releaseRetentionPolicy != nil {
		lifecycle.ReleaseRetentionPolicy = *releaseRetentionPolicy
	}

	tentacleRetentionPolicy := getRetentionPeriod(d, constTentacleRetentionPolicy)
	if tentacleRetentionPolicy != nil {
		lifecycle.TentacleRetentionPolicy = *tentacleRetentionPolicy
	}

	if attr, ok := d.GetOk(constPhase); ok {
		tfPhases := attr.([]interface{})

		for _, tfPhase := range tfPhases {
			phase := buildPhaseResource(tfPhase.(map[string]interface{}))
			lifecycle.Phases = append(lifecycle.Phases, phase)
		}
	}

	return lifecycle, nil
}

func getRetentionPeriod(d *schema.ResourceData, key string) *model.RetentionPeriod {
	attr, ok := d.GetOk(key)
	if ok {
		tfRetentionSettings := attr.(*schema.Set)
		if len(tfRetentionSettings.List()) == 1 {
			tfRetentionItem := tfRetentionSettings.List()[0].(map[string]interface{})
			retention := model.RetentionPeriod{
				Unit:           model.RetentionUnit(tfRetentionItem[constUnit].(string)),
				QuantityToKeep: int32(tfRetentionItem[constQuantityToKeep].(int)),
			}
			return &retention
		}
	}

	return nil
}

func buildPhaseResource(tfPhase map[string]interface{}) model.Phase {
	phase := model.Phase{
		Name:                               tfPhase[constName].(string),
		MinimumEnvironmentsBeforePromotion: int32(tfPhase[constMinimumEnvironmentsBeforePromotion].(int)),
		IsOptionalPhase:                    tfPhase[constIsOptionalPhase].(bool),
		AutomaticDeploymentTargets:         getSliceFromTerraformTypeList(tfPhase[constAutomaticDeploymentTargets]),
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

func resourceLifecycleRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Lifecycles.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingLifecycle, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constLifecycle, m)

	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)

	return nil
}

func resourceLifecycleUpdate(d *schema.ResourceData, m interface{}) error {
	lifecycle, err := buildLifecycleResource(d)
	if err != nil {
		return err
	}
	lifecycle.ID = d.Id() // set ID so Octopus API knows which lifecycle to update

	apiClient := m.(*client.Client)
	resource, err := apiClient.Lifecycles.Update(*lifecycle)
	if err != nil {
		return createResourceOperationError(errorUpdatingLifecycle, d.Id(), err)
	}

	d.SetId(resource.ID)

	return nil
}

func resourceLifecycleDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.Lifecycles.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingLifecycle, id, err)
	}

	d.SetId(constEmptyString)

	return nil
}
