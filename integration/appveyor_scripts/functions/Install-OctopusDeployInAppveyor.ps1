function Install-OctopusDeployInAppveyor
{
    [CmdletBinding()]
    Param
    (
        # Octopus Deploy Administrator Username
        [Parameter(Mandatory=$true)]
        [string]
        $OctopusAdministartorUser,

        # Octopus Deploy Administrator Password
        [Parameter(Mandatory=$true)]
        [string]
        $OctopusAdministartorPassword
    )

    # INSTALL PACKAGES
    choco install octopusdeploy -y --no-progress
    choco install octopusdeploy.tentacle -y --no-progress
    choco install octopustools -y --no-progress

    # Enable Firewall Rules
    New-NetFirewallRule -DisplayName 'Allow Access to MSSQL' -Direction Inbound -LocalPort 1433 -Protocol TCP -Action Allow
    New-NetFirewallRule -DisplayName 'HTTP(S) Inbound' -Profile 'Any' -Direction Inbound -Action Allow -Protocol TCP -LocalPort @('80', '443')

    # CONFIGURE OCTO SERVER
    $OctoExe = "C:\Program Files\Octopus Deploy\Octopus\Octopus.Server.exe"

    $installCommands = @(
        "create-instance --instance `"OctopusServer`" --config `"C:\Octopus\OctopusServer.config`""
        "database --instance `"OctopusServer`" --connectionString `"Data Source=localhost,1433\SQL2017;Initial Catalog=Octopus;Integrated Security=True`" --create --grant `"NT AUTHORITY\SYSTEM`""
        "configure --instance `"OctopusServer`" --upgradeCheck `"False`" --upgradeCheckWithStatistics `"False`" --webForceSSL `"False`" --webListenPrefixes `"http://localhost:80/`" --commsListenPort `"10943`" --serverNodeName `"$($env:COMPUTERNAME)`" --usernamePasswordIsEnabled `"True`""
        "service --instance `"OctopusServer`" --stop"
        "admin --instance `"OctopusServer`" --username `"$($OctopusAdministartorUser)`" --email `"admin@octo.com`" --password `"$($OctopusAdministartorPassword)`""
        "license --instance `"OctopusServer`" --licenseBase64 `"PExpY2Vuc2UgU2lnbmF0dXJlPSJDd0R1YUh2L2JveVBiS2tISnRqdjVBdmRWUjFWdG1zdktrSlZJQTJyM3ZhbDQ4d0lObThLbm1pUHlQRG1TYXNTKzl2OTlGUERNNlc0ZE92SjYvd2IzZz09Ij4KICA8TGljZW5zZWRUbz5WYWdyYW50PC9MaWNlbnNlZFRvPgogIDxMaWNlbnNlS2V5PjI2MDgwLTQ1MDc1LTU1NDIyLTI5NDU3PC9MaWNlbnNlS2V5PgogIDxWZXJzaW9uPjIuMDwhLS0gTGljZW5zZSBTY2hlbWEgVmVyc2lvbiAtLT48L1ZlcnNpb24+CiAgPFZhbGlkRnJvbT4yMDE3LTEyLTAzPC9WYWxpZEZyb20+CiAgPFZhbGlkVG8+MjAxOC0wMS0xNzwvVmFsaWRUbz4KICA8UHJvamVjdExpbWl0PlVubGltaXRlZDwvUHJvamVjdExpbWl0PgogIDxNYWNoaW5lTGltaXQ+VW5saW1pdGVkPC9NYWNoaW5lTGltaXQ+CiAgPFVzZXJMaW1pdD5VbmxpbWl0ZWQ8L1VzZXJMaW1pdD4KPC9MaWNlbnNlPg==`""
        "service --instance `"OctopusServer`" --install --reconfigure --start --dependOn `"MSSQL`$SQL2017`""
    )

    foreach ($command in $installCommands) {
        Write-Output "Running $($OctoExe) $($command)"
        Start-ProcessAdvanced -FilePath $OctoExe -ArgumentList $command
    }


}
