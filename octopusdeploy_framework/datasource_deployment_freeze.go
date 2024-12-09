package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deploymentfreezes"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

const deploymentFreezeDatasourceName = "deployment_freezes"

type recurringScheduleDatasourceModel struct {
	Type                types.String `tfsdk:"type"`
	Unit                types.Int64  `tfsdk:"unit"`
	EndType             types.String `tfsdk:"end_type"`
	EndOnDate           types.String `tfsdk:"end_on_date"`
	EndAfterOccurrences types.Int64  `tfsdk:"end_after_occurrences"`
	MonthlyScheduleType types.String `tfsdk:"monthly_schedule_type"`
	DateOfMonth         types.String `tfsdk:"date_of_month"`
	DayNumberOfMonth    types.String `tfsdk:"day_number_of_month"`
	DaysOfWeek          types.List   `tfsdk:"days_of_week"`
	DayOfWeek           types.String `tfsdk:"day_of_week"`
}

type deploymentFreezeDatasourceModel struct {
	ID                      types.String                      `tfsdk:"id"`
	Name                    types.String                      `tfsdk:"name"`
	Start                   types.String                      `tfsdk:"start"`
	End                     types.String                      `tfsdk:"end"`
	ProjectEnvironmentScope types.Map                         `tfsdk:"project_environment_scope"`
	RecurringSchedule       *recurringScheduleDatasourceModel `tfsdk:"recurring_schedule"`
}

type deploymentFreezesDatasourceModel struct {
	ID                types.String `tfsdk:"id"`
	IDs               types.List   `tfsdk:"ids"`
	PartialName       types.String `tfsdk:"partial_name"`
	ProjectIDs        types.List   `tfsdk:"project_ids"`
	EnvironmentIDs    types.List   `tfsdk:"environment_ids"`
	IncludeComplete   types.Bool   `tfsdk:"include_complete"`
	Status            types.String `tfsdk:"status"`
	Skip              types.Int64  `tfsdk:"skip"`
	Take              types.Int64  `tfsdk:"take"`
	DeploymentFreezes types.List   `tfsdk:"deployment_freezes"`
}

type deploymentFreezeDataSource struct {
	*Config
}

func (d *deploymentFreezeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = DataSourceConfiguration(req, resp)
}

func NewDeploymentFreezeDataSource() datasource.DataSource {
	return &deploymentFreezeDataSource{}
}

func (d *deploymentFreezeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(deploymentFreezeDatasourceName)
}

func (d *deploymentFreezeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.DeploymentFreezeSchema{}.GetDatasourceSchema()
}

func (d *deploymentFreezeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data deploymentFreezesDatasourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := deploymentfreezes.DeploymentFreezeQuery{
		IDs:             util.Ternary(data.IDs.IsNull(), []string{}, util.ExpandStringList(data.IDs)),
		PartialName:     data.PartialName.ValueString(),
		ProjectIds:      util.Ternary(data.ProjectIDs.IsNull(), []string{}, util.ExpandStringList(data.ProjectIDs)),
		EnvironmentIds:  util.Ternary(data.EnvironmentIDs.IsNull(), []string{}, util.ExpandStringList(data.EnvironmentIDs)),
		IncludeComplete: data.IncludeComplete.ValueBool(),
		Status:          data.Status.ValueString(),
		Skip:            int(data.Skip.ValueInt64()),
		Take:            int(data.Take.ValueInt64()),
	}

	util.DatasourceReading(ctx, "deployment freezes", query)

	existingFreezes, err := deploymentfreezes.Get(d.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load deployment freezes", err.Error())
		return
	}

	flattenedFreezes := []interface{}{}
	for _, freeze := range existingFreezes.DeploymentFreezes {
		flattenedFreeze, diags := mapFreezeToAttribute(ctx, freeze)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		flattenedFreezes = append(flattenedFreezes, flattenedFreeze)
	}

	data.ID = types.StringValue("Deployment Freezes " + time.Now().UTC().String())
	data.DeploymentFreezes, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: freezeObjectType()}, flattenedFreezes)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

var _ datasource.DataSource = &deploymentFreezeDataSource{}
var _ datasource.DataSourceWithConfigure = &deploymentFreezeDataSource{}

func mapFreezeToAttribute(ctx context.Context, freeze deploymentfreezes.DeploymentFreeze) (attr.Value, diag.Diagnostics) {
	projectScopes := make(map[string]attr.Value)
	for projectId, environmentScopes := range freeze.ProjectEnvironmentScope {
		projectScopes[projectId] = util.FlattenStringList(environmentScopes)
	}

	scopeType, diags := types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, projectScopes)
	if diags.HasError() {
		return nil, diags
	}

	attrs := map[string]attr.Value{
		"id":                        types.StringValue(freeze.ID),
		"name":                      types.StringValue(freeze.Name),
		"start":                     types.StringValue(freeze.Start.Format(time.RFC3339)),
		"end":                       types.StringValue(freeze.End.Format(time.RFC3339)),
		"project_environment_scope": scopeType,
	}

	if freeze.RecurringSchedule != nil {
		daysOfWeek, diags := types.ListValueFrom(ctx, types.StringType, freeze.RecurringSchedule.DaysOfWeek)
		if diags.HasError() {
			return nil, diags
		}

		endOnDate := types.StringNull()
		if freeze.RecurringSchedule.EndOnDate != nil {
			endOnDate = types.StringValue(freeze.RecurringSchedule.EndOnDate.Format(time.RFC3339))
		}

		endAfterOccurrences := types.Int64Null()
		if freeze.RecurringSchedule.EndAfterOccurrences != nil {
			endAfterOccurrences = types.Int64Value(int64(*freeze.RecurringSchedule.EndAfterOccurrences))
		}

		monthlyScheduleType := types.StringNull()
		if freeze.RecurringSchedule.MonthlyScheduleType != "" {
			monthlyScheduleType = types.StringValue(freeze.RecurringSchedule.MonthlyScheduleType)
		}

		dateOfMonth := types.StringNull()
		if freeze.RecurringSchedule.DateOfMonth != nil {
			dateOfMonth = types.StringValue(*freeze.RecurringSchedule.DateOfMonth)
		}

		dayNumberOfMonth := types.StringNull()
		if freeze.RecurringSchedule.DayNumberOfMonth != nil {
			dayNumberOfMonth = types.StringValue(*freeze.RecurringSchedule.DayNumberOfMonth)
		}

		dayOfWeek := types.StringNull()
		if freeze.RecurringSchedule.DayOfWeek != nil {
			dayOfWeek = types.StringValue(*freeze.RecurringSchedule.DayOfWeek)
		}

		scheduleAttrs := map[string]attr.Value{
			"type":                  types.StringValue(string(freeze.RecurringSchedule.Type)),
			"unit":                  types.Int64Value(int64(freeze.RecurringSchedule.Unit)),
			"end_type":              types.StringValue(string(freeze.RecurringSchedule.EndType)),
			"end_on_date":           endOnDate,
			"end_after_occurrences": endAfterOccurrences,
			"monthly_schedule_type": monthlyScheduleType,
			"date_of_month":         dateOfMonth,
			"day_number_of_month":   dayNumberOfMonth,
			"days_of_week":          daysOfWeek,
			"day_of_week":           dayOfWeek,
		}

		recurringSchedule, diags := types.ObjectValue(freezeRecurringScheduleObjectType(), scheduleAttrs)
		if diags.HasError() {
			return nil, diags
		}

		attrs["recurring_schedule"] = recurringSchedule
	} else {
		attrs["recurring_schedule"] = types.ObjectNull(freezeRecurringScheduleObjectType())
	}

	return types.ObjectValueMust(freezeObjectType(), attrs), diags
}

func freezeRecurringScheduleObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                  types.StringType,
		"unit":                  types.Int64Type,
		"end_type":              types.StringType,
		"end_on_date":           types.StringType,
		"end_after_occurrences": types.Int64Type,
		"monthly_schedule_type": types.StringType,
		"date_of_month":         types.StringType,
		"day_number_of_month":   types.StringType,
		"days_of_week":          types.ListType{ElemType: types.StringType},
		"day_of_week":           types.StringType,
	}
}

func freezeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                        types.StringType,
		"name":                      types.StringType,
		"start":                     types.StringType,
		"end":                       types.StringType,
		"project_environment_scope": types.MapType{ElemType: types.ListType{ElemType: types.StringType}},
		"recurring_schedule":        types.ObjectType{AttrTypes: freezeRecurringScheduleObjectType()},
	}
}
