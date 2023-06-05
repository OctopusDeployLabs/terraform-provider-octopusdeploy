data "octopusdeploy_accounts" "example" {
  partial_name = "AWS Account"
  skip         = 0
  take         = 1
}

output "data_lookup" {
  value = data.octopusdeploy_accounts.example.accounts[0].id
}