package util

import (
	"fmt"
	"html/template"
)

func Int64ToBTC(sats int64) string {
	return fmt.Sprintf("%.8f", float64(sats)/1e8)
}

func NoEscapeHTML(str string) template.HTML {
	return template.HTML(str)
}

func FloatBTCToSats(btc float64) int64 {
	return int64(btc * 1e8)
}
