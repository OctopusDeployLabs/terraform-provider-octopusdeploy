---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "octopusdeploy_step_template Resource - terraform-provider-octopusdeploy"
subcategory: ""
description: |-
  This resource manages step_templates in Octopus Deploy.
---

# octopusdeploy_step_template (Resource)

This resource manages step_templates in Octopus Deploy.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action_type` (String) The action type of the step template
- `name` (String) The name of this resource.
- `packages` (Attributes List) Package information for the step template (see [below for nested schema](#nestedatt--packages))
- `parameters` (Attributes List) List of parameters that can be used in Step Template. (see [below for nested schema](#nestedatt--parameters))
- `properties` (Map of String) Properties for the step template
- `step_package_id` (String) The ID of the step package

### Optional

- `community_action_template_id` (String) The ID of the community action template
- `description` (String) The description of this step_template.
- `space_id` (String) The space ID associated with this step_template.

### Read-Only

- `id` (String) The unique ID for this resource.
- `version` (Number) The version of the step template

<a id="nestedatt--packages"></a>
### Nested Schema for `packages`

Required:

- `feed_id` (String) ID of the feed.
- `name` (String) The name of this resource.
- `properties` (Attributes) Properties for the package. (see [below for nested schema](#nestedatt--packages--properties))

Optional:

- `acquisition_location` (String) Acquisition location for the package.
- `package_id` (String) The ID of the package to use.

Read-Only:

- `id` (String) The unique ID for this resource.

<a id="nestedatt--packages--properties"></a>
### Nested Schema for `packages.properties`

Required:

- `selection_mode` (String) The selection mode.

Optional:

- `extract` (String) If the package should extract.
- `package_parameter_name` (String) The name of the package parameter
- `purpose` (String) The purpose of this property.



<a id="nestedatt--parameters"></a>
### Nested Schema for `parameters`

Required:

- `id` (String) The id for the property.
- `name` (String) The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods. Example: `ServerName`

Optional:

- `default_value` (String) A default value for the parameter, if applicable. This can be a hard-coded value or a variable reference.
- `display_settings` (Map of String) The display settings for the parameter.
- `help_text` (String) The help presented alongside the parameter input.
- `label` (String) The label shown beside the parameter when presented in the deployment process. Example: `Server name`.

