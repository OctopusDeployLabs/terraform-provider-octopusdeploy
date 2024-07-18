package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &gitCredentialResource{}

type gitCredentialResource struct {
	*Config
}

type gitCredentialResourceModel struct {
	ID          types.String `tfsdk:"id"`
	SpaceID     types.String `tfsdk:"space_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
}

func NewGitCredentialResource() resource.Resource {
	return &gitCredentialResource{}
}

func (g *gitCredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.GitCredentialResourceName)
}

func (g *gitCredentialResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetGitCredentialResourceSchema()
}

func (g *gitCredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	g.Config = ResourceConfiguration(req, resp)
}

func (g *gitCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan gitCredentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gitCredential := expandGitCredential(&plan)
	createdCredential, err := credentials.Add(g.Client, gitCredential)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Git credential", err.Error())
		return
	}

	setGitCredential(ctx, &plan, createdCredential)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (g *gitCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state gitCredentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource, err := credentials.GetByID(g.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading Git credential", err.Error())
		return
	}

	setGitCredential(ctx, &state, resource)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (g *gitCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan gitCredentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource := expandGitCredential(&plan)
	updatedResource, err := credentials.Update(g.Client, resource)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Git credential", err.Error())
		return
	}

	setGitCredential(ctx, &plan, updatedResource)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (g *gitCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state gitCredentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := credentials.DeleteByID(g.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Git credential", err.Error())
		return
	}
}

func expandGitCredential(model *gitCredentialResourceModel) *credentials.Resource {
	password := core.NewSensitiveValue(model.Password.ValueString())
	name := model.Name.ValueString()
	username := model.Username.ValueString()

	usernamePassword := credentials.NewUsernamePassword(username, password)

	resource := credentials.NewResource(name, usernamePassword)
	resource.ID = model.ID.ValueString()
	resource.Description = model.Description.ValueString()
	resource.SpaceID = model.SpaceID.ValueString()

	return resource
}
func setGitCredential(ctx context.Context, model *gitCredentialResourceModel, resource *credentials.Resource) {
	if resource == nil {
		tflog.Warn(ctx, "Resource is nil in setGitCredential")
		return
	}

	model.ID = types.StringValue(resource.GetID())
	model.SpaceID = types.StringValue(resource.SpaceID)
	model.Name = types.StringValue(resource.GetName())
	model.Description = types.StringValue(resource.Description)

	if resource.Details != nil {
		model.Type = types.StringValue(string(resource.Details.Type()))

		if usernamePassword, ok := resource.Details.(*credentials.UsernamePassword); ok && usernamePassword != nil {
			model.Username = types.StringValue(usernamePassword.Username)
		} else {
			tflog.Debug(ctx, "Git credential is not of type UsernamePassword", map[string]interface{}{
				"type": resource.Details.Type(),
			})
		}
	} else {
		tflog.Warn(ctx, "Git credential details are nil")
	}

	tflog.Debug(ctx, "Git credential state set", map[string]interface{}{
		"id":          model.ID.ValueString(),
		"name":        model.Name.ValueString(),
		"space_id":    model.SpaceID.ValueString(),
		"type":        model.Type.ValueString(),
		"description": model.Description.ValueString(),
		"username":    model.Username.ValueString(),
	})
}
