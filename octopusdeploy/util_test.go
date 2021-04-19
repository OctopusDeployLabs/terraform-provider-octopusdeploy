package octopusdeploy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSliceFromTerraformTypeList(t *testing.T) {
	var list []interface{}
	slice := getSliceFromTerraformTypeList(list)

	assert.Nil(t, slice)

	list = []interface{}{}
	slice = getSliceFromTerraformTypeList(list)

	assert.Nil(t, slice)
}
