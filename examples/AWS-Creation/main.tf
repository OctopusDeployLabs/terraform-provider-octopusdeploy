provider "octopusdeploy" {
    address = var.serverURL
    apikey  = var.apiKey
    space   = var.space
}

resource "octopusdeploy_aws_account" "aw" {
    name = "awsaccount"
    access_key = var.accessKey
    secret_key = var.secretKey
}