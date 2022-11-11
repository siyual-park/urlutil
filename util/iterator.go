package util

import (
	"reflect"
)

func KeyTo[T any](value T, convert func(key any) any) T {
	if IsNil(value) {
		return value
	}
	res := keyTo(reflect.ValueOf(value), convert)
	return res.Interface().(T)
}

func ValueTo[T any](value T, convert func(key any) any) T {
	if IsNil(value) {
		return value
	}
	res := valueTo(reflect.ValueOf(value), convert)
	return res.Interface().(T)
}

func keyTo(value reflect.Value, convert func(key any) any) reflect.Value {
	if !value.IsValid() {
		return value
	}

	v := reflect.ValueOf(value.Interface())

	switch v.Kind() {
	case reflect.Map:
		res := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			v := keyTo(v.MapIndex(key), convert)
			res.SetMapIndex(reflect.ValueOf(convert(key.Interface())), v)
		}
		return res
	case reflect.Slice:
		res := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			res.Index(i).Set(keyTo(v.Index(i), convert))
		}
		return res
	case reflect.Array:
		res := reflect.New(v.Type())
		for i := 0; i < v.Len(); i++ {
			res.Elem().Index(i).Set(keyTo(v.Index(i), convert))
		}
		return res.Elem()
	case reflect.Pointer:
		if v.IsNil() {
			return v
		}
		e := v.Elem()
		res := reflect.New(e.Type())
		res.Elem().Set(keyTo(e, convert))
		return res
	}

	return value
}

func valueTo(value reflect.Value, convert func(key any) any) reflect.Value {
	if !value.IsValid() {
		return value
	}

	v := reflect.ValueOf(convert(value.Interface()))

	switch v.Kind() {
	case reflect.Map:
		res := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			v := valueTo(v.MapIndex(key), convert)
			res.SetMapIndex(key, v)
		}
		return res
	case reflect.Slice:
		res := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			res.Index(i).Set(valueTo(v.Index(i), convert))
		}
		return res
	case reflect.Array:
		res := reflect.New(v.Type())
		for i := 0; i < v.Len(); i++ {
			res.Elem().Index(i).Set(valueTo(v.Index(i), convert))
		}
		return res.Elem()
	case reflect.Pointer:
		if v.IsNil() {
			return v
		}
		e := v.Elem()
		res := reflect.New(e.Type())
		res.Elem().Set(valueTo(e, convert))
		return res
	}

	return v
}
