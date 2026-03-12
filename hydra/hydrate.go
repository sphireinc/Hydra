package hydra

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
)

func (h *Hydratable) hydrate(ctx context.Context, db any, whereClauses map[string]interface{}) error {
	if !h.isInitialized {
		return ErrNotInitialized
	}

	// Reflect on the actual value of `self` (which holds the parent struct)
	v := reflect.ValueOf(h.self)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("hydrated object must be a struct, got %T", h.self)
	}
	if !v.CanAddr() {
		return fmt.Errorf("hydrated object is not addressable")
	}

	tableName, err := h.resolveTableName()
	if err != nil {
		return err
	}

	meta, err := getTypeMetadata(h.self)
	if err != nil {
		return err
	}

	data, err := h.FetchContext(ctx, db, tableName, meta.Columns, whereClauses)
	if err != nil {
		return err
	}

	for _, field := range meta.Fields {
		value, ok := data[field.Column]
		if !ok {
			continue
		}

		fieldValue := v.Field(field.Index)
		if !fieldValue.CanSet() {
			return fmt.Errorf("field %s is not settable", field.Name)
		}

		if err := assignFieldValue(h.self, field, fieldValue, value); err != nil {
			return fmt.Errorf("field %s (tag %q): %w", field.Name, field.Column, err)
		}
	}

	return nil
}

func assignFieldValue(parent any, meta fieldMetadata, dst reflect.Value, src interface{}) error {
	if !dst.CanSet() {
		return fmt.Errorf("destination cannot be set")
	}

	if converted, handled, err := tryParentConverter(parent, meta, src); handled {
		if err != nil {
			return err
		}
		src = converted
	}

	if dst.CanAddr() {
		if converter, ok := dst.Addr().Interface().(HydraValueConverter); ok {
			return converter.HydraConvert(src)
		}
	}

	if src == nil {
		return assignNil(dst)
	}

	if dst.Kind() == reflect.Ptr {
		return assignPointerValue(parent, meta, dst, src)
	}

	if tryAssignDirect(dst, src) {
		return nil
	}

	switch dst.Kind() {
	case reflect.String:
		return assignString(dst, src)
	case reflect.Bool:
		return assignBool(dst, src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return assignInt(dst, src)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return assignUint(dst, src)
	case reflect.Float32, reflect.Float64:
		return assignFloat(dst, src)
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
		if tryAssignDirect(dst, src) {
			return nil
		}
		return fmt.Errorf("cannot assign %T to %s", src, dst.Type())
	case reflect.Interface:
		dst.Set(reflect.ValueOf(src))
		return nil
	default:
		return fmt.Errorf("unsupported destination kind %s for value %T", dst.Kind(), src)
	}
}

func tryParentConverter(parent any, meta fieldMetadata, src any) (any, bool, error) {
	provider, ok := parent.(HydraConverterProvider)
	if !ok {
		return nil, false, nil
	}

	converters := provider.HydraConverters()
	if converters == nil {
		return nil, false, nil
	}

	if converter, ok := converters[meta.Name]; ok {
		value, err := converter(src)
		return value, true, err
	}
	if converter, ok := converters[meta.Column]; ok {
		value, err := converter(src)
		return value, true, err
	}

	return nil, false, nil
}

func assignNil(dst reflect.Value) error {
	switch dst.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface:
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	default:
		return fmt.Errorf("cannot assign NULL to %s", dst.Type())
	}
}

func assignPointerValue(parent any, meta fieldMetadata, dst reflect.Value, src interface{}) error {
	if src == nil {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	elem := reflect.New(dst.Type().Elem())
	if err := assignFieldValue(parent, meta, elem.Elem(), src); err != nil {
		return err
	}
	dst.Set(elem)
	return nil
}

func tryAssignDirect(dst reflect.Value, src interface{}) bool {
	sv := reflect.ValueOf(src)
	if !sv.IsValid() {
		return false
	}
	if sv.Type().AssignableTo(dst.Type()) {
		dst.Set(sv)
		return true
	}
	if sv.Type().ConvertibleTo(dst.Type()) {
		dst.Set(sv.Convert(dst.Type()))
		return true
	}
	return false
}

func assignString(dst reflect.Value, src interface{}) error {
	switch v := src.(type) {
	case string:
		dst.SetString(v)
		return nil
	case []byte:
		dst.SetString(string(v))
		return nil
	default:
		return fmt.Errorf("cannot assign %T to %s", src, dst.Type())
	}
}

func assignBool(dst reflect.Value, src interface{}) error {
	switch v := src.(type) {
	case bool:
		dst.SetBool(v)
		return nil
	case int:
		dst.SetBool(v != 0)
		return nil
	case int8:
		dst.SetBool(v != 0)
		return nil
	case int16:
		dst.SetBool(v != 0)
		return nil
	case int32:
		dst.SetBool(v != 0)
		return nil
	case int64:
		dst.SetBool(v != 0)
		return nil
	case uint:
		dst.SetBool(v != 0)
		return nil
	case uint8:
		dst.SetBool(v != 0)
		return nil
	case uint16:
		dst.SetBool(v != 0)
		return nil
	case uint32:
		dst.SetBool(v != 0)
		return nil
	case uint64:
		dst.SetBool(v != 0)
		return nil
	case string:
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("cannot parse %q as bool", v)
		}
		dst.SetBool(parsed)
		return nil
	case []byte:
		parsed, err := strconv.ParseBool(string(v))
		if err != nil {
			return fmt.Errorf("cannot parse %q as bool", string(v))
		}
		dst.SetBool(parsed)
		return nil
	default:
		return fmt.Errorf("cannot assign %T to %s", src, dst.Type())
	}
}

func assignInt(dst reflect.Value, src interface{}) error {
	bits := dst.Type().Bits()
	switch v := src.(type) {
	case int:
		dst.SetInt(int64(v))
	case int8:
		dst.SetInt(int64(v))
	case int16:
		dst.SetInt(int64(v))
	case int32:
		dst.SetInt(int64(v))
	case int64:
		dst.SetInt(v)
	case uint:
		dst.SetInt(int64(v))
	case uint8:
		dst.SetInt(int64(v))
	case uint16:
		dst.SetInt(int64(v))
	case uint32:
		dst.SetInt(int64(v))
	case uint64:
		dst.SetInt(int64(v))
	case float32:
		dst.SetInt(int64(v))
	case float64:
		dst.SetInt(int64(v))
	case string:
		parsed, err := strconv.ParseInt(v, 10, bits)
		if err != nil {
			return fmt.Errorf("cannot parse %q as int", v)
		}
		dst.SetInt(parsed)
	case []byte:
		parsed, err := strconv.ParseInt(string(v), 10, bits)
		if err != nil {
			return fmt.Errorf("cannot parse %q as int", string(v))
		}
		dst.SetInt(parsed)
	default:
		return fmt.Errorf("cannot assign %T to %s", src, dst.Type())
	}
	return nil
}

func assignUint(dst reflect.Value, src interface{}) error {
	bits := dst.Type().Bits()
	switch v := src.(type) {
	case int:
		dst.SetUint(uint64(v))
	case int8:
		dst.SetUint(uint64(v))
	case int16:
		dst.SetUint(uint64(v))
	case int32:
		dst.SetUint(uint64(v))
	case int64:
		dst.SetUint(uint64(v))
	case uint:
		dst.SetUint(uint64(v))
	case uint8:
		dst.SetUint(uint64(v))
	case uint16:
		dst.SetUint(uint64(v))
	case uint32:
		dst.SetUint(uint64(v))
	case uint64:
		dst.SetUint(v)
	case float32:
		dst.SetUint(uint64(v))
	case float64:
		dst.SetUint(uint64(v))
	case string:
		parsed, err := strconv.ParseUint(v, 10, bits)
		if err != nil {
			return fmt.Errorf("cannot parse %q as uint", v)
		}
		dst.SetUint(parsed)
	case []byte:
		parsed, err := strconv.ParseUint(string(v), 10, bits)
		if err != nil {
			return fmt.Errorf("cannot parse %q as uint", string(v))
		}
		dst.SetUint(parsed)
	default:
		return fmt.Errorf("cannot assign %T to %s", src, dst.Type())
	}
	return nil
}

func assignFloat(dst reflect.Value, src interface{}) error {
	bits := dst.Type().Bits()
	switch v := src.(type) {
	case float32:
		dst.SetFloat(float64(v))
	case float64:
		dst.SetFloat(v)
	case int:
		dst.SetFloat(float64(v))
	case int8:
		dst.SetFloat(float64(v))
	case int16:
		dst.SetFloat(float64(v))
	case int32:
		dst.SetFloat(float64(v))
	case int64:
		dst.SetFloat(float64(v))
	case uint:
		dst.SetFloat(float64(v))
	case uint8:
		dst.SetFloat(float64(v))
	case uint16:
		dst.SetFloat(float64(v))
	case uint32:
		dst.SetFloat(float64(v))
	case uint64:
		dst.SetFloat(float64(v))
	case string:
		parsed, err := strconv.ParseFloat(v, bits)
		if err != nil {
			return fmt.Errorf("cannot parse %q as float", v)
		}
		dst.SetFloat(parsed)
	case []byte:
		parsed, err := strconv.ParseFloat(string(v), bits)
		if err != nil {
			return fmt.Errorf("cannot parse %q as float", string(v))
		}
		dst.SetFloat(parsed)
	default:
		return fmt.Errorf("cannot assign %T to %s", src, dst.Type())
	}
	return nil
}
