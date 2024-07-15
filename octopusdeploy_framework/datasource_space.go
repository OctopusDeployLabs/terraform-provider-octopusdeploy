package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type spaceDataSource struct {
	*Config
}

func NewSpaceDataSource() datasource.DataSource {
	return &spaceDataSource{}
}

func (*spaceDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("space")
}

func (*spaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides information about an existing space.",
		Attributes:  schemas.GetSpaceDatasourceSchema(),
	}
}

func (b *spaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	b.Config = DataSourceConfiguration(req, resp)
}

func (b *spaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data schemas.SpaceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// construct query
	query := spaces.SpacesQuery{PartialName: data.Name.ValueString()}
	spacesResult, err := spaces.Get(b.Client, query)

	if err != nil {
		resp.Diagnostics.AddError("unable to query spaces", err.Error())
		return
	}

	var matchedSpace *spaces.Space
	for _, spaceResult := range spacesResult.Items {
		if strings.EqualFold(spaceResult.Name, data.Name.ValueString()) {
			matchedSpace = spaceResult
		}
	}

	mapSpaceToState(ctx, &data, matchedSpace)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapSpaceToState(ctx context.Context, data *schemas.SpaceModel, space *spaces.Space) {
	data.ID = types.StringValue(space.ID)
	data.Description = types.StringValue(space.Description)
	data.Slug = types.StringValue(space.Slug)
	data.IsTaskQueueStopped = types.BoolValue(space.TaskQueueStopped)
	data.IsDefault = types.BoolValue(space.IsDefault)
	data.SpaceManagersTeamMembers, _ = types.ListValueFrom(ctx, types.StringType, space.SpaceManagersTeamMembers)
	data.SpaceManagersTeams, _ = types.ListValueFrom(ctx, types.StringType, space.SpaceManagersTeams)
}

func mapSpaceFromState(ctx context.Context, data *schemas.SpaceModel, space *spaces.Space) {
	space.ID = data.ID.ValueString()
	space.Name = data.Name.ValueString()
	space.Description = data.Description.ValueString()
	space.Slug = data.Slug.ValueString()
	space.IsDefault = data.IsDefault.ValueBool()
	space.TaskQueueStopped = data.IsTaskQueueStopped.ValueBool()

	for _, t := range data.SpaceManagersTeams.Elements() {
		space.SpaceManagersTeams = append(space.SpaceManagersTeams, t.String())
	}

	for _, t := range data.SpaceManagersTeamMembers.Elements() {
		space.SpaceManagersTeamMembers = append(space.SpaceManagersTeamMembers, t.String())
	}
}
