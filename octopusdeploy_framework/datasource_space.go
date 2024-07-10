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

func (*spaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("space")
}

func (*spaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides information about an existing space.",
		Attributes:  schemas.GetSpaceDatasourceSchema(),
	}
}

func (b *spaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	mapSpace(ctx, &data, matchedSpace)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapSpace(ctx context.Context, data *schemas.SpaceModel, matchedSpace *spaces.Space) {
	data.ID = types.StringValue(matchedSpace.ID)
	data.Description = types.StringValue(matchedSpace.Description)
	data.Slug = types.StringValue(matchedSpace.Slug)
	data.IsTaskQueueStopped = types.BoolValue(matchedSpace.TaskQueueStopped)
	data.IsDefault = types.BoolValue(matchedSpace.IsDefault)
	data.SpaceManagersTeamMembers, _ = types.ListValueFrom(ctx, types.StringType, matchedSpace.SpaceManagersTeamMembers)
	data.SpaceManagersTeams, _ = types.ListValueFrom(ctx, types.StringType, matchedSpace.SpaceManagersTeams)
}
