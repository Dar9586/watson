package types

import (
	"fmt"
	"reflect"
)

// ToValue converts an arbitrary value `v` into `Value` by applying the following rules:
//   * If v is any of int, int8, int16, int32, or int64, then v is converted to Int.
//   * If v is any of uint, uint8, uint16, uint32, or uint64, then v is converted to Uint.
//   * If v is float32 or float64, then v is converted to Float.
//   * If v is bool, then v is converted to Bool.
//   * If v is string, then v is converted to String.
//   * If v is a struct that implements `Marshaler`, then v is converted to Value by calling `v.MarshalWatson()`.
//   * If v is a struct that does not implement `Marshaler`, then v is converted to Object with its keys correspond to the fields of v.
//   * If v is a slice or an array, then v is converted to Array with its elements converted by these rules.
//   * If v is a map, then v is converted to Object with its elements converted by these rules.
//   * If v is a pointer, then v is converted to `Value` by converting `*v` with these rules.
//
// Note that you can configure struct fields by adding "watson" tag to fields.
// Tag must be like `watson:"name,flag1,flag2,...,flagN"`.
// If `ToValue` finds a field that has such tag, it uses `name` as a key of output instead of using the name of the field, or omits such field if `name` equals to "-".
//
// Currntly these flags are available:
//   omitempty      If the field is zero value, it will be omitted from the output.
//   inline         Inline the field. Currently the field must be a struct.
func ToValue(v interface{}) *Value {
	if v == nil {
		return NewNilValue()
	}
	switch v := v.(type) {
	case bool:
		return NewBoolValue(v)
	case int:
		return NewIntValue(int64(v))
	case int8:
		return NewIntValue(int64(v))
	case int16:
		return NewIntValue(int64(v))
	case int32:
		return NewIntValue(int64(v))
	case int64:
		return NewIntValue(v)
	case uint:
		return NewUintValue(uint64(v))
	case uint8:
		return NewUintValue(uint64(v))
	case uint16:
		return NewUintValue(uint64(v))
	case uint32:
		return NewUintValue(uint64(v))
	case uint64:
		return NewUintValue(uint64(v))
	case string:
		return NewStringValue([]byte(v))
	case float32:
		return NewFloatValue(float64(v))
	case float64:
		return NewFloatValue(v)
	}
	if marshaler, ok := v.(Marshaler); ok {
		return marshaler.MarshalWatson()
	}
	vv := reflect.ValueOf(v)
	return ToValueByReflection(vv)
}

// `ToValueByReflection` does almost the same thing as `ToValue`, but it always uses reflection.
func ToValueByReflection(v reflect.Value) *Value {
	if isMarshaler(v) {
		return marshalerToValueByReflection(v)
	} else if isIntFamily(v) {
		return intToValueByReflection(v)
	} else if isUintFamily(v) {
		return uintToValueByReflection(v)
	} else if isFloatFamily(v) {
		return floatToValueByReflection(v)
	} else if isBool(v) {
		return boolToValueByReflection(v)
	} else if isString(v) {
		return stringToValueByReflection(v)
	} else if isArray(v) {
		return sliceOrArrayToValueByReflection(v)
	} else if isStruct(v) {
		return structToValueByReflection(v)
	} else if isNil(v) {
		// Marshalers should be placed before nil so as to handle `MarshalWatson` correctly.
		return NewNilValue()
		// Maps, slices, and pointers should be placed after nil so as to convert nil into Nil correctly.
	} else if isPtr(v) {
		return reflectPtrToValue(v)
	} else if isMapConvertibleToValue(v) {
		return reflectMapToValue(v)
	} else if isSlice(v) {
		return sliceOrArrayToValueByReflection(v)
	}

	panic(fmt.Errorf("can't convert %s to *Value", v.Type().String()))
}

func intToValueByReflection(v reflect.Value) *Value {
	return NewIntValue(v.Int())
}

func uintToValueByReflection(v reflect.Value) *Value {
	return NewUintValue(v.Uint())
}

func floatToValueByReflection(v reflect.Value) *Value {
	return NewFloatValue(v.Float())
}

func boolToValueByReflection(v reflect.Value) *Value {
	return NewBoolValue(v.Bool())
}

func stringToValueByReflection(v reflect.Value) *Value {
	return NewStringValue([]byte(v.String()))
}

func reflectMapToValue(v reflect.Value) *Value {
	obj := map[string]*Value{}
	iter := v.MapRange()
	for iter.Next() {
		k := iter.Key().String()
		v := iter.Value()
		if v.CanInterface() {
			obj[k] = ToValue(v.Interface())
		} else {
			obj[k] = ToValueByReflection(v)
		}
	}
	return NewObjectValue(obj)
}

func sliceOrArrayToValueByReflection(v reflect.Value) *Value {
	arr := []*Value{}
	size := v.Len()
	for i := 0; i < size; i++ {
		elem := v.Index(i)
		if elem.CanInterface() {
			arr = append(arr, ToValue(elem.Interface()))
		} else {
			arr = append(arr, ToValueByReflection(elem))
		}
	}
	return NewArrayValue(arr)
}

func reflectPtrToValue(v reflect.Value) *Value {
	elem := v.Elem()
	if elem.CanInterface() {
		return ToValue(elem.Interface())
	} else {
		return ToValueByReflection(elem)
	}
}

func structToValueByReflection(v reflect.Value) *Value {
	obj := map[string]*Value{}
	addFields(obj, v)
	return NewObjectValue(obj)
}

func addFields(obj map[string]*Value, v reflect.Value) {
	size := v.NumField()
	t := v.Type()
	for i := 0; i < size; i++ {
		field := t.Field(i)
		tag := parseTag(&field)
		if tag.ShouldAlwaysOmit() {
			continue
		}
		name := tag.Key()
		elem := v.Field(i)
		if tag.OmitEmpty() && elem.IsZero() {
			continue
		}
		if tag.Inline() {
			addFields(obj, elem)
		} else if elem.CanInterface() {
			obj[name] = ToValue(elem.Interface())
		} else {
			obj[name] = ToValueByReflection(elem)
		}
	}
}

func marshalerToValueByReflection(v reflect.Value) *Value {
	marshal := v.MethodByName("MarshalWatson")
	ret := marshal.Call([]reflect.Value{})
	return ret[0].Interface().(*Value)
}
