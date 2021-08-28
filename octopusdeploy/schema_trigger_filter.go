package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTriggerFilter(d map[string]interface{}) *octopusdeploy.TriggerFilter {
	triggerFilter := octopusdeploy.NewTriggerFilter()
	if v, ok := d["id"]; ok {
		triggerFilter.ID = v.(string)
	}

	if v, ok := d["filter_type"]; ok {
		triggerFilter.FilterType = v.(string)
	}

	switch triggerFilter.FilterType {
	case "ContinuousDailySchedule":
		expandContinuousDailyScheduledTriggerFilter(triggerFilter, d)
	case "CronExpressionSchedule":
		expandCronScheduledTriggerFilter(triggerFilter, d)
	case "DaysPerMonthSchedule":
		expandDaysPerMonthScheduledTriggerFilter(triggerFilter, d)
	case "MachineFilter":
		expandMachineTriggerFilter(triggerFilter, d)
	case "OnceDailySchedule":
		expandOnceDailyScheduledTriggerFilter(triggerFilter, d)
	}
	return triggerFilter
}

func flattenTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter) map[string]interface{} {
	if triggerFilter == nil {
		return nil
	}
	flattened := make(map[string]interface{})

	flattened["filter_type"] = triggerFilter.FilterType
	flattened["id"] = triggerFilter.ID
	switch triggerFilter.FilterType {
	case "ContinuousDailySchedule":
		flattenContinuousDailyScheduledTriggerFilter(triggerFilter, flattened)
	case "CronExpressionSchedule":
		flattenCronScheduledTriggerFilter(triggerFilter, flattened)
	case "DaysPerMonthSchedule":
		flattenDaysPerMonthScheduledTriggerFilter(triggerFilter, flattened)
	case "MachineFilter":
		flattenMachineTriggerFilter(triggerFilter, flattened)
	case "OnceDailySchedule":
		flattenOnceDailyScheduledTriggerFilter(triggerFilter, flattened)
	}
	return flattened
}

func getTriggerFilterSchema() (*schema.Schema, *schema.Resource) {
	element := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": getIDSchema(),
			"filter_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	triggerFilterSchema := &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem:     element,
		Optional: true,
	}
	return triggerFilterSchema, element
}

func expandMachineTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, d map[string]interface{}) {
	if v, ok := d["environment_ids"]; ok {
		triggerFilter.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}
	if v, ok := d["roles"]; ok {
		triggerFilter.Roles = getSliceFromTerraformTypeList(v)
	}
	if v, ok := d["event_groups"]; ok {
		triggerFilter.EventGroups = getSliceFromTerraformTypeList(v)
	}
	if v, ok := d["event_categories"]; ok {
		triggerFilter.EventCategories = getSliceFromTerraformTypeList(v)
	}
}

func flattenMachineTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, flattened map[string]interface{}) {
	flattened["environment_ids"] = triggerFilter.EnvironmentIDs
	flattened["roles"] = triggerFilter.Roles
	flattened["event_groups"] = triggerFilter.EventGroups
	flattened["event_categories"] = triggerFilter.EventCategories
}

func getMachineTriggerFilterSchema() *schema.Schema {
	schema, element := getTriggerFilterSchema()
	addMachineTriggerFilterSchema(element)
	return schema
}

func addMachineTriggerFilterSchema(element *schema.Resource) {
	element.Schema["environment_ids"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	element.Schema["roles"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	element.Schema["event_groups"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	element.Schema["event_categories"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func expandScheduledTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, d map[string]interface{}) {
	if v, ok := d["timezone"]; ok {
		triggerFilter.Timezone = v.(string)
	}
}

func flattenScheduledTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, flattened map[string]interface{}) {
	flattened["timezone"] = triggerFilter.Timezone
}

func addScheduledTriggerFilterSchema(element *schema.Resource) {
	element.Schema["timezone"] = &schema.Schema{
		Type:     schema.TypeString,
		Default:  "UTC",
		Optional: true,
	}
}

func expandCronScheduledTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, d map[string]interface{}) {
	expandScheduledTriggerFilter(triggerFilter, d)
	if v, ok := d["cron_expression"]; ok {
		triggerFilter.CronExpression = v.(string)
	}
}

func flattenCronScheduledTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, flattened map[string]interface{}) {
	flattenScheduledTriggerFilter(triggerFilter, flattened)
	flattened["cron_expression"] = triggerFilter.CronExpression
}

func getCronScheduledTriggerFilterSchema() *schema.Schema {
	schema, element := getTriggerFilterSchema()
	addScheduledTriggerFilterSchema(element)
	addCronScheduledTriggerFilterSchema(element)
	return schema
}

func addCronScheduledTriggerFilterSchema(element *schema.Resource) {
	element.Schema["cron_expression"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
}

func expandContinuousDailyScheduledTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, d map[string]interface{}) {
	expandScheduledTriggerFilter(triggerFilter, d)
	if v, ok := d["run_after"]; ok {
		triggerFilter.RunAfter = v.(string)
	}
	if v, ok := d["run_until"]; ok {
		triggerFilter.RunUntil = v.(string)
	}
	if v, ok := d["interval"]; ok {
		triggerFilter.Interval = v.(string)
	}
	if v, ok := d["hour_interval"]; ok {
		triggerFilter.HourInterval = v.(int)
	}
	if v, ok := d["minute_interval"]; ok {
		triggerFilter.MinuteInterval = v.(int)
	}
	if v, ok := d["days_of_week"]; ok {
		triggerFilter.DaysOfWeek = getSliceFromTerraformTypeList(v)
	}
}

func flattenContinuousDailyScheduledTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, flattened map[string]interface{}) {
	flattenScheduledTriggerFilter(triggerFilter, flattened)
	flattened["run_after"] = triggerFilter.RunAfter
	flattened["run_until"] = triggerFilter.RunUntil
	flattened["interval"] = triggerFilter.Interval
	flattened["hour_interval"] = triggerFilter.HourInterval
	flattened["minute_interval"] = triggerFilter.MinuteInterval
	flattened["days_of_week"] = triggerFilter.DaysOfWeek
}

func getContinuousDailyScheduledTriggerFilterSchema() *schema.Schema {
	schema, element := getTriggerFilterSchema()
	addScheduledTriggerFilterSchema(element)
	addContinuousDailyScheduledTriggerFilterSchema(element)
	return schema
}

func addContinuousDailyScheduledTriggerFilterSchema(element *schema.Resource) {
	element.Schema["run_after"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	element.Schema["run_util"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	element.Schema["interval"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	element.Schema["hour_interval"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
	}
	element.Schema["minute_internal"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	element.Schema["days_of_week"] = &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	}
}

func expandDaysPerMonthScheduledTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, d map[string]interface{}) {
	expandScheduledTriggerFilter(triggerFilter, d)
	if v, ok := d["start_time"]; ok {
		triggerFilter.StartTime = v.(string)
	}
	if v, ok := d["monthly_schedule_type"]; ok {
		triggerFilter.MonthlyScheduleType = v.(string)
	}
	if v, ok := d["date_of_month"]; ok {
		triggerFilter.DateOfMonth = v.(string)
	}
	if v, ok := d["day_number_of_month"]; ok {
		triggerFilter.DayNumberOfMonth = v.(string)
	}
	if v, ok := d["day_of_week"]; ok {
		triggerFilter.DayOfWeek = v.(string)
	}
}

func flattenDaysPerMonthScheduledTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, flattened map[string]interface{}) {
	flattenScheduledTriggerFilter(triggerFilter, flattened)
	flattened["start_time"] = triggerFilter.StartTime
	flattened["monthly_schedule_type"] = triggerFilter.MonthlyScheduleType
	flattened["date_of_month"] = triggerFilter.DateOfMonth
	flattened["day_number_of_month"] = triggerFilter.DayNumberOfMonth
	flattened["day_of_week"] = triggerFilter.DayOfWeek
}

func getDaysPerMonthScheduledTriggerFilterSchema() *schema.Schema {
	schema, element := getTriggerFilterSchema()
	addScheduledTriggerFilterSchema(element)
	addDaysPerMonthScheduledTriggerFilterSchema(element)
	return schema
}

func addDaysPerMonthScheduledTriggerFilterSchema(element *schema.Resource) {
	element.Schema["start_time"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	element.Schema["monthly_schedule_type"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	element.Schema["date_of_month"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	element.Schema["day_number_of_month"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	element.Schema["day_of_week"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
}

func expandOnceDailyScheduledTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, d map[string]interface{}) {
	expandScheduledTriggerFilter(triggerFilter, d)
	if v, ok := d["start_time"]; ok {
		triggerFilter.StartTime = v.(string)
	}
	if v, ok := d["days_of_week"]; ok {
		triggerFilter.DaysOfWeek = getSliceFromTerraformTypeList(v)
	}
}
func flattenOnceDailyScheduledTriggerFilter(triggerFilter *octopusdeploy.TriggerFilter, flattened map[string]interface{}) {
	flattenScheduledTriggerFilter(triggerFilter, flattened)
	flattened["start_time"] = triggerFilter.StartTime
	flattened["days_of_week"] = triggerFilter.DaysOfWeek
}

func getOnceDailyScheduledTriggerFilterSchema() *schema.Schema {
	schema, element := getTriggerFilterSchema()
	addScheduledTriggerFilterSchema(element)
	addOnceDailyScheduledTriggerFilterSchema(element)
	return schema
}

func addOnceDailyScheduledTriggerFilterSchema(element *schema.Resource) {
	element.Schema["start_time"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	element.Schema["days_of_week"] = &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Required: true,
	}
}
