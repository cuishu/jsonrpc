package jsonrpc

import (
	"errors"
	"reflect"
)

func convert(input interface{}, shouldBe reflect.Kind) (interface{}, error) {
	switch shouldBe {
	case reflect.Int:
		if v, ok := input.(float64); ok {
			return int(v), nil
		}
		return int(0), errors.New("type Error")
	case reflect.Int8:
		if v, ok := input.(float64); ok {
			return int8(v), nil
		}
		return int8(0), errors.New("type Error")
	case reflect.Int16:
		if v, ok := input.(float64); ok {
			return int16(v), nil
		}
		return int16(0), errors.New("type Error")
	case reflect.Int32:
		if v, ok := input.(float64); ok {
			return int32(v), nil
		}
		return int32(0), errors.New("type Error")
	case reflect.Int64:
		if v, ok := input.(float64); ok {
			return int64(v), nil
		}
		return int64(0), errors.New("type Error")
	case reflect.Uint8:
		if v, ok := input.(float64); ok {
			return uint8(v), nil
		}
		return uint8(0), errors.New("type Error")
	case reflect.Uint16:
		if v, ok := input.(float64); ok {
			return uint16(v), nil
		}
		return uint16(0), errors.New("type Error")
	case reflect.Uint32:
		if v, ok := input.(float64); ok {
			return uint32(v), nil
		}
		return uint32(0), errors.New("type Error")
	case reflect.Uint64:
		if v, ok := input.(float64); ok {
			return uint64(v), nil
		}
		return uint64(0), errors.New("type Error")
	case reflect.Uint:
		if v, ok := input.(float64); ok {
			return uint(v), nil
		}
		return uint(0), errors.New("type Error")
	case reflect.String:
		if v, ok := input.(string); ok {
			return v, nil
		}
		return "", errors.New("type Error")
	case reflect.Bool:
		if v, ok := input.(bool); ok {
			return v, nil
		}
		return false, errors.New("type Error")
	case reflect.Float32:
		if v, ok := input.(float32); ok {
			return v, nil
		}
		return float32(0), errors.New("type Error")
	case reflect.Float64:
		if v, ok := input.(float64); ok {
			return v, nil
		}
		return float64(0), errors.New("type Error")
	case reflect.Slice:
		if reflect.ValueOf(input).Kind() == reflect.Slice {
			return input, nil
		}
		return nil, errors.New("type Error")
	case reflect.Struct:
		if reflect.ValueOf(input).Kind() == reflect.Struct {
			return input, nil
		}
		return input, errors.New("type Error")
	case reflect.Ptr:
		if reflect.ValueOf(input).Kind() == reflect.Ptr {
			return input, nil
		}
		return nil, errors.New("type Error")
	case reflect.Map:
		if reflect.ValueOf(input).Kind() == reflect.Map {
			return input, nil
		}
		return nil, errors.New("type Error")
	}
	return input, errors.New("unsupport type.")
}
