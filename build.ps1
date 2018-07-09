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
	$buildOutputPath = "$($binary.goos)-$($binary.goarch)"

	if ($binary.goos -eq "windows") {
		$fileExtension = ".exe"
	}

	New-Item -Path $buildOutputPath -Force

	Start-ProcessAdvanced -FilePath 'go' -ArgumentList "build -o $($buildOutputPath)/terraform-provider-octopusdeploy$($fileExtension)" -EnvironmentKeyValues @{ GOOS = $binary.goos; GOARCH = $binary.goarch; } -Verbose
}
