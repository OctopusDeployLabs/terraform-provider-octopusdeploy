package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
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

func (*spacesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("spaces")
}

func (*spacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// request
			"ids":          schemas.GetQueryIDsDatasourceSchema(),
			"partial_name": schemas.GetQueryPartialNameDatasourceSchema(),
			"skip":         schemas.GetQuerySkipDatasourceSchema(),
			"take":         schemas.GetQueryTakeDatasourceSchema(),

			// response
			"id": schemas.GetIdDatasourceSchema(true),
			"spaces": schema.ListNestedAttribute{
				Computed: true,
				Optional: false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: schemas.GetSpacesDatasourceSchema(),
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
		Skip:        schemas.GetNumber(data.Skip),
		Take:        schemas.GetNumber(data.Take),
	}

	util.DatasourceReading(ctx, "spaces", query)

	existingSpaces, err := spaces.Get(b.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load spaces", err.Error())
		return
	}

	var mappedSpaces []schemas.SpaceModel
	for _, space := range existingSpaces.Items {
		var s schemas.SpaceModel
		mapSpaceToState(ctx, &s, space)
		mappedSpaces = append(mappedSpaces, s)
	}

	util.DatasourceResultCount(ctx, "spaces", len(mappedSpaces))

	data.Spaces, _ = types.ListValueFrom(ctx, schemas.GetSpaceTypeAttributes(), mappedSpaces)
	data.ID = types.StringValue("Spaces " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
