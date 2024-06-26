package strutil

import (
	"strings"

	"github.com/xyzbit/pkg/convertor"
)

const (
	varPatternLeft  = "${"
	varPatternRight = "}"
)

// SprintfDict 根据pattern和dict格式化字符串。
// pattern是一个字符串，包含${key}形式的变量占位符。
// dict map[string]string 被替换的变量。
// 例如，SprintfDict（“你好，$｛name｝！”，map[string]string｛“name”：“Alice”｝）返回“你好，Alice！”。
// 如果pattern包含一个不在dict中的key，它将保持不变。
// 如果dict包含一个不在pattern中的key，它将被忽略。
func SprintfDict(pattern string, dict map[string]string) string {
	return SprintfVar(pattern, "", dict)
}

func SprintfDictAny(pattern string, dict map[string]interface{}) string {
	dictStr := make(map[string]string, len(dict))
	for k, v := range dict {
		dictStr[k] = convertor.ToString(v)
	}
	return SprintfVar(pattern, "", dictStr)
}

func SprintfDictBytes(pattern string, dict map[string]string) []byte {
	str := SprintfVar(pattern, "", dict)
	return convertor.Str2Bytes(str)
}

// SprintfVar 根据pattern和dict格式化字符串。
// pattern是一个字符串，包含${key}形式的变量占位符。
// dict map[string]string 被替换的变量。
// keyPrefix key前缀，dict所有key将会加上前缀keyPrefix再进行替换。
// 例如，SprintfDict（“你好，$｛name｝！”，map[string]string｛“name”：“Alice”｝）返回“你好，Alice！”。
// 如果pattern包含一个不在dict中的key，它将保持不变。
// 如果dict包含一个不在pattern中的key，它将被忽略。
func SprintfVar(pattern string, keyPrefix string, dict map[string]string) string {
	result := pattern
	for key, value := range dict {
		result = ProcessVar(result, keyPrefix+key, value)
	}
	return result
}

func ProcessVar(pattern, key, val string) string {
	varPattern := varPatternLeft + key + varPatternRight
	return strings.Replace(pattern, varPattern, val, -1)
}

// CheckHasVar 检查字符串是否有占位符
func CheckHasVar(str string) bool {
	return strings.Contains(str, "${") && strings.Contains(str, "}")
}

// RemoveBraces A function that takes a string with ${} and returns a string without them
func RemoveBraces(s string) string {
	// Create a new empty string
	result := ""
	// Loop through each character in the input string
	for i := 0; i < len(s); i++ {
		// Get the current character
		c := s[i]
		// If the character is $, check the next character
		if c == '$' && i+1 < len(s) {
			// If the next character is {, skip it and move to the next one
			if s[i+1] == '{' {
				i++
				continue
			}
		}
		// If the character is }, skip it and move to the next one
		if c == '}' {
			continue
		}
		// If the character is a space, skip it and move to the next one
		if c == ' ' {
			continue
		}
		// Otherwise, append the character to the result string
		result += string(c)
	}
	// Return the result string
	return result
}
