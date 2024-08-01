resource "octopusdeploy_space" "octopus_space_test" {
  name                  = var.octopus_space_name
  is_default            = false
  is_task_queue_stopped = false
  description           = var.octopus_space_description
  space_managers_teams  = ["teams-administrators"]
}

output "octopus_space_id" {
  value = octopusdeploy_space.octopus_space_test.id
}

variable "octopus_space_name" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The name of the new space"
  default     = "Test"
}

variable "octopus_space_description" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The description of the new space"
  default     = "My test space"
}
