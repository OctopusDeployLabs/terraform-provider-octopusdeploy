data "octopusdeploy_machine_policies" "default_machine_policy" {
  ids          = null
  partial_name = "Default Machine Policy"
  skip         = 0
  take         = 1
}

resource octopusdeploy_kubernetes_cluster_deployment_target test_eks {
  cluster_url                       = "https://cluster"
  environments                      = [octopusdeploy_environment.test_environment.id]
  name                              = "Test"
  roles                             = ["eks"]
  cluster_certificate               = ""
  cluster_certificate_path          = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
  machine_policy_id                 = data.octopusdeploy_machine_policies.default_machine_policy.machine_policies[0].id
  namespace                         = ""
  skip_tls_verification             = true
  tenant_tags                       = []
  tenanted_deployment_participation = "Untenanted"
  tenants                           = []
  thumbprint                        = ""
  uri                               = ""

  container {
    feed_id = ""
    image   = ""
  }

  endpoint {
    communication_style    = "Kubernetes"
  }

  pod_authentication {
    token_path = "/var/run/secrets/kubernetes.io/serviceaccount/token"
  }
}