function Get-OctopusDeployApiKey
{
    [CmdletBinding()]
    Param
    (
        # Octopus Deploy Server Url
        [Parameter(Mandatory=$true)]
        [string]
        $OctopusUrl,

        # Octopus Deploy Username to get API Key
        [Parameter(Mandatory=$true)]
        [string]
        $Username,

        # Octopus Deploy Password to get API Key
        [Parameter(Mandatory=$true)]
        [string]
        $Password,

        # The purpose of the API Key to store as a description on the server
        [Parameter(Mandatory=$false)]
        [string]
        $Purpose="Powershell"
    )

    $ErrorActionPreference = 'Stop'

    #Adding libraries. Make sure to modify these paths acording to your environment setup.
    Add-Type -Path "C:\Program Files\Octopus Deploy\Octopus\Newtonsoft.Json.dll"
    Add-Type -Path "C:\Program Files\Octopus Deploy\Octopus\Octopus.Client.dll"

    Write-Verbose "Attempting to connect to Octopus Server"

    #Creating a connection
    $endpoint = new-object Octopus.Client.OctopusServerEndpoint $OctopusUrl
    $repository = new-object Octopus.Client.OctopusRepository $endpoint

    #Creating login object
    $LoginObj = New-Object Octopus.Client.Model.LoginCommand
    $LoginObj.Username = $Username
    $LoginObj.Password = $Password

    #Loging in to Octopus
    $repository.Users.SignIn($LoginObj)

    #Getting current user logged in
    $UserObj = $repository.Users.GetCurrent()

    #Creating API Key for user. This automatically gets saved to the database.
    $ApiObj = $repository.Users.CreateApiKey($UserObj, $Purpose)

    return $ApiObj.ApiKey
}
