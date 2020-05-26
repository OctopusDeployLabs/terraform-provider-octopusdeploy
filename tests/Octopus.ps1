function Connect-ToOctopus() {
    param (
        [string]$url,
        [string]$username = "admin",
        [string]$password = "Password01!"
    )
    $endpoint = New-Object Octopus.Client.OctopusServerEndpoint $url
    $repository = New-Object Octopus.Client.OctopusRepository $endpoint
    $LoginObj = New-Object Octopus.Client.Model.LoginCommand
    $LoginObj.Username = $username
    $LoginObj.Password = $password
    Invoke-ScriptBlockWithRetries {$repository.Users.SignIn($LoginObj)} | Out-Null
    return $repository
}

function Invoke-CommandWithRetries
{
    [CmdletBinding()]
    param (
        [Parameter(Mandatory = $True)]
        [string]$Command,
        [Array]$Arguments,
        [bool]$TrustExitCode = $True,
        [int]$RetrySleepSeconds = 10,
        [int]$MaxAttempts = 10,
        [bool]$PrintCommand = $True,
        [bool]$PrintOutput = $False,
        [int[]]$AllowedReturnValues = @(),
        [string]$ErrorMessage
    )

    Process
    {
        $attempt = 0
        while ($true)
        {
            Write-Host $(if ($PrintCommand) { "Executing: $Command $Arguments" }
            else { "Executing command..." })

            try
            {
                $output = & $Command $Arguments 2>&1
                if ($PrintOutput) { Write-Host $output }

                $stderr = $output | where { $_ -is [System.Management.Automation.ErrorRecord] }
                if (($LASTEXITCODE -eq 0 -or $AllowedReturnValues -contains $LASTEXITCODE) -and ($TrustExitCode -or !($stderr)))
                {
                    Write-Host "Command executed successfully"
                    return $output
                }

                Write-Host "Command failed with exit code ($LASTEXITCODE) and stderr: $stderr" -ForegroundColor Yellow
            }
            catch
            {
                Write-Host "Command failed with exit code ($LASTEXITCODE), exception ($_) and stderr: $stderr" -ForegroundColor Yellow
            }

            if ($attempt -eq $MaxAttempts)
            {
                $ex = new-object System.Management.Automation.CmdletInvocationException "All retry attempts exhausted $ErrorMessage"
                $category = [System.Management.Automation.ErrorCategory]::LimitsExceeded
                $errRecord = new-object System.Management.Automation.ErrorRecord $ex, "CommandFailed", $category, $Command
                $psCmdlet.WriteError($errRecord)
                return $output
            }

            $attempt++;
            Write-Host "Retrying test execution [#$attempt/$MaxAttempts] in $RetrySleepSeconds seconds..."
            Start-Sleep -s $RetrySleepSeconds
        }
    }
}

function Invoke-ScriptBlockWithRetries {
    [CmdletBinding()]
    param (
        [parameter(Mandatory, ValueFromPipeline)]
        [ValidateNotNullOrEmpty()]
        [scriptblock] $ScriptBlock,
        [int] $RetryCount = 3,
        [int] $TimeoutInSecs = 30,
        [string] $SuccessMessage = "",
        [string] $FailureMessage = ""
    )

    process {
        $Attempt = 1

        do {
            try {
                $PreviousPreference = $ErrorActionPreference
                $ErrorActionPreference = 'Stop'
                Invoke-Command -ScriptBlock $ScriptBlock -OutVariable Result | Out-Null
                $ErrorActionPreference = $PreviousPreference

                # flow control will execute the next line only if the command in the scriptblock executed without any errors
                # if an error is thrown, flow control will go to the 'catch' block
                if (-not [string]::IsNullOrEmpty($SuccessMessage)) {
                    Write-Host "$SuccessMessage `n"
                }
                return $result
            }
            catch {
                if ($Attempt -gt $RetryCount) {
                    if (-not [string]::IsNullOrEmpty($FailureMessage)) {
                        Write-Host "$FailureMessage! Error was $($_.exception.Message). Total retry attempts: $RetryCount"
                    }
                    throw $_.exception
                }
                else {
                    if (-not [string]::IsNullOrEmpty($FailureMessage)) {
                        Write-Host "[$Attempt/$RetryCount] $FailureMessage. Error was $($_.exception.Message). Retrying in $TimeoutInSecs seconds..."
                    }
                    Start-Sleep -Seconds $TimeoutInSecs
                    $Attempt = $Attempt + 1
                }
            }
        }
        While ($true)
    }
}
