package schemas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type daysOfWeekValidator struct{}

func NewDaysOfWeekValidator() daysOfWeekValidator {
	return daysOfWeekValidator{}
}

func (v daysOfWeekValidator) Description(ctx context.Context) string {
	return "validates that days of the week are valid and in correct order (Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday)"
}

func (v daysOfWeekValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v daysOfWeekValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	validDays := map[string]int{
		"Sunday":    0,
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
	}

	var days []string
	req.ConfigValue.ElementsAs(ctx, &days, false)

	for i := 1; i < len(days); i++ {
		currentDay := days[i]
		previousDay := days[i-1]

		currentPos, currentExists := validDays[currentDay]
		previousPos, previousExists := validDays[previousDay]

		if !currentExists {
			resp.Diagnostics.AddError(
				"Invalid day of week",
				fmt.Sprintf("'%s' is not a valid day of week. Must be one of: Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday", currentDay),
			)
			return
		}

		if !previousExists {
			resp.Diagnostics.AddError(
				"Invalid day of week",
				fmt.Sprintf("'%s' is not a valid day of week. Must be one of: Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday", previousDay),
			)
			return
		}

		if currentPos <= previousPos {
			resp.Diagnostics.AddError(
				"Invalid day order",
				fmt.Sprintf("Days of the week must be in order (Sunday through Saturday). Found '%s' after '%s'", currentDay, previousDay),
			)
			return
		}
	}
}

type recurringScheduleValidator struct{}

func NewRecurringScheduleValidator() recurringScheduleValidator {
	return recurringScheduleValidator{}
}

func (v recurringScheduleValidator) Description(_ context.Context) string {
	return "validates that required fields are set based on the schedule type"
}

func (v recurringScheduleValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v recurringScheduleValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var schedule struct {
		Type                types.String      `tfsdk:"type"`
		Unit                types.Int64       `tfsdk:"unit"`
		UtcOffsetInMinutes  types.Int64       `tfsdk:"utc_offset_in_minutes"`
		EndType             types.String      `tfsdk:"end_type"`
		EndOnDate           timetypes.RFC3339 `tfsdk:"end_on_date"`
		EndAfterOccurrences types.Int64       `tfsdk:"end_after_occurrences"`
		MonthlyScheduleType types.String      `tfsdk:"monthly_schedule_type"`
		DateOfMonth         types.String      `tfsdk:"date_of_month"`
		DayNumberOfMonth    types.String      `tfsdk:"day_number_of_month"`
		DaysOfWeek          types.List        `tfsdk:"days_of_week"`
		DayOfWeek           types.String      `tfsdk:"day_of_week"`
	}

	resp.Diagnostics.Append(req.ConfigValue.As(ctx, &schedule, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	scheduleType := schedule.Type.ValueString()

	switch scheduleType {
	case "Daily":
		// Daily only requires type and unit which are already marked as required

	case "Weekly":
		if schedule.DaysOfWeek.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("days_of_week"),
				"Missing Required Field",
				"days_of_week must be set when schedule type is DaysPerWeek",
			)
		}

	case "Monthly":
		if schedule.MonthlyScheduleType.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("monthly_schedule_type"),
				"Missing Required Field",
				"monthly_schedule_type must be set when schedule type is DaysPerMonth",
			)
			return
		}

		monthlyType := schedule.MonthlyScheduleType.ValueString()
		switch monthlyType {
		case "DateOfMonth":
			if schedule.DateOfMonth.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("date_of_month"),
					"Missing Required Field",
					"date_of_month must be set when monthly_schedule_type is DateOfMonth",
				)
			}

		case "DayOfMonth":
			if schedule.DayNumberOfMonth.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("day_number_of_month"),
					"Missing Required Field",
					"day_number_of_month must be set when monthly_schedule_type is DayOfMonth",
				)
			} else {
				dayNum := schedule.DayNumberOfMonth.ValueString()
				validDayNums := map[string]bool{
					"1": true,
					"2": true,
					"3": true,
					"4": true,
					"L": true,
				}
				if !validDayNums[dayNum] {
					resp.Diagnostics.AddAttributeError(
						path.Root("day_number_of_month"),
						"Invalid Day Number",
						fmt.Sprintf("day_number_of_month must be one of: 1, 2, 3, 4, L, got: %s", dayNum),
					)
				}
			}
			if schedule.DayOfWeek.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("day_of_week"),
					"Missing Required Field",
					"day_of_week must be set when monthly_schedule_type is DayOfMonth",
				)
			}

		default:
			resp.Diagnostics.AddAttributeError(
				path.Root("monthly_schedule_type"),
				"Invalid Monthly Schedule Type",
				fmt.Sprintf("monthly_schedule_type must be either DateOfMonth or DayOfMonth, got: %s", monthlyType),
			)
		}

	case "Annually":
		// Annually only requires type and unit which are already marked as required

	default:
		resp.Diagnostics.AddAttributeError(
			path.Root("type"),
			"Invalid Schedule Type",
			fmt.Sprintf("type must be one of: Daily, Weekly, Monthly, Annually, got: %s", scheduleType),
		)
	}
}
