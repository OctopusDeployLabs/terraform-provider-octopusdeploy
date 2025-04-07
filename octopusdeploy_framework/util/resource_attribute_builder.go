package util

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	case *schema.Int32Attribute:
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
	case *schema.Int32Attribute:
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
	case *schema.Int32Attribute:
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
	case *schema.Int32Attribute:
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
	case *schema.Int32Attribute:
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

func (b *AttributeBuilder[T]) Default(defaultValue interface{}) *AttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		if strDefault, ok := defaultValue.(string); ok {
			a.Default = stringdefault.StaticString(strDefault)
		}
	case *schema.BoolAttribute:
		if boolDefault, ok := defaultValue.(bool); ok {
			a.Default = booldefault.StaticBool(boolDefault)
		}
	case *schema.Int64Attribute:
		if intDefault, ok := defaultValue.(int64); ok {
			a.Default = int64default.StaticInt64(intDefault)
		}
	case *schema.Int32Attribute:
		if intDefault, ok := defaultValue.(int32); ok {
			a.Default = int32default.StaticInt32(intDefault)
		}
	case *schema.NumberAttribute:
	case *schema.Float64Attribute:
		if floatDefault, ok := defaultValue.(float64); ok {
			a.Default = float64default.StaticFloat64(floatDefault)
		}
	case *schema.ListAttribute:
		if defaultList, ok := defaultValue.(defaults.List); ok {
			a.Default = defaultList
		}
	case *schema.SetAttribute:
		if defaultSet, ok := defaultValue.(defaults.Set); ok {
			a.Default = defaultSet
		}
	case *schema.MapAttribute:
		if defaultMap, ok := defaultValue.(defaults.Map); ok {
			a.Default = defaultMap
		}
	}
	return b
}

// DefaultEmpty sets the default value of an attribute to an empty collection.
//
// This method applies only to ListAttribute, SetAttribute or MapAttribute types.
func (b *AttributeBuilder[T]) DefaultEmpty() *AttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.ListAttribute:
		a.Default = listdefault.StaticValue(types.ListValueMust(a.ElementType, []attr.Value{}))
	case *schema.SetAttribute:
		a.Default = setdefault.StaticValue(types.SetValueMust(a.ElementType, []attr.Value{}))
	case *schema.MapAttribute:
		a.Default = mapdefault.StaticValue(types.MapValueMust(a.ElementType, map[string]attr.Value{}))
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
	case *schema.Int32Attribute:
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

func (b *AttributeBuilder[T]) PlanModifiers(modifiers ...any) *AttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		if stringModifiers, ok := convertToTypedSlice[planmodifier.String](modifiers); ok {
			a.PlanModifiers = append(a.PlanModifiers, stringModifiers...)
		}
	case *schema.BoolAttribute:
		if boolModifiers, ok := convertToTypedSlice[planmodifier.Bool](modifiers); ok {
			a.PlanModifiers = append(a.PlanModifiers, boolModifiers...)
		}
	case *schema.Int64Attribute:
		if int64Modifiers, ok := convertToTypedSlice[planmodifier.Int64](modifiers); ok {
			a.PlanModifiers = append(a.PlanModifiers, int64Modifiers...)
		}
	case *schema.Int32Attribute:
		if int32Modifiers, ok := convertToTypedSlice[planmodifier.Int32](modifiers); ok {
			a.PlanModifiers = append(a.PlanModifiers, int32Modifiers...)
		}
	case *schema.Float64Attribute:
		if float64Modifiers, ok := convertToTypedSlice[planmodifier.Float64](modifiers); ok {
			a.PlanModifiers = append(a.PlanModifiers, float64Modifiers...)
		}
	case *schema.ListAttribute:
		if listModifiers, ok := convertToTypedSlice[planmodifier.List](modifiers); ok {
			a.PlanModifiers = append(a.PlanModifiers, listModifiers...)
		}
	case *schema.SetAttribute:
		if setModifiers, ok := convertToTypedSlice[planmodifier.Set](modifiers); ok {
			a.PlanModifiers = append(a.PlanModifiers, setModifiers...)
		}
	case *schema.MapAttribute:
		if mapModifiers, ok := convertToTypedSlice[planmodifier.Map](modifiers); ok {
			a.PlanModifiers = append(a.PlanModifiers, mapModifiers...)
		}
	}
	return b
}
func (b *AttributeBuilder[T]) Validators(validators ...any) *AttributeBuilder[T] {
	switch a := any(&b.attr).(type) {
	case *schema.StringAttribute:
		if stringValidators, ok := convertToTypedSlice[validator.String](validators); ok {
			a.Validators = append(a.Validators, stringValidators...)
		}
	case *schema.BoolAttribute:
		if boolValidators, ok := convertToTypedSlice[validator.Bool](validators); ok {
			a.Validators = append(a.Validators, boolValidators...)
		}
	case *schema.Int64Attribute:
		if int64Validators, ok := convertToTypedSlice[validator.Int64](validators); ok {
			a.Validators = append(a.Validators, int64Validators...)
		}
	case *schema.Int32Attribute:
		if int32Validators, ok := convertToTypedSlice[validator.Int32](validators); ok {
			a.Validators = append(a.Validators, int32Validators...)
		}
	case *schema.Float64Attribute:
		if float64Validators, ok := convertToTypedSlice[validator.Float64](validators); ok {
			a.Validators = append(a.Validators, float64Validators...)
		}
	case *schema.ListAttribute:
		if listValidators, ok := convertToTypedSlice[validator.List](validators); ok {
			a.Validators = append(a.Validators, listValidators...)
		}
	case *schema.SetAttribute:
		if setValidators, ok := convertToTypedSlice[validator.Set](validators); ok {
			a.Validators = append(a.Validators, setValidators...)
		}
	case *schema.MapAttribute:
		if mapValidators, ok := convertToTypedSlice[validator.Map](validators); ok {
			a.Validators = append(a.Validators, mapValidators...)
		}
	}
	return b
}

func convertToTypedSlice[T any](slice []any) ([]T, bool) {
	typedSlice := make([]T, 0, len(slice))
	for _, item := range slice {
		if typed, ok := item.(T); ok {
			typedSlice = append(typedSlice, typed)
		} else {
			return nil, false
		}
	}
	return typedSlice, true
}
func ResourceString() *AttributeBuilder[schema.StringAttribute] {
	return NewAttributeBuilder[schema.StringAttribute]()
}

func ResourceBool() *AttributeBuilder[schema.BoolAttribute] {
	return NewAttributeBuilder[schema.BoolAttribute]()
}

func ResourceInt32() *AttributeBuilder[schema.Int32Attribute] {
	return NewAttributeBuilder[schema.Int32Attribute]()
}

func ResourceInt64() *AttributeBuilder[schema.Int64Attribute] {
	return NewAttributeBuilder[schema.Int64Attribute]()
}

func ResourceFloat64() *AttributeBuilder[schema.Float64Attribute] {
	return NewAttributeBuilder[schema.Float64Attribute]()
}

func ResourceNumber() *AttributeBuilder[schema.NumberAttribute] {
	return NewAttributeBuilder[schema.NumberAttribute]()
}

func ResourceList(elementType attr.Type) *AttributeBuilder[schema.ListAttribute] {
	return NewAttributeBuilder[schema.ListAttribute]().ElementType(elementType)
}

func ResourceSet(elementType attr.Type) *AttributeBuilder[schema.SetAttribute] {
	return NewAttributeBuilder[schema.SetAttribute]().ElementType(elementType)
}
func ResourceMap(elementType attr.Type) *AttributeBuilder[schema.MapAttribute] {
	return NewAttributeBuilder[schema.MapAttribute]().ElementType(elementType)
}

func ResourceObject(attributeTypes map[string]attr.Type) *AttributeBuilder[schema.ObjectAttribute] {
	return NewAttributeBuilder[schema.ObjectAttribute]().AttributeTypes(attributeTypes)
}
