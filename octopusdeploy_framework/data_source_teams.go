package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/teams"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"time"
)

type teamDataSource struct {
	*Config
}

type teamsDataSourceModel struct {
	ID            types.String                      `tfsdk:"id"`
	IDs           types.List                        `tfsdk:"ids"`
	IncludeSystem types.Bool                        `tfsdk:"include_system"`
	PartialName   types.String                      `tfsdk:"partial_name"`
	Spaces        types.List                        `tfsdk:"spaces"`
	Skip          types.Int64                       `tfsdk:"skip"`
	Take          types.Int64                       `tfsdk:"take"`
	Teams         []schemas.TeamTypeDatasourceModel `tfsdk:"teams"`
}

var _ datasource.DataSource = &teamDataSource{}

func NewTeamsDataSource() datasource.DataSource {
	return &teamDataSource{}
}

func (t *teamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("teams")
}

func (t *teamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.TeamSchema{}.GetDatasourceSchema()
}

func (t *teamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	t.Config = DataSourceConfiguration(req, resp)
}

func (t *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data teamsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := teams.TeamsQuery{
		IDs:           util.GetIds(data.IDs),
		IncludeSystem: data.IncludeSystem.ValueBool(),
		PartialName:   data.PartialName.ValueString(),
		Spaces:        util.ExpandStringList(data.Spaces),
		Skip:          util.GetNumber(data.Skip),
		Take:          util.GetNumber(data.Take),
	}

	util.DatasourceReading(ctx, "teams", query)

	existingTeams, err := t.Client.Teams.Get(query)

	if err != nil {
		diag.FromErr(err)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("users returned from API: %#v", existingTeams))

	data.Teams = make([]schemas.TeamTypeDatasourceModel, 0, len(existingTeams.Items))
	for _, team := range existingTeams.Items {
		data.Teams = append(data.Teams, schemas.MapToTeamsDatasourceModel(team))
	}

	util.DatasourceResultCount(ctx, "teams", len(data.Teams))

	data.ID = types.StringValue("Teams " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
