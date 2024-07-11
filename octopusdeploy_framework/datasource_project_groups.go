package project_group

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeployv6/config"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeployv6/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

type projectGroupsDataSource struct {
	*config.Config
}

type projectGroupsDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	SpaceID       types.String `tfsdk:"space_id"`
	IDs           types.List   `tfsdk:"ids"`
	PartialName   types.String `tfsdk:"partial_name"`
	Skip          types.Int64  `tfsdk:"skip"`
	Take          types.Int64  `tfsdk:"take"`
	ProjectGroups types.List   `tfsdk:"project_groups"`
	//ProjectGroups []projectGroupTypeResourceModel `tfsdk:"project_groups"`
}

func NewProjectGroupsDataSource() datasource.DataSource {
	return &projectGroupsDataSource{}
}

func getNestedGroupAttributes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                  types.StringType,
		"space_id":            types.StringType,
		"name":                types.StringType,
		"retention_policy_id": types.StringType,
		"description":         types.StringType,
	}
}

func (p *projectGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "groups datasource Metadata")
	resp.TypeName = "octopusdeploy_project_groups"
}

func (p *projectGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "groups datasource Schema")
	description := "project group"
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// request
			"space_id":     util.GetSpaceIdDatasourceSchema(description),
			"ids":          util.GetQueryIDsDatasourceSchema(),
			"partial_name": util.GetQueryPartialNameDatasourceSchema(),
			"skip":         util.GetQuerySkipDatasourceSchema(),
			"take":         util.GetQueryTakeDatasourceSchema(),

			// response
			"id": util.GetIdDatasourceSchema(),
		},
		Blocks: map[string]schema.Block{
			"project_groups": schema.ListNestedBlock{
				Description: "A list of project groups that match the filter(s).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":       util.GetIdResourceSchema(),
						"space_id": util.GetSpaceIdResourceSchema(description),
						"name":     util.GetNameResourceSchema(true),
						"retention_policy_id": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "The ID of the retention policy associated with this project group.",
						},
						"description": util.GetDescriptionResourceSchema(description),
					},
				},
			},
		},
	}
}

func (p *projectGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "groups datasource Configure")
	p.Config = dataSourceConfiguration(req, resp)
}

func (p *projectGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "groups datasource Read")
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

	existingProjectGroups, err := projectgroups.Get(p.Client, spaceID, query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load project groups", err.Error())
		return
	}

	var newGroups []projectGroupTypeResourceModel
	for _, projectGroup := range existingProjectGroups.Items {
		tflog.Debug(ctx, "loaded group "+projectGroup.Name)
		var g projectGroupTypeResourceModel
		g.ID = types.StringValue(projectGroup.ID)
		g.SpaceID = types.StringValue(projectGroup.SpaceID)
		g.Name = types.StringValue(projectGroup.Name)
		g.RetentionPolicyID = types.StringValue(projectGroup.RetentionPolicyID)
		g.Description = types.StringValue(projectGroup.Description)
		newGroups = append(newGroups, g)
	}

	//groups, _ := types.ObjectValueFrom(ctx, types.ObjectType{AttrTypes: getNestedGroupAttributes()}, newGroups)
	for _, projectGroup := range newGroups {
		tflog.Debug(ctx, "mapped group "+projectGroup.Name.ValueString())
	}
	g, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: getNestedGroupAttributes()}, newGroups)

	data.ProjectGroups = g
	data.ID = types.StringValue("ProjectGroups " + time.Now().UTC().String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func dataSourceConfiguration(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *config.Config {
	if req.ProviderData == nil {
		return nil
	}

	config, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return config
}
