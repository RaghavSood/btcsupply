package webui

import (
	"net/http"

	"github.com/RaghavSood/blockreward"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) HalvingSchedule(c *gin.Context) {
	schedule := blockreward.SubsidySchedule(blockreward.BitcoinMainnet)

	tmpl := templates.New()
	err := tmpl.Render(c.Writer, "halving_schedule.tmpl", map[string]interface{}{
		"Title":    "Halving Schedule",
		"Schedule": schedule,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
