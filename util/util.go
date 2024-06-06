package util

import (
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/RaghavSood/btcsupply/notes"
	"github.com/RaghavSood/btcsupply/types"
	"github.com/rs/zerolog/log"
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

func IsScriptInNotes(script string, noteList []notes.Note) bool {
	for _, note := range noteList {
		log.Debug().
			Strs("pathElements", note.PathElements).
			Str("script", script).
			Msg("Checking note")
		if note.PathElements[0] == script {
			return true
		}
	}
	return false
}

func FutureBlock(height int64, lastBlock int64) types.Block {
	blocksTill := height - lastBlock
	timeToMine := blocksTill * 10 * int64(time.Minute)
	durationToMine := time.Duration(timeToMine)
	estimatedMineTime := time.Now().Add(time.Duration(durationToMine))

	return types.Block{
		BlockHeight:    height,
		BlockHash:      "Unknown",
		BlockTimestamp: estimatedMineTime,
	}
}

func BlockHeightString(height int64) string {
	return fmt.Sprintf("%d", height)
}

func TimeDisclaimer(target time.Time) string {
	if target.After(time.Now()) {
		duration := target.Sub(time.Now()).Round(time.Minute)
		return fmt.Sprintf(" (Estimated to be mined in %s)", duration.String())
	}
	return ""
}
