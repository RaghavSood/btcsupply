package webui

import (
	"net/http"

	"github.com/RaghavSood/btcsupply/templates"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) OpReturns(c *gin.Context) {
	opReturns, err := w.db.GetOpReturnSummaries(500)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recent op returns")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "opreturns.tmpl", map[string]interface{}{
		"Title":     "Top OP_RETURNs",
		"OpReturns": opReturns,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (w *WebUI) OpReturn(c *gin.Context) {
	opReturn := c.Param("opreturn")

	opReturnSummary, err := w.db.GetOpReturnSummary(opReturn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get op return summary")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	burnTransactions, err := w.db.GetTransactionLossSummaryForScript(opReturn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get transactions for script")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "opreturn.tmpl", map[string]interface{}{
		"Title":        "OP_RETURN " + opReturn,
		"OpReturn":     opReturnSummary,
		"Transactions": burnTransactions,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
