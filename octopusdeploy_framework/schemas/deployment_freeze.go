package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
						Description: "List of days of the week for weekly schedules",
						Optional:    true,
						ElementType: types.StringType,
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
						"recurring_schedule": datasourceSchema.SingleNestedAttribute{
							Computed: true,
							Optional: false,
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
				Optional: false,
				Computed: true,
			},
		},
	}
}

var _ EntitySchema = &DeploymentFreezeSchema{}
