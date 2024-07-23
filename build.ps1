# Set variables
$VERSION = "0.7.514"
$BINARY = "terraform-provider-octopusdeploy.exe"
$HOSTNAME = "octopus.com"
$NAMESPACE = "com"
$NAME = "octopusdeploy"
$OS_ARCH = "windows_amd64"

# Build the provider
go build -o $BINARY

# Create the plugin directory if it doesn't exist
$pluginDir = "$env:APPDATA\terraform.d\plugins\$HOSTNAME\$NAMESPACE\$NAME\$VERSION\$OS_ARCH"
New-Item -ItemType Directory -Force -Path $pluginDir

# Move the binary to the plugin directory
Move-Item -Force $BINARY $pluginDir

Write-Host "Provider installed successfully to $pluginDir"