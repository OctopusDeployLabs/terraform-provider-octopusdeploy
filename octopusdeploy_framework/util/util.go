package util

import "fmt"

func GetProviderName() string {
	return "octopusdeploy"
}

func GetTypeName(name string) string {
	return fmt.Sprintf("%s_%s", GetProviderName(), name)
}
