

Templates:


## Datasource


```golang
import (
	"context"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// datasource model
type blahsDataSource struct {
	*Config
}

// query model
type blahsModel struct {
	ID          types.String `tfsdk:"id"`
	SpaceID     types.String `tfsdk:"space_id"`
	IDs         types.List   `tfsdk:"ids"`
	PartialName types.String `tfsdk:"partial_name"`
	Skip        types.Int64  `tfsdk:"skip"`
	Take        types.Int64  `tfsdk:"take"`
	Blahs       types.List   `tfsdk:"blahs"`
}

// new datasource
func NewBlahsDataSource() datasource.DataSource {
	return &blahsDataSource{}
}

func (*blahsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = ProviderTypeName + "_blahs"
}

func (*blahsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// request
			"space_id":     util.GetSpaceIdDatasourceSchema("blahs"),
			"ids":          util.GetQueryIDsDatasourceSchema(),
			"partial_name": util.GetQueryPartialNameDatasourceSchema(),
			"skip":         util.GetQuerySkipDatasourceSchema(),
			"take":         util.GetQueryTakeDatasourceSchema(),

			// response
			"id": util.GetIdDatasourceSchema(),
		},
		Blocks: map[string]schema.Block{
			"blahs": schema.ListNestedBlock{
				Description: "blahs description",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":       util.GetIdResourceSchema(),
						"space_id": util.GetSpaceIdResourceSchema("blahs"),
						"name":     util.GetNameResourceSchema(true),
						...
					},
				},
			},
		},
	}
}

func (b *blahsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	b.Config = DataSourceConfiguration(req, resp)
}

func (b *blahsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data blahsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	// construct query
	
	// map response model back to state
	/*
	
	for _, b  :-= range resp.Items {
		map from api model to internal model
	}
	
	
	 */
	
	// set state
	data.Blahs = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: getNestedGroupAttrs()}, )
	data.ID = // something, usually a string concatenation specific to this resource 
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

Resource:


```golang
package octopusdeploy_framework

import (
	"context"
	
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type blahResource struct {
	*Config
}

type resourceModel struct {
	ID      string `tfsdk:"id"`
	Name    string `tfsdk:"name"`
	SpaceID string `tfsdk:"space_id"`
}

// check that basic Resource interface has been implemented - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/resource#Resource
var _ resource.Resource = &resourceModel{}

// check that the ResourceWithImportState has been implented - https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/resource#ResourceWithImportState
var _ resource.ResourceWithImportState = &resourceModel{}

func NewBlahResource() resource.Resource {
	return &blahResource{}
}

func (b *blahResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (b *blahResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	b.Config = ResourceConfiguration(req, resp)
}

func (b *blahResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("resource name")
}

func (b *blahResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	....
}

func (b *blahResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	....
}

func (b *blahResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data, state model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	....
	
	// get the resource from the api
	result = client.get(state.ID.ValueString())
	
	// update the resource
	result.Name = plan.Name.ValueString()
	...
	
	// update api
	_, err = resources.Update(b.Client, result)

    // map result back to state

	
}

func (b *blahResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	// delete by id 
	if err := b.Client.Resource.DeleteById(data.ID.ValueString(); err != nil {
		resp.Diagnostics.AddError("unable to delete resource", err.Error())
		return
	}
		
	tflog.Info("resource deleted")
}

// only reqired if the resource is supporting Import
func (s *blahResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
```

# Gotchas

## SpaceID

Most resources can specify the SpaceID on the provider configuration, or on the resource it self. When doing API calls, you need to pass the space ID to the client function from the resource if specified or from the config/client if not specified

## Diags

A lot of built-in functions return a tuple with a result and a `diag.Diagnotic`. If any errors occured during the function call, you can check and ideally return either the diagnostics or just return from the current function.

Example:

Within one of the resource interface functions

```golang
resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
if resp.Diagnostics.HasError() {
    return
}
```

Returning a specific error

```golang
if err := b.Client.Resource.DeleteById(data.ID.ValueString(); err != nil {
    resp.Diagnostics.AddError("unable to delete resource", err.Error())
    return
}
```

