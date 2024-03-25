package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
)

func flattenProjectScheduledTrigger(projectScheduledTrigger *triggers.ProjectTrigger) map[string]interface{} {
	flattenedProjectScheduledTrigger := map[string]interface{}{}
	flattenedProjectScheduledTrigger["id"] = projectScheduledTrigger.GetID()
	flattenedProjectScheduledTrigger["name"] = projectScheduledTrigger.Name
	flattenedProjectScheduledTrigger["description"] = projectScheduledTrigger.Description
	flattenedProjectScheduledTrigger["project_id"] = projectScheduledTrigger.ProjectID
	flattenedProjectScheduledTrigger["space_id"] = projectScheduledTrigger.SpaceID
	flattenedProjectScheduledTrigger["is_disabled"] = projectScheduledTrigger.IsDisabled

	actionType := projectScheduledTrigger.Action.GetActionType()
	if actionType == actions.DeployLatestRelease {
		deployLatestReleaseAction := projectScheduledTrigger.Action.(*actions.DeployLatestReleaseAction)
		flattenedProjectScheduledTrigger["deploy_latest_release_action"] = []map[string]interface{}{
			{
				"source_environment_id":      deployLatestReleaseAction.SourceEnvironments[0],
				"destination_environment_id": deployLatestReleaseAction.DestinationEnvironment,
				"should_redeploy":            deployLatestReleaseAction.ShouldRedeploy,
			},
		}

	} else if actionType == actions.DeployNewRelease {
		deployNewReleaseAction := projectScheduledTrigger.Action.(*actions.DeployNewReleaseAction)
		flattenedProjectScheduledTrigger["deploy_new_release_action"] = []map[string]interface{}{
			{
				"destination_environment_id": deployNewReleaseAction.Environment,
				"git_reference":              deployNewReleaseAction.VersionControlReference.GitRef,
			},
		}
	}
	// TODO throw error if action is nil

	filterType := projectScheduledTrigger.Filter.GetFilterType()
	if filterType == filters.OnceDailySchedule {
		onceDailyScheduleFilter := projectScheduledTrigger.Filter.(*filters.OnceDailyScheduledTriggerFilter)
		days := onceDailyScheduleFilter.Days
		parsedDays := make([]string, len(days))
		for i := range days {
			// TODO handle error case
			parsedDays[i] = filters.Weekday.String(days[i])
		}
		flattenedProjectScheduledTrigger["once_daily_schedule"] = []map[string]interface{}{
			{
				"start_time":   onceDailyScheduleFilter.Start.String(),
				"days_of_week": flattenArray(parsedDays),
			},
		}
		flattenedProjectScheduledTrigger["timezone"] = onceDailyScheduleFilter.TimeZone

	} else if filterType == filters.ContinuousDailySchedule {
		continuousDailyScheduleFilter := projectScheduledTrigger.Filter.(*filters.ContinuousDailyScheduledTriggerFilter)
		days := continuousDailyScheduleFilter.Days
		parsedDays := make([]string, len(days))
		for i := range days {
			// TODO handle error case
			parsedDays[i] = filters.Weekday.String(days[i])
		}

		flattenedProjectScheduledTrigger["continuous_daily_schedule"] = []map[string]interface{}{
			{
				"interval": continuousDailyScheduleFilter.Interval.String(),
				// TODO handle the types
				"hour_interval":   continuousDailyScheduleFilter.HourInterval,
				"minute_interval": continuousDailyScheduleFilter.MinuteInterval,
				"run_after":       continuousDailyScheduleFilter.RunAfter.Format(filters.RFC3339NanoNoZone),
				"run_until":       continuousDailyScheduleFilter.RunUntil.Format(filters.RFC3339NanoNoZone),
				"days_of_week":    flattenArray(parsedDays),
			},
		}
		flattenedProjectScheduledTrigger["timezone"] = continuousDailyScheduleFilter.TimeZone
	} else if filterType == filters.DaysPerMonthSchedule {
		daysPerMonthScheduleFilter := projectScheduledTrigger.Filter.(*filters.DaysPerMonthScheduledTriggerFilter)
		flattenedProjectScheduledTrigger["days_per_month_schedule"] = []map[string]interface{}{
			{
				"start_time":            daysPerMonthScheduleFilter.Start.Format(filters.RFC3339NanoNoZone),
				"monthly_schedule_type": daysPerMonthScheduleFilter.MonthlySchedule.String(),
				"date_of_month":         daysPerMonthScheduleFilter.DateOfMonth,
				"day_number_of_month":   daysPerMonthScheduleFilter.DayNumberOfMonth,
				"day_of_week":           filters.Weekday.String(*daysPerMonthScheduleFilter.Day),
			},
		}
		flattenedProjectScheduledTrigger["timezone"] = daysPerMonthScheduleFilter.TimeZone
	} else if filterType == filters.CronExpressionSchedule {
		cronExpressionScheduleFilter := projectScheduledTrigger.Filter.(*filters.CronScheduledTriggerFilter)
		flattenedProjectScheduledTrigger["cron_expression_schedule"] = []map[string]interface{}{
			{
				"cron_expression": cronExpressionScheduleFilter.CronExpression,
			},
		}
		flattenedProjectScheduledTrigger["timezone"] = cronExpressionScheduleFilter.TimeZone
	}

	return flattenedProjectScheduledTrigger
}

func expandProjectScheduledTrigger(projectScheduledTrigger *schema.ResourceData, client *client.Client) (*triggers.ProjectTrigger, error) {
	name := projectScheduledTrigger.Get("name").(string)
	description := projectScheduledTrigger.Get("description").(string)
	isDisabled := projectScheduledTrigger.Get("is_disabled").(bool)
	timezone := projectScheduledTrigger.Get("timezone").(string)

	projectId := projectScheduledTrigger.Get("project_id").(string)
	spaceId := projectScheduledTrigger.Get("space_id").(string)
	// TODO handle error
	project, _ := projects.GetByID(client, spaceId, projectId)

	var action actions.ITriggerAction = nil
	var filter filters.ITriggerFilter = nil

	// Action configuration
	if attr, ok := projectScheduledTrigger.GetOk("deploy_latest_release_action"); ok {
		deployLatestReleaseActionList := attr.(*schema.Set).List()
		deployLatestReleaseActionMap := deployLatestReleaseActionList[0].(map[string]interface{})
		deploymentAction := actions.NewDeployLatestReleaseAction(
			deployLatestReleaseActionMap["destination_environment_id"].(string),
			deployLatestReleaseActionMap["should_redeploy"].(bool),
			[]string{deployLatestReleaseActionMap["source_environment_id"].(string)},
			"",
		)

		// Might need to add some validation here.
		deploymentAction.Channel = projectScheduledTrigger.Get("channel_id").(string)
		deploymentAction.Tenants = expandArray(projectScheduledTrigger.Get("tenant_ids").([]interface{}))

		action = deploymentAction
	}

	if attr, ok := projectScheduledTrigger.GetOk("deploy_new_release_action"); ok {
		// TODO Blow up if action isn't nil anymore

		deployNewReleaseActionList := attr.(*schema.Set).List()
		deployNewReleaseActionMap := deployNewReleaseActionList[0].(map[string]interface{})
		deploymentAction := actions.NewDeployNewReleaseAction(
			deployNewReleaseActionMap["destination_environment_id"].(string),
			"",
			&actions.VersionControlReference{GitRef: deployNewReleaseActionMap["git_reference"].(string)},
		)

		// Might need to add some validation here.
		deploymentAction.Channel = projectScheduledTrigger.Get("channel_id").(string)
		deploymentAction.Tenants = expandArray(projectScheduledTrigger.Get("tenant_ids").([]interface{}))
		action = deploymentAction
	}

	// Filter configuration
	if attr, ok := projectScheduledTrigger.GetOk("once_daily_schedule"); ok {
		// TODO Blow up if filter isn't nil anymore

		onceDailyScheduleFilterList := attr.(*schema.Set).List()
		onceDailyScheduleFilterMap := onceDailyScheduleFilterList[0].(map[string]interface{})

		// TODO handle error case
		startTime, _ := time.Parse(filters.RFC3339NanoNoZone, onceDailyScheduleFilterMap["start_time"].(string))
		days := expandArray(onceDailyScheduleFilterMap["days_of_week"].([]interface{}))

		parsedDays := make([]filters.Weekday, len(days))
		for i := range days {
			// TODO handle error case
			parsedDays[i], _ = filters.WeekdayString(days[i])
		}

		onceDailyScheduleFilter := filters.NewOnceDailyScheduledTriggerFilter(parsedDays, startTime)
		onceDailyScheduleFilter.TimeZone = timezone
		filter = onceDailyScheduleFilter
	}

	if attr, ok := projectScheduledTrigger.GetOk("continuous_daily_schedule"); ok {
		// TODO Blow up if filter isn't nil anymore

		continuousDailyScheduleFilterList := attr.(*schema.Set).List()
		continuousDailyScheduleFilterMap := continuousDailyScheduleFilterList[0].(map[string]interface{})

		// TODO handle error case
		interval, _ := filters.DailyScheduledIntervalString(continuousDailyScheduleFilterMap["interval"].(string))
		runAfter, err := time.Parse(filters.RFC3339NanoNoZone, continuousDailyScheduleFilterMap["run_after"].(string))

		if err != nil {
			return nil, err
		}

		runUntil, err := time.Parse(filters.RFC3339NanoNoZone, continuousDailyScheduleFilterMap["run_until"].(string))

		if err != nil {
			return nil, err
		}

		days := expandArray(continuousDailyScheduleFilterMap["days_of_week"].([]interface{}))
		// TODO move into helper function
		parsedDays := make([]filters.Weekday, len(days))
		for i := range days {
			// TODO handle error case
			parsedDays[i], _ = filters.WeekdayString(days[i])
		}

		continuousDailyScheduleFilter := filters.NewContinuousDailyScheduledTriggerFilter(parsedDays, timezone)
		continuousDailyScheduleFilter.Interval = &interval
		continuousDailyScheduleFilter.RunAfter = &runAfter
		continuousDailyScheduleFilter.RunUntil = &runUntil

		if interval == filters.OnceHourly {
			hourInterval := int16(continuousDailyScheduleFilterMap["hour_interval"].(int))
			continuousDailyScheduleFilter.HourInterval = &hourInterval
		} else if interval == filters.OnceEveryMinute {
			minuteInterval := int16(continuousDailyScheduleFilterMap["minute_interval"].(int))
			continuousDailyScheduleFilter.MinuteInterval = &minuteInterval
		}

		filter = continuousDailyScheduleFilter
	}

	if attr, ok := projectScheduledTrigger.GetOk("days_per_month_schedule"); ok {
		// TODO Blow up if filter isn't nil anymore
		daysPerMonthScheduleFilterList := attr.(*schema.Set).List()
		daysPerMonthScheduleFilterMap := daysPerMonthScheduleFilterList[0].(map[string]interface{})

		startTime, _ := time.Parse(filters.RFC3339NanoNoZone, daysPerMonthScheduleFilterMap["start_time"].(string))
		monthlyScheduleType, _ := filters.MonthlyScheduleString(daysPerMonthScheduleFilterMap["monthly_schedule_type"].(string))

		daysPerMonthScheduleFilter := filters.NewDaysPerMonthScheduledTriggerFilter(monthlyScheduleType, startTime)

		daysPerMonthScheduleFilter.DateOfMonth = daysPerMonthScheduleFilterMap["date_of_month"].(string)
		daysPerMonthScheduleFilter.DayNumberOfMonth = daysPerMonthScheduleFilterMap["day_number_of_month"].(string)

		dayOfWeek, _ := filters.WeekdayString(daysPerMonthScheduleFilterMap["day_of_week"].(string))
		daysPerMonthScheduleFilter.Day = &dayOfWeek

		daysPerMonthScheduleFilter.TimeZone = timezone
		filter = daysPerMonthScheduleFilter
	}

	if attr, ok := projectScheduledTrigger.GetOk("cron_expression_schedule"); ok {
		cronExpressionScheduleFilterList := attr.(*schema.Set).List()
		cronExpressionScheduleFilterMap := cronExpressionScheduleFilterList[0].(map[string]interface{})

		cronExpression := cronExpressionScheduleFilterMap["cron_expression"].(string)
		cronExpressionScheduleFilter := filters.NewCronScheduledTriggerFilter(cronExpression, timezone)

		filter = cronExpressionScheduleFilter
	}

	// NewProjectTrigger doesn't set the description
	projectTriggerToCreate := triggers.NewProjectTrigger(name, description, isDisabled, project, action, filter)

	projectTriggerToCreate.Description = description
	return projectTriggerToCreate, nil
}

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
			Description:      "The space ID where this trigger's project exists.",
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"deploy_latest_release_action": {
			Description: "Configuration for deploying the latest release. Can not be used with 'deploy_new_release_action'.",
			Optional:    true,
			Type:        schema.TypeSet,
			Elem:        &schema.Resource{Schema: getDeployLatestReleaseActionSchema()},
			MaxItems:    1,
		},
		"deploy_new_release_action": {
			Description: "Configuration for deploying a new release. Can not be used with 'deploy_latest_release_action'.",
			Optional:    true,
			Type:        schema.TypeSet,
			Elem:        &schema.Resource{Schema: getDeployNewReleaseActionSchema()},
			MaxItems:    1,
		},
		"channel_id": {
			Description:      "The channel ID to use when creating the release. Will use the default channel if left blank.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"tenant_ids": {
			Description: "The IDs of the tenants to deploy to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		//"tenant_tags": {}, Not actually sure if this is used
		"is_disabled": {
			Description: "Indicates whether the trigger is disabled.",
			Optional:    true,
			Default:     false,
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

// DeployLatestRelease
func getDeployLatestReleaseActionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"source_environment_id": {
			Required:         true,
			Description:      "The environment ID to use when selecting the release to deploy from.",
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"destination_environment_id": {
			Required:         true,
			Description:      "The environment ID to deploy the selected release to.",
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"should_redeploy": {
			Description: "Enable to re-deploy to the deployment targets even if they are already up-to-date with the current deployment.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
	}
}

// DeployNewRelease
func getDeployNewReleaseActionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"destination_environment_id": {
			Required:         true,
			Description:      "The environment ID to deploy the selected release to.",
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"git_reference": {
			Optional:    true,
			Description: "The git reference to use when creating the release. Can be a branch, tag, or commit hash.",
			Type:        schema.TypeString,
		},
	}
}

// OnceDailySchedule
func getOnceDailyScheduleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start_time": {
			Required:    true,
			Description: "The time of day to start the trigger.",
			Type:        schema.TypeString,
			//ValidateFunc: validation.IsRFC3339Time,
		},
		"days_of_week": {
			Required:    true,
			Description: "The days of the week to run the trigger.",
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			//ValidateDiagFunc: validation.ToDiagFunc(validation.All(
			//	validation.IsDayOfTheWeek(true))),
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
			Required:    true,
			Description: "The time of day to start the trigger.",
			Type:        schema.TypeString,
			//ValidateFunc: validation.IsRFC3339Time,
			// TODO handle supressing diff when the time component is the same
			//DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//	if strings.ToLower(old) == strings.ToLower(new) {
			//		return true
			//	}
			//	return false
			//},
		},
		"run_until": {
			Required:    true,
			Description: "The time of day to end the trigger.",
			Type:        schema.TypeString,
			//ValidateFunc: validation.IsRFC3339Time,
		},
		"days_of_week": {
			Required:    true,
			Description: "The days of the week to run the trigger.",
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			//ValidateDiagFunc: validation.ToDiagFunc(validation.All(
			//	validation.IsDayOfTheWeek(true))),
		},
	}
}

// DaysPerMonthSchedule
func getDaysPerMonthScheduleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start_time": {
			Required:    true,
			Description: "The time of day to start the trigger.",
			Type:        schema.TypeString,
			//ValidateFunc: validation.IsRFC3339Time,
		},
		"monthly_schedule_type": {
			Required:    true,
			Description: "The type of monthly schedule to run the trigger",
			Type:        schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"DateOfMonth",
				"DayOfMonth",
			}, true)),
		},
		"date_of_month": {
			Optional:    true,
			Description: "Which date of the month to run the trigger. String number between 1 - 31 Incl. or L for the last day of the month.",
			Type:        schema.TypeString,
		},
		"day_number_of_month": {
			Optional:    true,
			Description: "Which ordinal day of the week to run the trigger on. String number between 1 - 4 Incl. or L for the last occurrence of day_of_week for the month.",
			Type:        schema.TypeString,
		},
		"day_of_week": {
			Optional:         true,
			Description:      "Which day of the week to run the trigger on.",
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsDayOfTheWeek(true)),
			// TODO supress if monthly schedule type is not DayOfMonth
		},
	}
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
