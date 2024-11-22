package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/proxies"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &machineProxyResource{}

type machineProxyResource struct {
	*Config
}

func NewMachineProxyResource() resource.Resource {
	return &machineProxyResource{}
}

func (r *machineProxyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.MachineProxyResourceName)
}

func (r *machineProxyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.MachineProxySchema{}.GetResourceSchema()
}

func (r *machineProxyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *machineProxyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.MachineProxyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	machineProxy := mapMachineProxyModelToRequest(&plan)
	createdProxy, err := proxies.Add(r.Client, machineProxy)
	if err != nil {
		resp.Diagnostics.AddError("Error creating machine proxy", err.Error())
		return
	}

	proxyModel := mapMachineProxyRequestToModel(createdProxy, &plan)

	diags := resp.State.Set(ctx, proxyModel)
	resp.Diagnostics.Append(diags...)
}

func (r *machineProxyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.MachineProxyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	machineProxy, err := proxies.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "machine proxy"); err != nil {
			resp.Diagnostics.AddError("Error reading machine proxy", err.Error())
		}
		return
	}

	proxyModel := mapMachineProxyRequestToModel(machineProxy, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, proxyModel)...)
}

func (r *machineProxyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.MachineProxyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	existingProxy, err := proxies.GetByID(r.Client, plan.SpaceID.ValueString(), plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving machine proxy", err.Error())
		return
	}

	updatedProxy := mapMachineProxyModelToRequest(&plan)
	updatedProxy.ID = existingProxy.ID
	updatedProxy.Links = existingProxy.Links

	updatedProxy, err = proxies.Update(r.Client, updatedProxy)
	if err != nil {
		resp.Diagnostics.AddError("Error updating machine proxy", err.Error())
		return
	}

	proxyModel := mapMachineProxyRequestToModel(updatedProxy, &plan)

	diags := resp.State.Set(ctx, proxyModel)
	resp.Diagnostics.Append(diags...)
}

func (r *machineProxyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.MachineProxyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := proxies.DeleteByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting machine proxy", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func mapMachineProxyModelToRequest(model *schemas.MachineProxyResourceModel) *proxies.Proxy {
	password := core.NewSensitiveValue(model.Password.ValueString())
	proxy := proxies.NewProxy(model.Name.ValueString(), model.Host.ValueString(), model.Username.ValueString(), password)
	proxy.SpaceID = model.SpaceID.ValueString()
	portNumber := model.Port.ValueInt32()
	proxy.Port = int(portNumber)
	return proxy
}

func mapMachineProxyRequestToModel(proxy *proxies.Proxy, state *schemas.MachineProxyResourceModel) *schemas.MachineProxyResourceModel {
	proxyModel := &schemas.MachineProxyResourceModel{
		SpaceID:  types.StringValue(proxy.SpaceID),
		Name:     types.StringValue(proxy.Name),
		Host:     types.StringValue(proxy.Host),
		Username: types.StringValue(proxy.Username),
		Password: types.StringValue(state.Password.ValueString()),
		Port:     types.Int32Value(int32(proxy.Port)),
	}
	proxyModel.ID = types.StringValue(proxy.ID)

	return proxyModel
}
