package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	resp.Schema = schemas.TagSetSchema{}.GetDatasourceSchema()
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
		IDs:         util.GetIds(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}

	util.DatasourceReading(ctx, "tag sets", query)

	spaceID := data.SpaceID.ValueString()

	existingTagSets, err := tagsets.Get(t.Client, spaceID, query)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, t.Config.SystemInfo, "Unable to query tag sets", err.Error())
		return
	}

	util.DatasourceResultCount(ctx, "tag sets", len(existingTagSets.Items))

	data.TagSets = flattenTagSets(ctx, existingTagSets.Items)
	data.ID = types.StringValue(fmt.Sprintf("TagSets-%s", time.Now().UTC().String()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenTagSets(ctx context.Context, tagSets []*tagsets.TagSet) types.List {
	if len(tagSets) == 0 {
		emptyList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.GetTagSetAttrTypes()}, tagSets)
		return emptyList
	}

	tfList := make([]attr.Value, len(tagSets))
	for i, tagSet := range tagSets {
		tfList[i] = types.ObjectValueMust(schemas.GetTagSetAttrTypes(), map[string]attr.Value{
			"id":          types.StringValue(tagSet.ID),
			"name":        types.StringValue(tagSet.Name),
			"description": types.StringValue(tagSet.Description),
			"sort_order":  types.Int64Value(int64(tagSet.SortOrder)),
			"space_id":    types.StringValue(tagSet.SpaceID),
		})
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: schemas.GetTagSetAttrTypes()}, tfList)
}
