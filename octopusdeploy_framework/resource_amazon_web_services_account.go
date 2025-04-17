package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type amazonWebServicesAccountResource struct {
	*Config
}

func NewAmazonWebServicesAccountResource() resource.Resource {
	return &amazonWebServicesAccountResource{}
}

var _ resource.ResourceWithImportState = &amazonWebServicesAccountResource{}

func (r *amazonWebServicesAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("aws_account")
}

func (r *amazonWebServicesAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.AmazonWebServicesAccountSchema{}.GetResourceSchema()
}

func (r *amazonWebServicesAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *amazonWebServicesAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.AmazonWebServicesAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newAccount, err := accounts.NewAmazonWebServicesAccount(data.Name.ValueString(), data.AccessKey.ValueString(), core.NewSensitiveValue(data.SecretKey.ValueString()))

	createdAccount, err := accounts.Add(r.Config.Client, newAccount)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Amazon Web Services account", err.Error())
		return
	}

	data.ID = types.StringValue(createdAccount.GetID())
	data.SpaceId = types.StringValue(createdAccount.GetSpaceID())

	tenantsList := make([]attr.Value, len(createdAccount.GetTenantIDs()))
	for i, tenantId := range createdAccount.GetTenantIDs() {
		tenantsList[i] = types.StringValue(tenantId)
	}

	data.Tenants, _ = types.ListValue(types.StringType, tenantsList)

	data.TenantedDeploymentParticipation = types.StringValue(string(createdAccount.GetTenantedDeploymentMode()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *amazonWebServicesAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	return
}

func (r *amazonWebServicesAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

func (r *amazonWebServicesAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	return
}

func (*amazonWebServicesAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
