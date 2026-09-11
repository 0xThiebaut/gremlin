package flatten

import (
	"encoding"
	"fmt"
	"math/big"
	"reflect"
	"strings"
)

func Unmarshal(props []Property, v any) error {
	value := reflect.ValueOf(v)
	for k := value.Kind(); k == reflect.Pointer || k == reflect.Interface; k = value.Kind() {
		value = value.Elem()
	}
	return unmarshalFields(props, value)
}

func unmarshalFields(props []Property, structure reflect.Value) error {
	// require structures
	if k, t := structure.Kind(), structure.Type(); k != reflect.Struct {
		return fmt.Errorf("invalid kind %s for type %s", k.String(), t.String())
	}
	fields := make(map[string]reflect.Value)
	// loop fields
	for field, value := range structure.Fields() {
		// ignore internal or filtered fields
		if !field.IsExported() {
			continue
		}
		// propagate anonymous fields
		if field.Anonymous {
			if err := unmarshalFields(props, value); err != nil {
				return err
			}
			continue
		}
		// prioritize text marshalers
	UNMARSHAL:
		if _, ok := reflect.TypeAssert[encoding.TextUnmarshaler](value); ok {
			fields[field.Name] = value
			continue
		}
		switch k, t := value.Kind(), value.Type(); k {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Bool:
			fields[field.Name] = value
		case reflect.Struct:
			var sub []Property

			prefix := field.Name + "."
			for _, property := range props {
				if label, ok := strings.CutPrefix(property.Label, prefix); ok {
					sub = append(sub, Property{
						Label: label,
						Value: property.Value,
					})
				}
			}

			if err := unmarshalFields(sub, value); err != nil {
				return err
			}
		case reflect.Interface, reflect.Pointer:
			value = value.Elem()
			goto UNMARSHAL
		default:
			return fmt.Errorf("unsupported kind %s for type %s", k.String(), t.String())
		}
	}

	// loop properties
	for _, property := range props {
		// ignore unknown fields
		field, ok := fields[property.Label]
		if !ok {
			continue
		}

		// error if the field cannot be set
		if !field.CanSet() {
			return fmt.Errorf("cannot set %s", property.Label)
		}

		// prioritize TextUnmarshalers
		value := reflect.ValueOf(property.Value)
		if value.Kind() == reflect.String {
			if unmarshaler, ok := reflect.TypeAssert[encoding.TextUnmarshaler](field); ok {
				if err := unmarshaler.UnmarshalText([]byte(value.String())); err != nil {
					return err
				}
				continue
			}
		}

		// Map similar or compatible kinds
		if field.Kind() == value.Kind() {
			field.Set(value)
			continue
		} else if field.CanInt() && value.CanInt() {
			field.SetInt(value.Int())
			continue
		} else if field.CanUint() && value.CanUint() {
			field.SetUint(value.Uint())
			continue
		}

		// Handle big-ints
		if b, ok := reflect.TypeAssert[*big.Int](value); ok && b.IsUint64() && field.CanUint() {
			field.SetUint(b.Uint64())
			continue
		}

		return fmt.Errorf("unsupported kind %s for type %T", field.Kind().String(), property.Value)
	}
	return nil
}
