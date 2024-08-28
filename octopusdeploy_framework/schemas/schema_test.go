package schemas

import (
	"fmt"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	GitCredentialSchema{},
	ArtifactoryGenericFeedSchema{},
	AwsElasticContainerRegistrySchema{},
	FeedsSchema{},
	GitHubRepositoryFeedSchema{},
	SpaceSchema{},
	ScriptModuleSchema{},
	LibraryVariableSetSchema{},
	ActionTemplateParameterSchema{},
	MavenFeedSchema{},
	ProjectGroupSchema{},
	TagSchema{},
	UsernamePasswordAccountSchema{},
	VariableSchema{},
}

func TestDatasourceSchemaDefinitionIsUsingCorrectTypes(t *testing.T) {
	for _, schema := range testableSchemas {
		datasourceTest(t, schema)
	}
}

func TestResourceSchemaDefinitionIsUsingCorrectTypes(t *testing.T) {
	for _, schema := range testableSchemas {
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
		if !strings.HasSuffix(reflect.TypeOf(attr).PkgPath(), "resource/schema") {
			require.Fail(t, fmt.Sprintf("%s in %s must be from the resource schema", name, schemaName))
		}
	}
}

func typeName(i interface{}) string {
	return fmt.Sprintf("%T", i)
}

func datasourceTest(t *testing.T, schema EntitySchema) {
	dataSourceSchema := schema.GetDatasourceSchema()
	schemaName := typeName(schema)

	checkDatasourceAttributes(t, schemaName, dataSourceSchema.Attributes)
	checkDatasourceBlocks(t, schemaName, dataSourceSchema.Blocks)
}

func checkDatasourceAttributes(t *testing.T, schemaName string, attributes map[string]datasourceSchema.Attribute) {
	for name, attr := range attributes {
		switch attrType := attr.(type) {
		case datasourceSchema.ListNestedAttribute:
			checkDatasourceAttributes(t, schemaName, attrType.NestedObject.Attributes)
		default:
			if !strings.HasSuffix(reflect.TypeOf(attr).PkgPath(), "datasource/schema") {
				require.Fail(t, fmt.Sprintf("%s in %s must be from the data source schema", name, schemaName))
			}
		}
	}
}

func checkDatasourceBlocks(t *testing.T, schemaName string, blocks map[string]datasourceSchema.Block) {
	for _, block := range blocks {
		switch b := block.(type) {
		case datasourceSchema.ListNestedBlock:
			checkDatasourceAttributes(t, schemaName, b.NestedObject.Attributes)
			checkDatasourceBlocks(t, schemaName, b.NestedObject.Blocks)
		case datasourceSchema.SetNestedBlock:
			checkDatasourceAttributes(t, schemaName, b.NestedObject.Attributes)
			checkDatasourceBlocks(t, schemaName, b.NestedObject.Blocks)
		case datasourceSchema.SingleNestedBlock:
			checkDatasourceAttributes(t, schemaName, b.Attributes)
			checkDatasourceBlocks(t, schemaName, b.Blocks)
		}
	}
}
