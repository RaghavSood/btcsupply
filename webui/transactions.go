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
		"Title":        "Burn Transactions",
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

	txSummary, err := w.db.GetTransactionLossSummaryForTxid(hash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transaction summary")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	notePointers, hasNulldata, scripts := transaction.NotePointers()

	for _, loss := range losses {
		notePointers = append(notePointers, notes.NotePointer{
			NoteType:     notes.Output,
			PathElements: []string{loss.TxID, fmt.Sprintf("%d", loss.Vout)},
		})

		if loss.Vout == -1 {
			notePointers = append(notePointers, notes.NotePointer{
				NoteType:     notes.Coinbase,
				PathElements: []string{"coinbase"},
			})
		}
	}
	burnScripts, err := w.db.GetBurnScriptsByScripts(scripts)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get burn scripts")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Ensure the common nulldata note shows on every nulldata tx page
	if hasNulldata {
		notePointers = append(notePointers, notes.NotePointer{
			NoteType:     notes.NullData,
			PathElements: []string{"nulldata"},
		})
	}

	seenGroups := make(map[string]bool)
	for _, burnScript := range burnScripts {
		if !seenGroups[burnScript.ScriptGroup] {
			seenGroups[burnScript.ScriptGroup] = true
			notePointers = append(notePointers, notes.NotePointer{
				NoteType:     notes.ScriptGroup,
				PathElements: []string{burnScript.ScriptGroup},
			})
		}
	}

	notes := notes.ReadNotes(notePointers)

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "transaction.tmpl", map[string]interface{}{
		"Title":       fmt.Sprintf("Transaction %s", hash),
		"Losses":      losses,
		"Transaction": transaction,
		"Summary":     txSummary,
		"Notes":       notes,
		"BurnScripts": burnScripts,
		"Block":       block,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
