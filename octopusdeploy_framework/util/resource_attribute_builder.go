package util

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type AttributeBuilder[T any] struct {
	attr T
}

func NewAttributeBuilder[T any]() *AttributeBuilder[T] {
	return &AttributeBuilder[T]{}
}

func (b *AttributeBuilder[T]) Optional() *AttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.Optional = true
	case *schema.BoolAttribute:
		a.Optional = true
	case *schema.Int64Attribute:
		a.Optional = true
	case *schema.Float64Attribute:
		a.Optional = true
	case *schema.NumberAttribute:
		a.Optional = true
	case *schema.ListAttribute:
		a.Optional = true
	case *schema.SetAttribute:
		a.Optional = true
	case *schema.MapAttribute:
		a.Optional = true
	case *schema.ObjectAttribute:
		a.Optional = true
	}
	return b
}

func (b *AttributeBuilder[T]) Deprecated(deprecationMessage string) *AttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.DeprecationMessage = deprecationMessage
	case *schema.BoolAttribute:
		a.DeprecationMessage = deprecationMessage
	case *schema.Int64Attribute:
		a.DeprecationMessage = deprecationMessage
	case *schema.Float64Attribute:
		a.DeprecationMessage = deprecationMessage
	case *schema.NumberAttribute:
		a.DeprecationMessage = deprecationMessage
	case *schema.ListAttribute:
		a.DeprecationMessage = deprecationMessage
	case *schema.SetAttribute:
		a.DeprecationMessage = deprecationMessage
	case *schema.MapAttribute:
		a.DeprecationMessage = deprecationMessage
	case *schema.ObjectAttribute:
		a.DeprecationMessage = deprecationMessage
	}
	return b
}

func (b *AttributeBuilder[T]) Computed() *AttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.Computed = true
	case *schema.BoolAttribute:
		a.Computed = true
	case *schema.Int64Attribute:
		a.Computed = true
	case *schema.Float64Attribute:
		a.Computed = true
	case *schema.NumberAttribute:
		a.Computed = true
	case *schema.ListAttribute:
		a.Computed = true
	case *schema.SetAttribute:
		a.Computed = true
	case *schema.MapAttribute:
		a.Computed = true
	case *schema.ObjectAttribute:
		a.Computed = true
	}
	return b
}

func (b *AttributeBuilder[T]) Required() *AttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.Required = true
	case *schema.BoolAttribute:
		a.Required = true
	case *schema.Int64Attribute:
		a.Required = true
	case *schema.Float64Attribute:
		a.Required = true
	case *schema.NumberAttribute:
		a.Required = true
	case *schema.ListAttribute:
		a.Required = true
	case *schema.SetAttribute:
		a.Required = true
	case *schema.MapAttribute:
		a.Required = true
	case *schema.ObjectAttribute:
		a.Required = true
	}
	return b
}

func (b *AttributeBuilder[T]) Description(desc string) *AttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.Description = desc
	case *schema.BoolAttribute:
		a.Description = desc
	case *schema.Int64Attribute:
		a.Description = desc
	case *schema.Float64Attribute:
		a.Description = desc
	case *schema.NumberAttribute:
		a.Description = desc
	case *schema.ListAttribute:
		a.Description = desc
	case *schema.SetAttribute:
		a.Description = desc
	case *schema.MapAttribute:
		a.Description = desc
	case *schema.ObjectAttribute:
		a.Description = desc
	}
	return b
}

func (b *AttributeBuilder[T]) Sensitive() *AttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		a.Sensitive = true
	case *schema.BoolAttribute:
		a.Sensitive = true
	case *schema.Int64Attribute:
		a.Sensitive = true
	case *schema.Float64Attribute:
		a.Sensitive = true
	case *schema.NumberAttribute:
		a.Sensitive = true
	case *schema.ListAttribute:
		a.Sensitive = true
	case *schema.SetAttribute:
		a.Sensitive = true
	case *schema.MapAttribute:
		a.Sensitive = true
	case *schema.ObjectAttribute:
		a.Sensitive = true
	}
	return b
}

func (b *AttributeBuilder[T]) ElementType(elementType attr.Type) *AttributeBuilder[T] {
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

func (b *AttributeBuilder[T]) AttributeTypes(attributeTypes map[string]attr.Type) *AttributeBuilder[T] {
	if a, ok := any(&b.attr).(*schema.ObjectAttribute); ok {
		a.AttributeTypes = attributeTypes
	}
	return b
}

func (b *AttributeBuilder[T]) Build() T {
	return b.attr
}

func String() *AttributeBuilder[schema.StringAttribute] {
	return NewAttributeBuilder[schema.StringAttribute]()
}

func Bool() *AttributeBuilder[schema.BoolAttribute] {
	return NewAttributeBuilder[schema.BoolAttribute]()
}

func Int64() *AttributeBuilder[schema.Int64Attribute] {
	return NewAttributeBuilder[schema.Int64Attribute]()
}

func Float64() *AttributeBuilder[schema.Float64Attribute] {
	return NewAttributeBuilder[schema.Float64Attribute]()
}

func Number() *AttributeBuilder[schema.NumberAttribute] {
	return NewAttributeBuilder[schema.NumberAttribute]()
}

func List(elementType attr.Type) *AttributeBuilder[schema.ListAttribute] {
	return NewAttributeBuilder[schema.ListAttribute]().ElementType(elementType)
}

func Set(elementType attr.Type) *AttributeBuilder[schema.SetAttribute] {
	return NewAttributeBuilder[schema.SetAttribute]().ElementType(elementType)
}
