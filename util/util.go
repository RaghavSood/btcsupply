package util

import (
	"encoding/hex"
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
	estimatedMineTime := time.Now().Add(time.Duration(durationToMine)).Truncate(time.Second)

	return types.Block{
		BlockHeight:    height,
		BlockHash:      "Unknown",
		BlockTimestamp: estimatedMineTime,
	}
}

func BlockHeightString(height int64) string {
	return fmt.Sprintf("%d", height)
}

func PrettyDuration(t time.Time, significantPlaces int) string {
	now := time.Now()
	var d time.Duration
	var prefix string
	var suffix string

	if now.After(t) {
		d = time.Since(t)
		prefix = "~"
		suffix = " ago"
	} else {
		d = time.Until(t)
		prefix = "in ~"
	}

	const (
		secondsInMinute = 60
		secondsInHour   = 60 * secondsInMinute
		secondsInDay    = 24 * secondsInHour
		secondsInMonth  = 30 * secondsInDay
		secondsInYear   = 12 * secondsInMonth
	)

	totalSeconds := int(d.Seconds())
	years := totalSeconds / secondsInYear
	months := (totalSeconds % secondsInYear) / secondsInMonth
	days := (totalSeconds % secondsInMonth) / secondsInDay
	hours := (totalSeconds % secondsInDay) / secondsInHour
	minutes := (totalSeconds % secondsInHour) / secondsInMinute
	seconds := totalSeconds % secondsInMinute

	parts := []struct {
		value int
		unit  string
	}{
		{years, "year"},
		{months, "month"},
		{days, "day"},
		{hours, "hour"},
		{minutes, "minute"},
		{seconds, "second"},
	}

	result := ""
	count := 0
	for _, part := range parts {
		if part.value > 0 {
			if count > 0 {
				result += ", "
			}
			result += fmt.Sprintf("%d %s", part.value, part.unit)
			if part.value > 1 {
				result += "s"
			}
			count++
			if significantPlaces > 0 && count >= significantPlaces {
				break
			}
		}
	}

	if result == "" {
		result = "0 seconds"
	}

	return fmt.Sprintf("%s%s%s", prefix, result, suffix)
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

func HexToASCII(hexstring string) string {
	// Drop OP_RETURN
	if hexstring[:2] == "6a" {
		hexstring = hexstring[2:]
	}

	bytes, err := hex.DecodeString(hexstring)
	if err != nil {
		fmt.Printf("Error decoding hex string: %w\n", err)
		return ""
	}

	return string(bytes)
}
