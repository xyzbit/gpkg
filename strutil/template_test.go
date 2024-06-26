package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSprintfDict(t *testing.T) {
	// 创建一个字典
	dict := map[string]string{
		"name": "Alice",
		"age":  "18",
	}
	// 使用SprintfDict来格式化字符串
	s := SprintfDict("Hello, ${name}. You are ${age} years old.", dict)
	assert.Equal(t, "Hello, Alice. You are 18 years old.", s)
}

func TestSprintfVar(t *testing.T) {
	// 创建一个字典
	dict := map[string]string{
		"name": "Alice",
		"age":  "18",
	}
	// 使用SprintfDict来格式化字符串
	s := SprintfVar("Hello, ${global.name}. You are ${global.age} years old.", "global.", dict)
	assert.Equal(t, "Hello, Alice. You are 18 years old.", s)
}
