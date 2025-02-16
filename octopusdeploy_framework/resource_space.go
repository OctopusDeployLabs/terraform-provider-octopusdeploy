package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const spaceManagersTeamIDPrefix = "teams-spacemanagers-"

type spaceResource struct {
	*Config
}

var _ resource.Resource = &spaceResource{}
var _ resource.ResourceWithImportState = &spaceResource{}

func NewSpaceResource() resource.Resource {
	return &spaceResource{}
}

func (s *spaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.SpaceSchema{}.GetResourceSchema()
}

func (s *spaceResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("space")
}

func (s *spaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	s.Config = ResourceConfiguration(req, resp)
}

func (s *spaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.SpaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	util.Create(ctx, "space", data)

	newSpace := spaces.NewSpace(data.Name.ValueString())
	newSpace.Slug = data.Slug.ValueString()
	newSpace.Description = data.Description.ValueString()
	newSpace.IsDefault = data.IsDefault.ValueBool()
	newSpace.TaskQueueStopped = data.IsTaskQueueStopped.ValueBool()

	convertedTeams, diags := util.SetToStringArray(ctx, data.SpaceManagersTeams)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	newSpace.SpaceManagersTeams = convertedTeams

	convertedTeamMembers, diags := util.SetToStringArray(ctx, data.SpaceManagersTeamMembers)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	newSpace.SpaceManagersTeamMembers = convertedTeamMembers

	tflog.Debug(ctx, fmt.Sprintf("creating space %#v", newSpace))

	createdSpace, err := s.Client.Spaces.Add(newSpace)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, s.Config.SystemInfo, "unable to create new space", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("resulting space %#v", createdSpace))

	// the result of a space add operation seems to not return the correct values for teams, but the get does
	createdSpace, _ = spaces.GetByID(s.Client, createdSpace.ID)

	if data.IsTaskQueueStopped.ValueBool() == true {
		// a space can't have a stopped task queue via the create, need to do a subsequent update
		createdSpace.TaskQueueStopped = true
		_, err = spaces.Update(s.Client, createdSpace)
		if err != nil {
			util.AddDiagnosticError(resp.Diagnostics, s.Config.SystemInfo, "Error updating space task queue", err.Error())
		}
		createdSpace, _ = spaces.GetByID(s.Client, createdSpace.ID)
		tflog.Debug(ctx, fmt.Sprintf("resulting space after setting task queue stopped %#v", createdSpace))
	}

	data.ID = types.StringValue(createdSpace.ID)
	data.Description = types.StringValue(createdSpace.Description)
	data.Slug = types.StringValue(createdSpace.Slug)
	data.IsTaskQueueStopped = types.BoolValue(createdSpace.TaskQueueStopped)
	data.IsDefault = types.BoolValue(createdSpace.IsDefault)

	data.SpaceManagersTeamMembers, diags = types.SetValueFrom(ctx, types.StringType, createdSpace.SpaceManagersTeamMembers)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	effectiveTeams := s.getEffectiveTeams(ctx, createdSpace)
	data.SpaceManagersTeams = types.SetValueMust(types.StringType, effectiveTeams)

	tflog.Debug(ctx, fmt.Sprintf("state space %#v", data))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	util.Created(ctx, "space")
}

func (s *spaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.SpaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading space (%s)", data.ID))

	spaceResult, err := spaces.GetByID(s.Client, data.ID.ValueString())

	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "space"); err != nil {
			util.AddDiagnosticError(resp.Diagnostics, s.Config.SystemInfo, "unable to query spaces", err.Error())
		}
		return
	}

	data.Name = types.StringValue(spaceResult.Name)
	data.Description = types.StringValue(spaceResult.Description)
	data.Slug = types.StringValue(spaceResult.Slug)
	data.IsTaskQueueStopped = types.BoolValue(spaceResult.TaskQueueStopped)
	data.IsDefault = types.BoolValue(spaceResult.IsDefault)
	data.SpaceManagersTeamMembers, _ = types.SetValueFrom(ctx, types.StringType, spaceResult.SpaceManagersTeamMembers)

	effectiveTeams := s.getEffectiveTeams(ctx, spaceResult)
	data.SpaceManagersTeams = types.SetValueMust(types.StringType, effectiveTeams)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Info(ctx, fmt.Sprintf("space read (%s)", data.ID))
}

func (s *spaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// read plan and state
	var plan, state schemas.SpaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get existing resource from api
	spaceResult, err := spaces.GetByID(s.Client, state.ID.ValueString())
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, s.Config.SystemInfo, "unable to query spaces", err.Error())
		return
	}

	// update the api resource
	spaceResult.Name = plan.Name.ValueString()
	spaceResult.Slug = plan.Slug.ValueString()
	spaceResult.Description = plan.Description.ValueString()
	spaceResult.IsDefault = plan.IsDefault.ValueBool()
	spaceResult.TaskQueueStopped = plan.IsTaskQueueStopped.ValueBool()

	convertedTeams, diags := util.SetToStringArray(ctx, plan.SpaceManagersTeams)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	spaceResult.SpaceManagersTeams = addSpaceManagers(spaceResult.ID, convertedTeams)

	convertedTeamMembers, diags := util.SetToStringArray(ctx, plan.SpaceManagersTeamMembers)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	spaceResult.SpaceManagersTeamMembers = convertedTeamMembers

	// push to api
	tflog.Debug(ctx, fmt.Sprintf("update: spaceResult before update: %#v", spaceResult))
	_, err = spaces.Update(s.Client, spaceResult)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, s.Config.SystemInfo, "unable to update space", err.Error())
		return
	}

	// refresh from the api
	updatedSpace, _ := spaces.GetByID(s.Client, state.ID.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("Update: updatedSpace: %+v", updatedSpace))

	// update the plan for the managers teams
	plan.ID = types.StringValue(spaceResult.ID)
	plan.Name = types.StringValue(spaceResult.Name)
	plan.Description = types.StringValue(spaceResult.Description)
	plan.Slug = types.StringValue(spaceResult.Slug)
	plan.IsTaskQueueStopped = types.BoolValue(spaceResult.TaskQueueStopped)
	plan.IsDefault = types.BoolValue(spaceResult.IsDefault)
	plan.SpaceManagersTeamMembers, _ = types.SetValueFrom(ctx, types.StringType, spaceResult.SpaceManagersTeamMembers)

	effectiveTeams := s.getEffectiveTeams(ctx, updatedSpace)
	plan.SpaceManagersTeams = types.SetValueMust(types.StringType, effectiveTeams)

	// save plan to state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (s *spaceResource) getEffectiveTeams(ctx context.Context, updatedSpace *spaces.Space) []attr.Value {
	effectiveTeamMembers := make([]attr.Value, 0)
	for _, t := range removeSpaceManagers(ctx, updatedSpace.SpaceManagersTeams) {
		effectiveTeamMembers = append(effectiveTeamMembers, types.StringValue(t))
	}
	return effectiveTeamMembers
}

func (s *spaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.SpaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("deleting space (%s)", data.ID.ValueString()))

	space, err := spaces.GetByID(s.Client, data.ID.ValueString())
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, s.Config.SystemInfo, "unable to read space", err.Error())
	}

	space.TaskQueueStopped = true

	_, err = spaces.Update(s.Client, space)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, s.Config.SystemInfo, "unable to stop task queue", err.Error())
		return
	}

	if err := s.Client.Spaces.DeleteByID(data.ID.ValueString()); err != nil {
		util.AddDiagnosticError(resp.Diagnostics, s.Config.SystemInfo, "unable to delete space", err.Error())
		return
	}

	tflog.Info(ctx, "space deleted ")
}

func (s *spaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func removeSpaceManagers(ctx context.Context, teamIDs []string) []string {
	if len(teamIDs) == 0 {
		return teamIDs
	}
	var newSlice []string

	for _, v := range teamIDs {
		if !strings.Contains(v, spaceManagersTeamIDPrefix) {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

func addSpaceManagers(spaceID string, teamIDs []string) []string {
	var newSlice []string
	if util.GetStringOrEmpty(spaceID) != "" {
		newSlice = append(newSlice, spaceManagersTeamIDPrefix+spaceID)
	}
	newSlice = append(newSlice, teamIDs...)
	return newSlice
}
