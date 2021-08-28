package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandProjectTriggerFilter(d map[string]interface{}) *octopusdeploy.TriggerFilter {
	var vrd map[string]interface{}
	if v, ok := d["cron"]; ok && len(v.([]interface{})) > 0 {
		vrd = v.([]interface{})[0].(map[string]interface{})
		vrd["filter_type"] = "CronExpressionSchedule"
	}
	if v, ok := d["once_daily"]; ok && len(v.([]interface{})) > 0 {
		vrd = v.([]interface{})[0].(map[string]interface{})
		vrd["filter_type"] = "OnceDailySchedule"
	}
	if v, ok := d["days_per_month"]; ok && len(v.([]interface{})) > 0 {
		vrd = v.([]interface{})[0].(map[string]interface{})
		vrd["filter_type"] = "DaysPerMonthSchedule"
	}
	if v, ok := d["continuous_daily"]; ok && len(v.([]interface{})) > 0 {
		vrd = v.([]interface{})[0].(map[string]interface{})
		vrd["filter_type"] = "ContinuousDailySchedule"
	}
	if v, ok := d["machine"]; ok && len(v.([]interface{})) > 0 {
		vrd = v.([]interface{})[0].(map[string]interface{})
		vrd["filter_type"] = "MachineFilter"
	}
	return expandTriggerFilter(vrd)
}

func flattenProjectTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter) map[string]interface{} {
	filter_type := ""
	switch triggerFilter.FilterType {
	case "CronExpressionSchedule":
		filter_type = "cron"
	case "ContinuousDailySchedule":
		filter_type = "continuous_daily"
	case "DaysPerMonthSchedule":
		filter_type = "day_per_month"
	case "MachineFilter":
		filter_type = "machine"
	case "OnceDailySchedule":
		filter_type = "once_daily"
	}
	return map[string]interface{}{
		filter_type: []interface{}{flattenTriggerFilter(triggerFilter)},
	}
}

func getProjectTriggerFilterSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cron":             getCronScheduledTriggerFilterSchema(),
				"once_daily":       getOnceDailyScheduledTriggerFilterSchema(),
				"days_per_month":   getDaysPerMonthScheduledTriggerFilterSchema(),
				"continuous_daily": getContinuousDailyScheduledTriggerFilterSchema(),
				"machine":          getMachineTriggerFilterSchema(),
			},
		},
	}
}

func expandProjectTriggerAction(d map[string]interface{}) *octopusdeploy.TriggerAction {
	var vrd map[string]interface{}
	if v, ok := d["auto_deploy"]; ok && len(v.([]interface{})) > 0 {
		vrd = v.([]interface{})[0].(map[string]interface{})
		vrd["action_type"] = "AutoDeploy"
	}
	if v, ok := d["deploy_latest"]; ok && len(v.([]interface{})) > 0 {
		vrd = v.([]interface{})[0].(map[string]interface{})
		vrd["action_type"] = "DeployLatestRelease"
	}
	if v, ok := d["deploy_new"]; ok && len(v.([]interface{})) > 0 {
		vrd = v.([]interface{})[0].(map[string]interface{})
		vrd["action_type"] = "DeployNewRelease"
	}
	return expandTriggerAction(vrd)
}

func flattenProjectTriggerAction(triggerAction *octopusdeploy.TriggerAction) map[string]interface{} {
	action_type := ""
	switch triggerAction.ActionType {
	case "AutoDeploy":
		action_type = "auto_deploy"
	case "DeployLatestRelease":
		action_type = "deploy_latest"
	case "DeployNewRelease":
		action_type = "deploy_new"
	}
	return map[string]interface{}{
		action_type: []interface{}{flattenTriggerAction(triggerAction)},
	}
}

func getProjectTriggerActionSchema() *schema.Schema {
	autoDeploy := getAutoDeployTriggerActionSchema()
	autoDeploy.ConflictsWith = []string{"filter.0.cron", "filter.0.once_daily", "filter.0.days_per_month", "filter.0.continuous_daily"}
	autoDeploy.RequiredWith = []string{"filter.0.machine"}
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_deploy":   autoDeploy,
				"deploy_latest": getDeployLatestReleaseTriggerActionSchema(),
				"deploy_new":    getDeployNewReleaseTriggerActionSchema(),
			},
		},
	}
}

func getProjectTriggerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema(),
		"id":          getIDSchema(),
		"name":        getNameSchema(true),
		"space_id":    getSpaceIDSchema(),
		"is_disabled": {
			Type:        schema.TypeBool,
			Description: "Whether the trigger is disabled or not.",
			Default:     false,
			Optional:    true,
		},
		"project_id": {
			Type:        schema.TypeString,
			Description: "The ID of the project to attach the trigger.",
			Required:    true,
		},
		"action": makeMutuallyExclusive(getProjectTriggerActionSchema(), "action"),
		"filter": makeMutuallyExclusive(getProjectTriggerFilterSchema(), "filter"),
	}
}

func expandProjectTrigger(d *schema.ResourceData) *octopusdeploy.ProjectTrigger {
	projectTrigger := octopusdeploy.NewProjectTrigger()
	projectTrigger.ID = d.Id()
	projectTrigger.Name = d.Get("name").(string)
	projectTrigger.ProjectID = d.Get("project_id").(string)

	if v, ok := d.GetOk("description"); ok {
		projectTrigger.Description = v.(string)
	}
	if v, ok := d.GetOk("space_id"); ok {
		projectTrigger.SpaceID = v.(string)
	}
	projectTrigger.Action = expandProjectTriggerAction(d.Get("action").([]interface{})[0].(map[string]interface{}))
	projectTrigger.Filter = expandProjectTriggerFilter(d.Get("filter").([]interface{})[0].(map[string]interface{}))
	return projectTrigger
}

func setProjectTrigger(ctx context.Context, d *schema.ResourceData, projectTrigger *octopusdeploy.ProjectTrigger) error {
	d.Set("description", projectTrigger.Description)

	d.Set("name", projectTrigger.Name)
	d.Set("project_id", projectTrigger.ProjectID)
	d.Set("space_id", projectTrigger.SpaceID)
	d.Set("is_disabled", projectTrigger.IsDisabled)

	d.SetId(projectTrigger.GetID())

	if err := d.Set("action", []interface{}{flattenProjectTriggerAction(projectTrigger.Action)}); err != nil {
		return fmt.Errorf("error setting action: %s", err)
	}

	if err := d.Set("filter", []interface{}{flattenProjectTriggerFilter(projectTrigger.Filter)}); err != nil {
		return fmt.Errorf("error setting filter: %s", err)
	}

	return nil
}
