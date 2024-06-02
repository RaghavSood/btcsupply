package webui

import (
	"fmt"
	"net/http"

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
	router.GET("/block/:hash", w.Block)

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

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "index.tmpl", map[string]interface{}{
		"Title":  "Home",
		"Losses": losses,
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
	hash := c.Param("hash")

	block, err := w.db.GetBlock(hash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	losses, err := w.db.GetBlockLosses(hash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block losses")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "block.tmpl", map[string]interface{}{
		"Title":  fmt.Sprintf("Block %d - %s", block.BlockHeight, hash),
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

	var noteIDs []string
	for _, loss := range losses {
		noteIDs = append(noteIDs, fmt.Sprintf("output:%s:%d", loss.TxID, loss.Vout))
	}

	for _, vout := range transaction.Vout {
		noteIDs = append(noteIDs, fmt.Sprintf("address:%s", vout.ScriptPubKey.Hex))
	}

	rawNotes, err := w.db.GetLossNotes(noteIDs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get loss notes")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var renderedNotes []types.LossNote
	for _, note := range rawNotes {
		renderedNote := types.LossNote{
			NoteID:      note.NoteID,
			Description: notes.RenderNote(note),
			CreatedAt:   note.CreatedAt,
			Version:     note.Version,
		}
		renderedNotes = append(renderedNotes, renderedNote)
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "transaction.tmpl", map[string]interface{}{
		"Title":       "Transaction",
		"Losses":      losses,
		"Transaction": transaction,
		"Notes":       renderedNotes,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
