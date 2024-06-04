package util

import (
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/RaghavSood/btcsupply/types"
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

func PrettyPrintJSON(data interface{}) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(b), nil
}

func RevaluePriceWithAdjustedSupply(expectedSupply, circulatingSupply *types.BigInt, currentPrice float64) float64 {
	expectedSupplyFloat, _ := expectedSupply.BigFloat().Float64()
	circulatingSupplyFloat, _ := circulatingSupply.BigFloat().Float64()
	return currentPrice * (expectedSupplyFloat / circulatingSupplyFloat)
}
