package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FeedsSchema struct{}

var _ EntitySchema = FeedsSchema{}

func (f FeedsSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{}
}

func (f FeedsSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing feeds.",
		Attributes: map[string]datasourceSchema.Attribute{
			"feed_type": datasourceSchema.StringAttribute{
				Description: "A filter to search by feed type. Valid feed types are `AwsElasticContainerRegistry`, `BuiltIn`, `Docker`, `GitHub`, `Helm`, `Maven`, `NuGet`, or `OctopusProject`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"AwsElasticContainerRegistry",
						"BuiltIn",
						"Docker",
						"GitHub",
						"Helm",
						"Maven",
						"NuGet",
						"OctopusProject"),
				},
			},
			"ids":          GetQueryIDsDatasourceSchema(),
			"name":         GetNameDatasourceSchema(false),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take":         GetQueryTakeDatasourceSchema(),
			"space_id":     GetSpaceIdDatasourceSchema("feeds", false),

			// response
			"id": GetIdDatasourceSchema(true),
			"feeds": datasourceSchema.ListNestedAttribute{
				Computed: true,
				Optional: false,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"feed_type": datasourceSchema.StringAttribute{
							Description: "A filter to search by feed type. Valid feed types are `AwsElasticContainerRegistry`, `BuiltIn`, `Docker`, `GitHub`, `Helm`, `Maven`, `NuGet`, or `OctopusProject`.",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"AwsElasticContainerRegistry",
									"BuiltIn",
									"Docker",
									"GitHub",
									"Helm",
									"Maven",
									"NuGet",
									"OctopusProject"),
							},
						},
						"feed_uri": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"id": GetIdDatasourceSchema(true),
						"is_enhanced_mode": datasourceSchema.BoolAttribute{
							Computed: true,
						},
						"name": GetReadonlyNameDatasourceSchema(),
						"password": datasourceSchema.StringAttribute{
							Description: "The password associated with this resource.",
							Sensitive:   true,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"package_acquisition_location_options": datasourceSchema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"region": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"registry_path": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"secret_key": datasourceSchema.StringAttribute{
							Computed:  true,
							Sensitive: true,
						},
						"space_id": GetSpaceIdDatasourceSchema("feeds", true),
						"username": datasourceSchema.StringAttribute{
							Description: "The username associated with this resource.",
							Sensitive:   true,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"delete_unreleased_packages_after_days": datasourceSchema.Int64Attribute{
							Computed: true,
						},
						"access_key": datasourceSchema.StringAttribute{
							Computed:    true,
							Description: "The AWS access key to use when authenticating against Amazon Web Services.",
						},
						"api_version": datasourceSchema.StringAttribute{
							Computed: true,
						},
						"download_attempts": datasourceSchema.Int64Attribute{
							Description: "The number of times a deployment should attempt to download a package from this feed before failing.",
							Computed:    true,
						},
						"download_retry_backoff_seconds": datasourceSchema.Int64Attribute{
							Description: "The number of seconds to apply as a linear back off between download attempts.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

type FeedsDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Feeds       types.List   `tfsdk:"feeds"`
	FeedType    types.String `tfsdk:"feed_type"`
	IDs         types.List   `tfsdk:"ids"`
	Name        types.String `tfsdk:"name"`
	PartialName types.String `tfsdk:"partial_name"`
	Skip        types.Int64  `tfsdk:"skip"`
	Take        types.Int64  `tfsdk:"take"`
	SpaceID     types.String `tfsdk:"space_id"`
}
