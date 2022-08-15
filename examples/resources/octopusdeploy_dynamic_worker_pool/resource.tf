resource "octopusdeploy_dynamic_worker_pool" "example" {
  description = "Description for the dynamic worker pool."
  is_default  = true
  name        = "Test Dynamic Worker Pool (OK to Delete)"
  sort_order  = 5
  worker_type = "UbuntuDefault"
}
