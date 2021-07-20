package octopusdeploy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSliceFromTerraformTypeList(t *testing.T) {
	var list []interface{}
	slice := getSliceFromTerraformTypeList(list)
	require.Nil(t, slice)

	list = []interface{}{}
	slice = getSliceFromTerraformTypeList(list)
	require.Nil(t, slice)

	randomNumber := 0
	slice = getSliceFromTerraformTypeList(randomNumber)
	require.Nil(t, slice)

	errList := []interface{}{}
	errList = append(errList, nil)
	slice = getSliceFromTerraformTypeList(errList)
	require.Nil(t, slice)
}
