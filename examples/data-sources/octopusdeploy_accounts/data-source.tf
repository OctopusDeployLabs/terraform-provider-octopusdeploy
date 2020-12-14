data "octopusdeploy_accounts" "example" {
  account_type = "UsernamePassword"
  ids          = ["Accounts-123", "Accounts-321"]
  partial_name = "Defau"
  skip         = 5
  take         = 100
}