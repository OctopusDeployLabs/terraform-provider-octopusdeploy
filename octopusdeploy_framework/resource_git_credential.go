package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
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
	SpaceID     types.String `tfsdk:"space_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`

	schemas.ResourceModel
}

func NewGitCredentialResource() resource.Resource {
	return &gitCredentialResource{}
}

func (g *gitCredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.GitCredentialResourceName)
}

func (g *gitCredentialResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GitCredentialSchema{}.GetResourceSchema()
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

	tflog.Debug(ctx, "Creating Git credential", map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	gitCredential := expandGitCredential(&plan)
	createdResponse, err := credentials.Add(g.Client, gitCredential)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, g.Config.SystemInfo, "Error creating Git credential", err.Error())
		return
	}

	createdGitCredential, err := credentials.GetByID(g.Client, gitCredential.SpaceID, createdResponse.ID)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, g.Config.SystemInfo, "Error retrieving created Git credential", err.Error())
		return
	}

	if createdGitCredential == nil {
		resp.Diagnostics.AddError("Error creating Git credential", "Created resource is nil")
		return
	}

	setGitCredential(ctx, &plan, createdGitCredential)

	tflog.Debug(ctx, "Git credential created", map[string]interface{}{
		"id":          plan.ID.ValueString(),
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (g *gitCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state gitCredentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gitCredential, err := credentials.GetByID(g.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "git credential"); err != nil {
			util.AddDiagnosticError(&resp.Diagnostics, g.Config.SystemInfo, "Error reading Git credential", err.Error())
		}
		return
	}

	setGitCredential(ctx, &state, gitCredential)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (g *gitCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan gitCredentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gitCredential := expandGitCredential(&plan)
	updatedResource, err := credentials.Update(g.Client, gitCredential)
	if err != nil {
		util.AddDiagnosticError(&resp.Diagnostics, g.Config.SystemInfo, "Error updating Git credential", err.Error())
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
		util.AddDiagnosticError(&resp.Diagnostics, g.Config.SystemInfo, "Error deleting Git credential", err.Error())
		return
	}
}

func expandGitCredential(model *gitCredentialResourceModel) *credentials.Resource {
	if model == nil {
		tflog.Error(context.Background(), "Model is nil in expandGitCredential")
		return nil
	}

	password := core.NewSensitiveValue(model.Password.ValueString())
	name := model.Name.ValueString()
	username := model.Username.ValueString()

	usernamePassword := credentials.NewUsernamePassword(username, password)

	gitCredential := credentials.NewResource(name, usernamePassword)

	// Only set these if they're not empty
	if !model.ID.IsNull() {
		gitCredential.ID = model.ID.ValueString()
	}
	if !model.Description.IsNull() {
		gitCredential.Description = model.Description.ValueString()
	}
	if !model.SpaceID.IsNull() {
		gitCredential.SpaceID = model.SpaceID.ValueString()
	}

	tflog.Debug(context.Background(), "Expanded Git credential", map[string]interface{}{
		"id":          gitCredential.ID,
		"name":        gitCredential.Name,
		"description": gitCredential.Description,
		"space_id":    gitCredential.SpaceID,
		"username":    username,
		// Don't log the password
	})

	return gitCredential
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

	tflog.Debug(ctx, "Setting Git credential state", map[string]interface{}{
		"id":          resource.GetID(),
		"space_id":    resource.SpaceID,
		"name":        resource.GetName(),
		"description": resource.Description,
	})

	if usernamePassword, ok := resource.Details.(*credentials.UsernamePassword); ok && usernamePassword != nil {
		model.Username = types.StringValue(usernamePassword.Username)
		// Note: We don't set the password here as it's sensitive and not returned by the API
	} else {
		tflog.Debug(ctx, "Git credential is not of type UsernamePassword", map[string]interface{}{
			"type": resource.Details.Type(),
		})
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
