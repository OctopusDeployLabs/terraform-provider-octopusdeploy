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

echo "TF_ACC=true" | Out-File -FilePath $env:GITHUB_ENV -Encoding utf8 -Append
echo "OCTOPUS_APIKEY=$($ApiObj.ApiKey)" | Out-File -FilePath $env:GITHUB_ENV -Encoding utf8 -Append
