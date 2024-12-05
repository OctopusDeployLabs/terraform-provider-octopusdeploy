package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deploymentfreezes"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

const deploymentFreezeDatasourceName = "deployment_freezes"

type deploymentFreezesModel struct {
	ID                types.String `tfsdk:"id"`
	IDs               types.List   `tfsdk:"ids"`
	PartialName       types.String `tfsdk:"partial_name"`
	ProjectIDs        types.List   `tfsdk:"project_ids"`
	EnvironmentIDs    types.List   `tfsdk:"environment_ids"`
	IncludeComplete   types.Bool   `tfsdk:"include_complete"`
	Status            types.String `tfsdk:"status"`
	Skip              types.Int64  `tfsdk:"skip"`
	Take              types.Int64  `tfsdk:"take"`
	DeploymentFreezes types.List   `tfsdk:"deployment_freezes"`
}

type deploymentFreezeDataSource struct {
	*Config
}

func (d *deploymentFreezeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = DataSourceConfiguration(req, resp)
}

func NewDeploymentFreezeDataSource() datasource.DataSource {
	return &deploymentFreezeDataSource{}
}

func (d *deploymentFreezeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(deploymentFreezeDatasourceName)
}

func (d *deploymentFreezeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.DeploymentFreezeSchema{}.GetDatasourceSchema()
}

func (d *deploymentFreezeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data deploymentFreezesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := deploymentfreezes.DeploymentFreezeQuery{
		IDs:             util.Ternary(data.IDs.IsNull(), []string{}, util.ExpandStringList(data.IDs)),
		PartialName:     data.PartialName.ValueString(),
		ProjectIds:      util.Ternary(data.ProjectIDs.IsNull(), []string{}, util.ExpandStringList(data.ProjectIDs)),
		EnvironmentIds:  util.Ternary(data.EnvironmentIDs.IsNull(), []string{}, util.ExpandStringList(data.EnvironmentIDs)),
		IncludeComplete: data.IncludeComplete.ValueBool(),
		Status:          data.Status.ValueString(),
		Skip:            int(data.Skip.ValueInt64()),
		Take:            int(data.Take.ValueInt64()),
	}

	util.DatasourceReading(ctx, "deployment freezes", query)

	existingFreezes, err := deploymentfreezes.Get(d.Client, query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load deployment freezes", err.Error())
		return
	}

	flattenedFreezes := []interface{}{}
	for _, freeze := range existingFreezes.DeploymentFreezes {
		flattenedFreeze, diags := mapFreezeToAttribute(ctx, freeze)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		flattenedFreezes = append(flattenedFreezes, flattenedFreeze)
	}

	data.ID = types.StringValue("Deployment Freezes " + time.Now().UTC().String())
	data.DeploymentFreezes, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: freezeObjectType()}, flattenedFreezes)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

var _ datasource.DataSource = &deploymentFreezeDataSource{}
var _ datasource.DataSourceWithConfigure = &deploymentFreezeDataSource{}

func mapFreezeToAttribute(ctx context.Context, freeze deploymentfreezes.DeploymentFreeze) (attr.Value, diag.Diagnostics) {
	projectScopes := make(map[string]attr.Value)
	for projectId, environmentScopes := range freeze.ProjectEnvironmentScope {
		projectScopes[projectId] = util.FlattenStringList(environmentScopes)
	}

	scopeType, diags := types.MapValueFrom(ctx, types.ListType{ElemType: types.StringType}, projectScopes)
	if diags.HasError() {
		return nil, diags
	}

	return types.ObjectValueMust(freezeObjectType(), map[string]attr.Value{
		"id":                        types.StringValue(freeze.ID),
		"name":                      types.StringValue(freeze.Name),
		"start":                     types.StringValue(freeze.Start.Format(time.RFC3339)),
		"end":                       types.StringValue(freeze.End.Format(time.RFC3339)),
		"project_environment_scope": scopeType,
	}), diags
}

func freezeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                        types.StringType,
		"name":                      types.StringType,
		"start":                     types.StringType,
		"end":                       types.StringType,
		"project_environment_scope": types.MapType{ElemType: types.ListType{ElemType: types.StringType}},
	}
}
