package webui

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/RaghavSood/blockreward"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/RaghavSood/btcsupply/types"
	"github.com/RaghavSood/btcsupply/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) HalvingSchedule(c *gin.Context) {
	schedule := blockreward.SubsidySchedule(blockreward.BitcoinMainnet)
	latestBlock, err := w.db.GetLatestBlock()
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to get latest block")
	}

	lastHalvingBlock := (latestBlock.BlockHeight / 210000) * 210000
	nextHalvingBlock := ((latestBlock.BlockHeight / 210000) + 1) * 210000

	currentSubsidy := types.FromMathBigInt(big.NewInt(blockreward.SubsidyAtHeight(blockreward.BitcoinMainnet, latestBlock.BlockHeight)))
	nextSubsidy := types.FromMathBigInt(big.NewInt(blockreward.SubsidyAtHeight(blockreward.BitcoinMainnet, nextHalvingBlock)))

	scheduleHeights := make([]int64, len(schedule))
	for i, s := range schedule {
		scheduleHeights[i] = s.Height
	}

	halvingBlocks, err := w.db.GetBlocksByHeights(scheduleHeights)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get halving blocks")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	blocks := make(map[int64]types.Block)
	for _, block := range halvingBlocks {
		blocks[block.BlockHeight] = block
	}

	for _, s := range schedule {
		if _, ok := blocks[s.Height]; !ok {
			blocks[s.Height] = util.FutureBlock(s.Height, latestBlock.BlockHeight)
		}
	}

	var emissionCurveHeights []int64
	var emissionCurveRewards []*types.BigInt
	var emissionCurveSupply []*types.BigInt

	// Calculate the supply and rewards in 1000 block intervals until the
	// reward reaches 0
	wasLower := true
	for i := int64(0); i < 7140001; i += 1000 {
		if wasLower {
			if i >= latestBlock.BlockHeight {
				wasLower = false
				emissionCurveHeights = append(emissionCurveHeights, latestBlock.BlockHeight)
				emissionCurveRewards = append(emissionCurveRewards, types.FromMathBigInt(big.NewInt(blockreward.SubsidyAtHeight(blockreward.BitcoinMainnet, latestBlock.BlockHeight))))
				emissionCurveSupply = append(emissionCurveSupply, types.FromMathBigInt(big.NewInt(blockreward.SupplyAtHeight(blockreward.BitcoinMainnet, latestBlock.BlockHeight))))
			}
		}
		emissionCurveHeights = append(emissionCurveHeights, i)
		emissionCurveRewards = append(emissionCurveRewards, types.FromMathBigInt(big.NewInt(blockreward.SubsidyAtHeight(blockreward.BitcoinMainnet, i))))
		emissionCurveSupply = append(emissionCurveSupply, types.FromMathBigInt(big.NewInt(blockreward.SupplyAtHeight(blockreward.BitcoinMainnet, i))))
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "halving_schedule.tmpl", map[string]interface{}{
		"Title":    "Bitcoin Halving Schedule",
		"Desc":     fmt.Sprintf("Track Bitcoin's Halving Schedule and Emission Curve. Next halving is at block #%d.", nextHalvingBlock),
		"Schedule": schedule,
		"Blocks":   blocks,
		"Curve": map[string]interface{}{
			"Heights": emissionCurveHeights,
			"Rewards": emissionCurveRewards,
			"Supply":  emissionCurveSupply,
		},
		"LatestBlock":      latestBlock,
		"CurrentSubsidy":   currentSubsidy,
		"NextSubsidy":      nextSubsidy,
		"LastHalvingBlock": lastHalvingBlock,
		"NextHalvingBlock": nextHalvingBlock,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
