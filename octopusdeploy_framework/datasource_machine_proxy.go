package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/proxies"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var _ datasource.DataSource = &machineProxyDataSource{}

type machineProxyDataSource struct {
	*Config
}

func NewMachineProxyDataSource() datasource.DataSource {
	return &machineProxyDataSource{}
}

func (p *machineProxyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.MachineProxyDataSourceName)
}

func (p *machineProxyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.MachineProxySchema{}.GetDatasourceSchema()
}

func (p *machineProxyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	p.Config = DataSourceConfiguration(req, resp)
}

func (p *machineProxyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data schemas.MachineProxyDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := proxies.ProxiesQuery{
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}

	util.DatasourceReading(ctx, "machine proxies", query)

	if !data.IDs.IsNull() {
		var ids []string
		resp.Diagnostics.Append(data.IDs.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		query.IDs = ids
	}

	spaceID := data.SpaceID.ValueString()

	proxiesData, err := proxies.Get(p.Client, spaceID, query)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, p.Config.SystemInfo, "Unable to query proxies", err.Error())
		return
	}

	util.DatasourceResultCount(ctx, "proxies", len(proxiesData.Items))

	data.Proxies = make([]schemas.ProxyDatasourceModel, 0, len(proxiesData.Items))
	for _, proxy := range proxiesData.Items {
		proxyModel := mapMachineProxyRequestToModel(proxy, &schemas.MachineProxyResourceModel{})
		data.Proxies = append(data.Proxies, schemas.ProxyDatasourceModel{
			ID:       proxyModel.ID,
			SpaceID:  proxyModel.SpaceID,
			Name:     proxyModel.Name,
			Host:     proxyModel.Host,
			Username: proxyModel.Username,
			Port:     proxyModel.Port,
		})
	}

	data.ID = types.StringValue(fmt.Sprintf("Proxies-%s", time.Now().UTC().String()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
