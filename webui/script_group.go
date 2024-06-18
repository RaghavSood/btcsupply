package webui

import (
	"fmt"
	"net/http"

	"github.com/RaghavSood/btcsupply/notes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) ScriptGroup(c *gin.Context) {
	scriptGroup := c.Param("scriptgroup")

	groupSummary, err := w.db.GetScriptGroupSummary(scriptGroup)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get script group summary")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	burnScriptSummaries, err := w.db.GetBurnScriptSummariesForGroup(scriptGroup)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get script summaries for group")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	groupNotePointer := notes.NotePointer{
		NoteType:     notes.ScriptGroup,
		PathElements: []string{scriptGroup},
	}

	notes := notes.ReadNotes([]notes.NotePointer{groupNotePointer})

	w.renderTemplate(c, "scriptgroup.tmpl", map[string]interface{}{
		"Title":            fmt.Sprintf("Script Group %s", scriptGroup),
		"Desc":             fmt.Sprintf("%s BTC lost across %d transactions in this group.", groupSummary.TotalLoss.SatoshisToBTC(true), groupSummary.Transactions),
		"OGImage":          fmt.Sprintf("https://burned.money/ogimage/scriptgroup-%s.png", scriptGroup),
		"BurnTransactions": burnScriptSummaries,
		"GroupSummary":     groupSummary,
		"Notes":            notes,
	})
}

func (w *WebUI) ScriptGroups(c *gin.Context) {
	scriptGroups, err := w.db.GetScriptGroupSummaries(500)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get script groups")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	w.renderTemplate(c, "scriptgroups.tmpl", map[string]interface{}{
		"Title":        "Script Groups",
		"Desc":         "View the top 500 groups that have burned Bitcoin.",
		"ScriptGroups": scriptGroups,
	})
}
