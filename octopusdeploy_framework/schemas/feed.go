package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FlattenFeed(feed *feeds.FeedResource) attr.Value {
	return types.ObjectValueMust(FeedObjectType(), map[string]attr.Value{
		"access_key":                            types.StringValue(feed.AccessKey),
		"api_version":                           types.StringValue(feed.APIVersion),
		"delete_unreleased_packages_after_days": types.Int64Value(int64(feed.DeleteUnreleasedPackagesAfterDays)),
		"download_attempts":                     types.Int64Value(int64(feed.DownloadAttempts)),
		"download_retry_backoff_seconds":        types.Int64Value(int64(feed.DownloadRetryBackoffSeconds)),
		"feed_type":                             types.StringValue(string(feed.FeedType)),
		"feed_uri":                              types.StringValue(feed.FeedURI),
		"id":                                    types.StringValue(feed.GetID()),
		"is_enhanced_mode":                      types.BoolValue(feed.EnhancedMode),
		"name":                                  types.StringValue(feed.Name),
		"package_acquisition_location_options":  types.ListValueMust(types.StringType, util.ToValueSlice(feed.PackageAcquisitionLocationOptions)),
		"region":                                types.StringValue(feed.Region),
		"registry_path":                         types.StringValue(feed.RegistryPath),
		"space_id":                              types.StringValue(feed.SpaceID),
		"username":                              types.StringValue(feed.Username),
		// Password and secret key are sensitive values that are not returned from the API.
		// Here we map empty values to keep the behaviour consistent with the SDK.
		"password":   types.StringValue(""),
		"secret_key": types.StringValue(""),
	})
}

func FeedObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"access_key":                            types.StringType,
		"api_version":                           types.StringType,
		"delete_unreleased_packages_after_days": types.Int64Type,
		"download_attempts":                     types.Int64Type,
		"download_retry_backoff_seconds":        types.Int64Type,
		"feed_type":                             types.StringType,
		"feed_uri":                              types.StringType,
		"id":                                    types.StringType,
		"is_enhanced_mode":                      types.BoolType,
		"name":                                  types.StringType,
		"package_acquisition_location_options":  types.ListType{ElemType: types.StringType},
		"region":                                types.StringType,
		"registry_path":                         types.StringType,
		"space_id":                              types.StringType,
		"username":                              types.StringType,
		"password":                              types.StringType,
		"secret_key":                            types.StringType,
	}
}

func GetFeedsDataSourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
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
	}
}

func GetFeedDataSourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
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
