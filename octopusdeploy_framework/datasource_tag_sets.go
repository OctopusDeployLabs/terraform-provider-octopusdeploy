package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var _ datasource.DataSource = &tagSetsDataSource{}

type tagSetsDataSource struct {
	*Config
}

func NewTagSetsDataSource() datasource.DataSource {
	return &tagSetsDataSource{}
}

func (t *tagSetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.TagSetDataSourceName)
}

func (t *tagSetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.GetTagSetDataSourceSchema()
}

func (t *tagSetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	t.Config = DataSourceConfiguration(req, resp)
}

func (t *tagSetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data schemas.TagSetDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := tagsets.TagSetsQuery{
		IDs:         schemas.GetIds(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}
	spaceID := data.SpaceID.ValueString()

	existingTagSets, err := tagsets.Get(t.Client, spaceID, query)
	if err != nil {
		resp.Diagnostics.AddError("Unable to query tag sets", err.Error())
		return
	}

	data.TagSets = flattenTagSets(ctx, existingTagSets.Items)

	data.ID = types.StringValue(fmt.Sprintf("TagSets-%s", time.Now().UTC().String()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenTagSets(ctx context.Context, tagSets []*tagsets.TagSet) types.List {
	tfList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.GetTagSetAttrTypes()}, tagSets)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: schemas.GetTagSetAttrTypes()})
	}
	return tfList
}
