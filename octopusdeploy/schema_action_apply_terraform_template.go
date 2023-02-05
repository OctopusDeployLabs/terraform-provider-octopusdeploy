package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func addTerraformTemplateAdvancedOptionsSchema(element *schema.Resource) {
	element.Schema["advanced_options"] = &schema.Schema{
		Description: "Optional advanced options for Terraform",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allow_additional_plugin_downloads": {
					Default:  true,
					Optional: true,
					Type:     schema.TypeBool,
				},
				"apply_parameters": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"init_parameters": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"plugin_cache_directory": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"workspace": {
					Optional: true,
					Type:     schema.TypeString,
				},
			},
		},
		MaxItems: 1,
		Required: true,
		Type:     schema.TypeSet,
	}
}

func addTerraformTemplateAwsAccountSchema(element *schema.Resource) {
	element.Schema["aws_account"] = &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"region": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"role": {
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"arn": {
								Optional: true,
								Type:     schema.TypeString,
							},
							"external_id": {
								Optional: true,
								Type:     schema.TypeString,
							},
							"role_session_name": {
								Optional: true,
								Type:     schema.TypeString,
							},
							"session_duration": {
								Optional: true,
								Type:     schema.TypeInt,
							},
						},
					},
					MaxItems: 1,
					Optional: true,
					Type:     schema.TypeSet,
				},
				"variable": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"use_instance_role": {
					Optional: true,
					Type:     schema.TypeBool,
				},
			},
		},
		MaxItems: 1,
		Optional: true,
		Type:     schema.TypeSet,
	}
}

func addTerraformTemplateAzureAccountSchema(element *schema.Resource) {
	element.Schema["azure_account"] = &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"variable": {
					Optional: true,
					Type:     schema.TypeString,
				},
			},
		},
		MaxItems: 1,
		Optional: true,
		Type:     schema.TypeSet,
	}
}

func addTerraformTemplateGoogleAccountSchema(element *schema.Resource) {
	element.Schema["google_cloud_account"] = &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"variable": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"use_vm_service_account": {
					Optional:    true,
					Description: "When running in a Compute Engine virtual machine, use the associated VM service account",
					Type:        schema.TypeBool,
				},
				"project": {
					Optional:    true,
					Description: "This sets GOOGLE_PROJECT environment variable",
					Type:        schema.TypeString,
				},
				"region": {
					Optional:    true,
					Description: "This sets GOOGLE_REGION environment variable",
					Type:        schema.TypeString,
				},
				"zone": {
					Optional:    true,
					Description: "This sets GOOGLE_ZONE environment variable",
					Type:        schema.TypeString,
				},
				"service_account_emails": {
					Optional:    true,
					Description: "This sets GOOGLE_IMPERSONATE_SERVICE_ACCOUNT environment variable",
					Type:        schema.TypeString,
				},
				"impersonate_service_account": {
					Optional:    true,
					Description: "Impersonate service accounts",
					Type:        schema.TypeBool,
				},
			},
		},
		MaxItems: 1,
		Optional: true,
		Type:     schema.TypeSet,
	}
}

func addTerraformTemplateParametersSchema(element *schema.Resource) {
	element.Schema["template_parameters"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
}

func addTerraformInlineTemplateSchema(element *schema.Resource) {
	element.Schema["inline_template"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
}

func addTerraformTemplateSchema(element *schema.Resource) {
	element.Schema["template"] = &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"additional_variable_files": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"directory": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"run_automatic_file_substitution": {
					Optional: true,
					Type:     schema.TypeBool,
				},
				"target_files": {
					Optional: true,
					Type:     schema.TypeString,
				},
			},
		},
		MaxItems: 1,
		Optional: true,
		Type:     schema.TypeSet,
	}
}

func flattenTerraformTemplate(properties map[string]core.PropertyValue) []interface{} {
	if len(properties) == 0 {
		return nil
	}

	flattenedMap := map[string]interface{}{}

	for k, v := range properties {
		switch k {
		case "Octopus.Action.Terraform.FileSubstitution":
			flattenedMap["target_files"] = v.Value
		case "Octopus.Action.Terraform.RunAutomaticFileSubstitution":
			runAutomaticFileSubstitution, _ := strconv.ParseBool(v.Value)
			flattenedMap["run_automatic_file_substitution"] = runAutomaticFileSubstitution
		case "Octopus.Action.Terraform.TemplateDirectory":
			flattenedMap["directory"] = v.Value
		case "Octopus.Action.Terraform.VarFiles":
			flattenedMap["additional_variable_files"] = v.Value
		}
	}

	return []interface{}{flattenedMap}
}

func expandApplyTerraformTemplateAction(flattenedAction map[string]interface{}) *deployments.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.TerraformApply"

	if v, ok := flattenedAction["template"]; ok {
		template := v.(*schema.Set).List()[0].(map[string]interface{})

		if v, ok := template["additional_variable_files"]; ok {
			action.Properties["Octopus.Action.Terraform.VarFiles"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := template["directory"]; ok {
			action.Properties["Octopus.Action.Terraform.TemplateDirectory"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := template["run_automatic_file_substitution"]; ok {
			runAutomaticFileSubstitution := v.(bool)
			action.Properties["Octopus.Action.Terraform.RunAutomaticFileSubstitution"] = core.NewPropertyValue(cases.Title(language.Und, cases.NoLower).String(strconv.FormatBool(runAutomaticFileSubstitution)), false)
		}

		if v, ok := template["target_files"]; ok {
			action.Properties["Octopus.Action.Terraform.FileSubstitution"] = core.NewPropertyValue(v.(string), false)
		}
	}

	if v, ok := flattenedAction["template_parameters"]; ok {
		action.Properties["Octopus.Action.Terraform.TemplateParameters"] = core.NewPropertyValue(v.(string), false)
	}

	if v, ok := flattenedAction["inline_template"]; ok {
		action.Properties["Octopus.Action.Terraform.Template"] = core.NewPropertyValue(v.(string), false)
	}

	if _, ok := flattenedAction["primary_package"]; ok {
		action.Properties["Octopus.Action.Script.ScriptSource"] = core.NewPropertyValue("Package", false)
	} else {
		action.Properties["Octopus.Action.Script.ScriptSource"] = core.NewPropertyValue("Inline", false)
	}

	if v, ok := flattenedAction["advanced_options"]; ok && len(v.(*schema.Set).List()) > 0 {
		advancedOptions := v.(*schema.Set).List()[0].(map[string]interface{})

		if v, ok := advancedOptions["allow_additional_plugin_downloads"]; ok {
			allowPluginDownloads := v.(bool)
			action.Properties["Octopus.Action.Terraform.AllowPluginDownloads"] = core.NewPropertyValue(cases.Title(language.Und, cases.NoLower).String(strconv.FormatBool(allowPluginDownloads)), false)
		}

		if v, ok := advancedOptions["apply_parameters"]; ok {
			action.Properties["Octopus.Action.Terraform.AdditionalActionParams"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := advancedOptions["init_parameters"]; ok {
			action.Properties["Octopus.Action.Terraform.AdditionalInitParams"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := advancedOptions["plugin_cache_directory"]; ok {
			action.Properties["Octopus.Action.Terraform.PluginsDirectory"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := advancedOptions["workspace"]; ok {
			action.Properties["Octopus.Action.Terraform.Workspace"] = core.NewPropertyValue(v.(string), false)
		}
	}

	if v, ok := flattenedAction["aws_account"]; ok && len(v.(*schema.Set).List()) > 0 {
		action.Properties["Octopus.Action.Terraform.ManagedAccount"] = core.NewPropertyValue("AWS", false)

		awsAccount := v.(*schema.Set).List()[0].(map[string]interface{})

		if v, ok := awsAccount["region"]; ok {
			action.Properties["Octopus.Action.Aws.Region"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := awsAccount["role"]; ok && len(v.(*schema.Set).List()) > 0 {
			action.Properties["Octopus.Action.Aws.AssumeRole"] = core.NewPropertyValue("True", false)

			role := v.(*schema.Set).List()[0].(map[string]interface{})

			if v, ok := role["arn"]; ok {
				action.Properties["Octopus.Action.Aws.AssumedRoleArn"] = core.NewPropertyValue(v.(string), false)
			}

			if v, ok := role["external_id"]; ok {
				action.Properties["Octopus.Action.Aws.AssumeRoleExternalId"] = core.NewPropertyValue(v.(string), false)
			}

			if v, ok := role["role_session_name"]; ok {
				action.Properties["Octopus.Action.Aws.AssumedRoleSession"] = core.NewPropertyValue(v.(string), false)
			}

			if v, ok := role["session_duration"]; ok {
				action.Properties["Octopus.Action.Aws.AssumeRoleSessionDurationSeconds"] = core.NewPropertyValue(strconv.Itoa(v.(int)), false)
			}
		}

		if v, ok := awsAccount["variable"]; ok {
			action.Properties["Octopus.Action.AwsAccount.Variable"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := awsAccount["use_instance_role"]; ok {
			action.Properties["Octopus.Action.AwsAccount.UseInstanceRole"] = core.NewPropertyValue(cases.Title(language.Und, cases.NoLower).String(strconv.FormatBool(v.(bool))), false)
		}
	}

	if v, ok := flattenedAction["azure_account"]; ok && len(v.(*schema.Set).List()) > 0 {
		action.Properties["Octopus.Action.Terraform.AzureAccount"] = core.NewPropertyValue("True", false)

		azureAccount := v.(*schema.Set).List()[0].(map[string]interface{})

		if v, ok := azureAccount["variable"]; ok {
			action.Properties["Octopus.Action.AzureAccount.Variable"] = core.NewPropertyValue(v.(string), false)
		}
	}

	if v, ok := flattenedAction["google_cloud_account"]; ok && len(v.(*schema.Set).List()) > 0 {
		action.Properties["Octopus.Action.Terraform.GoogleCloudAccount"] = core.NewPropertyValue("True", false)

		googleAccount := v.(*schema.Set).List()[0].(map[string]interface{})

		if v, ok := googleAccount["variable"]; ok {
			action.Properties["Octopus.Action.GoogleCloudAccount.Variable"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := googleAccount["use_vm_service_account"]; ok {
			action.Properties["Octopus.Action.GoogleCloud.UseVMServiceAccount"] = core.NewPropertyValue(cases.Title(language.Und, cases.NoLower).String(strconv.FormatBool(v.(bool))), false)
		}

		if v, ok := googleAccount["impersonate_service_account"]; ok {
			action.Properties["Octopus.Action.GoogleCloud.ImpersonateServiceAccount"] = core.NewPropertyValue(cases.Title(language.Und, cases.NoLower).String(strconv.FormatBool(v.(bool))), false)
		}

		if v, ok := googleAccount["service_account_emails"]; ok {
			action.Properties["Octopus.Action.GoogleCloud.ServiceAccountEmails"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := googleAccount["zone"]; ok {
			action.Properties["Octopus.Action.GoogleCloud.Zone"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := googleAccount["region"]; ok {
			action.Properties["Octopus.Action.GoogleCloud.Region"] = core.NewPropertyValue(v.(string), false)
		}

		if v, ok := googleAccount["project"]; ok {
			action.Properties["Octopus.Action.GoogleCloud.Project"] = core.NewPropertyValue(v.(string), false)
		}
	}

	return action
}

func flattenTerraformTemplateAdvancedOptions(properties map[string]core.PropertyValue) []interface{} {
	if len(properties) == 0 {
		return nil
	}

	flattenedMap := map[string]interface{}{}

	for k, v := range properties {
		switch k {
		case "Octopus.Action.Terraform.AdditionalActionParams":
			flattenedMap["apply_parameters"] = v.Value
		case "Octopus.Action.Terraform.AdditionalInitParams":
			flattenedMap["init_parameters"] = v.Value
		case "Octopus.Action.Terraform.AllowPluginDownloads":
			allowPluginDownloads, _ := strconv.ParseBool(v.Value)
			flattenedMap["allow_additional_plugin_downloads"] = allowPluginDownloads
		case "Octopus.Action.Terraform.PluginsDirectory":
			flattenedMap["plugin_cache_directory"] = v.Value
		case "Octopus.Action.Terraform.Workspace":
			flattenedMap["workspace"] = v.Value
		}
	}

	return []interface{}{flattenedMap}
}

func flattenTerraformTemplateGoogleAccount(properties map[string]core.PropertyValue) []interface{} {
	if len(properties) == 0 {
		return nil
	}

	flattenedMap := map[string]interface{}{}

	for k, v := range properties {
		switch k {
		case "Octopus.Action.GoogleCloudAccount.Variable":
			flattenedMap["variable"] = v.Value
		case "Octopus.Action.GoogleCloud.ServiceAccountEmails":
			flattenedMap["service_account_emails"] = v.Value
		case "Octopus.Action.GoogleCloud.Project":
			flattenedMap["project"] = v.Value
		case "Octopus.Action.GoogleCloud.Zone":
			flattenedMap["zone"] = v.Value
		case "Octopus.Action.GoogleCloud.Region":
			flattenedMap["region"] = v.Value
		case "Octopus.Action.GoogleCloud.UseVMServiceAccount":
			useVMServiceAccount, _ := strconv.ParseBool(v.Value)
			flattenedMap["use_vm_service_account"] = useVMServiceAccount
		case "Octopus.Action.GoogleCloud.ImpersonateServiceAccount":
			impersonateServiceAccount, _ := strconv.ParseBool(v.Value)
			flattenedMap["impersonate_service_account"] = impersonateServiceAccount
		}
	}

	return []interface{}{flattenedMap}
}

func flattenTerraformTemplateAwsAccount(properties map[string]core.PropertyValue) []interface{} {
	if len(properties) == 0 {
		return nil
	}

	flattenedMap := map[string]interface{}{}

	for k, v := range properties {
		switch k {
		case "Octopus.Action.Aws.AssumeRole":
			if v.Value == "True" {
				flattenedMap["role"] = flattenTerraformTemplateAwsRole(properties)
			}
		case "Octopus.Action.Aws.Region":
			flattenedMap["region"] = v.Value
		case "Octopus.Action.AwsAccount.Variable":
			flattenedMap["variable"] = v.Value
		case "Octopus.Action.AwsAccount.UseInstanceRole":
			useInstanceRole, _ := strconv.ParseBool(v.Value)
			flattenedMap["use_instance_role"] = useInstanceRole
		}
	}

	return []interface{}{flattenedMap}
}

func flattenTerraformTemplateAwsRole(properties map[string]core.PropertyValue) []interface{} {
	if len(properties) == 0 {
		return nil
	}

	flattenedMap := map[string]interface{}{}

	for k, v := range properties {
		switch k {
		case "Octopus.Action.Aws.AssumedRoleArn":
			flattenedMap["arn"] = v.Value
		case "Octopus.Action.Aws.AssumeRoleExternalId":
			flattenedMap["external_id"] = v.Value
		case "Octopus.Action.Aws.AssumedRoleSession":
			flattenedMap["role_session_name"] = v.Value
		case "Octopus.Action.Aws.AssumeRoleSessionDurationSeconds":
			duration, _ := strconv.ParseInt(v.Value, 10, 32)
			flattenedMap["session_duration"] = duration
		}
	}

	return []interface{}{flattenedMap}
}

func flattenTerraformTemplateAzureAccount(properties map[string]core.PropertyValue) []interface{} {
	if len(properties) == 0 {
		return nil
	}

	flattenedMap := map[string]interface{}{}

	if v, ok := properties["Octopus.Action.AzureAccount.Variable"]; ok {
		flattenedMap["variable"] = v.Value
	}

	return []interface{}{flattenedMap}
}

func flattenApplyTerraformTemplateAction(action *deployments.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	for k, v := range action.Properties {
		switch k {
		case "Octopus.Action.RunOnServer":
			runOnServer, _ := strconv.ParseBool(v.Value)
			flattenedAction["run_on_server"] = runOnServer
		case "Octopus.Action.Script.ScriptSource":
			if v.Value == "Package" {
				flattenedAction["template"] = flattenTerraformTemplate(action.Properties)
			}
		case "Octopus.Action.Terraform.AzureAccount":
			if v.Value == "True" {
				flattenedAction["azure_account"] = flattenTerraformTemplateAzureAccount(action.Properties)
			}
		case "Octopus.Action.Terraform.GoogleCloudAccount":
			if v.Value == "True" {
				flattenedAction["google_cloud_account"] = flattenTerraformTemplateGoogleAccount(action.Properties)
			}
		case "Octopus.Action.Terraform.ManagedAccount":
			if v.Value == "AWS" {
				flattenedAction["aws_account"] = flattenTerraformTemplateAwsAccount(action.Properties)
			}
		case "Octopus.Action.Terraform.Template":
			flattenedAction["inline_template"] = v.Value
		case "Octopus.Action.Terraform.TemplateParameters":
			flattenedAction["template_parameters"] = v.Value
		}
	}

	flattenedAction["advanced_options"] = flattenTerraformTemplateAdvancedOptions(action.Properties)

	return flattenedAction
}

func getApplyTerraformTemplateActionSchema() *schema.Schema {
	actionSchema, element := getActionSchema()
	addExecutionLocationSchema(element)
	addTerraformTemplateAdvancedOptionsSchema(element)
	addTerraformTemplateAwsAccountSchema(element)
	addTerraformTemplateAzureAccountSchema(element)
	addTerraformTemplateGoogleAccountSchema(element)
	addTerraformTemplateParametersSchema(element)
	addTerraformTemplateSchema(element)
	addTerraformInlineTemplateSchema(element)
	addPrimaryPackageSchema(element, false)

	return actionSchema
}
