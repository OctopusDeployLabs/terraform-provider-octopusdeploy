package octopusdeploy

import (
	"fmt"
)

func createInvalidParameterError(methodName string, parameterName string) error {
	return fmt.Errorf("%s: invalid input parameter, %s", methodName, parameterName)
}

func nameIsNil(methodName string) error {
	return fmt.Errorf("%s: Name is nil", methodName)
}
