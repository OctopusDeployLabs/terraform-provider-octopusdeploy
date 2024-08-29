package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
