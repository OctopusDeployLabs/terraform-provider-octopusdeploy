package octopusdeploy

import (
	"fmt"
)

func createInvalidParameterError(methodName string, parameterName string) error {
	return fmt.Errorf("%s: invalid input parameter, %s", methodName, parameterName)
}

func createResourceOperationError(text string, resourceID string, err error) error {
	return fmt.Errorf(text, resourceID, err)
}

func nameIsNil(methodName string) error {
	return fmt.Errorf("%s: name is nil", methodName)
}
