package webui

import (
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/RaghavSood/blockreward"
	btypes "github.com/RaghavSood/btcsupply/bitcoinrpc/types"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/RaghavSood/btcsupply/types"
	"github.com/RaghavSood/btcsupply/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) Block(c *gin.Context) {
	identifier := c.Param("identifier")

	block, err := w.getBlockOrFutureBlock(identifier)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	losses, err := w.db.GetTransactionLossSummaryForBlock(identifier)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block losses")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	blockSummary, err := w.db.GetBlockLossSummary(identifier)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Err(err).Msg("Failed to get block loss summary")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var indexStats types.IndexStatistics
	var txOutSetInfo types.TxOutSetInfo
	var blockStats btypes.BlockStats
	if !block.IsFutureBlock() {
		indexStats, err = w.statsForHeight(block.BlockHeight)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get block statistics")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		txOutSetInfo, err = w.db.GetTxOutSetInfo(block.BlockHash)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get txoutset info")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		blockStats, err = w.db.GetBlockStats(block.BlockHash)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get block stats")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	theoreticalSubsidy := blockreward.SubsidyAtHeight(blockreward.BitcoinMainnet, block.BlockHeight)
	theoreticalSupply := blockreward.SupplyAtHeight(blockreward.BitcoinMainnet, block.BlockHeight)

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "block.tmpl", map[string]interface{}{
		"Title":              fmt.Sprintf("Block %d - %s", block.BlockHeight, block.BlockHash),
		"Block":              block,
		"BlockSummary":       blockSummary,
		"Losses":             losses,
		"TheoreticalSubsidy": types.FromMathBigInt(big.NewInt(theoreticalSubsidy)),
		"TheoreticalSupply":  types.FromMathBigInt(big.NewInt(theoreticalSupply)),
		"Stats":              indexStats,
		"CoinStats":          txOutSetInfo,
		"BlockStats":         blockStats,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (w *WebUI) Blocks(c *gin.Context) {
	blocks, err := w.db.GetLossyBlocks(500)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get lossy blocks")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "blocks.tmpl", map[string]interface{}{
		"Title":  "Blocks",
		"Blocks": blocks,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (w *WebUI) getBlockOrFutureBlock(identifier string) (types.Block, error) {
	block, err := w.db.GetBlock(identifier)
	if err != nil && err != sql.ErrNoRows {
		return types.Block{}, fmt.Errorf("Failed to get block: %w", err)
	}

	if err == sql.ErrNoRows {
		latestBlock, err := w.db.GetLatestBlock()
		if err != nil {
			return types.Block{}, fmt.Errorf("Failed to get latest block: %w", err)
		}

		int64Identifier, err := strconv.ParseInt(identifier, 10, 64)
		if err != nil {
			return types.Block{}, fmt.Errorf("Failed to parse identifier: %w", err)
		}

		block = util.FutureBlock(int64Identifier, latestBlock.BlockHeight)
	}

	return block, nil
}
