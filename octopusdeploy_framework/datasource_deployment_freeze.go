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

type deploymentFreezesDatasourceModel struct {
	ID                types.String `tfsdk:"id"`
	IDs               types.List   `tfsdk:"ids"`
	PartialName       types.String `tfsdk:"partial_name"`
	ProjectIDs        types.List   `tfsdk:"project_ids"`
	EnvironmentIDs    types.List   `tfsdk:"environment_ids"`
	TenantIDs         types.List   `tfsdk:"tenant_ids"`
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
		TenantIds:       util.Ternary(data.TenantIDs.IsNull(), []string{}, util.ExpandStringList(data.TenantIDs)),
		IncludeComplete: data.IncludeComplete.ValueBool(),
		Status:          data.Status.ValueString(),
		Skip:            int(data.Skip.ValueInt64()),
		Take:            int(data.Take.ValueInt64()),
	}

	util.DatasourceReading(ctx, "deployment freezes", query)

	existingFreezes, err := deploymentfreezes.Get(d.Client, query)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, d.Config.SystemInfo, "unable to load deployment freezes", err.Error())
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
	var diags diag.Diagnostics

	projectScopes := make(map[string]attr.Value)
	for projectId, environmentScopes := range freeze.ProjectEnvironmentScope {
		projectScopes[projectId] = util.FlattenStringList(environmentScopes)
	}

	scopeType, scopeDiags := types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, projectScopes)
	if scopeDiags.HasError() {
		diags.Append(scopeDiags...)
		return nil, diags
	}

	tenantScopes := make([]attr.Value, 0)
	for _, scope := range freeze.TenantProjectEnvironmentScope {
		tenantScope, tDiags := types.ObjectValue(tenantScopeObjectType(), map[string]attr.Value{
			"tenant_id":      types.StringValue(scope.TenantId),
			"project_id":     types.StringValue(scope.ProjectId),
			"environment_id": types.StringValue(scope.EnvironmentId),
		})
		if tDiags.HasError() {
			diags.Append(tDiags...)
			return nil, diags
		}
		tenantScopes = append(tenantScopes, tenantScope)
	}

	tenantScopesList, tsDiags := types.ListValue(
		types.ObjectType{AttrTypes: tenantScopeObjectType()},
		tenantScopes,
	)
	if tsDiags.HasError() {
		diags.Append(tsDiags...)
		return nil, diags
	}

	attrs := map[string]attr.Value{
		"id":                               types.StringValue(freeze.ID),
		"name":                             types.StringValue(freeze.Name),
		"start":                            types.StringValue(freeze.Start.Format(time.RFC3339)),
		"end":                              types.StringValue(freeze.End.Format(time.RFC3339)),
		"project_environment_scope":        scopeType,
		"tenant_project_environment_scope": tenantScopesList,
	}

	if freeze.RecurringSchedule != nil {
		daysOfWeek, daysDiags := types.ListValueFrom(ctx, types.StringType, freeze.RecurringSchedule.DaysOfWeek)
		if daysDiags.HasError() {
			diags.Append(daysDiags...)
			return nil, diags
		}

		endOnDate := types.StringNull()
		if freeze.RecurringSchedule.EndOnDate != nil {
			endOnDate = types.StringValue(freeze.RecurringSchedule.EndOnDate.Format(time.RFC3339))
		}

		endAfterOccurrences := types.Int64Value(int64(freeze.RecurringSchedule.EndAfterOccurrences))

		monthlyScheduleType := types.StringNull()
		if freeze.RecurringSchedule.MonthlyScheduleType != "" {
			monthlyScheduleType = types.StringValue(freeze.RecurringSchedule.MonthlyScheduleType)
		}

		dateOfMonth := types.StringValue(freeze.RecurringSchedule.DateOfMonth)

		dayNumberOfMonth := types.StringValue(freeze.RecurringSchedule.DayNumberOfMonth)

		dayOfWeek := types.StringValue(freeze.RecurringSchedule.DayOfWeek)

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

		recurringSchedule, rsDiags := types.ObjectValue(freezeRecurringScheduleObjectType(), scheduleAttrs)
		if rsDiags.HasError() {
			diags.Append(rsDiags...)
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

func tenantScopeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"tenant_id":      types.StringType,
		"project_id":     types.StringType,
		"environment_id": types.StringType,
	}
}

func freezeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                               types.StringType,
		"name":                             types.StringType,
		"start":                            types.StringType,
		"end":                              types.StringType,
		"project_environment_scope":        types.MapType{ElemType: types.ListType{ElemType: types.StringType}},
		"tenant_project_environment_scope": types.ListType{ElemType: types.ObjectType{AttrTypes: tenantScopeObjectType()}},
		"recurring_schedule":               types.ObjectType{AttrTypes: freezeRecurringScheduleObjectType()},
	}
}
