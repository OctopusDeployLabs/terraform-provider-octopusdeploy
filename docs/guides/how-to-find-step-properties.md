---
page_title: "How to Find Step Properties"
subcategory: "Guides"
---

# How to Find Step Properties

When defining step properties for a deployment or runbook process it can be challenging to discover the required key and
values for a step. This guides goal is to help with the discovery of step properties.

To discover properties for a step its best to first define deployment process within the UI including the steps you wish
to use within the provider. Once the deployment process is defined and the step are configured do the following:

1. Click the ellipsis to the right of the `save` button on the Process Editor page
1. From the menu select `Download as JSON`
1. Open the downloaded JSON file in an editor of your choice
1. Within the JSON you will find a key `Steps`, this is an array of step objects. The objects are in the run order.
1. Locate the key `Properties` within a step object, this maps to the key/value pairs for the `properties` field on the `process_step` resource in HCL.
1. Locate the key `Actions` within a step object, this is and array of actions. Every Step object will have at lest one action, the first action is always related to the step object while additional actions are used to define child steps.
1. Locate the key `Properties` within the first action object of the array, this maps the key/value pairs for the `execution_properties` field on the `process_step` resource in HCL.
1. The following action objects within the array, if any, will map to the `property` key to the `execution_properties` field on the `process_child_step` resource within HCL.


Note that in the future we plan to improve the discoverability of step properties.
