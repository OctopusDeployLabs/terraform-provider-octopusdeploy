resource "octopusdeploy_kubernetes_cluster_deployment_target" "k8s-target" {
  cluster_url                       = "https://example.com:1234/"
  environments                      = ["Environments-123", "Environment-321"]
  name                              = "Kubernetes Cluster Deployment Target (OK to Delete)"
  roles                             = ["Development Team", "System Administrators"]
  tenanted_deployment_participation = "Untenanted"

  aws_account_authentication {
    account_id   = "Accounts-123"
    cluster_name = "cluster-name"
  }
}
