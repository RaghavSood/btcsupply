package webui

import (
	"fmt"
	"strings"

	"github.com/RaghavSood/btcsupply/address"
	"github.com/gin-gonic/gin"
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

	w.renderTemplate(c, "search.tmpl", map[string]interface{}{
		"Title": "Search",
		"Desc":  "Search the Bitcoin Blockchain for blocks, transactions, scripts, and addresses that burn BTC.",
		"Query": q,
	})
}
