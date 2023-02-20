resource "octopusdeploy_aws_elastic_container_registry" "feed_ecr" {
  name                                 = "ECR"
  access_key                           = var.feed_ecr_access_key
  secret_key                           = var.feed_ecr_secret_key
  package_acquisition_location_options = ["ExecutionTarget", "NotAcquired"]
  region                               = "us-east-1"
}
variable "feed_ecr_access_key" {
  type        = string
  nullable    = false
  sensitive   = true
  description = "The access key used for the ECR feed"
}
variable "feed_ecr_secret_key" {
  type        = string
  nullable    = false
  sensitive   = true
  description = "The secret key used for the ECR feed"
}