$ErrorActionPreference = 'Stop'

# Load Functions
$functionFolder = Get-ChildItem -Path (Join-Path -Path $PWD -ChildPath 'integration\appveyor_scripts\functions')
foreach ($function in $functionFolder) { . $function.FullName }

Start-ProcessAdvanced -FilePath 'go' -ArgumentList "test -v -timeout 120s ./octopusdeploy/..." -EnvironmentKeyValues @{ TF_ACC = 1 } -Verbose
