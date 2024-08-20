package schemas

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"strings"
	"testing"
)

var testableSchemas = []EntitySchema{
	LifecycleSchema{},
	NugetFeedSchema{},
}

func TestSchemaDefinitionIsUsingCorrectTypes(t *testing.T) {
	for _, schema := range testableSchemas {
		datasourceTest(t, schema)
		resourceTest(t, schema)
	}
}

func resourceTest(t *testing.T, schema EntitySchema) {
	entitySchema := schema.GetResourceSchema()
	schemaName := typeName(schema)
	for name, attr := range entitySchema.Attributes {
		if strings.HasSuffix(reflect.TypeOf(attr).PkgPath(), "datasource/schema") {
			require.Fail(t, fmt.Sprintf("%s in %s must be from the resource schema", name, schemaName))
		}
	}
}

func typeName(i interface{}) string {
	return fmt.Sprintf("%T", i)
}

func datasourceTest(t *testing.T, schema EntitySchema) {
	entitySchema := schema.GetDatasource()
	schemaName := typeName(schema)
	for name, attr := range entitySchema.Attributes {
		if strings.HasSuffix(reflect.TypeOf(attr).PkgPath(), "resource/schema") {
			require.Fail(t, fmt.Sprintf("%s in %s must be from the data source schema", name, schemaName))
		}
	}
}
