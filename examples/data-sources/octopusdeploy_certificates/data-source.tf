data "octopusdeploy_certificates" "example" {
  archived     = false
  ids          = ["Certificates-123", "Certificates-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}
