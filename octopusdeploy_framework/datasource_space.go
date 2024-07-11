package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type spaceDataSource struct {
	*Config
}

type spaceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Slug                     types.String `tfsdk:"slug"`
	Description              types.String `tfsdk:"description"`
	IsDefault                types.Bool   `tfsdk:"is_default"`
	SpaceManagersTeams       types.List   `tfsdk:"space_managers_teams"`
	SpaceManagersTeamMembers types.List   `tfsdk:"space_managers_team_members"`
	IsTaskQueueStopped       types.Bool   `tfsdk:"is_task_queue_stopped"`
}

func NewSpaceDataSource() datasource.DataSource {
	return &spaceDataSource{}
}

func (*spaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = ProviderTypeName + "_space"
}

func (*spaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides information about an existing space.",
		Attributes:  getSpaceSchema(),
	}
}

func (b *spaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	b.Config = DataSourceConfiguration(req, resp)
}

func (b *spaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data spaceModel
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

func mapSpace(ctx context.Context, data *spaceModel, matchedSpace *spaces.Space) {
	data.ID = types.StringValue(matchedSpace.ID)
	data.Description = types.StringValue(matchedSpace.Description)
	data.Slug = types.StringValue(matchedSpace.Slug)
	data.IsTaskQueueStopped = types.BoolValue(matchedSpace.TaskQueueStopped)
	data.IsDefault = types.BoolValue(matchedSpace.IsDefault)
	data.SpaceManagersTeamMembers, _ = types.ListValueFrom(ctx, types.StringType, matchedSpace.SpaceManagersTeamMembers)
	data.SpaceManagersTeams, _ = types.ListValueFrom(ctx, types.StringType, matchedSpace.SpaceManagersTeams)
}

func getSpaceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":          util.GetIdDatasourceSchema(),
		"description": util.GetDescriptionDatasourceSchema("space"),
		"name":        util.GetNameDatasourceWithMaxLengthSchema(true, 20),
		"slug":        util.GetSlugDatasourceSchema("space"),
		"space_managers_teams": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of team IDs designated to be managers of this space.",
			Optional:    true,
			Computed:    true,
		},
		"space_managers_team_members": schema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of user IDs designated to be managers of this space.",
			Optional:    true,
			Computed:    true,
		},
		"is_task_queue_stopped": schema.BoolAttribute{
			Description: "Specifies the status of the task queue for this space.",
			Optional:    true,
		},
		"is_default": schema.BoolAttribute{
			Description: "Specifies if this space is the default space in Octopus.",
			Optional:    true,
		},
	}
}
