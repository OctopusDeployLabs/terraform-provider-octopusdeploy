package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

type spacesDataSource struct {
	*Config
}

type spacesModel struct {
	ID          types.String `tfsdk:"id"`
	IDs         types.List   `tfsdk:"ids"`
	PartialName types.String `tfsdk:"partial_name"`
	Skip        types.Int64  `tfsdk:"skip"`
	Take        types.Int64  `tfsdk:"take"`
	Spaces      types.List   `tfsdk:"spaces"`
}

func NewSpacesDataSource() datasource.DataSource {
	return &spacesDataSource{}
}

func (*spacesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("spaces")
}

func (*spacesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// request
			"ids":          util.GetQueryIDsDatasourceSchema(),
			"partial_name": util.GetQueryPartialNameDatasourceSchema(),
			"skip":         util.GetQuerySkipDatasourceSchema(),
			"take":         util.GetQueryTakeDatasourceSchema(),

			// response
			"id": util.GetIdDatasourceSchema(),
		},
		Blocks: map[string]schema.Block{
			"spaces": schema.ListNestedBlock{
				Description: "Provides information about existing spaces.",
				NestedObject: schema.NestedBlockObject{
					Attributes: getSpaceSchema(),
				},
			},
		},
	}
}

func (b *spacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	b.Config = DataSourceConfiguration(req, resp)
}

func (b *spacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data spacesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := spaces.SpacesQuery{
		IDs:         util.GetIds(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        util.GetNumber(data.Skip),
		Take:        util.GetNumber(data.Take),
	}

	existingSpaces, err := spaces.Get(b.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load spaces", err.Error())
		return
	}

	var mappedSpaces []spaceModel
	for _, space := range existingSpaces.Items {
		var s spaceModel
		mapSpace(ctx, &s, space)
		mappedSpaces = append(mappedSpaces, s)
	}

	data.Spaces, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":                          types.StringType,
		"name":                        types.StringType,
		"slug":                        types.StringType,
		"description":                 types.StringType,
		"is_default":                  types.BoolType,
		"space_managers_teams":        types.ListType{ElemType: types.StringType},
		"space_managers_team_members": types.ListType{ElemType: types.StringType},
		"is_task_queue_stopped":       types.BoolType}},
		mappedSpaces)
	data.ID = types.StringValue("Spaces " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
