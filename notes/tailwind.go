package notes

import (
	"fmt"
	"strings"
)

var fullMatch = map[string]string{}

var prefixMatch = map[string]string{
	"h3": "text-lg font-bold",
}

func wrapForTailwind(content string) string {
	for k, v := range fullMatch {
		content = strings.ReplaceAll(content, tag(k), tagWithClass(k, v))
	}

	for k, v := range prefixMatch {
		fmt.Println(k, v)
		content = strings.ReplaceAll(content, "<"+k, "<"+k+" class=\""+v+"\"")
	}

	return content
}

func tag(value string) string {
	return "<" + value + ">"
}

func tagWithClass(tag string, class string) string {
	return "<" + tag + " class=\"" + class + "\">"
}
