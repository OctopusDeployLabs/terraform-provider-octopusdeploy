resource "octopusdeploy_azure_container_registry" "example" {
  name          = "Test Azure Container Registry (OK to Delete)"
  feed_uri      = "https://test-azure.azurecr.io"
  username      = "username"
  password      = "password"
}

resource "octopusdeploy_azure_container_register" "example_with_oidc" {
    name          = "Test Azure Container Registry (OK to Delete)"
    feed_uri      = "https://test-azure.azurecr.io"
    oidc_authentication = {
      client_id     = "client_id"
      tenant_id     = "tenant_id"
      audience      = "audience"
      subject_keys = ["feed", "space"]
    } 
}