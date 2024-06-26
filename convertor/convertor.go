package convertor

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/bytedance/sonic"
)

// Str2Bytes 字符串转[]byte
func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Bytes2Str []byte转字符串
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ToString input的值转成字符串
func ToString(input interface{}) string {
	if input == nil {
		return ""
	}
	switch v := input.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case float64:
		ft, _ := input.(float64)
		return strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft, _ := input.(float32)
		return strconv.FormatFloat(float64(ft), 'f', -1, 32)
	case int:
		return strconv.Itoa(v)
	case uint:
		return strconv.Itoa(int(v))
	case int8:
		return strconv.Itoa(int(v))
	case uint8:
		return strconv.Itoa(int(v))
	case int16:
		return strconv.Itoa(int(v))
	case uint16:
		return strconv.Itoa(int(v))
	case int32:
		return strconv.Itoa(int(v))
	case uint32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case []byte:
		return Bytes2Str(v)
	case fmt.Stringer:
		return v.String()
	case error:
		return v.Error()
	default:
		if newValue, err := sonic.Marshal(input); err == nil {
			return string(newValue)
		} else {
			return ""
		}
	}
}

// ToInt convert value to int64 value, if input is not numerical, return 0 and error.
func ToInt(value any) (int64, error) {
	v := reflect.ValueOf(value)

	var result int64
	err := fmt.Errorf("ToInt: invalid value type %T", value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		result = v.Int()
		return result, nil
	case uint, uint8, uint16, uint32, uint64:
		result = int64(v.Uint())
		return result, nil
	case float32, float64:
		result = int64(v.Float())
		return result, nil
	case string:
		result, err = strconv.ParseInt(v.String(), 0, 64)
		if err != nil {
			result = 0
		}
		return result, err
	default:
		return result, err
	}
}
