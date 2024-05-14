data "octopusdeploy_kubernetes_agent_deployment_targets" "kubernetes_agent_deployment_targets" {
   deployment_id   = "Deployments-123"
   environments    = ["Environments-123", "Environments-321"]
   health_statuses = ["HasWarnings"]
   ids             = ["Machines-123", "Machines-321"]
   is_disabled     = false
   name            = "Kubernetes Agent"
   partial_name    = "Kubernetes Age"
   roles           = ["Roles-123", "Roles-321"]
   shell_names     = []
   skip            = 5
   take            = 100
   tenant_tags     = ["TagSet1/Tag"]
   tenants         = ["Tenants-123"]
   thumbprint      = "96203ED84246201C26A2F4360D7CBC36AC1D232D"
}
