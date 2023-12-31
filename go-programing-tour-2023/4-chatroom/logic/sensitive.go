package logic

import (
	"strings"

	"github.com/dapings/examples/go-programing-tour-2023/chatroom/global"
)

// FilterSensitive 过滤敏感信息。
func FilterSensitive(content string) string {
	for _, word := range global.SensitiveWords {
		content = strings.ReplaceAll(content, word, "**")
	}

	return content
}
