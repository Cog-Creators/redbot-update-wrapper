package logutils

import (
	"log/slog"
	"reflect"
)

type anyLogValue struct {
	value reflect.Value
}

func (lv anyLogValue) LogValue() slog.Value {
	var value reflect.Value
	switch lv.value.Kind() {
	case reflect.Pointer, reflect.Interface:
		if lv.value.IsNil() {
			return slog.AnyValue(nil)
		}
		value = lv.value.Elem()
	default:
		value = lv.value
	}

	t := value.Type()
	if t.Kind() == reflect.Map && t.Key().Kind() == reflect.String {
		attrs := []slog.Attr{}
		for k, v := range value.Seq2() {
			attrs = append(attrs, slog.Any(k.String(), anyLogValue{v}))
		}
		return slog.GroupValue(attrs...)
	}
	if t.Kind() != reflect.Struct {
		return slog.AnyValue(value.Interface())
	}

	attrs := []slog.Attr{}
	for field, v := range value.Fields() {
		attrs = append(attrs, slog.Any(field.Name, anyLogValue{v}))
	}
	return slog.GroupValue(attrs...)
}

type structLogValue[T any] struct {
	value T
}

func NewStructLogValue[T any](v *T) slog.LogValuer {
	return structLogValue[*T]{v}
}

func (lv structLogValue[T]) LogValue() slog.Value {
	return anyLogValue{reflect.ValueOf(lv.value)}.LogValue()
}
