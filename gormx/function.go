package gormx

import (
	"fmt"
	"strings"

	"github.com/xyzbit/gpkg/convertor"
)

var (
	Upper = UpperFunc{}
	Lower = LowerFunc{}
)

type Function interface {
	Name() string
	Expression(params ...string) string
	ConvertVal(v any) string
}

type UpperFunc struct{}

func (u UpperFunc) Name() string {
	return "UPPER"
}

func (u UpperFunc) Expression(params ...string) string {
	return fmt.Sprintf("UPPER(`%s`)", params[0])
}

func (u UpperFunc) ConvertVal(v any) string {
	return strings.ToUpper(convertor.ToString(v))
}

type LowerFunc struct{}

func (l LowerFunc) Name() string {
	return "LOWER"
}

func (l LowerFunc) Expression(params ...string) string {
	return fmt.Sprintf("LOWER(`%s`)", params[0])
}

func (l LowerFunc) ConvertVal(v any) string {
	return strings.ToLower(convertor.ToString(v))
}
