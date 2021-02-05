data "octopusdeploy_lifecycles" "example" {
  ids          = ["Lifecycles-123", "Lifecycles-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
