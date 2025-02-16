package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deploymentfreezes"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

const deploymentFreezeResourceName = "deployment_freeze"

type recurringScheduleModel struct {
	Type                types.String      `tfsdk:"type"`
	Unit                types.Int64       `tfsdk:"unit"`
	EndType             types.String      `tfsdk:"end_type"`
	EndOnDate           timetypes.RFC3339 `tfsdk:"end_on_date"`
	EndAfterOccurrences types.Int64       `tfsdk:"end_after_occurrences"`
	MonthlyScheduleType types.String      `tfsdk:"monthly_schedule_type"`
	DateOfMonth         types.String      `tfsdk:"date_of_month"`
	DayNumberOfMonth    types.String      `tfsdk:"day_number_of_month"`
	DaysOfWeek          types.List        `tfsdk:"days_of_week"`
	DayOfWeek           types.String      `tfsdk:"day_of_week"`
}

type deploymentFreezeModel struct {
	Name              types.String            `tfsdk:"name"`
	Start             timetypes.RFC3339       `tfsdk:"start"`
	End               timetypes.RFC3339       `tfsdk:"end"`
	RecurringSchedule *recurringScheduleModel `tfsdk:"recurring_schedule"`
	schemas.ResourceModel
}

func getStringPointer(s types.String) *string {
	if s.IsNull() {
		return nil
	}
	value := s.ValueString()
	return &value
}

type deploymentFreezeResource struct {
	*Config
}

var _ resource.Resource = &deploymentFreezeResource{}

func NewDeploymentFreezeResource() resource.Resource {
	return &deploymentFreezeResource{}
}

func (f *deploymentFreezeResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(deploymentFreezeResourceName)
}

func (f *deploymentFreezeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.DeploymentFreezeSchema{}.GetResourceSchema()
}

func (f *deploymentFreezeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	f.Config = ResourceConfiguration(req, resp)
}

func (f *deploymentFreezeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state *deploymentFreezeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentFreeze, err := deploymentfreezes.GetById(f.Config.Client, state.GetID())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "deployment freeze"); err != nil {
			util.AddDiagnosticError(&resp.Diagnostics, f.Config.SystemInfo, "unable to load deployment freeze", err.Error())
		}
		return
	}

	if deploymentFreeze.Name != state.Name.ValueString() {
		state.Name = types.StringValue(deploymentFreeze.Name)
	}

	mapToState(ctx, state, deploymentFreeze)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (f *deploymentFreezeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan *deploymentFreezeModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentFreeze, diags := mapFromState(plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	createdFreeze, err := deploymentfreezes.Add(f.Config.Client, deploymentFreeze)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, f.Config.SystemInfo, "error while creating deployment freeze", err.Error())
		return
	}

	diags.Append(mapToState(ctx, plan, createdFreeze)...)
	if diags.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (f *deploymentFreezeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan *deploymentFreezeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	existingFreeze, err := deploymentfreezes.GetById(f.Config.Client, plan.ID.ValueString())
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, f.Config.SystemInfo, "unable to load deployment freeze", err.Error())
		return
	}

	updatedFreeze, diags := mapFromState(plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Preserve both project and tenant scopes from the existing freeze
	updatedFreeze.ProjectEnvironmentScope = existingFreeze.ProjectEnvironmentScope
	updatedFreeze.TenantProjectEnvironmentScope = existingFreeze.TenantProjectEnvironmentScope

	updatedFreeze.SetID(existingFreeze.GetID())
	updatedFreeze.Links = existingFreeze.Links

	updatedFreeze, err = deploymentfreezes.Update(f.Config.Client, updatedFreeze)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, f.Config.SystemInfo, "error while updating deployment freeze", err.Error())
		return
	}

	diags.Append(mapToState(ctx, plan, updatedFreeze)...)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (f *deploymentFreezeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state *deploymentFreezeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	freeze, err := deploymentfreezes.GetById(f.Config.Client, state.GetID())
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, f.Config.SystemInfo, "unable to load deployment freeze", err.Error())
		return
	}

	err = deploymentfreezes.Delete(f.Config.Client, freeze)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, f.Config.SystemInfo, "unable to delete deployment freeze", err.Error())
	}

	resp.State.RemoveResource(ctx)
}
func mapFromState(state *deploymentFreezeModel) (*deploymentfreezes.DeploymentFreeze, diag.Diagnostics) {
	start, diags := state.Start.ValueRFC3339Time()
	if diags.HasError() {
		return nil, diags
	}
	start = start.UTC()

	end, diags := state.End.ValueRFC3339Time()
	if diags.HasError() {
		return nil, diags
	}
	end = end.UTC()

	freeze := deploymentfreezes.DeploymentFreeze{
		Name:  state.Name.ValueString(),
		Start: &start,
		End:   &end,
	}

	if state.RecurringSchedule != nil {
		var daysOfWeek []string

		if !state.RecurringSchedule.DaysOfWeek.IsNull() {
			diags.Append(state.RecurringSchedule.DaysOfWeek.ElementsAs(context.TODO(), &daysOfWeek, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		freeze.RecurringSchedule = &deploymentfreezes.RecurringSchedule{
			Type:                deploymentfreezes.RecurringScheduleType(state.RecurringSchedule.Type.ValueString()),
			Unit:                int(state.RecurringSchedule.Unit.ValueInt64()),
			EndType:             deploymentfreezes.RecurringScheduleEndType(state.RecurringSchedule.EndType.ValueString()),
			EndAfterOccurrences: getOptionalIntValue(state.RecurringSchedule.EndAfterOccurrences),
			MonthlyScheduleType: getOptionalString(state.RecurringSchedule.MonthlyScheduleType),
			DateOfMonth:         getOptionalString(state.RecurringSchedule.DateOfMonth),
			DayNumberOfMonth:    getOptionalString(state.RecurringSchedule.DayNumberOfMonth),
			DaysOfWeek:          daysOfWeek,
			DayOfWeek:           getOptionalString(state.RecurringSchedule.DayOfWeek),
		}

		if !state.RecurringSchedule.EndOnDate.IsNull() {
			date, diagsDate := state.RecurringSchedule.EndOnDate.ValueRFC3339Time()
			if diagsDate.HasError() {
				diags.Append(diagsDate...)
				return nil, diags
			}
			freeze.RecurringSchedule.EndOnDate = &date
		}
	}

	freeze.ID = state.ID.String()
	return &freeze, nil
}
func mapToState(ctx context.Context, state *deploymentFreezeModel, deploymentFreeze *deploymentfreezes.DeploymentFreeze) diag.Diagnostics {
	state.ID = types.StringValue(deploymentFreeze.ID)
	state.Name = types.StringValue(deploymentFreeze.Name)

	updatedStart, diags := calculateStateTime(ctx, state.Start, *deploymentFreeze.Start)
	if diags.HasError() {
		return diags
	}
	state.Start = updatedStart

	updatedEnd, diags := calculateStateTime(ctx, state.End, *deploymentFreeze.End)
	if diags.HasError() {
		return diags
	}
	state.End = updatedEnd

	if deploymentFreeze.RecurringSchedule != nil {
		var daysOfWeek types.List
		if len(deploymentFreeze.RecurringSchedule.DaysOfWeek) > 0 {
			elements := make([]attr.Value, len(deploymentFreeze.RecurringSchedule.DaysOfWeek))
			for i, day := range deploymentFreeze.RecurringSchedule.DaysOfWeek {
				elements[i] = types.StringValue(day)
			}

			var listDiags diag.Diagnostics
			daysOfWeek, listDiags = types.ListValue(types.StringType, elements)
			if listDiags.HasError() {
				diags.Append(listDiags...)
				return diags
			}
		} else {
			daysOfWeek = types.ListNull(types.StringType)
		}

		state.RecurringSchedule = &recurringScheduleModel{
			Type:                types.StringValue(string(deploymentFreeze.RecurringSchedule.Type)),
			Unit:                types.Int64Value(int64(deploymentFreeze.RecurringSchedule.Unit)),
			EndType:             types.StringValue(string(deploymentFreeze.RecurringSchedule.EndType)),
			DaysOfWeek:          daysOfWeek,
			MonthlyScheduleType: mapOptionalStringValue(deploymentFreeze.RecurringSchedule.MonthlyScheduleType),
		}

		if deploymentFreeze.RecurringSchedule.EndOnDate != nil {
			state.RecurringSchedule.EndOnDate = timetypes.NewRFC3339TimeValue(*deploymentFreeze.RecurringSchedule.EndOnDate)
		} else {
			state.RecurringSchedule.EndOnDate = timetypes.NewRFC3339Null()
		}

		state.RecurringSchedule.EndAfterOccurrences = mapOptionalIntValue(deploymentFreeze.RecurringSchedule.EndAfterOccurrences)
		state.RecurringSchedule.DateOfMonth = mapOptionalStringValue(deploymentFreeze.RecurringSchedule.DateOfMonth)
		state.RecurringSchedule.DayNumberOfMonth = mapOptionalStringValue(deploymentFreeze.RecurringSchedule.DayNumberOfMonth)
		state.RecurringSchedule.DayOfWeek = mapOptionalStringValue(deploymentFreeze.RecurringSchedule.DayOfWeek)
	}

	return nil
}

func calculateStateTime(ctx context.Context, stateValue timetypes.RFC3339, updatedValue time.Time) (timetypes.RFC3339, diag.Diagnostics) {
	stateTime, diags := stateValue.ValueRFC3339Time()
	if diags.HasError() {
		return timetypes.RFC3339{}, diags
	}
	stateTimeUTC := timetypes.NewRFC3339TimeValue(stateTime.UTC())
	updatedValueUTC := updatedValue.UTC()
	valuesAreEqual, diags := stateTimeUTC.StringSemanticEquals(ctx, timetypes.NewRFC3339TimeValue(updatedValueUTC))
	if diags.HasError() {
		return timetypes.NewRFC3339Null(), diags
	}

	if valuesAreEqual {
		return stateValue, diags
	}

	location := stateTime.Location()
	newValue := timetypes.NewRFC3339TimeValue(updatedValueUTC.In(location))
	return newValue, diags
}

func getOptionalStringPointer(value types.String) *string {
	if value.IsNull() {
		return nil
	}
	str := value.ValueString()
	return &str
}
func mapOptionalStringValue(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}
func getOptionalIntValue(value types.Int64) int {
	if value.IsNull() {
		return 0
	}
	return int(value.ValueInt64())
}

func mapOptionalIntValue(value int) types.Int64 {
	if value == 0 {
		return types.Int64Null()
	}
	return types.Int64Value(int64(value))
}

func getOptionalString(value types.String) string {
	if value.IsNull() {
		return ""
	}
	return value.ValueString()
}
