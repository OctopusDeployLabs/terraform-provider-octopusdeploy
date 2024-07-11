package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

const spaceManagersTeamIDPrefix = "teams-spacemanagers-"

type spaceResource struct {
	*Config
}

func NewSpaceResource() resource.Resource {
	return &spaceResource{}
}

func (s *spaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource manages spaces in Octopus Deploy.",
		Attributes:  schemas.GetSpaceResourceSchema(),
	}
}

func (s *spaceResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("space")
}

func (s *spaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	s.Config = ResourceConfiguration(req, resp)
}

func (s *spaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.SpaceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	//tflog.Info(ctx, fmt.Sprintf("creating space: %#v", data))
	util.Create(ctx, "space", data)

	newSpace := spaces.NewSpace(data.Name.ValueString())
	newSpace.Slug = data.Slug.ValueString()
	newSpace.Description = data.Description.ValueString()
	newSpace.IsDefault = data.IsDefault.ValueBool()
	newSpace.TaskQueueStopped = data.IsTaskQueueStopped.ValueBool()

	teams := make([]types.String, 0, len(data.SpaceManagersTeams.Elements()))
	diags := data.SpaceManagersTeams.ElementsAs(ctx, &teams, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	convertTeams := make([]string, 0)
	for _, t := range teams {
		convertTeams = append(convertTeams, t.ValueString())
	}
	newSpace.SpaceManagersTeams = convertTeams

	teamMembers := make([]types.String, 0, len(data.SpaceManagersTeamMembers.Elements()))
	diags = data.SpaceManagersTeamMembers.ElementsAs(ctx, &teamMembers, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	convertTeamMembers := make([]string, 0)
	for _, t := range teamMembers {
		convertTeamMembers = append(convertTeamMembers, t.ValueString())
	}
	newSpace.SpaceManagersTeamMembers = convertTeamMembers

	tflog.Debug(ctx, fmt.Sprintf("creating space %#v", newSpace))

	createdSpace, err := s.Client.Spaces.Add(newSpace)
	if err != nil {
		resp.Diagnostics.AddError("unable to create new space", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("resulting space %#v", createdSpace))

	data.ID = types.StringValue(createdSpace.ID)
	data.Description = types.StringValue(createdSpace.Description)
	data.Slug = types.StringValue(createdSpace.Slug)
	data.IsTaskQueueStopped = types.BoolValue(createdSpace.TaskQueueStopped)
	data.IsDefault = types.BoolValue(createdSpace.IsDefault)
	data.SpaceManagersTeamMembers, _ = types.SetValueFrom(ctx, types.StringType, createdSpace.SpaceManagersTeamMembers)
	data.SpaceManagersTeams, _ = types.SetValueFrom(ctx, types.StringType, removeSpaceManagers(ctx, createdSpace.SpaceManagersTeams))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	//tflog.Info(ctx, fmt.Sprintf("space created (%s)", data.ID.ValueString()))
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
		resp.Diagnostics.AddError("unable to query spaces", err.Error())
		return
	}

	data.ID = types.StringValue(spaceResult.ID)
	data.Description = types.StringValue(spaceResult.Description)
	data.Slug = types.StringValue(spaceResult.Slug)
	data.IsTaskQueueStopped = types.BoolValue(spaceResult.TaskQueueStopped)
	data.IsDefault = types.BoolValue(spaceResult.IsDefault)
	data.SpaceManagersTeamMembers, _ = types.SetValueFrom(ctx, types.StringType, spaceResult.SpaceManagersTeamMembers)
	data.SpaceManagersTeams, _ = types.SetValueFrom(ctx, types.StringType, removeSpaceManagers(ctx, spaceResult.SpaceManagersTeams))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Info(ctx, fmt.Sprintf("space read (%s)", data.ID))
}

func (s *spaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schemas.SpaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceResult, err := spaces.GetByID(s.Client, data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("unable to query spaces", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("update: spaceResult before update: %#v", spaceResult))
	updatedSpace, err := spaces.Update(s.Client, spaceResult)
	if err != nil {
		resp.Diagnostics.AddError("unable to update space", err.Error())
		return
	}

	updatedSpace, err = spaces.GetByID(s.Client, data.ID.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("Update: updatedSpace: %+v", updatedSpace))

	data.ID = types.StringValue(updatedSpace.ID)
	data.Description = types.StringValue(updatedSpace.Description)
	data.Slug = types.StringValue(updatedSpace.Slug)
	data.IsTaskQueueStopped = types.BoolValue(updatedSpace.TaskQueueStopped)
	data.IsDefault = types.BoolValue(updatedSpace.IsDefault)
	data.SpaceManagersTeamMembers, _ = types.SetValueFrom(ctx, types.StringType, updatedSpace.SpaceManagersTeamMembers)
	data.SpaceManagersTeams, _ = types.SetValueFrom(ctx, types.StringType, removeSpaceManagers(ctx, updatedSpace.SpaceManagersTeams))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
		resp.Diagnostics.AddError("unable to read space", err.Error())
	}

	space.TaskQueueStopped = true

	_, err = spaces.Update(s.Client, space)
	if err != nil {
		resp.Diagnostics.AddError("unable to stop task queue", err.Error())
		return
	}

	if err := s.Client.Spaces.DeleteByID(data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete space", err.Error())
		return
	}

	tflog.Info(ctx, "space deleted ")
}

func removeSpaceManagers(ctx context.Context, teamIDs []string) []string {
	if len(teamIDs) == 0 {
		return teamIDs
	}
	var newSlice []string

	tflog.Debug(ctx, fmt.Sprintf("before removeSpaceManagers: %#v", teamIDs))
	for _, v := range teamIDs {
		if !strings.Contains(v, spaceManagersTeamIDPrefix) {
			newSlice = append(newSlice, v)
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("after removeSpaceManagers: %#v", newSlice))
	return newSlice
}
