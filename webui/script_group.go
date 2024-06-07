package webui

import (
	"net/http"

	"github.com/RaghavSood/btcsupply/templates"
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

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "scriptgroup.tmpl", map[string]interface{}{
		"Title":            "Script Group",
		"BurnTransactions": burnScriptSummaries,
		"GroupSummary":     groupSummary,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (w *WebUI) ScriptGroups(c *gin.Context) {
	scriptGroups, err := w.db.GetScriptGroupSummaries(500)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get script groups")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "scriptgroups.tmpl", map[string]interface{}{
		"Title":        "Script Groups",
		"ScriptGroups": scriptGroups,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
