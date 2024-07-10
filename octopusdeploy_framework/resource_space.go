package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type spaceResource struct {
	*Config
}

func NewSpaceResource() resource.Resource {
	return &spaceResource{}
}

func (b *spaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides information about an existing space.",
		Attributes:  schemas.GetSpaceResourceSchema(),
	}
}

func (b *spaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("space")
}

func (b *spaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.SpaceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating space: %#v", data))

	newSpace := spaces.NewSpace(data.Name.ValueString())
	newSpace.Slug = data.Slug.ValueString()
	newSpace.Description = data.Description.ValueString()
	newSpace.IsDefault = data.IsDefault.ValueBool()
	newSpace.TaskQueueStopped = data.IsTaskQueueStopped.ValueBool()
	data.SpaceManagersTeams.ElementsAs(ctx, newSpace.SpaceManagersTeams, false)
	data.SpaceManagersTeamMembers.ElementsAs(ctx, newSpace.SpaceManagersTeamMembers, false)

	createdSpace, err := b.Client.Spaces.Add(newSpace)
	if err != nil {
		resp.Diagnostics.AddError("unable to create new space", err.Error())
	}

}

func (b *spaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (b *spaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (b *spaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}
