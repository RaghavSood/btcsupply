package webui

import (
	"fmt"
	"net/http"

	"github.com/RaghavSood/btcsupply/notes"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/RaghavSood/btcsupply/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) Transactions(c *gin.Context) {
	recentTransactions, err := w.db.GetTransactionLossSummary(500)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recent transactions")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "transactions.tmpl", map[string]interface{}{
		"Title":        "Transactions",
		"Transactions": recentTransactions,
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

	block, err := w.getBlockOrFutureBlock(util.BlockHeightString(rawTransaction.BlockHeight))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get block")
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

	hasNulldata := false
	for _, vout := range transaction.Vout {
		noteType := notes.Script

		if vout.ScriptPubKey.Type == "nulldata" {
			hasNulldata = true
			noteType = notes.NullData
		}

		notePointers = append(notePointers, notes.NotePointer{
			NoteType:     noteType,
			PathElements: []string{vout.ScriptPubKey.Hex},
		})
	}

	// Ensure the common nulldata note shows on every nulldata tx page
	if hasNulldata {
		notePointers = append(notePointers, notes.NotePointer{
			NoteType:     notes.NullData,
			PathElements: []string{"nulldata"},
		})
	}

	notes := notes.ReadNotes(notePointers)

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "transaction.tmpl", map[string]interface{}{
		"Title":       "Transaction",
		"Losses":      losses,
		"Transaction": transaction,
		"Notes":       notes,
		"Block":       block,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
