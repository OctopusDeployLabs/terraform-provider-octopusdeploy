package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTriggerAction(d map[string]interface{}) *octopusdeploy.TriggerAction {
	triggerAction := octopusdeploy.NewTriggerAction()

	if v, ok := d["id"]; ok {
		triggerAction.ID = v.(string)
	}
	if v, ok := d["action_type"]; ok {
		triggerAction.ActionType = v.(string)
	}

	switch triggerAction.ActionType {
	case "AutoDeploy":
		expandAutoDeployTriggerAction(triggerAction, d)
	case "DeployLatestRelease":
		expandDeployLatestReleaseTriggerAction(triggerAction, d)
	case "DeployNewRelease":
		expandDeployNewReleaseTriggerAction(triggerAction, d)
	case "RunRunbook":
		expandRunRunbookTriggerAction(triggerAction, d)
	}
	return triggerAction
}

func flattenTriggerAction(triggerAction *octopusdeploy.TriggerAction) map[string]interface{} {
	if triggerAction == nil {
		return nil
	}
	flattened := make(map[string]interface{})

	flattened["action_type"] = triggerAction.ActionType
	flattened["id"] = triggerAction.ID

	switch triggerAction.ActionType {
	case "AutoDeploy":
		flattenAutoDeployTriggerAction(triggerAction, flattened)
	case "DeployLatestRelease":
		flattenDeployLatestReleaseTriggerAction(triggerAction, flattened)
	case "DeployNewRelease":
		flattenDeployNewReleaseTriggerAction(triggerAction, flattened)
	case "RunRunbook":
		flattenRunRunbookTriggerAction(triggerAction, flattened)
	}
	return flattened
}

func getTriggerActionSchema() (*schema.Schema, *schema.Resource) {
	element := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": getIDSchema(),
			"action_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	triggerActionSchema := &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem:     element,
		Optional: true,
	}
	return triggerActionSchema, element
}

func flattenScopedDeploymentTriggerAction(triggerAction *octopusdeploy.TriggerAction, flattened map[string]interface{}) {
	flattened["channel_id"] = triggerAction.ChannelID
	flattened["tenant_ids"] = triggerAction.TenantIDs
	flattened["tenant_tags"] = triggerAction.TenantTags
}
func expandScopedDeploymentTriggerAction(triggerAction *octopusdeploy.TriggerAction, d map[string]interface{}) {
	if v, ok := d["channel_id"]; ok {
		triggerAction.ChannelID = v.(string)
	}
	if v, ok := d["tenant_ids"]; ok {
		triggerAction.TenantIDs = getSliceFromTerraformTypeList(v)
	}
	if v, ok := d["tenant_tags"]; ok {
		triggerAction.TenantTags = getSliceFromTerraformTypeList(v)
	}
}
func getScopedDeploymentTriggerActionSchema() *schema.Schema {
	schema, element := getTriggerActionSchema()
	addScopedDeploymentTriggerActionSchema(element)
	return schema
}

func addScopedDeploymentTriggerActionSchema(element *schema.Resource) {
	element.Schema["channel_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	element.Schema["tenant_ids"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	element.Schema["tenant_tags"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func flattenAutoDeployTriggerAction(triggerAction *octopusdeploy.TriggerAction, flattened map[string]interface{}) {
	flattened["should_redeploy"] = triggerAction.ShouldRedeployWhenMachineHasBeenDeployedTo
}
func expandAutoDeployTriggerAction(triggerAction *octopusdeploy.TriggerAction, d map[string]interface{}) {
	if v, ok := d["should_redeploy"]; ok {
		triggerAction.ShouldRedeployWhenMachineHasBeenDeployedTo = v.(bool)
	}
}

func getAutoDeployTriggerActionSchema() *schema.Schema {
	schema, element := getTriggerActionSchema()
	addAutoDeployTriggerActionSchema(element)
	return schema
}

func addAutoDeployTriggerActionSchema(element *schema.Resource) {
	element.Schema["should_redeploy"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Enable to re-deploy to the deployment targets even if they are already up-to-date with the current deployment.",
		Required:    true,
	}
}

func flattenRunRunbookTriggerAction(triggerAction *octopusdeploy.TriggerAction, flattened map[string]interface{}) {
	flattened["runbook_id"] = triggerAction.RunbookID
	flattened["environment_ids"] = triggerAction.SourceEnvironmentIDs
	flattened["tenant_ids"] = triggerAction.TenantIDs
	flattened["tenant_tagss"] = triggerAction.TenantTags
}
func expandRunRunbookTriggerAction(triggerAction *octopusdeploy.TriggerAction, d map[string]interface{}) {
	if v, ok := d["runbook_id"]; ok {
		triggerAction.RunbookID = v.(string)
	}
	if v, ok := d["environment_ids"]; ok {
		triggerAction.SourceEnvironmentIDs = getSliceFromTerraformTypeList(v)
	}
	if v, ok := d["tenant_ids"]; ok {
		triggerAction.TenantIDs = getSliceFromTerraformTypeList(v)
	}
	if v, ok := d["tenant_tags"]; ok {
		triggerAction.TenantTags = getSliceFromTerraformTypeList(v)
	}
}

func getRunRunbookTriggerActionSchema() *schema.Schema {
	schema, element := getTriggerActionSchema()
	addRunRunbookTriggerActionSchema(element)
	return schema
}
func addRunRunbookTriggerActionSchema(element *schema.Resource) {
	element.Schema["runbook_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	element.Schema["environment_ids"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	element.Schema["tenant_ids"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
	element.Schema["tenant_tags"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func flattenDeployNewReleaseTriggerAction(triggerAction *octopusdeploy.TriggerAction, flattened map[string]interface{}) {
	flattenScopedDeploymentTriggerAction(triggerAction, flattened)
	flattened["variables"] = triggerAction.Variables
	flattened["destination_environment_id"] = triggerAction.DestinationEnvironmentID
}
func expandDeployNewReleaseTriggerAction(triggerAction *octopusdeploy.TriggerAction, d map[string]interface{}) {
	expandScopedDeploymentTriggerAction(triggerAction, d)
	if v, ok := d["variables"]; ok {
		triggerAction.Variables = v.(string)
	}
	if v, ok := d["destination_environment_id"]; ok {
		triggerAction.DestinationEnvironmentID = v.(string)
	}
}

func getDeployNewReleaseTriggerActionSchema() *schema.Schema {
	schema, element := getTriggerActionSchema()
	addScopedDeploymentTriggerActionSchema(element)
	addDeployNewReleaseTriggerActionSchema(element)
	return schema
}
func addDeployNewReleaseTriggerActionSchema(element *schema.Resource) {
	element.Schema["variables"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	element.Schema["destination_environment_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
}

func flattenDeployLatestReleaseTriggerAction(triggerAction *octopusdeploy.TriggerAction, flattened map[string]interface{}) {
	flattenScopedDeploymentTriggerAction(triggerAction, flattened)
	flattened["variables"] = triggerAction.Variables
	flattened["destination_environment_id"] = triggerAction.DestinationEnvironmentID
	flattened["source_environment_ids"] = triggerAction.SourceEnvironmentIDs
	flattened["should_redeploy"] = triggerAction.RedeployCurrent
}
func expandDeployLatestReleaseTriggerAction(triggerAction *octopusdeploy.TriggerAction, d map[string]interface{}) {
	expandScopedDeploymentTriggerAction(triggerAction, d)
	if v, ok := d["variables"]; ok {
		triggerAction.Variables = v.(string)
	}
	if v, ok := d["destination_environment_id"]; ok {
		triggerAction.DestinationEnvironmentID = v.(string)
	}
	if v, ok := d["source_environment_ids"]; ok {
		triggerAction.SourceEnvironmentIDs = getSliceFromTerraformTypeList(v)
	}
	if v, ok := d["should_redeploy"]; ok {
		triggerAction.RedeployCurrent = v.(bool)
	}
}

func getDeployLatestReleaseTriggerActionSchema() *schema.Schema {
	schema, element := getTriggerActionSchema()
	addScopedDeploymentTriggerActionSchema(element)
	addDeployLatestReleaseTriggerActionSchema(element)
	return schema
}

func addDeployLatestReleaseTriggerActionSchema(element *schema.Resource) {
	element.Schema["variables"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	element.Schema["source_environment_ids"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	element.Schema["destination_environment_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	element.Schema["should_redeploy"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Enable to re-deploy to the deployment targets even if they are already up-to-date with the current deployment.",
		Default:     true,
		Optional:    true,
	}
}
