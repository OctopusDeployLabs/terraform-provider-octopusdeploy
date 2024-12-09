package schemas

import (
	"context"
	"fmt"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type DeploymentFreezeSchema struct{}

func (d DeploymentFreezeSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":    GetIdResourceSchema(),
			"name":  GetNameResourceSchema(true),
			"start": GetDateTimeResourceSchema("The start time of the freeze, must be RFC3339 format", true),
			"end":   GetDateTimeResourceSchema("The end time of the freeze, must be RFC3339 format", true),
			"recurring_schedule": resourceSchema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]resourceSchema.Attribute{
					"type": resourceSchema.StringAttribute{
						Description: "Type of recurring schedule (OnceDaily, DaysPerWeek, DaysPerMonth, Annually)",
						Required:    true,
					},
					"unit": resourceSchema.Int64Attribute{
						Description: "The unit value for the schedule",
						Required:    true,
					},
					"end_type": resourceSchema.StringAttribute{
						Description: "When the recurring schedule should end (Never, OnDate, AfterOccurrences)",
						Required:    true,
					},
					"end_on_date": GetDateTimeResourceSchema("The date when the recurring schedule should end", false),
					"end_after_occurrences": resourceSchema.Int64Attribute{
						Description: "Number of occurrences after which the schedule should end",
						Optional:    true,
					},
					"monthly_schedule_type": resourceSchema.StringAttribute{
						Description: "Type of monthly schedule (DayOfMonth, DateOfMonth)",
						Optional:    true,
					},
					"date_of_month": resourceSchema.StringAttribute{
						Description: "The date of the month for monthly schedules",
						Optional:    true,
					},
					"day_number_of_month": resourceSchema.StringAttribute{
						Description: "The day number of the month for monthly schedules",
						Optional:    true,
					},
					"days_of_week": resourceSchema.ListAttribute{
						Description: "List of days of the week for weekly schedules. Must follow order: Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday",
						Optional:    true,
						ElementType: types.StringType,
						Validators: []validator.List{
							NewDaysOfWeekValidator(),
						},
					},
					"day_of_week": resourceSchema.StringAttribute{
						Description: "The day of the week for monthly schedules",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (d DeploymentFreezeSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about deployment freezes",
		Attributes: map[string]datasourceSchema.Attribute{
			"id":           GetIdDatasourceSchema(true),
			"ids":          GetQueryIDsDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take":         GetQueryTakeDatasourceSchema(),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"project_ids": datasourceSchema.ListAttribute{
				Description: "A filter to search by a list of project IDs",
				ElementType: types.StringType,
				Optional:    true,
			},
			"tenant_ids": datasourceSchema.ListAttribute{
				Description: "A filter to search by a list of tenant IDs",
				ElementType: types.StringType,
				Optional:    true,
			},
			"environment_ids": datasourceSchema.ListAttribute{
				Description: "A filter to search by a list of environment IDs",
				ElementType: types.StringType,
				Optional:    true,
			},
			"include_complete": GetBooleanDatasourceAttribute("Include deployment freezes that completed, default is true", true),
			"status": datasourceSchema.StringAttribute{
				Description: "Filter by the status of the deployment freeze, value values are Expired, Active, Scheduled (case-insensitive)",
				Optional:    true,
			},
			"deployment_freezes": datasourceSchema.ListNestedAttribute{
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id":   GetIdDatasourceSchema(true),
						"name": GetReadonlyNameDatasourceSchema(),
						"start": datasourceSchema.StringAttribute{
							Description: "The start time of the freeze",
							Optional:    false,
							Computed:    true,
						},
						"end": datasourceSchema.StringAttribute{
							Description: "The end time of the freeze",
							Optional:    false,
							Computed:    true,
						},
						"project_environment_scope": datasourceSchema.MapAttribute{
							ElementType: types.ListType{ElemType: types.StringType},
							Description: "The project environment scope of the deployment freeze",
							Optional:    false,
							Computed:    true,
						},
						"tenant_project_environment_scope": datasourceSchema.ListNestedAttribute{
							Description: "The tenant project environment scope of the deployment freeze",
							Optional:    false,
							Computed:    true,
							NestedObject: datasourceSchema.NestedAttributeObject{
								Attributes: map[string]datasourceSchema.Attribute{
									"tenant_id": datasourceSchema.StringAttribute{
										Description: "The tenant ID",
										Computed:    true,
									},
									"project_id": datasourceSchema.StringAttribute{
										Description: "The project ID",
										Computed:    true,
									},
									"environment_id": datasourceSchema.StringAttribute{
										Description: "The environment ID",
										Computed:    true,
									},
								},
							},
						},
						"recurring_schedule": datasourceSchema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]datasourceSchema.Attribute{
								"type": datasourceSchema.StringAttribute{
									Description: "Type of recurring schedule (OnceDaily, DaysPerWeek, DaysPerMonth, Annually)",
									Computed:    true,
								},
								"unit": datasourceSchema.Int64Attribute{
									Description: "The unit value for the schedule",
									Computed:    true,
								},
								"end_type": datasourceSchema.StringAttribute{
									Description: "When the recurring schedule should end (Never, OnDate, AfterOccurrences)",
									Computed:    true,
								},
								"end_on_date": datasourceSchema.StringAttribute{
									Description: "The date when the recurring schedule should end",
									Computed:    true,
								},
								"end_after_occurrences": datasourceSchema.Int64Attribute{
									Description: "Number of occurrences after which the schedule should end",
									Computed:    true,
								},
								"monthly_schedule_type": datasourceSchema.StringAttribute{
									Description: "Type of monthly schedule (DayOfMonth, DateOfMonth)",
									Computed:    true,
								},
								"date_of_month": datasourceSchema.StringAttribute{
									Description: "The date of the month for monthly schedules",
									Computed:    true,
								},
								"day_number_of_month": datasourceSchema.StringAttribute{
									Description: "The day number of the month for monthly schedules",
									Computed:    true,
								},
								"days_of_week": datasourceSchema.ListAttribute{
									Description: "List of days of the week for weekly schedules",
									Computed:    true,
									ElementType: types.StringType,
								},
								"day_of_week": datasourceSchema.StringAttribute{
									Description: "The day of the week for monthly schedules",
									Computed:    true,
								},
							},
						},
					},
				},
				Computed: true,
			},
		},
	}
}

var _ EntitySchema = &DeploymentFreezeSchema{}
