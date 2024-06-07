package webui

import (
	"net/http"

	"github.com/RaghavSood/btcsupply/templates"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) Scripts(c *gin.Context) {
	topScripts, err := w.db.GetBurnScriptSummary(500)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top scripts")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "scripts.tmpl", map[string]interface{}{
		"Title":   "Top Scripts",
		"Scripts": topScripts,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
