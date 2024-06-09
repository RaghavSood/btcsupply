package util

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"runtime/debug"
	"time"

	"github.com/RaghavSood/btcsupply/notes"
	"github.com/RaghavSood/btcsupply/types"
)

func Int64ToBTC(sats int64) *types.BigInt {
	bigInt := big.NewInt(sats)
	return types.FromMathBigInt(bigInt)
}

func NoEscapeHTML(str string) template.HTML {
	return template.HTML(str)
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

func IsScriptInBurnScripts(script string, burnScripts []types.BurnScript) bool {
	for _, burnScript := range burnScripts {
		if burnScript.Script == script {
			return true
		}
	}

	return false
}

func IsScriptInNotes(script string, noteList []notes.Note) bool {
	for _, note := range noteList {
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

func GitCommit() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		settings := buildInfo.Settings
		for _, s := range settings {
			if s.Key == "vcs.revision" {
				return s.Value
			}
		}
	}

	return "unknown"
}
