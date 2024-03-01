package octopusdeploy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenDeploymentAction(action *deployments.DeploymentAction) map[string]interface{} {
	flattenedDeploymentAction := flattenAction(action)

	if len(action.ActionType) > 0 {
		flattenedDeploymentAction["action_type"] = action.ActionType
	}

	if len(action.WorkerPool) > 0 {
		flattenedDeploymentAction["worker_pool_id"] = action.WorkerPool
	}

	if len(action.WorkerPoolVariable) > 0 {
		flattenedDeploymentAction["worker_pool_variable"] = action.WorkerPoolVariable
	}

	if v, ok := action.Properties["Octopus.Action.RunOnServer"]; ok {
		runOnServer, _ := strconv.ParseBool(v.Value)
		flattenedDeploymentAction["run_on_server"] = runOnServer
	}

	return flattenedDeploymentAction
}

func flattenAction(action *deployments.DeploymentAction) map[string]interface{} {
	if action == nil {
		return nil
	}

	flattenedAction := map[string]interface{}{
		"can_be_used_for_project_versioning": action.CanBeUsedForProjectVersioning,
		"is_disabled":                        action.IsDisabled,
		"is_required":                        action.IsRequired,
	}

	if len(action.Channels) > 0 {
		flattenedAction["channels"] = action.Channels
	}

	if len(action.Condition) > 0 {
		flattenedAction["condition"] = action.Condition
	}

	if action.Container != nil {
		flattenedAction["container"] = flattenContainer(action.Container)
	}

	if len(action.Environments) > 0 {
		flattenedAction["environments"] = action.Environments
	}

	if len(action.ExcludedEnvironments) > 0 {
		flattenedAction["excluded_environments"] = action.ExcludedEnvironments
	}

	if len(action.ID) > 0 {
		flattenedAction["id"] = action.ID
	}

	if len(action.Name) > 0 {
		flattenedAction["name"] = action.Name
	}

	if len(action.Notes) > 0 {
		flattenedAction["notes"] = action.Notes
	}

	if len(action.Properties) > 0 {
		flattenedAction["properties"] = flattenProperties(action.Properties)
	}

	if len(action.TenantTags) > 0 {
		flattenedAction["tenant_tags"] = action.TenantTags
	}

	if v, ok := action.Properties["Octopus.Action.EnabledFeatures"]; ok {
		flattenedAction["features"] = strings.Split(v.Value, ",")
	}

	if v, ok := action.Properties["Octopus.Action.Template.Id"]; ok {
		actionTemplate := map[string]interface{}{
			"id": v.Value,
		}

		if v, ok := action.Properties["Octopus.Action.Template.Version"]; ok {
			version, _ := strconv.Atoi(v.Value)
			actionTemplate["version"] = version
		}

		flattenedAction["action_template"] = []interface{}{actionTemplate}
	}

	if len(action.Packages) > 0 {
		flattenedPackageReferences := []interface{}{}
		for _, packageReference := range action.Packages {
			flattenedPackageReference := flattenPackageReference(packageReference)
			if len(packageReference.Name) == 0 {
				flattenedAction["primary_package"] = []interface{}{flattenedPackageReference}
				// TODO: consider these properties
				// actionProperties["Octopus.Action.Package.DownloadOnTentacle"] = packageReference.AcquisitionLocation
				// flattenedAction["properties"] = actionProperties
				continue
			}

			if v, ok := packageReference.Properties["Extract"]; ok {
				extractDuringDeployment, _ := strconv.ParseBool(v)
				flattenedPackageReference["extract_during_deployment"] = extractDuringDeployment
			}

			flattenedPackageReferences = append(flattenedPackageReferences, flattenedPackageReference)
		}
		flattenedAction["package"] = flattenedPackageReferences
	}

	return flattenedAction
}

func getDeploymentActionSchema() *schema.Schema {
	actionSchema, element := getActionSchema()
	addActionTypeSchema(element)
	addExecutionLocationSchema(element)
	addWorkerPoolSchema(element)
	addWorkerPoolVariableSchema(element)
	addPackagesSchema(element, false)

	return actionSchema
}

func getActionSchema() (*schema.Schema, *schema.Resource) {
	element := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action_template": {
				Computed:    true,
				Description: "Represents the template that is associated with this action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"community_action_template_id": {
							Computed: true,
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Computed: true,
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeSet,
			},
			"can_be_used_for_project_versioning": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeBool,
			},
			"channels": {
				Computed:    true,
				Description: "The channels associated with this deployment action.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"condition": {
				Computed:    true,
				Description: "The condition associated with this deployment action.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"container": {
				Computed:    true,
				Description: "The deployment action container associated with this deployment action.",
				Elem:        &schema.Resource{Schema: getDeploymentActionContainerSchema()},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"environments": {
				Computed:    true,
				Description: "The environments within which this deployment action will run.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"excluded_environments": {
				Computed:    true,
				Description: "The environments that this step will be skipped in",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"features": {
				Computed:    true,
				Description: "A list of enabled features for this action.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Type:        schema.TypeList,
			},
			"id": getIDSchema(),
			"is_disabled": {
				Default:     false,
				Description: "Indicates the disabled status of this deployment action.",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"is_required": {
				Default:     false,
				Description: "Indicates the required status of this deployment action.",
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"name": getNameSchema(true),
			"notes": {
				Description: "The notes associated with this deployment action.",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"package": getPackageSchema(false),
			"properties": {
				Computed:         true,
				Description:      "The properties associated with this deployment action.",
				Elem:             &schema.Schema{Type: schema.TypeString},
				Optional:         true,
				Type:             schema.TypeMap,
				ValidateDiagFunc: warnIfIncludesRunOnServer(),
			},
			"sort_order": {
				Description: "Order used by terraform to ensure correct ordering of actions. This property must be either omitted from all actions, or provided on all actions",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
			},
			"tenant_tags": getTenantTagsSchema(),
		},
	}

	actionSchema := &schema.Schema{
		Computed: false,
		Elem:     element,
		Optional: true,
		Type:     schema.TypeList,
	}

	return actionSchema, element
}

func warnIfIncludesRunOnServer() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		keys := v.(map[string]interface{})
		const key = "Octopus.Action.RunOnServer"

		if _, ok := keys[key]; ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       fmt.Sprintf("\"%s\" is defined in properties", key),
				Detail:        "Please update your template to specify \"run_on_server\" under the action instead.",
				AttributePath: append(path, cty.IndexStep{Key: cty.StringVal(key)}),
			})
		}

		return diags
	}
}

func addExecutionLocationSchema(element *schema.Resource) {
	element.Schema["run_on_server"] = &schema.Schema{
		Default:     false,
		Description: "Whether this step runs on a worker or on the target",
		Optional:    true,
		Type:        schema.TypeBool,
	}
}

func addActionTypeSchema(element *schema.Resource) {
	element.Schema["action_type"] = &schema.Schema{
		Description:      "The type of action",
		Required:         true,
		Type:             schema.TypeString,
		ValidateDiagFunc: warnIfActionTypeExists(),
	}
}

func warnIfActionTypeExists() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		value := v.(string)

		var actionType string
		switch value {
		case "Octopus.KubernetesDeploySecret":
			actionType = "deploy_kubernetes_secret_action"
		case "Octopus.KubernetesRunScript":
			actionType = "run_kubectl_script_action"
		case "Octopus.Manual":
			actionType = "manual_intervention_action"
		case "Octopus.Script":
			actionType = "run_script_action"
		case "Octopus.TentaclePackage":
			actionType = "deploy_package_action"
		case "Octopus.TerraformApply":
			actionType = "apply_terraform_template_action"
		case "Octopus.WindowsService":
			actionType = "deploy_windows_service_action"
		}

		if len(actionType) > 0 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       "This action type has its own explicit action type",
				Detail:        fmt.Sprintf("Please use the new \"%s\" instead of the generic \"action\".", actionType),
				AttributePath: path,
			})
		}

		return diags
	}
}

func addPropertiesSchema(element *schema.Resource, deprecated string) {
	element.Schema["properties"] = &schema.Schema{
		Computed:    true,
		Description: "The properties associated with this deployment action.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeMap,
	}

	if len(deprecated) > 0 {
		element.Schema["properties"].Deprecated = deprecated
	}
}

func addWorkerPoolSchema(element *schema.Resource) {
	element.Schema["worker_pool_id"] = &schema.Schema{
		Description: "The worker pool associated with this deployment action.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func addWorkerPoolVariableSchema(element *schema.Resource) {
	element.Schema["worker_pool_variable"] = &schema.Schema{
		Description: "The worker pool variable associated with this deployment action.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func expandAction(flattenedAction map[string]interface{}) *deployments.DeploymentAction {
	if len(flattenedAction) == 0 {
		return nil
	}

	if _, ok := flattenedAction["name"].(string); !ok {
		return nil
	}
	name := flattenedAction["name"].(string)

	var actionType string
	if v, ok := flattenedAction["action_type"].(string); ok {
		actionType = v
	}

	action := deployments.NewDeploymentAction(name, actionType)

	// Retain the existing id as older channel rules might still be referencing the action by id
	if v, ok := flattenedAction["id"]; ok {
		action.ID = v.(string)
	}

	// expand properties first
	if v, ok := flattenedAction["properties"]; ok {
		action.Properties = expandProperties(v)
	}

	if v, ok := flattenedAction["can_be_used_for_project_versioning"]; ok {
		action.CanBeUsedForProjectVersioning = v.(bool)
	}

	if v, ok := flattenedAction["channels"]; ok {
		action.Channels = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedAction["condition"]; ok {
		action.Condition = v.(string)
	}

	if v, ok := flattenedAction["container"]; ok {
		action.Container = expandContainer(v)
	}

	if v, ok := flattenedAction["environments"]; ok {
		action.Environments = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedAction["excluded_environments"]; ok {
		action.ExcludedEnvironments = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedAction["features"]; ok {
		action.Properties["Octopus.Action.EnabledFeatures"] = core.NewPropertyValue(strings.Join(getSliceFromTerraformTypeList(v), ","), false)
	}

	if v, ok := flattenedAction["is_disabled"]; ok {
		action.IsDisabled = v.(bool)
	}

	if v, ok := flattenedAction["is_required"]; ok {
		action.IsRequired = v.(bool)
	}

	if v, ok := flattenedAction["notes"]; ok {
		action.Notes = v.(string)
	}

	if v, ok := flattenedAction["run_on_server"]; ok {
		runOnServer := v.(bool)
		action.Properties["Octopus.Action.RunOnServer"] = core.NewPropertyValue(cases.Title(language.Und, cases.NoLower).String(strconv.FormatBool(runOnServer)), false)
	}

	if v, ok := flattenedAction["action_template"]; ok {
		templateList := v.(*schema.Set).List()
		if len(templateList) > 0 {
			template := templateList[0].(map[string]interface{})
			action.Properties["Octopus.Action.Template.Id"] = core.NewPropertyValue(template["id"].(string), false)
			version := strconv.Itoa(template["version"].(int))
			action.Properties["Octopus.Action.Template.Version"] = core.NewPropertyValue(version, false)
		}
	}

	if v, ok := flattenedAction["tenant_tags"]; ok {
		action.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedAction["worker_pool_id"]; ok {
		action.WorkerPool = v.(string)
	}

	if v, ok := flattenedAction["worker_pool_variable"]; ok {
		action.WorkerPoolVariable = v.(string)
	}

	if v, ok := flattenedAction["primary_package"]; ok {
		primaryPackages := v.([]interface{})
		for _, primaryPackage := range primaryPackages {
			primaryPackageReference := expandPackageReference(primaryPackage.(map[string]interface{}))

			switch primaryPackageReference.AcquisitionLocation {
			case "Server":
				action.Properties["Octopus.Action.Package.DownloadOnTentacle"] = core.NewPropertyValue("False", false)
			default:
				action.Properties["Octopus.Action.Package.DownloadOnTentacle"] = core.NewPropertyValue(primaryPackageReference.AcquisitionLocation, false)
			}

			if len(primaryPackageReference.PackageID) > 0 {
				action.Properties["Octopus.Action.Package.PackageId"] = core.NewPropertyValue(primaryPackageReference.PackageID, false)
			}

			if len(primaryPackageReference.FeedID) > 0 {
				action.Properties["Octopus.Action.Package.FeedId"] = core.NewPropertyValue(primaryPackageReference.FeedID, false)
			}

			action.Packages = append(action.Packages, primaryPackageReference)
		}
	}

	if v, ok := flattenedAction["package"]; ok {
		packageReferences := v.([]interface{})
		for _, packageReference := range packageReferences {
			action.Packages = append(action.Packages, expandPackageReference(packageReference.(map[string]interface{})))
		}
	}

	return action
}
