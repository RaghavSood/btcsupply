package webui

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/RaghavSood/blockreward"
	"github.com/RaghavSood/btcsupply/notes"
	"github.com/RaghavSood/btcsupply/static"
	"github.com/RaghavSood/btcsupply/storage"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/RaghavSood/btcsupply/types"
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

	indexStats, err := w.statsForHeight(-1)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get index statistics")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

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

	block, err := w.db.GetBlock(identifier)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	losses, err := w.db.GetBlockLosses(block.BlockHash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block losses")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	indexStats, err := w.statsForHeight(block.BlockHeight)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block statistics")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	txOutSetInfo, err := w.db.GetTxOutSetInfo(block.BlockHash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get txoutset info")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	blockStats, err := w.db.GetBlockStats(block.BlockHash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block stats")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	theoreticalSubsidy := blockreward.SubsidyAtHeight(blockreward.BitcoinMainnet, block.BlockHeight)

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "block.tmpl", map[string]interface{}{
		"Title":              fmt.Sprintf("Block %d - %s", block.BlockHeight, block.BlockHash),
		"Block":              block,
		"Losses":             losses,
		"TheoreticalSubsidy": types.FromMathBigInt(big.NewInt(theoreticalSubsidy)),
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

func (w *WebUI) Transaction(c *gin.Context) {
	hash := c.Param("hash")

	losses, err := w.db.GetTransactionLosses(hash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transaction losses")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rawTransaction, err := w.db.GetTransaction(hash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transaction")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	transaction, err := rawTransaction.TransactionDetail()
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse transaction")
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
