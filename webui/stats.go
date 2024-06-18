package webui

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/RaghavSood/blockreward"
	"github.com/RaghavSood/btcsupply/prices"
	"github.com/RaghavSood/btcsupply/types"
	"github.com/RaghavSood/btcsupply/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) Stats(c *gin.Context) {
	heightLossSummary, err := w.db.GetHeightLossSummary()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get height loss summary")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	theoreticalSupplySummary := make([]types.HeightLossSummary, len(heightLossSummary))
	for i, summary := range heightLossSummary {
		theoreticalSupplySummary[i].BlockHeight = summary.BlockHeight
		theoreticalSupplySummary[i].TotalLoss = types.FromMathBigInt(big.NewInt(blockreward.SupplyAtHeight(blockreward.BitcoinMainnet, summary.BlockHeight)))
	}

	w.renderTemplate(c, "stats.tmpl", map[string]interface{}{
		"Title":             "Stats",
		"HeightLossSummary": heightLossSummary,
		"TheoreticalSupply": theoreticalSupplySummary,
	})
}

func (w *WebUI) statsForHeight(height int64) (types.IndexStatistics, error) {
	indexStats, err := w.db.GetIndexStatistics(height)
	if err != nil {
		return types.IndexStatistics{}, fmt.Errorf("Failed to get index statistics: %w", err)
	}

	btcPrice, err := prices.GetBTCUSDPrice()
	if err != nil {
		return types.IndexStatistics{}, fmt.Errorf("Failed to get BTC price: %w", err)
	}

	indexStats.PlannedSupply = types.FromMathBigInt(big.NewInt(blockreward.SupplyAtHeight(blockreward.BitcoinMainnet, indexStats.LastBlockHeight)))
	indexStats.CirculatingSupply = types.FromMathBigInt(big.NewInt(0).Sub(indexStats.PlannedSupply.BigInt(), indexStats.BurnedSupply.BigInt()))
	indexStats.CurrentPrice = btcPrice
	indexStats.AdjustedPrice = util.RevaluePriceWithAdjustedSupply(indexStats.PlannedSupply, indexStats.CirculatingSupply, btcPrice)
	indexStats.PriceChange = indexStats.AdjustedPrice - indexStats.CurrentPrice

	return indexStats, nil
}
