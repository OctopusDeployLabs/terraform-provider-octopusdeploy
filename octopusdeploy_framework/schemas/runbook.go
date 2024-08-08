package schemas

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
			RunbookSchemaAttributeNames.ID: util.GetIdResourceSchema(),
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
			RunbookSchemaAttributeNames.Description: util.GetDescriptionResourceSchema(RunbookResourceDescription),
			RunbookSchemaAttributeNames.ProjectID: resourceSchema.StringAttribute{
				Description: "The project that this runbook belongs to.",
				Required:    true,
			},
			RunbookSchemaAttributeNames.RunbookProcessID: resourceSchema.StringAttribute{
				Description: "The runbook process ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.PublishedRunbookSnapshotID: resourceSchema.StringAttribute{
				Description: "The published snapshot ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.SpaceID: util.GetSpaceIdResourceSchema(RunbookResourceDescription),
			RunbookSchemaAttributeNames.MultiTenancyMode: resourceSchema.StringAttribute{
				Description: fmt.Sprintf("The tenanted deployment mode of the runbook. Valid modes are %s", strings.Join(util.Map(tenantedDeploymentModes, func(item string) string { return fmt.Sprintf("`%s`", item) }), ", ")),
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(tenantedDeploymentModes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.EnvironmentScope: resourceSchema.StringAttribute{
				Description: "Determines how the runbook is scoped to environments.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(environmentScopeTypes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.Environments: resourceSchema.ListAttribute{
				Description: fmt.Sprintf("When %s is set to \"%s\", this is the list of environments the runbook can be run against.", RunbookSchemaAttributeNames.EnvironmentScope, environmentScopeNames.Specified),
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.DefaultGuidedFailureMode: resourceSchema.StringAttribute{
				Description: "Sets the runbook guided failure mode.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(defaultGuidedFailureModes...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			RunbookSchemaAttributeNames.ForcePackageDownload: resourceSchema.BoolAttribute{
				Description: "Whether to force packages to be re-downloaded or not.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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

func (data *RunbookTypeResourceModel) RefreshFromApiResponse(ctx context.Context, runbook *runbooks.Runbook) diag.Diagnostics {
	var diags diag.Diagnostics

	if runbook == nil {
		return diags
	}

	data.ID = types.StringValue(runbook.ID)
	data.Name = types.StringValue(runbook.Name)
	data.ProjectID = types.StringValue(runbook.ProjectID)
	data.Description = types.StringValue(runbook.Description)
	data.RunbookProcessID = types.StringValue(runbook.RunbookProcessID)
	data.PublishedRunbookSnapshotID = types.StringValue(runbook.PublishedRunbookSnapshotID)
	data.SpaceID = types.StringValue(runbook.SpaceID)
	data.MultiTenancyMode = types.StringValue(string(runbook.MultiTenancyMode))
	data.EnvironmentScope = types.StringValue(runbook.EnvironmentScope)
	data.Environments = util.FlattenStringList(runbook.Environments)
	data.DefaultGuidedFailureMode = types.StringValue(runbook.DefaultGuidedFailureMode)
	data.ForcePackageDownload = types.BoolValue(runbook.ForcePackageDownload)
	if !data.ConnectivityPolicy.IsNull() {
		result, d := types.ListValueFrom(
			ctx,
			types.ObjectType{AttrTypes: GetConnectivityPolicyObjectType()},
			[]attr.Value{MapFromConnectivityPolicy(runbook.ConnectivityPolicy)},
		)
		diags.Append(d...)
		data.ConnectivityPolicy = result
	} /*else {
		data.ConnectivityPolicy = types.ListValueMust(
			types.ObjectType{AttrTypes: GetConnectivityPolicyObjectType()},
			[]attr.Value{MapFromConnectivityPolicy(GetDefaultConnectivityPolicy())},
		)
	}*/
	if !data.RunRetentionPolicy.IsNull() {
		result, d := types.ListValueFrom(
			ctx,
			types.ObjectType{AttrTypes: GetRunbookRetentionPeriodObjectType()},
			[]attr.Value{MapFromRunbookRetentionPeriod(runbook.RunRetentionPolicy)},
		)
		diags.Append(d...)
		data.RunRetentionPolicy = result
	}

	return diags
}
