resource "octopusdeploy_static_worker_pool" "workerpool_docker" {
  name        = "Docker"
  description = "A test worker pool"
  is_default  = false
  sort_order  = 3
}

resource "octopusdeploy_static_worker_pool" "workerpool_example" {
  name        = "Example"
  description = "My other example worker pool"
  is_default  = false
  sort_order  = 4
}
