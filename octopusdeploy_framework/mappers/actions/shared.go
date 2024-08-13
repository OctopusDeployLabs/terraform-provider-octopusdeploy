package actions

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

func mapBaseDeploymentActionToState(ctx context.Context, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	newAction["can_be_used_for_project_versioning"] = types.BoolValue(action.CanBeUsedForProjectVersioning)
	newAction["is_disabled"] = types.BoolValue(action.IsDisabled)
	newAction["is_required"] = types.BoolValue(action.IsRequired)
	newAction["channels"] = util.FlattenStringList(action.Channels)
	newAction["condition"] = types.StringValue(action.Condition)
	newAction["container"] = mapContainerToState(action.Container)
	newAction["environments"] = util.FlattenStringList(action.Environments)
	newAction["excluded_environments"] = util.FlattenStringList(action.ExcludedEnvironments)
	newAction["id"] = types.StringValue(action.ID)
	newAction["name"] = types.StringValue(action.Name)
	newAction["slug"] = types.StringValue(action.Slug)
	newAction["notes"] = types.StringValue(action.Notes)

	updatedProperties, diags := MapPropertiesToState(ctx, action.Properties)
	if diags.HasError() {
		return diags
	}
	newAction["properties"] = updatedProperties

	newAction["tenant_tags"] = util.FlattenStringList(action.TenantTags)

	if v, ok := action.Properties["Octopus.Action.EnabledFeatures"]; ok {
		newAction["features"] = util.FlattenStringList(strings.Split(v.Value, ","))
	} else {
		newAction["features"] = types.ListNull(types.StringType)
	}

	attrTypes := map[string]attr.Type{"id": types.StringType, "version": types.StringType}
	if v, ok := action.Properties["Octopus.Action.Template.Id"]; ok {
		actionTemplate := map[string]attr.Value{
			"id": types.StringValue(v.Value),
		}

		if v, ok := action.Properties["Octopus.Action.Template.Version"]; ok {
			actionTemplate["version"] = types.StringValue(v.Value)
		}

		list := make([]attr.Value, 1)
		list[0] = types.ObjectValueMust(attrTypes, actionTemplate)

		newAction["action_template"] = types.ListValueMust(types.ObjectType{AttrTypes: attrTypes}, list)
	} else {
		newAction["action_template"] = types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	hasPackageReference := false
	if len(action.Packages) > 0 {
		var packageReferences []attr.Value
		for _, packageReference := range action.Packages {
			packageReferenceAttribute, diags := mapPackageReferenceToState(ctx, packageReference)
			if diags.HasError() {
				return diags
			}
			if len(packageReference.Name) == 0 {

				newAction["primary_package"] = types.ListValueMust(types.ObjectType{AttrTypes: GetPackageReferenceAttrTypes(true)}, []attr.Value{types.ObjectValueMust(GetPackageReferenceAttrTypes(true), packageReferenceAttribute)})
				// TODO: consider these properties
				// actionProperties["Octopus.Action.Package.DownloadOnTentacle"] = packageReference.AcquisitionLocation
				// flattenedAction["properties"] = actionProperties
			} else {
				packageReferences = append(packageReferences, types.ObjectValueMust(GetPackageReferenceAttrTypes(false), packageReferenceAttribute))
				newAction["package"] = types.ListValueMust(types.ObjectType{AttrTypes: GetPackageReferenceAttrTypes(false)}, packageReferences)
				hasPackageReference = true
			}
		}
	} else {
		newAction["primary_package"] = types.ListNull(types.ObjectType{AttrTypes: GetPackageReferenceAttrTypes(true)})
	}

	if !hasPackageReference {
		newAction["package"] = types.ListNull(types.ObjectType{AttrTypes: GetPackageReferenceAttrTypes(false)})
	}

	if len(action.GitDependencies) > 0 {
		var gitDepenedencyList []attr.Value
		gitDepenedencyList = append(gitDepenedencyList, types.ObjectValueMust(GetGitDependencyAttrTypes(), mapGitDependencyToState(action.GitDependencies[0])))
		newAction["git_dependency"] = types.SetValueMust(types.ObjectType{AttrTypes: GetGitDependencyAttrTypes()}, gitDepenedencyList)
	} else {
		newAction["git_dependency"] = types.SetNull(types.ObjectType{AttrTypes: GetGitDependencyAttrTypes()})
	}

	return nil
}

func mapContainerToState(container *deployments.DeploymentActionContainer) types.List {
	attributeTypes := map[string]attr.Type{"feed_id": types.StringType, "image": types.StringType}
	if container == nil || (container.Image == "" && container.FeedID == "") {
		return types.ListNull(types.ObjectType{AttrTypes: attributeTypes})
	}

	list := make([]attr.Value, 0)
	containerAttributes := map[string]attr.Value{
		"feed_id": types.StringValue(container.FeedID),
		"image":   types.StringValue(container.Image),
	}

	list = append(list, types.ObjectValueMust(attributeTypes, containerAttributes))
	return types.ListValueMust(types.ObjectType{AttrTypes: attributeTypes}, list)
}

func mapPropertyToStateString(action *deployments.DeploymentAction, actionState map[string]attr.Value, propertyName string, attrName string) {
	if v, ok := action.Properties[propertyName]; ok {
		actionState[attrName] = types.StringValue(v.Value)
	} else {
		actionState[attrName] = types.StringValue("")
	}
}

func GetGitDependencyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"repository_uri":      types.StringType,
		"default_branch":      types.StringType,
		"git_credential_type": types.StringType,
		"file_path_filters":   types.ListType{ElemType: types.StringType},
		"git_credential_id":   types.StringType,
	}
}

func mapGitDependencyToState(gitDependency *gitdependencies.GitDependency) map[string]attr.Value {
	return map[string]attr.Value{
		"repository_uri":      types.StringValue(gitDependency.RepositoryUri),
		"default_branch":      types.StringValue(gitDependency.DefaultBranch),
		"git_credential_type": types.StringValue(gitDependency.GitCredentialType),
		"file_path_filters":   util.FlattenStringList(gitDependency.FilePathFilters),
		"git_credential_id":   types.StringValue(gitDependency.GitCredentialId),
	}
}

func GetPackageReferenceAttrTypes(isPrimaryPackage bool) map[string]attr.Type {
	attrTypes := map[string]attr.Type{
		"acquisition_location": types.StringType,
		"feed_id":              types.StringType,
		"id":                   types.StringType,
		"package_id":           types.StringType,
		"properties":           types.MapType{types.StringType},
	}

	if !isPrimaryPackage {
		attrTypes["name"] = types.StringType
		attrTypes["extract_during_deployment"] = types.BoolType
	}

	return attrTypes
}

func mapPackageReferenceToState(ctx context.Context, packageReference *packages.PackageReference) (map[string]attr.Value, diag.Diagnostics) {
	properties, diags := types.MapValueFrom(ctx, types.StringType, packageReference.Properties)
	if diags.HasError() {
		return nil, diags
	}

	reference := map[string]attr.Value{
		"acquisition_location": types.StringValue(packageReference.AcquisitionLocation),
		"feed_id":              types.StringValue(packageReference.FeedID),
		"id":                   types.StringValue(packageReference.ID),
		"package_id":           types.StringValue(packageReference.PackageID),
		"properties":           properties,
	}

	if v, ok := packageReference.Properties["Extract"]; ok {
		if len(packageReference.Name) > 0 {
			extractDuringDeployment, _ := strconv.ParseBool(v)
			reference["extract_during_deployment"] = types.BoolValue(extractDuringDeployment)
			reference["name"] = types.StringValue(packageReference.Name)
		}
	}

	return reference, nil
}

func MapPropertiesToState(ctx context.Context, properties map[string]core.PropertyValue) (types.Map, diag.Diagnostics) {
	if properties == nil || len(properties) == 0 {
		return types.MapNull(types.StringType), nil
	}

	stateMap := make(map[string]attr.Value)
	for key, value := range properties {
		if !value.IsSensitive {
			stateMap[key] = types.StringValue(value.Value)
		}
	}

	return types.MapValueFrom(ctx, types.StringType, stateMap)
}

func GetActionAttributes(actionAttribute attr.Value) map[string]attr.Value {
	actionAttrsList := actionAttribute.(types.List)
	if actionAttrsList.IsNull() {
		return nil
	}

	actionAttrsElements := actionAttrsList.Elements()
	if len(actionAttrsElements) == 0 {
		return nil
	}

	return actionAttrsElements[0].(types.Object).Attributes()
}

func GetBaseAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := GetActionAttributes(actionAttribute)
	if actionAttrs == nil {
		return nil
	}

	var name string
	util.SetString(actionAttrs, "name", &name)

	var actionType string
	util.SetString(actionAttrs, "action_type", &actionType)

	action := deployments.NewDeploymentAction(name, actionType)
	util.SetString(actionAttrs, "id", &action.ID)
	util.SetString(actionAttrs, "condition", &action.Condition)
	util.SetBool(actionAttrs, "is_disabled", &action.IsDisabled)
	util.SetBool(actionAttrs, "is_required", &action.IsRequired)
	util.SetString(actionAttrs, "notes", &action.Notes)
	action.Channels = util.ExpandStringList(actionAttrs["channels"].(types.List))

	action.Container = getContainer(actionAttrs)

	action.Environments = getArray(actionAttrs, "environments")
	action.ExcludedEnvironments = getArray(actionAttrs, "excluded_environments")

	// TODO map properties from state
	for k, v := range actionAttrs["properties"].(types.Map).Elements() {
		action.Properties[k] = core.NewPropertyValue(v.(types.String).ValueString(), false)
	}

	features := getArray(actionAttrs, "features")
	if features != nil {
		action.Properties["Octopus.Action.EnabledFeatures"] = core.NewPropertyValue(strings.Join(features, ","), false)
	}

	if v, ok := actionAttrs["run_on_server"]; ok {
		runOnServer := v.(types.Bool).ValueBool()
		action.Properties["Octopus.Action.RunOnServer"] = core.NewPropertyValue(formatBoolForActionProperty(runOnServer), false)
	}

	util.SetString(actionAttrs, "slug", &action.Slug)

	tenantTags := getArray(actionAttrs, "tenant_tags")
	if tenantTags != nil {
		action.TenantTags = tenantTags
	}

	util.SetString(actionAttrs, "worker_pool_id", &action.WorkerPool)
	util.SetString(actionAttrs, "worker_pool_variable", &action.WorkerPoolVariable)

	setActionTemplate(actionAttrs, action)
	setPrimaryPackage(actionAttrs, action)

	for key, attr := range actionAttrs {
		if key == "package" {
			for _, p := range attr.(types.List).Elements() {
				pkg := getPackageReference(p.(types.Object).Attributes())
				action.Packages = append(action.Packages, pkg)
			}
		}

		if key == "git_dependency" && len(attr.(types.Set).Elements()) > 0 {
			for _, gd := range attr.(types.Set).Elements() {
				gitDependency := getGitDependency(gd.(types.Object).Attributes())
				action.GitDependencies = append(action.GitDependencies, gitDependency)
			}

		}
	}

	// Polyfill the Kubernetes Object status check to default to true if not specified for Kubernetes steps
	switch actionType {
	case "Octopus.KubernetesDeployContainers":
		fallthrough
	case "Octopus.KubernetesDeployRawYaml":
		fallthrough
	case "Octopus.KubernetesDeployService":
		fallthrough
	case "Octopus.KubernetesDeployIngress":
		fallthrough
	case "Octopus.KubernetesDeployConfigMap":
		fallthrough
	case "Octopus.Kustomize":
		if _, exists := action.Properties["Octopus.Action.Kubernetes.ResourceStatusCheck"]; !exists {
			action.Properties["Octopus.Action.Kubernetes.ResourceStatusCheck"] = core.NewPropertyValue(formatBoolForActionProperty(true), false)
		}
		break
	}

	return action
}

func setActionTemplate(attrs map[string]attr.Value, action *deployments.DeploymentAction) {
	actionTemplate := getAttributesForSingleElementList(attrs, "action_template")
	if actionTemplate != nil {
		if id, ok := actionTemplate["id"]; ok {
			action.Properties["Octopus.Action.Template.Id"] = core.NewPropertyValue(id.(types.String).ValueString(), false)
		}

		if v, ok := actionTemplate["version"]; ok {
			action.Properties["Octopus.Action.Template.Version"] = core.NewPropertyValue(v.(types.String).ValueString(), false)
		}
	}
}

func getAttributesForSingleElementList(attrs map[string]attr.Value, s string) map[string]attr.Value {
	if a, ok := attrs[s]; ok {
		list := a.(types.List)
		if len(list.Elements()) > 0 {
			return list.Elements()[0].(types.Object).Attributes()
		}
	}

	return nil
}

func getArray(attrs map[string]attr.Value, s string) []string {
	if a, ok := attrs[s]; ok {
		list := a.(types.List)
		return util.GetStringSlice(list)
	}

	return nil
}

func formatBoolForActionProperty(b bool) string {
	return cases.Title(language.Und, cases.NoLower).String(strconv.FormatBool(b))
}

func getContainer(attrs map[string]attr.Value) *deployments.DeploymentActionContainer {
	if c, ok := attrs["container"]; ok {
		if c == nil || c.IsNull() || c.IsUnknown() {
			return nil
		}

		containerAttrs := c.(types.List).Elements()[0].(types.Object).Attributes()
		actionContainer := &deployments.DeploymentActionContainer{}
		util.SetString(containerAttrs, "feed_id", &actionContainer.FeedID)
		util.SetString(containerAttrs, "image", &actionContainer.Image)
		return actionContainer
	}

	return nil
}

func getGitDependency(gitAttrs map[string]attr.Value) *gitdependencies.GitDependency {
	gitDependency := &gitdependencies.GitDependency{}
	util.SetString(gitAttrs, "repository_uri", &gitDependency.RepositoryUri)
	util.SetString(gitAttrs, "default_branch", &gitDependency.DefaultBranch)
	util.SetString(gitAttrs, "git_credential_type", &gitDependency.GitCredentialType)
	util.SetString(gitAttrs, "git_credential_id", &gitDependency.GitCredentialId)
	gitDependency.FilePathFilters = getArray(gitAttrs, "file_path_filters")
	return gitDependency
}

func getPackageReference(attrs map[string]attr.Value) *packages.PackageReference {
	pkg := &packages.PackageReference{Properties: map[string]string{}}
	util.SetString(attrs, "acquisition_location", &pkg.AcquisitionLocation)
	util.SetString(attrs, "feed_id", &pkg.FeedID)
	util.SetString(attrs, "name", &pkg.Name)
	util.SetString(attrs, "package_id", &pkg.PackageID)

	var extractDuringDeployment bool
	util.SetBool(attrs, "extract_during_deployment", &extractDuringDeployment)
	pkg.Properties["Extract"] = formatBoolForActionProperty(extractDuringDeployment)

	if properties := attrs["properties"]; properties != nil {
		propertyMap := properties.(types.Map).Elements()
		for k, v := range propertyMap {
			pkg.Properties[k] = v.(types.String).ValueString()
		}
	}

	return pkg
}

func setPrimaryPackage(attrs map[string]attr.Value, action *deployments.DeploymentAction) {
	primaryPackageAttributes := getAttributesForSingleElementList(attrs, "primary_package")
	if primaryPackageAttributes == nil {
		return
	}

	primaryPackageReference := getPackageReference(primaryPackageAttributes)
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

func mapAttributeToProperty(action *deployments.DeploymentAction, attrs map[string]attr.Value, attributeName string, propertyName string) {
	var value string
	util.SetString(attrs, attributeName, &value)
	if value != "" {
		action.Properties[propertyName] = core.NewPropertyValue(value, false)
	}
}

func mapBooleanAttributeToProperty(action *deployments.DeploymentAction, attrs map[string]attr.Value, attributeName string, propertyName string) {
	if v, ok := attrs[attributeName]; ok {
		b := v.(types.Bool).ValueBool()
		action.Properties[propertyName] = core.NewPropertyValue(formatBoolForActionProperty(b), false)
	}
}

func mapPropertyToStateBool(action *deployments.DeploymentAction, actionState map[string]attr.Value, propertyName string, attrName string, defaultValue bool) {
	if v, ok := action.Properties[propertyName]; ok {
		parsedValue, _ := strconv.ParseBool(v.Value)
		actionState[attrName] = types.BoolValue(parsedValue)
	} else {
		actionState[attrName] = types.BoolValue(defaultValue)
	}
}

func ensureFeatureIsEnabled(action *deployments.DeploymentAction, feature string) {
	const enabledFeatures = "Octopus.Action.EnabledFeatures"
	if len(action.Properties[enabledFeatures].Value) == 0 {
		action.Properties[enabledFeatures] = core.NewPropertyValue(feature, false)
	} else {
		// fixing https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/issues/641
		currentFeatures := action.Properties[enabledFeatures].Value
		if !strings.Contains(currentFeatures, feature) {
			action.Properties[enabledFeatures] = core.NewPropertyValue(currentFeatures+","+feature, false)
		}
	}
}
