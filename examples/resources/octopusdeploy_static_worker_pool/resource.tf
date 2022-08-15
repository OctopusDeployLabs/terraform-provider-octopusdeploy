resource "octopusdeploy_static_worker_pool" "example" {
  description = "Description for the static worker pool."
  is_default  = true
  name        = "Test Static Worker Pool (OK to Delete)"
  sort_order  = 5
}
