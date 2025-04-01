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
)

const projectDeploymentFreezeResourceName = "project_deployment_freeze"

type projectDeploymentFreezeModel struct {
	OwnerID        types.String `tfsdk:"owner_id"`
	EnvironmentIDs types.List   `tfsdk:"environment_ids"`
	deploymentFreezeModel
}

type projectDeploymentFreezeResource struct {
	*Config
}

var _ resource.Resource = &projectDeploymentFreezeResource{}

func NewProjectDeploymentFreezeResource() resource.Resource {
	return &projectDeploymentFreezeResource{}
}

func (f *projectDeploymentFreezeResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(projectDeploymentFreezeResourceName)
}

func (f *projectDeploymentFreezeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProjectDeploymentFreezeSchema{}.GetResourceSchema()
}

func (f *projectDeploymentFreezeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	f.Config = ResourceConfiguration(req, resp)
}

func (f *projectDeploymentFreezeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state *projectDeploymentFreezeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentFreeze, err := deploymentfreezes.GetById(f.Config.Client, state.GetID())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "project deployment freeze"); err != nil {
			resp.Diagnostics.AddError("unable to load project deployment freeze", err.Error())
		}
		return
	}

	mapFromProjectDeploymentFreezeToState(ctx, state, deploymentFreeze)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (f *projectDeploymentFreezeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan *projectDeploymentFreezeModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentFreeze, diags := mapFromStateToProjectDeploymentFreeze(plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	createdFreeze, err := deploymentfreezes.Add(f.Config.Client, deploymentFreeze)
	if err != nil {
		resp.Diagnostics.AddError("error while creating project deployment freeze", err.Error())
		return
	}

	diags.Append(mapFromProjectDeploymentFreezeToState(ctx, plan, createdFreeze)...)
	if diags.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (f *projectDeploymentFreezeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan *projectDeploymentFreezeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	existingFreeze, err := deploymentfreezes.GetById(f.Config.Client, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load project deployment freeze", err.Error())
		return
	}

	updatedFreeze, diags := mapFromStateToProjectDeploymentFreeze(plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Preserve both tenant scopes from the existing freeze
	updatedFreeze.TenantProjectEnvironmentScope = existingFreeze.TenantProjectEnvironmentScope

	updatedFreeze.SetID(existingFreeze.GetID())
	updatedFreeze.Links = existingFreeze.Links

	updatedFreeze, err = deploymentfreezes.Update(f.Config.Client, updatedFreeze)
	if err != nil {
		resp.Diagnostics.AddError("error while updating project deployment freeze", err.Error())
		return
	}

	diags.Append(mapFromProjectDeploymentFreezeToState(ctx, plan, updatedFreeze)...)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (f *projectDeploymentFreezeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state *projectDeploymentFreezeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	freeze, err := deploymentfreezes.GetById(f.Config.Client, state.GetID())
	if err != nil {
		resp.Diagnostics.AddError("unable to load project deployment freeze", err.Error())
		return
	}

	err = deploymentfreezes.Delete(f.Config.Client, freeze)
	if err != nil {
		resp.Diagnostics.AddError("unable to delete project deployment freeze", err.Error())
	}

	resp.State.RemoveResource(ctx)
}
func mapFromStateToProjectDeploymentFreeze(state *projectDeploymentFreezeModel) (*deploymentfreezes.DeploymentFreeze, diag.Diagnostics) {
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
		OwnerId: state.OwnerID.ValueString(),
		Name:    state.Name.ValueString(),
		Start:   &start,
		End:     &end,
	}

	if !state.EnvironmentIDs.IsNull() && !state.EnvironmentIDs.IsUnknown() {
		projectEnvironments := map[string][]string{
			state.OwnerID.ValueString(): util.ExpandStringList(state.EnvironmentIDs),
		}
		freeze.ProjectEnvironmentScope = projectEnvironments
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
			EndAfterOccurrences: util.GetOptionalIntValue(state.RecurringSchedule.EndAfterOccurrences),
			MonthlyScheduleType: util.GetOptionalString(state.RecurringSchedule.MonthlyScheduleType),
			DateOfMonth:         util.GetOptionalString(state.RecurringSchedule.DateOfMonth),
			DayNumberOfMonth:    util.GetOptionalString(state.RecurringSchedule.DayNumberOfMonth),
			DaysOfWeek:          daysOfWeek,
			DayOfWeek:           util.GetOptionalString(state.RecurringSchedule.DayOfWeek),
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

func mapFromProjectDeploymentFreezeToState(ctx context.Context, state *projectDeploymentFreezeModel, deploymentFreeze *deploymentfreezes.DeploymentFreeze) diag.Diagnostics {
	state.ID = types.StringValue(deploymentFreeze.ID)
	state.OwnerID = types.StringValue(deploymentFreeze.OwnerId)
	state.Name = types.StringValue(deploymentFreeze.Name)

	updatedStart, diags := util.CalculateStateTime(ctx, state.Start, *deploymentFreeze.Start)
	if diags.HasError() {
		return diags
	}
	state.Start = updatedStart

	updatedEnd, diags := util.CalculateStateTime(ctx, state.End, *deploymentFreeze.End)
	if diags.HasError() {
		return diags
	}
	state.End = updatedEnd

	if deploymentFreeze.ProjectEnvironmentScope != nil {
		state.EnvironmentIDs = util.FlattenStringList(deploymentFreeze.ProjectEnvironmentScope[state.OwnerID.ValueString()])
	}

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
			MonthlyScheduleType: util.MapOptionalStringValue(deploymentFreeze.RecurringSchedule.MonthlyScheduleType),
		}

		if deploymentFreeze.RecurringSchedule.EndOnDate != nil {
			state.RecurringSchedule.EndOnDate = timetypes.NewRFC3339TimeValue(*deploymentFreeze.RecurringSchedule.EndOnDate)
		} else {
			state.RecurringSchedule.EndOnDate = timetypes.NewRFC3339Null()
		}

		state.RecurringSchedule.EndAfterOccurrences = util.MapOptionalIntValue(deploymentFreeze.RecurringSchedule.EndAfterOccurrences)
		state.RecurringSchedule.DateOfMonth = util.MapOptionalStringValue(deploymentFreeze.RecurringSchedule.DateOfMonth)
		state.RecurringSchedule.DayNumberOfMonth = util.MapOptionalStringValue(deploymentFreeze.RecurringSchedule.DayNumberOfMonth)
		state.RecurringSchedule.DayOfWeek = util.MapOptionalStringValue(deploymentFreeze.RecurringSchedule.DayOfWeek)
	}

	return nil
}
