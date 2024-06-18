package webui

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/RaghavSood/btcsupply/notes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) Scripts(c *gin.Context) {
	topScripts, err := w.db.GetBurnScriptSummaries(500)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top scripts")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	w.renderTemplate(c, "scripts.tmpl", map[string]interface{}{
		"Title":   "Bitcoin Burn Scripts",
		"Desc":    "View the top 500 burn scripts on the Bitcoin blockchain.",
		"Scripts": topScripts,
	})
}

func (w *WebUI) Script(c *gin.Context) {
	script := c.Param("script")

	burnScriptSummary, err := w.db.GetBurnScriptSummary(script)
	if err != nil && err != sql.ErrNoRows {
		log.Error().Err(err).Msg("Failed to get script summary")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	scriptNotePointer := notes.NotePointer{
		NoteType:     notes.Script,
		PathElements: []string{script},
	}

	groupNotePointer := notes.NotePointer{
		NoteType:     notes.ScriptGroup,
		PathElements: []string{burnScriptSummary.Group},
	}

	notes := notes.ReadNotes([]notes.NotePointer{scriptNotePointer, groupNotePointer})

	burnTransactions, err := w.db.GetTransactionLossSummaryForScript(script)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transactions for script")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	w.renderTemplate(c, "script.tmpl", map[string]interface{}{
		"Title":         fmt.Sprintf("Burn Script %s", script),
		"Desc":          fmt.Sprintf("%s BTC burned in %d transactions for this script.", burnScriptSummary.TotalLoss.SatoshisToBTC(true), burnScriptSummary.Transactions),
		"OGImage":       fmt.Sprintf("https://burned.money/ogimage/script-%s.png", script),
		"ScriptSummary": burnScriptSummary,
		"Transactions":  burnTransactions,
		"Notes":         notes,
	})
}
