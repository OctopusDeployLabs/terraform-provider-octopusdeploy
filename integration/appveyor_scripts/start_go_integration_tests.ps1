$ErrorActionPreference = 'Stop'

# Load Functions
$functionFolder = Get-ChildItem -Path (Join-Path -Path $PWD -ChildPath 'integration\appveyor_scripts\functions')
foreach ($function in $functionFolder) { . $function.FullName }

$OCTOPUS_URL = "http://localhost"
$OCTOPUS_APIKEY = Get-OctopusDeployApiKey -OctopusUrl $OCTOPUS_URL -Username $env:TEST_OCTOPUS_USERNAME -Password $env:TEST_OCTOPUS_PASSWORD

Start-ProcessAdvanced -FilePath 'go' -ArgumentList "test -v -timeout 30s ./octopusdeploy/..." -EnvironmentKeyValues @{ OCTOPUS_URL = $OCTOPUS_URL; OCTOPUS_APIKEY = $OCTOPUS_APIKEY; TF_ACC = 1 } -Verbose
