. $PSScriptRoot\Octopus.ps1

Wait-ForOctopus

#Creating a connection
$repository = Connect-ToOctopus "http://localhost:8080"

#Creating login object
$LoginObj = New-Object Octopus.Client.Model.LoginCommand
$LoginObj.Username = "admin"
$LoginObj.Password = "Password01!"

#Loging in to Octopus
$repository.Users.SignIn($LoginObj)

#Getting current user logged in
$UserObj = $repository.Users.GetCurrent()

#Creating API Key for user. This automatically gets saved to the database.
$ApiObj = $repository.Users.CreateApiKey($UserObj, "Terraform tests")

#Save the API key so we can use it later
Set-Content -Path tests\octopus_api.txt -Value $ApiObj.ApiKey

Set-Content -Path tests\provider.tf -Value @"
provider "octopusdeploy" {
    address = "http://localhost:8080"
    apikey  = "$($ApiObj.ApiKey)"
}
"@
