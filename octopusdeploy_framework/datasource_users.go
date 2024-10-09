package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

type userDataSource struct {
	*Config
}

type usersDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	SpaceID types.String `tfsdk:"space_id"`
	IDs     types.List   `tfsdk:"ids"`
	Filter  types.String `tfsdk:"filter"`
	Skip    types.Int64  `tfsdk:"skip"`
	Take    types.Int64  `tfsdk:"take"`
	Users   types.List   `tfsdk:"users"`
}

func NewUsersDataSource() datasource.DataSource {
	return &userDataSource{}
}

func (u *userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("users")
}

func (u *userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.UserSchema{}.GetDatasourceSchema()
}

func (u *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	u.Config = DataSourceConfiguration(req, resp)
}

func (u *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data usersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := users.UsersQuery{
		IDs:    util.GetIds(data.IDs),
		Filter: data.Filter.ValueString(),
		Skip:   util.GetNumber(data.Skip),
		Take:   util.GetNumber(data.Take),
	}

	util.DatasourceReading(ctx, "users", query)

	existingUsers, err := users.Get(u.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load users", err.Error())
		return
	}

	mappedUsers := []schemas.UserTypeDatasourceModel{}
	tflog.Debug(ctx, fmt.Sprintf("users returned from API: %#v", existingUsers))
	for _, user := range existingUsers.Items {
		mappedUsers = append(mappedUsers, schemas.MapToUserDatasourceModel(user))
	}

	util.DatasourceResultCount(ctx, "users", len(mappedUsers))

	data.Users, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.UserObjectType()}, mappedUsers)
	data.ID = types.StringValue("Users " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
