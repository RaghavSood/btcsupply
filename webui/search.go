package webui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RaghavSood/btcsupply/address"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) Search(c *gin.Context) {
	q := c.Query("q")
	lowerQ := strings.ToLower(q)

	// Check if query is a block number or hash
	block, err := w.db.GetBlock(lowerQ)
	if err == nil {
		c.Redirect(302, fmt.Sprintf("/block/%d", block.BlockHeight))
		return
	}

	// Check if query is a transaction ID
	tx, err := w.db.GetTransaction(lowerQ)
	if err == nil {
		c.Redirect(302, fmt.Sprintf("/transaction/%s", tx.TxID))
		return
	}

	// Check if query is a script
	burnScript, err := w.db.GetBurnScript(lowerQ)
	if err == nil {
		c.Redirect(302, fmt.Sprintf("/script/%s", burnScript.Script))
		return
	}

	// Check if query is an address
	script, err := address.AddressToScript(q)
	if err == nil {
		burnScript, err = w.db.GetBurnScript(script)
		if err == nil {
			c.Redirect(302, fmt.Sprintf("/script/%s", burnScript.Script))
			return
		}
	}

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "search.tmpl", map[string]interface{}{
		"Title": "Search",
		"Desc":  "Search the Bitcoin Blockchain for blocks, transactions, scripts, and addresses that burn BTC.",
		"Query": q,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
