package util

import "fmt"

func Int64ToBTC(sats int64) string {
	return fmt.Sprintf("%.8f", float64(sats)/1e8)
}
