package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProjectDeploymentFreezeSchema struct{}

func (d ProjectDeploymentFreezeSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id": GetIdResourceSchema(),
			"owner_id": resourceSchema.StringAttribute{
				Description: "The Owner ID of the freeze",
				Required:    true,
			},
			"name":  GetNameResourceSchema(true),
			"start": GetDateTimeResourceSchema("The start time of the freeze, must be RFC3339 format", true),
			"end":   GetDateTimeResourceSchema("The end time of the freeze, must be RFC3339 format", true),
			"environment_ids": resourceSchema.ListAttribute{
				Description: "The environment IDs associated with this project deployment freeze scope",
				Optional:    true,
				ElementType: types.StringType,
			},
			"recurring_schedule": resourceSchema.SingleNestedAttribute{
				Optional: true,
				Validators: []validator.Object{
					NewRecurringScheduleValidator(),
				},
				Attributes: map[string]resourceSchema.Attribute{
					"type": resourceSchema.StringAttribute{
						Description: "Type of recurring schedule (Daily, Weekly, Monthly, Annually)",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("Daily", "Weekly", "Monthly", "Annually"),
						},
					},
					"unit": resourceSchema.Int64Attribute{
						Description: "The unit value for the schedule",
						Required:    true,
					},
					"end_type": resourceSchema.StringAttribute{
						Description: "When the recurring schedule should end (Never, OnDate, AfterOccurrences)",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("Never", "OnDate", "AfterOccurrences"),
						},
					},
					"end_on_date": GetDateTimeResourceSchema("The date when the recurring schedule should end", false),
					"end_after_occurrences": resourceSchema.Int64Attribute{
						Description: "Number of occurrences after which the schedule should end",
						Optional:    true,
					},
					"monthly_schedule_type": resourceSchema.StringAttribute{
						Description: "Type of monthly schedule (DayOfMonth, DateOfMonth)",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("DayOfMonth", "DateOfMonth"),
						},
					},
					"date_of_month": resourceSchema.StringAttribute{
						Description: "The date of the month for monthly schedules",
						Optional:    true,
					},
					"day_number_of_month": resourceSchema.StringAttribute{
						Description: "Specifies which weekday position in the month. Valid values: 1 (First), 2 (Second), 3 (Third), 4 (Fourth), L (Last). Used with day_of_week",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("1", "2", "3", "4", "L"),
						},
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
						Description: "The day of the week for monthly schedules when using DayOfMonth type",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"),
						},
					},
				},
			},
		},
	}
}

func (d ProjectDeploymentFreezeSchema) GetDatasourceSchema() datasourceSchema.Schema {
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
						"id": GetIdDatasourceSchema(true),
						"owner_id": datasourceSchema.StringAttribute{
							Optional:    false,
							Computed:    true,
							Description: "The Owner ID of the freeze.",
						},
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

var _ EntitySchema = &ProjectDeploymentFreezeSchema{}
