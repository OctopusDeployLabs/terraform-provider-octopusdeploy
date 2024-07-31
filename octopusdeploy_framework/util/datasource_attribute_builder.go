package util

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type DataSourceAttributeBuilder[T any] struct {
	attr T
}

func NewDataSourceAttributeBuilder[T any]() *DataSourceAttributeBuilder[T] {
	return &DataSourceAttributeBuilder[T]{}
}

func (b *DataSourceAttributeBuilder[T]) Optional() *DataSourceAttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.Optional = true
	case *schema.BoolAttribute:
		a.Optional = true
	case *schema.Int64Attribute:
		a.Optional = true
	case *schema.Float64Attribute:
		a.Optional = true
	case *schema.ListAttribute:
		a.Optional = true
	case *schema.SetAttribute:
		a.Optional = true
	case *schema.MapAttribute:
		a.Optional = true
	}
	return b
}

func (b *DataSourceAttributeBuilder[T]) Required() *DataSourceAttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.Required = true
	case *schema.BoolAttribute:
		a.Required = true
	case *schema.Int64Attribute:
		a.Required = true
	case *schema.Float64Attribute:
		a.Required = true
	case *schema.ListAttribute:
		a.Required = true
	case *schema.SetAttribute:
		a.Required = true
	case *schema.MapAttribute:
		a.Required = true
	}
	return b
}

func (b *DataSourceAttributeBuilder[T]) Computed() *DataSourceAttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.Computed = true
	case *schema.BoolAttribute:
		a.Computed = true
	case *schema.Int64Attribute:
		a.Computed = true
	case *schema.Float64Attribute:
		a.Computed = true
	case *schema.ListAttribute:
		a.Computed = true
	case *schema.SetAttribute:
		a.Computed = true
	case *schema.MapAttribute:
		a.Computed = true
	}
	return b
}

func (b *DataSourceAttributeBuilder[T]) Description(desc string) *DataSourceAttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.Description = desc
	case *schema.BoolAttribute:
		a.Description = desc
	case *schema.Int64Attribute:
		a.Description = desc
	case *schema.Float64Attribute:
		a.Description = desc
	case *schema.ListAttribute:
		a.Description = desc
	case *schema.SetAttribute:
		a.Description = desc
	case *schema.MapAttribute:
		a.Description = desc
	}
	return b
}

func (b *DataSourceAttributeBuilder[T]) Sensitive() *DataSourceAttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.Sensitive = true
	case *schema.BoolAttribute:
		a.Sensitive = true
	case *schema.Int64Attribute:
		a.Sensitive = true
	case *schema.Float64Attribute:
		a.Sensitive = true
	case *schema.ListAttribute:
		a.Sensitive = true
	case *schema.SetAttribute:
		a.Sensitive = true
	case *schema.MapAttribute:
		a.Sensitive = true
	}
	return b
}

func (b *DataSourceAttributeBuilder[T]) ElementType(elementType attr.Type) *DataSourceAttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.ListAttribute:
		a.ElementType = elementType
	case *schema.SetAttribute:
		a.ElementType = elementType
	case *schema.MapAttribute:
		a.ElementType = elementType
	}
	return b
}

func (b *DataSourceAttributeBuilder[T]) Build() T {
	return b.attr
}

func DataSourceString() *DataSourceAttributeBuilder[schema.StringAttribute] {
	return NewDataSourceAttributeBuilder[schema.StringAttribute]()
}

func DataSourceBool() *DataSourceAttributeBuilder[schema.BoolAttribute] {
	return NewDataSourceAttributeBuilder[schema.BoolAttribute]()
}

func DataSourceInt64() *DataSourceAttributeBuilder[schema.Int64Attribute] {
	return NewDataSourceAttributeBuilder[schema.Int64Attribute]()
}

func DataSourceFloat64() *DataSourceAttributeBuilder[schema.Float64Attribute] {
	return NewDataSourceAttributeBuilder[schema.Float64Attribute]()
}

func DataSourceList(elementType attr.Type) *DataSourceAttributeBuilder[schema.ListAttribute] {
	return NewDataSourceAttributeBuilder[schema.ListAttribute]().ElementType(elementType)
}

func DataSourceSet(elementType attr.Type) *DataSourceAttributeBuilder[schema.SetAttribute] {
	return NewDataSourceAttributeBuilder[schema.SetAttribute]().ElementType(elementType)
}

func DataSourceMap(elementType attr.Type) *DataSourceAttributeBuilder[schema.MapAttribute] {
	return NewDataSourceAttributeBuilder[schema.MapAttribute]().ElementType(elementType)
}
