package gormx

import (
	"testing"
)

func TestLower(t *testing.T) {
	t.Log(Lower.Name())
	t.Log(Lower.Expression("a"))
	t.Log(Lower.ConvertVal("ABc"))
}

func TestUpper(t *testing.T) {
	t.Log(Upper.Name())
	t.Log(Upper.Expression("a"))
	t.Log(Upper.ConvertVal("abc"))
}
