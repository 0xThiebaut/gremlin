package flatten

import (
	"encoding"
	"fmt"
	"reflect"
)

func Marshal(v any) (props []Property, err error) {
	return MarshalFilter(v, func(tag reflect.StructTag) bool {
		return true
	})
}

func MarshalFilter(v any, filter Filter) (props []Property, err error) {
	value := reflect.ValueOf(v)
	for k := value.Kind(); k == reflect.Pointer || k == reflect.Interface; k = value.Kind() {
		value = value.Elem()
	}
	return marshalFields(value, filter)
}

type Filter func(tag reflect.StructTag) bool

func marshalFields(structure reflect.Value, filter Filter) (props []Property, err error) {
	// require structures
	if k, t := structure.Kind(), structure.Type(); k != reflect.Struct {
		return props, fmt.Errorf("invalid kind %s for type %s", k.String(), t.String())
	}
	// loop fields
	for field, value := range structure.Fields() {
		// ignore internal or filtered fields
		if !field.IsExported() || !filter(field.Tag) {
			continue
		}
		// propagate anonymous fields
		if field.Anonymous {
			sub, err := marshalFields(value, filter)
			if err != nil {
				return props, err
			}
			props = append(props, sub...)
			continue
		}
		// prioritize text marshalers
	MARSHAL:
		if marshaler, ok := reflect.TypeAssert[encoding.TextMarshaler](value); ok {
			b, err := marshaler.MarshalText()
			if err != nil {
				return props, err
			}

			props = append(props, Property{
				Label: field.Name,
				Value: string(b),
			})
			continue
		}
		switch k, t := value.Kind(), value.Type(); k {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Bool:
			props = append(props, Property{
				Label: field.Name,
				Value: value.Interface(),
			})
		case reflect.Struct:
			sub, err := marshalFields(value, filter)
			if err != nil {
				return props, err
			}

			for _, prop := range sub {
				props = append(props, Property{
					Label: field.Name + "." + prop.Label,
					Value: prop.Value,
				})
			}
		case reflect.Interface, reflect.Pointer:
			value = value.Elem()
			goto MARSHAL
		default:
			return props, fmt.Errorf("unsupported kind %s for type %s", k.String(), t.String())
		}
	}
	return props, err
}
