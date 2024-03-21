package octopusdeploy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getProjectScheduledTriggerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": getNameSchema(true),
		"description": {
			Description: "A description of the trigger.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"project_id": {
			Description:      "The ID of the project to attach the trigger.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"space_id": {
			Required:         true,
			Description:      "The space ID associated with the project to attach the trigger.",
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"release_action":          {},
		"channel":                 {},
		"source_environment":      {},
		"destination_environment": {},
		"tenant_ids":              {},
		"tenant_tags":             {},
		"is_disabled": {
			Description: "Indicates whether the trigger is disabled.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"should_redeploy": {
			Description: "Enable to re-deploy to the deployment targets even if they are already up-to-date with the current deployment.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"timezone": {
			Description: "The timezone for the trigger.",
			Optional:    true,
			Default:     "UTC",
			Type:        schema.TypeString,
		},
		"once_daily_schedule": {
			Description: "The daily schedule for the trigger.",
			Optional:    true,
			Type:        schema.TypeSet,
			Elem:        &schema.Resource{Schema: getOnceDailyScheduleSchema()},
			MaxItems:    1,
		},
		"continuous_daily_schedule": {
			Description: "The daily schedule for the trigger.",
			Optional:    true,
			Type:        schema.TypeSet,
			Elem:        &schema.Resource{Schema: getContinuousDailyScheduleSchema()},
			MaxItems:    1,
		},
		"days_per_month_schedule": {
			Description: "The daily schedule for the trigger.",
			Optional:    true,
			Type:        schema.TypeSet,
			Elem:        &schema.Resource{Schema: getDaysPerMonthScheduleSchema()},
			MaxItems:    1,
		},
		"cron_expression_schedule": {
			Description: "The cron expression schedule for the trigger.",
			Optional:    true,
			Type:        schema.TypeSet,
			Elem:        &schema.Resource{Schema: getCronExpressionScheduleSchema()},
			MaxItems:    1,
		},
	}
}

// OnceDailySchedule
func getOnceDailyScheduleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start_time": {
			Required:     true,
			Description:  "The time of day to start the trigger.",
			Type:         schema.TypeString,
			ValidateFunc: validation.IsRFC3339Time,
		},
		"days_of_week": {
			Required:    true,
			Description: "The days of the week to run the trigger.",
			Type:        schema.TypeList,
			ValidateDiagFunc: validation.ToDiagFunc(validation.All(
				validation.IsDayOfTheWeek(true))),
		},
	}
}

// ContinuousDailySchedule
func getContinuousDailyScheduleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"interval": {
			Required:    true,
			Description: "The interval in minutes to run the trigger.",
			Type:        schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"OnceDaily",
				"OnceHourly",
				"OnceEveryMinute",
			}, true)),
		},
		"hour_interval": {
			Optional:         true,
			Description:      "How often to run the trigger in hours. Only used when the interval is set to 'OnceHourly'.",
			Type:             schema.TypeInt,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		},
		"minute_interval": {
			Optional:         true,
			Description:      "How often to run the trigger in minutes. Only used when the interval is set to 'OnceEveryMinute'.",
			Type:             schema.TypeInt,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		},
		"run_after": {
			Required:     true,
			Description:  "The time of day to start the trigger.",
			Type:         schema.TypeString,
			ValidateFunc: validation.IsRFC3339Time,
		},
		"run_until": {
			Required:     true,
			Description:  "The time of day to end the trigger.",
			Type:         schema.TypeString,
			ValidateFunc: validation.IsRFC3339Time,
		},
		"days_of_week": {
			Required:    true,
			Description: "The days of the week to run the trigger.",
			Type:        schema.TypeList,
			ValidateDiagFunc: validation.ToDiagFunc(validation.All(
				validation.IsDayOfTheWeek(true))),
		},
	}
}

// DaysPerMonthSchedule
func getDaysPerMonthScheduleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

// CronExpressionSchedule
func getCronExpressionScheduleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cron_expression": {
			Description: "The cron expression for the schedule.",
			Required:    true,
			Type:        schema.TypeString,
		},
	}
}
