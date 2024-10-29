resource "octopusdeploy_azure_container_registry" "example" {
  name          = "Test Azure Container Registry (OK to Delete)"
  feed_uri      = "https://test-azure.azurecr.io"
  username      = "username"
  password      = "password"
}
