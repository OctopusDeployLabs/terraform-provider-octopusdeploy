package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

type workersDataSource struct {
	*Config
}

func NewWorkersDataSource() datasource.DataSource {
	return &workersDataSource{}
}

func (*workersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("workers")
}

func (e *workersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	e.Config = DataSourceConfiguration(req, resp)
}

func (*workersDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.WorkersSchema{}.GetDatasourceSchema()
}

func (e *workersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data schemas.WorkersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := machines.WorkersQuery{
		Name:                data.Name.ValueString(),
		IDs:                 util.GetIds(data.IDs),
		PartialName:         data.PartialName.ValueString(),
		Skip:                util.GetNumber(data.Skip),
		Take:                util.GetNumber(data.Take),
		CommunicationStyles: util.ExpandStringList(data.CommunicationStyle),
		HealthStatuses:      util.ExpandStringList(data.HealthStatuses),
		WorkerPoolIDs:       util.ExpandStringList(data.WorkerPoolIDs),
		IsDisabled:          data.IsDisabled.ValueBool(),
		Thumbprint:          data.Thumbprint.ValueString(),
	}

	util.DatasourceReading(ctx, "workers", query)

	existingWorkers, err := workers.Get(e.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load workers", err.Error())
		return
	}

	util.DatasourceResultCount(ctx, "workers", len(existingWorkers.Items))

	workers := []interface{}{}
	for _, worker := range existingWorkers.Items {
		workers = append(workers, schemas.FlattenWorker(worker))
	}

	data.Workers, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.WorkerObjectType()}, workers)
	data.ID = types.StringValue("Workers " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
