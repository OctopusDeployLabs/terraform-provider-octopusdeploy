package schemas

import (
	"fmt"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/require"
	"reflect"
	"strings"
	"testing"
)

var testableSchemas = []EntitySchema{
	LifecycleSchema{},
	NugetFeedSchema{},
	ProjectSchema{},
	HelmFeedSchema{},
	DockerContainerRegistryFeedSchema{},
	EnvironmentSchema{},
	RunbookSchema{},
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

	checkResourceAttributes(t, schemaName, entitySchema.Attributes)
	checkResourceBlocks(t, schemaName, entitySchema.Blocks)
}

func checkResourceBlocks(t *testing.T, schemaName string, blocks map[string]resourceSchema.Block) {
	for _, block := range blocks {
		switch b := block.(type) {
		case resourceSchema.ListNestedBlock:
			checkResourceAttributes(t, schemaName, b.NestedObject.Attributes)
			checkResourceBlocks(t, schemaName, b.NestedObject.Blocks)
		case resourceSchema.SetNestedBlock:
			checkResourceAttributes(t, schemaName, b.NestedObject.Attributes)
			checkResourceBlocks(t, schemaName, b.NestedObject.Blocks)
		case resourceSchema.SingleNestedBlock:
			checkResourceAttributes(t, schemaName, b.Attributes)
			checkResourceBlocks(t, schemaName, b.Blocks)
		}
	}
}

func checkResourceAttributes(t *testing.T, schemaName string, attributes map[string]resourceSchema.Attribute) {
	for name, attr := range attributes {
		if strings.HasSuffix(reflect.TypeOf(attr).PkgPath(), "datasource/schema") {
			require.Fail(t, fmt.Sprintf("%s in %s must be from the resource schema", name, schemaName))
		}
	}
}

func typeName(i interface{}) string {
	return fmt.Sprintf("%T", i)
}

func datasourceTest(t *testing.T, schema EntitySchema) {
	dataSourceAttributes := schema.GetDatasourceSchemaAttributes()
	schemaName := typeName(schema)
	for name, attr := range dataSourceAttributes {
		if strings.HasSuffix(reflect.TypeOf(attr).PkgPath(), "resource/schema") {
			require.Fail(t, fmt.Sprintf("%s in %s must be from the data source schema", name, schemaName))
		}
	}

}
