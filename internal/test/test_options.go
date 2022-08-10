package test

import "github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

type TestOptions[T any] struct {
	Resource      *T
	LocalName     string
	QualifiedName string
	ResourceName  string
}

func NewTestOptions[T any](resource string) *TestOptions[T] {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	resourceName := "octopusdeploy_" + resource
	qualifiedName := resourceName + "." + localName

	return &TestOptions[T]{
		LocalName:     localName,
		QualifiedName: qualifiedName,
		ResourceName:  resourceName,
	}
}
