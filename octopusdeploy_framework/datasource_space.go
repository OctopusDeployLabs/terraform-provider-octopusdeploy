package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	resp.Schema = schemas.SpaceSchema{}.GetDatasourceSchema()
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

	util.DatasourceReading(ctx, "space", query)

	spacesResult, err := spaces.Get(b.Client, query)

	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, b.Config.SystemInfo, "unable to query spaces", err.Error())
		return
	}

	var matchedSpace *spaces.Space
	for _, spaceResult := range spacesResult.Items {
		if strings.EqualFold(spaceResult.Name, data.Name.ValueString()) {
			matchedSpace = spaceResult
		}
	}
	if matchedSpace == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("unable to find space with name %s", data.Name.ValueString()), "")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading space returned ID %s", matchedSpace.ID))

	mapSpaceToState(ctx, &data, matchedSpace)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapSpaceToState(ctx context.Context, data *schemas.SpaceModel, space *spaces.Space) {
	data.ID = types.StringValue(space.ID)
	data.Name = types.StringValue(space.Name)
	data.Description = types.StringValue(space.Description)
	data.Slug = types.StringValue(space.Slug)
	data.IsTaskQueueStopped = types.BoolValue(space.TaskQueueStopped)
	data.IsDefault = types.BoolValue(space.IsDefault)
	data.SpaceManagersTeamMembers, _ = types.SetValueFrom(ctx, types.StringType, space.SpaceManagersTeamMembers)
	data.SpaceManagersTeams, _ = types.SetValueFrom(ctx, types.StringType, space.SpaceManagersTeams)
}
