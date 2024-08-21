package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/scriptmodules"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"time"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type scriptModulesDataSource struct {
	*Config
}

func NewScriptModuleDataSource() datasource.DataSource {
	return &scriptModulesDataSource{}
}

func (l *scriptModulesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "script modules datasource Metadata")
	resp.TypeName = util.GetTypeName("script_modules")
}

func (l *scriptModulesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "script modules datasource Schema")
	resp.Schema = schemas.GetDatasourceScriptModuleSchema()
}

func (l *scriptModulesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "script modules datasource Configure")
	l.Config = DataSourceConfiguration(req, resp)
}

func (l *scriptModulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "script modules datasource Read")
	var data schemas.ScriptModuleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := variables.LibraryVariablesQuery{
		ContentType: "ScriptModule",
		IDs:         util.ExpandStringList(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading script modules data source query: %+v", query))

	spaceID := data.SpaceID.ValueString()
	existingScriptModules, err := scriptmodules.Get(l.Config.Client, spaceID, query)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read script modules, got error: %s", err))
		return
	}

	flattenedScriptModules := []attr.Value{}
	for _, scriptModule := range existingScriptModules.Items {
		flattenedScriptModules = append(flattenedScriptModules, schemas.FlattenScriptModule(scriptModule))
	}

	data.ScriptModules = types.ListValueMust(types.ObjectType{AttrTypes: schemas.ScriptModuleObjectType()},
		flattenedScriptModules)
	data.ID = types.StringValue("Script Modules " + time.Now().UTC().String())

	tflog.Debug(ctx, fmt.Sprintf("Read script modules data source returned %d items", len(existingScriptModules.Items)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
