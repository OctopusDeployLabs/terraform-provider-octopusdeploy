data "octopusdeploy_workers" "example" {
  communication_styles  = ["TentaclePassive"]
  health_statuses       = ["Unavailable"]
  ids                   = ["Workers-123"]
  name                  = "Exact name"
  partial_name          = "Test"
  skip                  = 5
  take                  = 100
  is_disabled           = true
}