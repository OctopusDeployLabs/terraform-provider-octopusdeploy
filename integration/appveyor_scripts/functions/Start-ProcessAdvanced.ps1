function Start-ProcessAdvanced
{
    [CmdletBinding()]
    Param
    (
        # Directory to start Docker Compose In
        [Parameter(Mandatory=$true)]
        [string]
        $FilePath,

        # Directory to start Docker Compose In
        [Parameter(Mandatory=$false)]
        [string]
        $WorkingDirectory=$PWD.Path,

        # Arguments to start Docker Compose With
        [Parameter(Mandatory=$false)]
        [string]
        $ArgumentList,

        # Environment Variables to pass to the process on startup
        [Parameter(Mandatory=$false)]
        [hashtable]
        $EnvironmentKeyValues,

        # Time to sleep between loops waiting for the process to finish
        [Parameter(Mandatory=$false)]
        [int]
        $SleepTime=1,

        # Enable a message which says how long you have been waiting for the process (in Minutes)
        [Parameter(Mandatory=$false)]
        [bool]
        $EnableWaitMessage=$false
    )

    $ErrorActionPreference = 'Stop'

    # Setup stdin\stdout redirection
    $StartInfo = New-Object System.Diagnostics.ProcessStartInfo -Property @{
                    FileName = $FilePath
                    Arguments = $ArgumentList
                    UseShellExecute = $false
                    RedirectStandardOutput = $true
                    RedirectStandardError = $true
                    WorkingDirectory = $WorkingDirectory
                }

    # Create new process
    $Process = New-Object System.Diagnostics.Process

    # Assign previously created StartInfo properties
    $Process.StartInfo = $StartInfo

    if ($EnvironmentKeyValues) {
        foreach ($envVariable in $EnvironmentKeyValues.GetEnumerator()) {
            Write-Verbose "Adding Environment Variable with Name: $($envVariable.Name)"
            $Process.StartInfo.Environment.Add($envVariable.Name, $envVariable.Value)
        }
    }

    # Register Object Events for stdin\stdout reading
    $OutEvent = Register-ObjectEvent -InputObject $Process -EventName OutputDataReceived -Action {
        Write-Host $Event.SourceEventArgs.Data
    }
    $ErrEvent = Register-ObjectEvent -InputObject $Process -EventName ErrorDataReceived -Action {
        Write-Host $Event.SourceEventArgs.Data
    }

    # Start process
    [void]$Process.Start()

    # Begin reading stdin\stdout
    $Process.BeginOutputReadLine()
    $Process.BeginErrorReadLine()

    if ($EnableWaitMessage) {
        # Start a timer to provide some feedback
        $stopwatch =  [system.diagnostics.stopwatch]::StartNew()
    }

    # Do something else while events are firing
    do
    {
        Start-Sleep -Seconds $SleepTime
        if ($EnableWaitMessage) {
            Write-Output "Currently waited $($stopwatch.Elapsed.Minutes) minutes for $($FilePath) to finish"
        }
    }
    while (!$Process.HasExited)

    # Unregister events
    $OutEvent.Name, $ErrEvent.Name |
		ForEach-Object {Unregister-Event -SourceIdentifier $_}

    [int]$exitCode = $Process.ExitCode
	Write-Verbose "Last Exit Code: $exitCode"

	if ($exitCode -gt 0){
        Write-Error "Exiting using the exit code: $($exitCode)"
		exit($exitCode)
	}
}
