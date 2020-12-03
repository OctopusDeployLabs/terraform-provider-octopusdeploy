resource "octopusdeploy_lifecycle" "example" {
  description = "This is the default lifecycle."
  name        = "Test Lifecycle (OK to Delete)"

  release_retention_policy {
    quantity_to_keep    = 1
    should_keep_forever = true
    unit                = "Days"
  }

  tentacle_retention_policy {
    quantity_to_keep    = 30
    should_keep_forever = false
    unit                = "Items"
  }

  phase {
    automatic_deployment_targets = ["Environments-321"]
    name                         = "foo"

    release_retention_policy {
      quantity_to_keep    = 1
      should_keep_forever = true
      unit                = "Days"
    }

    tentacle_retention_policy {
      quantity_to_keep    = 30
      should_keep_forever = false
      unit                = "Items"
    }
  }

  phase {
    is_optional_phase           = true
    name                        = "bar"
    optional_deployment_targets = ["Environments-321"]
  }
}
