package schemas

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const RunbookResourceDescription = "runbook"

var RunbookSchemaAttributeNames = struct {
	ID                         string
	Name                       string
	Description                string
	ProjectID                  string
	RunbookProcessID           string
	PublishedRunbookSnapshotID string
	SpaceID                    string
	MultiTenancyMode           string
	ConnectivityPolicy         string
	EnvironmentScope           string
	Environments               string
	DefaultGuidedFailureMode   string
	RetentionPolicy            string
	ForcePackageDownload       string
}{
	ID:                         "id",
	Name:                       "name",
	Description:                "description",
	ProjectID:                  "project_id",
	RunbookProcessID:           "runbook_process_id",
	PublishedRunbookSnapshotID: "published_runbook_snapshot_id",
	SpaceID:                    "space_id",
	MultiTenancyMode:           "multi_tenancy_mode",
	ConnectivityPolicy:         "connectivity_policy",
	EnvironmentScope:           "environment_scope",
	Environments:               "environments",
	DefaultGuidedFailureMode:   "default_guided_failure_mode",
	RetentionPolicy:            "retention_policy",
	ForcePackageDownload:       "force_package_download",
}

var tenantedDeploymentModeNames = struct {
	Untenanted           string
	TenantedOrUntenanted string
	Tenanted             string
}{
	Untenanted:           "Untenanted",
	TenantedOrUntenanted: "TenantedOrUntenanted",
	Tenanted:             "Tenanted",
}

var tenantedDeploymentModes = []string{
	tenantedDeploymentModeNames.Untenanted,
	tenantedDeploymentModeNames.TenantedOrUntenanted,
	tenantedDeploymentModeNames.Tenanted,
}

var environmentScopeNames = struct {
	All                   string
	Specified             string
	FromProjectLifecycles string
}{
	All:                   "All",
	Specified:             "Specified",
	FromProjectLifecycles: "FromProjectLifecycles",
}

var environmentScopeTypes = []string{
	environmentScopeNames.All,
	environmentScopeNames.Specified,
	environmentScopeNames.FromProjectLifecycles,
}

var defaultGuidedFailureModeNames = struct {
	EnvironmentDefault string
	Off                string
	On                 string
}{
	EnvironmentDefault: "EnvironmentDefault",
	Off:                "Off",
	On:                 "On",
}

var defaultGuidedFailureModes = []string{
	defaultGuidedFailureModeNames.EnvironmentDefault,
	defaultGuidedFailureModeNames.Off,
	defaultGuidedFailureModeNames.On,
}

type RunbookTypeResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	ProjectID                  types.String `tfsdk:"project_id"`
	Description                types.String `tfsdk:"description"`
	RunbookProcessID           types.String `tfsdk:"runbook_process_id"`
	PublishedRunbookSnapshotID types.String `tfsdk:"published_runbook_snapshot_id"`
	SpaceID                    types.String `tfsdk:"space_id"`
	MultiTenancyMode           types.String `tfsdk:"multi_tenancy_mode"`
	ConnectivityPolicy         types.List   `tfsdk:"connectivity_policy"`
	EnvironmentScope           types.String `tfsdk:"environment_scope"`
	Environments               types.List   `tfsdk:"environments"`
	DefaultGuidedFailureMode   types.String `tfsdk:"default_guided_failure_mode"`
	RunRetentionPolicy         types.List   `tfsdk:"retention_policy"`
	ForcePackageDownload       types.Bool   `tfsdk:"force_package_download"`
}

type RunbookRetentionPeriodModel struct {
	QuantityToKeep    types.String `tfsdk:"quantity_to_keep"`
	ShouldKeepForever types.Bool   `tfsdk:"should_keep_forever"`
}

type RunbookConnectivityPolicyModel struct {
	AllowDeploymentsToNoTargets types.Bool   `tfsdk:"allow_deployments_to_no_targets"`
	ExcludeUnhealthyTargets     types.Bool   `tfsdk:"exclude_unhealthy_targets"`
	SkipMachineBehavior         types.String `tfsdk:"skip_machine_behaviour"`
	TargetRoles                 types.List   `tfsdk:"target_roles"`
}

func GetRunbookResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(RunbookResourceDescription),
		Attributes: map[string]resourceSchema.Attribute{
			RunbookSchemaAttributeNames.ID: resourceSchema.StringAttribute{
				Description: "The unique ID for this runbook.",
				Computed:    true,
			},
			RunbookSchemaAttributeNames.Name: resourceSchema.StringAttribute{
				Description: "The name of the runbook in Octopus Deploy. This name must be unique.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`\S+`),
						"expected value to not be an empty string or whitespace",
					),
				},
			},
			RunbookSchemaAttributeNames.Description: resourceSchema.StringAttribute{
				Description: "The description of this runbook.",
				Optional:    true,
				Computed:    true,
			},
			RunbookSchemaAttributeNames.ProjectID: resourceSchema.StringAttribute{
				Description: "The project that this runbook belongs to.",
				Required:    true,
			},
			RunbookSchemaAttributeNames.RunbookProcessID: resourceSchema.StringAttribute{
				Description: "The runbook process ID.",
				Computed:    true,
			},
			RunbookSchemaAttributeNames.PublishedRunbookSnapshotID: resourceSchema.StringAttribute{
				Description: "The published snapshot ID.",
				Computed:    true,
			},
			RunbookSchemaAttributeNames.SpaceID: util.GetSpaceIdResourceSchema(RunbookResourceDescription),
			RunbookSchemaAttributeNames.MultiTenancyMode: resourceSchema.StringAttribute{
				Description: fmt.Sprintf("The tenanted deployment mode of the runbook. Valid modes are %s", strings.Join(util.Map(tenantedDeploymentModes, func(item string) string { return fmt.Sprintf("`%s`", item) }), ", ")),
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(tenantedDeploymentModes...),
				},
			},
			RunbookSchemaAttributeNames.EnvironmentScope: resourceSchema.StringAttribute{
				Description: "Determines how the runbook is scoped to environments.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(environmentScopeTypes...),
				},
			},
			RunbookSchemaAttributeNames.Environments: resourceSchema.ListAttribute{
				Description: fmt.Sprintf("When %s is set to \"%s\", this is the list of environments the runbook can be run agains.", RunbookSchemaAttributeNames.EnvironmentScope, environmentScopeNames.Specified),
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			RunbookSchemaAttributeNames.DefaultGuidedFailureMode: resourceSchema.StringAttribute{
				Description: "Sets the runbook guided failure mode.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(defaultGuidedFailureModes...),
				},
			},
			RunbookSchemaAttributeNames.ForcePackageDownload: resourceSchema.BoolAttribute{
				Description: "Whether to force packages to be re-downloaded or not.",
				Computed:    true,
				Optional:    true,
			},
		},
		Blocks: map[string]resourceSchema.Block{
			RunbookSchemaAttributeNames.ConnectivityPolicy: resourceSchema.ListNestedBlock{
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: getConnectivityPolicySchema(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			RunbookSchemaAttributeNames.RetentionPolicy: resourceSchema.ListNestedBlock{
				Description: "Sets the runbook retention policy.",
				NestedObject: resourceSchema.NestedBlockObject{
					Attributes: getRunbookRetentionPeriodSchema(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
		},
	}
}
