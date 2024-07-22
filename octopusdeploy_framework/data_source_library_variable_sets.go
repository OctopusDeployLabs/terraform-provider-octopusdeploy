package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/libraryvariablesets"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
)

type libraryVariableSetDataSource struct {
	*Config
}

type libraryVariableSetObjectType struct {
	ContentType         types.String `tfsdk:"content_type"`
	ID                  types.String `tfsdk:"id"`
	SpaceID             types.String `tfsdk:"space_id"`
	IDs                 types.List   `tfsdk:"ids"`
	PartialName         types.String `tfsdk:"partial_name"`
	Skip                types.Int64  `tfsdk:"skip"`
	Take                types.Int64  `tfsdk:"take"`
	LibraryVariableSets types.List   `tfsdk:"library_variable_sets"`
}

func NewLibraryVariableSetDataSource() datasource.DataSource {
	return &libraryVariableSetDataSource{}
}

func (l *libraryVariableSetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "library variable set Metadata")
	resp.TypeName = util.GetTypeName("library_variable_set")
}

func (l *libraryVariableSetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "library variable set Schema")
	resp.Schema = schemas.GetDatasourceLifecycleSchema()
}

func (l *libraryVariableSetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "library variable set datasource Configure")
	l.Config = DataSourceConfiguration(req, resp)
}

func (l *libraryVariableSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "library variable set datasource Read")
	var data libraryVariableSetObjectType
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := variables.LibraryVariablesQuery{
		ContentType: data.ContentType.ValueString(),
		IDs:         util.GetStringSlice(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}

	existingLibraryVariableSets, err := libraryvariablesets.Get(l.Config.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read library variable sets, got error: %s", err))
		return
	}

	flattenedLibraryVariableSets := []interface{}{}
	for _, libraryVariableSet := range existingLibraryVariableSets.Items {
		flattenedLibraryVariableSets = append(flattenedLibraryVariableSets, schemas.FlattenLibraryVariableSet(libraryVariableSet))
	}

	flattenedLibraryVariableSetsValues, _ := types.ListValueFrom(context.Background(),
		types.ObjectType{AttrTypes: schemas.GetLibraryVariableSetObjectType()},
		[]any{flattenedLibraryVariableSets})
	
	data.LibraryVariableSets = flattenedLibraryVariableSetsValues

	data.ID = types.StringValue("Library Variables Sets " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
