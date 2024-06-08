package webui

import (
	"math/big"
	"net/http"

	"github.com/RaghavSood/blockreward"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/RaghavSood/btcsupply/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) HalvingSchedule(c *gin.Context) {
	schedule := blockreward.SubsidySchedule(blockreward.BitcoinMainnet)

	var emissionCurveHeights []int64
	var emissionCurveRewards []*types.BigInt
	var emissionCurveSupply []*types.BigInt

	// Calculate the supply and rewards in 1000 block intervals until the
	// reward reaches 0
	for i := int64(0); i < 7140001; i += 1000 {
		emissionCurveHeights = append(emissionCurveHeights, i)
		emissionCurveRewards = append(emissionCurveRewards, types.FromMathBigInt(big.NewInt(blockreward.SubsidyAtHeight(blockreward.BitcoinMainnet, i))))
		emissionCurveSupply = append(emissionCurveSupply, types.FromMathBigInt(big.NewInt(blockreward.SupplyAtHeight(blockreward.BitcoinMainnet, i))))
	}

	tmpl := templates.New()
	err := tmpl.Render(c.Writer, "halving_schedule.tmpl", map[string]interface{}{
		"Title":    "Bitcoin Halving Schedule",
		"Schedule": schedule,
		"Curve": map[string]interface{}{
			"Heights": emissionCurveHeights,
			"Rewards": emissionCurveRewards,
			"Supply":  emissionCurveSupply,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
