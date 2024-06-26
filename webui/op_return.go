package webui

import (
	"fmt"
	"net/http"

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

	w.renderTemplate(c, "opreturns.tmpl", map[string]interface{}{
		"Title":     "Top OP_RETURNs",
		"Desc":      "The top 500 OP_RETURN losses in Bitcoin history.",
		"OpReturns": opReturns,
	})
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

	w.renderTemplate(c, "opreturn.tmpl", map[string]interface{}{
		"Title":        "OP_RETURN " + opReturn,
		"Desc":         fmt.Sprintf("%s BTC burned in %d transactions with this OP_RETURN.", opReturnSummary.TotalLoss.SatoshisToBTC(true), opReturnSummary.Transactions),
		"OGImage":      fmt.Sprintf("https://burned.money/ogimage/opreturn-%s.png", opReturn),
		"OpReturn":     opReturnSummary,
		"Transactions": burnTransactions,
	})
}
