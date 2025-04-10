resource "octopusdeploy_aws_elastic_container_registry" "example" {
  access_key = "access-key"
  name       = "Test AWS Elastic Container Registry (OK to Delete)"
  region     = "us-east-1"
  secret_key = "secret-key"
}

resource "octopusdeploy_aws_elastic_container_registry" "example_with_oidc" {
  name       = "Test AWS Elastic Container Registry with OIDC (OK to Delete)"
  region     = "us-east-1"
  oidc_authentication = {
    session_duration = 3600
    role_arn = "role_arn_value"
    subject_keys = [ "feed", "space" ]
  }
}
