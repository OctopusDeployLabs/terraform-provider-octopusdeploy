package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type gitTriggerResource struct {
	*Config
}

func NewGitTriggerResource() resource.Resource {
	return &gitTriggerResource{}
}

var _ resource.ResourceWithImportState = &gitTriggerResource{}

func (r *gitTriggerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("git_trigger")
}

func (r *gitTriggerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GitTriggerSchema{}.GetResourceSchema()
}

func (r *gitTriggerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *gitTriggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *gitTriggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.GitTriggerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	gitTriggerSources := convertListToGitTriggerSources(data.Sources)

	action := actions.NewCreateReleaseAction(data.ChannelId.ValueString())
	filter := filters.NewGitTriggerFilter(gitTriggerSources)

	client := r.Config.Client

	project, err := projects.GetByID(client, data.SpaceId.ValueString(), data.ProjectId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("error finding project", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating Git trigger: %s", data.Name.ValueString()))

	projectTrigger := triggers.NewProjectTrigger(data.Name.ValueString(), data.Description.ValueString(), data.IsDisabled.ValueBool(), project, action, filter)

	createdGitTrigger, err := client.ProjectTriggers.Add(projectTrigger)

	if err != nil {
		resp.Diagnostics.AddError("unable to create Git trigger", err.Error())
		return
	}

	data.ID = types.StringValue(createdGitTrigger.GetID())
	data.Name = types.StringValue(createdGitTrigger.Name)
	data.ProjectId = types.StringValue(createdGitTrigger.ProjectID)
	data.SpaceId = types.StringValue(createdGitTrigger.SpaceID)
	data.IsDisabled = types.BoolValue(createdGitTrigger.IsDisabled)
	data.Sources = convertGitTriggerSourcesToList(createdGitTrigger.Filter.(*filters.GitTriggerFilter).Sources)

	tflog.Info(ctx, fmt.Sprintf("Git trigger created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *gitTriggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.GitTriggerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading Git Trigger (%s)", data.ID))

	client := r.Config.Client

	gitTrigger, err := client.ProjectTriggers.GetByID(data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "error retrieving Git Trigger"); err != nil {
			resp.Diagnostics.AddError("unable to load Git Trigger", err.Error())
		}
		return
	}

	data.ID = types.StringValue(gitTrigger.GetID())
	data.Name = types.StringValue(gitTrigger.Name)
	data.ProjectId = types.StringValue(gitTrigger.ProjectID)
	data.SpaceId = types.StringValue(gitTrigger.SpaceID)
	data.IsDisabled = types.BoolValue(gitTrigger.IsDisabled)
	data.Sources = convertGitTriggerSourcesToList(gitTrigger.Filter.(*filters.GitTriggerFilter).Sources)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *gitTriggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *schemas.GitTriggerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating Git Trigger '%s'", data.ID.ValueString()))

	client := r.Config.Client

	gitTrigger, err := client.ProjectTriggers.GetByID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load Git Trigger", err.Error())
		return
	}

	gitTriggerSources := convertListToGitTriggerSources(data.Sources)
	action := actions.NewCreateReleaseAction(data.ChannelId.ValueString())
	filter := filters.NewGitTriggerFilter(gitTriggerSources)
	project, err := projects.GetByID(client, data.SpaceId.ValueString(), data.ProjectId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("error finding project", err.Error())
		return
	}

	updatedGitTrigger := triggers.NewProjectTrigger(data.Name.ValueString(), data.Description.ValueString(), data.IsDisabled.ValueBool(), project, action, filter)
	updatedGitTrigger.ID = gitTrigger.ID

	updatedGitTrigger, err = client.ProjectTriggers.Update(updatedGitTrigger)
	tflog.Info(ctx, fmt.Sprintf("Git Trigger updated (%s)", data.ID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *gitTriggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.GitTriggerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.Client

	if err := client.ProjectTriggers.DeleteByID(data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete Git Trigger", err.Error())
		return
	}
}

func convertListToGitTriggerSources(list types.List) []filters.GitTriggerSource {
	var gitTriggerSources []filters.GitTriggerSource

	for _, elem := range list.Elements() {
		obj := elem.(types.Object)
		attrs := obj.Attributes()

		deploymentActionSlug := attrs["deployment_action_slug"].(types.String).ValueString()
		gitDependencyName := attrs["git_dependency_name"].(types.String).ValueString()
		includeFilePaths := convertToStringSlice(attrs["include_file_paths"].(types.List))
		excludeFilePaths := convertToStringSlice(attrs["exclude_file_paths"].(types.List))

		gitTriggerSource := filters.GitTriggerSource{
			DeploymentActionSlug: deploymentActionSlug,
			GitDependencyName:    gitDependencyName,
			IncludeFilePaths:     includeFilePaths,
			ExcludeFilePaths:     excludeFilePaths,
		}

		gitTriggerSources = append(gitTriggerSources, gitTriggerSource)
	}

	return gitTriggerSources
}

func convertToStringSlice(list types.List) []string {
	var result []string
	for _, elem := range list.Elements() {
		result = append(result, elem.(types.String).ValueString())
	}
	return result
}

func convertGitTriggerSourcesToList(gitTriggerSources []filters.GitTriggerSource) types.List {
	var elements []attr.Value

	for _, source := range gitTriggerSources {
		attributes := map[string]attr.Value{
			"deployment_action_slug": types.StringValue(source.DeploymentActionSlug),
			"git_dependency_name":    types.StringValue(source.GitDependencyName),
			"include_file_paths":     convertStringSliceToList(source.IncludeFilePaths),
			"exclude_file_paths":     convertStringSliceToList(source.ExcludeFilePaths),
		}
		objectValue, _ := types.ObjectValue(sourcesObjectType(), attributes)
		elements = append(elements, objectValue)
	}

	listValue, _ := types.ListValue(types.ObjectType{AttrTypes: sourcesObjectType()}, elements)
	return listValue
}

func convertStringSliceToList(strings []string) types.List {
	var elements []attr.Value

	for _, str := range strings {
		elements = append(elements, types.StringValue(str))
	}

	listValue, _ := types.ListValue(types.StringType, elements)
	return listValue
}

func sourcesObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"deployment_action_slug": types.StringType,
		"git_dependency_name":    types.StringType,
		"include_file_paths":     types.ListType{ElemType: types.StringType},
		"exclude_file_paths":     types.ListType{ElemType: types.StringType},
	}
}
