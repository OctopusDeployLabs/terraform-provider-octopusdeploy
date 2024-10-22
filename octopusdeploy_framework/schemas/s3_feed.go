package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type S3FeedSchema struct{}

func (m S3FeedSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages a OCI Registry feed in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"use_machine_credentials": GetRequiredBooleanResourceAttribute("When true will use credentials configured on the worker"),
			"access_key":              GetOptionalStringResourceSchema("The AWS access key to use when authenticating against Amazon Web Services"),
			"secret_key":              GetSensitiveResourceSchema("The AWS secret key to use when authenticating against Amazon Web Services.", false),
			"id":                      GetIdResourceSchema(),
			"name":                    GetNameResourceSchema(true),
			"password":                GetPasswordResourceSchema(false),
			"space_id":                GetSpaceIdResourceSchema("AWS S3 Bucket Feed"),
			"username":                GetUsernameResourceSchema(false),
		},
	}
}

func (m S3FeedSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

var _ EntitySchema = S3FeedSchema{}

type S3FeedTypeResourceModel struct {
	UseMachineCredentials types.Bool   `tfsdk:"use_machine_credentials"`
	AccessKey             types.String `tfsdk:"access_key"`
	SecretKey             types.String `tfsdk:"secret_key"`
	Name                  types.String `tfsdk:"name"`
	Password              types.String `tfsdk:"password"`
	SpaceID               types.String `tfsdk:"space_id"`
	Username              types.String `tfsdk:"username"`

	ResourceModel
}
