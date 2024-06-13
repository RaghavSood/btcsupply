package webui

import (
	"fmt"
	"net/http"

	"github.com/RaghavSood/btcsupply/address"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) Search(c *gin.Context) {
	q := c.Query("q")
	log.Debug().Str("query", q).Msg("Search query")

	// Check if query is a block number or hash
	block, err := w.db.GetBlock(q)
	if err == nil {
		c.Redirect(302, fmt.Sprintf("/block/%d", block.BlockHeight))
		return
	}

	// Check if query is a transaction ID
	tx, err := w.db.GetTransaction(q)
	if err == nil {
		c.Redirect(302, fmt.Sprintf("/transaction/%s", tx.TxID))
		return
	}

	// Check if query is a script
	burnScript, err := w.db.GetBurnScript(q)
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
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
