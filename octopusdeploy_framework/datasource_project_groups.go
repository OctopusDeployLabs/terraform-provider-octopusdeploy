package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

type projectGroupsDataSource struct {
	*Config
}

type projectGroupsDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	SpaceID       types.String `tfsdk:"space_id"`
	IDs           types.List   `tfsdk:"ids"`
	PartialName   types.String `tfsdk:"partial_name"`
	Skip          types.Int64  `tfsdk:"skip"`
	Take          types.Int64  `tfsdk:"take"`
	ProjectGroups types.List   `tfsdk:"project_groups"`
}

func NewProjectGroupsDataSource() datasource.DataSource {
	return &projectGroupsDataSource{}
}

func getNestedGroupAttributes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"space_id":    types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
	}
}

func (p *projectGroupsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("project_groups")
}

func (p *projectGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.ProjectGroupSchema{}.GetDatasourceSchema()
}

func (p *projectGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	p.Config = DataSourceConfiguration(req, resp)
}

func (p *projectGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data projectGroupsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ids = make([]string, 0, len(data.IDs.Elements()))
	for _, id := range data.IDs.Elements() {
		ids = append(ids, id.String())
	}

	skip := 0
	if !data.Skip.IsNull() {
		skip = int(data.Skip.ValueInt64())
	}

	take := 0
	if !data.Take.IsNull() {
		take = int(data.Take.ValueInt64())
	}

	query := projectgroups.ProjectGroupsQuery{
		IDs:         ids,
		PartialName: data.PartialName.ValueString(),
		Skip:        skip,
		Take:        take,
	}
	spaceID := data.SpaceID.ValueString()

	util.DatasourceReading(ctx, "project groups", query)

	existingProjectGroups, err := projectgroups.Get(p.Client, spaceID, query)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, p.Config.SystemInfo, "unable to load project groups", err.Error())
		return
	}

	newGroups := []schemas.ProjectGroupTypeResourceModel{}
	for _, projectGroup := range existingProjectGroups.Items {
		tflog.Debug(ctx, "loaded group "+projectGroup.Name)
		var g schemas.ProjectGroupTypeResourceModel
		g.ID = types.StringValue(projectGroup.ID)
		g.SpaceID = types.StringValue(projectGroup.SpaceID)
		g.Name = types.StringValue(projectGroup.Name)
		g.Description = types.StringValue(projectGroup.Description)
		newGroups = append(newGroups, g)
	}

	util.DatasourceResultCount(ctx, "project groups", len(newGroups))

	for _, projectGroup := range newGroups {
		tflog.Debug(ctx, "mapped group "+projectGroup.Name.ValueString())
	}
	g, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: getNestedGroupAttributes()}, newGroups)

	data.ProjectGroups = g
	data.ID = types.StringValue("ProjectGroups " + time.Now().UTC().String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
