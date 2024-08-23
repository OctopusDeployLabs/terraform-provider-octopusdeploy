package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var _ datasource.DataSource = &projectsDataSource{}

type projectsDataSource struct {
	*Config
}

type projectsDataSourceModel struct {
	ID                  types.String           `tfsdk:"id"`
	SpaceID             types.String           `tfsdk:"space_id"`
	ClonedFromProjectID types.String           `tfsdk:"cloned_from_project_id"`
	IDs                 types.List             `tfsdk:"ids"`
	IsClone             types.Bool             `tfsdk:"is_clone"`
	Name                types.String           `tfsdk:"name"`
	PartialName         types.String           `tfsdk:"partial_name"`
	Skip                types.Int64            `tfsdk:"skip"`
	Take                types.Int64            `tfsdk:"take"`
	Projects            []projectResourceModel `tfsdk:"projects"`
}

func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

func (p *projectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProjectDataSourceName)
}

func (p *projectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema.Schema{
		Description: "Provides information about existing Octopus Deploy projects.",
		Attributes:  schemas.ProjectSchema{}.GetDatasourceSchemaAttributes(),
	}
}

func (p *projectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	p.Config = DataSourceConfiguration(req, resp)
}

func (p *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data projectsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := projects.ProjectsQuery{
		ClonedFromProjectID: data.ClonedFromProjectID.ValueString(),
		IsClone:             data.IsClone.ValueBool(),
		Name:                data.Name.ValueString(),
		PartialName:         data.PartialName.ValueString(),
		Skip:                int(data.Skip.ValueInt64()),
		Take:                int(data.Take.ValueInt64()),
	}

	util.DatasourceReading(ctx, "projects", query)

	if !data.IDs.IsNull() {
		var ids []string
		resp.Diagnostics.Append(data.IDs.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		query.IDs = ids
	}

	spaceID := data.SpaceID.ValueString()

	existingProjects, err := projects.Get(p.Client, spaceID, query)
	if err != nil {
		resp.Diagnostics.AddError("Unable to query projects", err.Error())
		return
	}

	util.DatasourceResultCount(ctx, "projects", len(existingProjects.Items))

	data.Projects = make([]projectResourceModel, 0, len(existingProjects.Items))
	for _, project := range existingProjects.Items {
		flattenedProject, diags := flattenProject(ctx, project, &projectResourceModel{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Projects = append(data.Projects, *flattenedProject)
	}

	data.ID = types.StringValue(fmt.Sprintf("Projects-%s", time.Now().UTC().String()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
