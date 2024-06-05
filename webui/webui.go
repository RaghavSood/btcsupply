package webui

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/RaghavSood/blockreward"
	"github.com/RaghavSood/btcsupply/notes"
	"github.com/RaghavSood/btcsupply/prices"
	"github.com/RaghavSood/btcsupply/static"
	"github.com/RaghavSood/btcsupply/storage"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/RaghavSood/btcsupply/types"
	"github.com/RaghavSood/btcsupply/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type WebUI struct {
	db storage.Storage
}

func NewWebUI(db storage.Storage) *WebUI {
	return &WebUI{
		db: db,
	}
}

func (w *WebUI) Serve() {
	router := gin.Default()

	router.GET("/", w.Index)

	router.GET("/blocks", w.Blocks)
	router.GET("/block/:identifier", w.Block)

	router.GET("/transaction/:hash", w.Transaction)

	router.GET("/schedule", w.HalvingSchedule)

	router.StaticFS("/static", http.FS(static.Static))

	router.Run(":8080")
}

func (w *WebUI) Index(c *gin.Context) {
	losses, err := w.db.GetRecentLosses(10)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recent losses")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	indexStats, err := w.db.GetIndexStatistics()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get index statistics")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	btcPrice, err := prices.GetBTCUSDPrice()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get BTC price")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	indexStats.PlannedSupply = types.FromMathBigInt(big.NewInt(blockreward.SupplyAtHeight(blockreward.BitcoinMainnet, indexStats.LastBlockHeight)))
	indexStats.CirculatingSupply = types.FromMathBigInt(big.NewInt(0).Sub(indexStats.PlannedSupply.BigInt(), indexStats.BurnedSupply.BigInt()))
	indexStats.CurrentPrice = btcPrice
	indexStats.AdjustedPrice = util.RevaluePriceWithAdjustedSupply(indexStats.PlannedSupply, indexStats.CirculatingSupply, btcPrice)
	indexStats.PriceChange = indexStats.AdjustedPrice - indexStats.CurrentPrice

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "index.tmpl", map[string]interface{}{
		"Title":  "Home",
		"Losses": losses,
		"Stats":  indexStats,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

}

func (w *WebUI) Blocks(c *gin.Context) {
	blocks, err := w.db.GetLossyBlocks(10)
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

func (w *WebUI) Block(c *gin.Context) {
	identifier := c.Param("identifier")

	blockHash, blockHeight, err := w.db.GetBlockIdentifiers(identifier)
	if err != nil {
		log.Error().
			Err(err).
			Str("identifier", identifier).
			Msg("Failed to get block identifiers")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	block, err := w.db.GetBlock(blockHash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	losses, err := w.db.GetBlockLosses(blockHash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block losses")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "block.tmpl", map[string]interface{}{
		"Title":  fmt.Sprintf("Block %d - %s", blockHeight, blockHash),
		"Block":  block,
		"Losses": losses,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (w *WebUI) Transaction(c *gin.Context) {
	hash := c.Param("hash")

	losses, err := w.db.GetTransactionLosses(hash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transaction losses")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	transaction, err := w.db.GetTransactionDetail(hash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transaction")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var notePointers []notes.NotePointer
	for _, loss := range losses {
		notePointers = append(notePointers, notes.NotePointer{
			NoteType:     notes.Output,
			PathElements: []string{loss.TxID, fmt.Sprintf("%d", loss.Vout)},
		})

	}

	for _, vout := range transaction.Vout {
		notePointers = append(notePointers, notes.NotePointer{
			NoteType:     notes.Script,
			PathElements: []string{vout.ScriptPubKey.Hex},
		})
	}

	notes := notes.ReadNotes(notePointers)

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "transaction.tmpl", map[string]interface{}{
		"Title":       "Transaction",
		"Losses":      losses,
		"Transaction": transaction,
		"Notes":       notes,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
