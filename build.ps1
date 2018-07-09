Param
(
    # Param1 help description
    [Parameter(Mandatory=$true)]
    [string]
    $BuildVersion
)


. ".\integration\appveyor_scripts\functions\Start-ProcessAdvanced.ps1"

$binaries = @(
	@{
		goos = "darwin"
		goarch = "386"
	},
	@{
		goos = "darwin"
		goarch = "amd64"
	},
	@{
		goos = "linux"
		goarch = "386"
	},
	@{
		goos = "linux"
		goarch = "amd64"
	},
	@{
		goos = "windows"
		goarch = "386"
	},
	@{
		goos = "windows"
		goarch = "amd64"
	}
)


foreach ($binary in $binaries){
    $buildName = "$($binary.goos)-$($binary.goarch)"
	$buildOutputPath = "build\$($buildName)"

	if ($binary.goos -eq "windows") {
		$fileExtension = ".exe"
	}

	New-Item -Path $buildOutputPath -ItemType Directory -Force

	Start-ProcessAdvanced -FilePath 'go' -ArgumentList "build -o $($buildOutputPath)/terraform-provider-octopusdeploy$($fileExtension)" -EnvironmentKeyValues @{ GOOS = $binary.goos; GOARCH = $binary.goarch; GOPATH = 'C:\gopath' } -Verbose

    Push-Location $buildOutputPath

    Compress-Archive -Path . -DestinationPath "terraform-provider-octopusdeploy-$($buildName)-$($BuildVersion).zip"

    Move-Item "terraform-provider-octopusdeploy-$($buildName)-$($BuildVersion).zip" ..

    Pop-Location
}
