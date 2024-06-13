package webui

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/snabb/sitemap"
)

func (w *WebUI) SitemapIndexTransactions(c *gin.Context) {
	transactionCount, err := w.db.GetTransactionCount()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	pages := (transactionCount / 10000) + 1

	si := sitemap.NewSitemapIndex()
	for i := 0; i < pages; i++ {
		result := "https://burned.money/sitemap/transactions/" + strconv.Itoa(i) + ".xml"
		si.Add(&sitemap.URL{Loc: result})
	}

	c.XML(200, si)
}

func (w *WebUI) SitemapTransactions(c *gin.Context) {
	index := c.Param("index")

	parts := strings.Split(index, ".")
	if len(parts) != 2 || parts[1] != "xml" {
		c.AbortWithStatus(404)
		return
	}

	page, err := strconv.Atoi(parts[0])
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	if page < 0 {
		c.AbortWithStatus(400)
		return
	}

	pageSize := 10000
	offset := page * pageSize
	txids, err := w.db.GetTransactionTxids(pageSize, offset)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	si := sitemap.New()
	for _, txid := range txids {
		result := "https://burned.money/transaction/" + txid
		si.Add(&sitemap.URL{Loc: result})
	}

	c.XML(200, si)
}

func (w *WebUI) SitemapIndexBlocks(c *gin.Context) {
	si := sitemap.NewSitemapIndex()
	for i := 0; i < 200; i++ {
		result := "https://burned.money/sitemap/blocks/" + strconv.Itoa(i) + ".xml"
		si.Add(&sitemap.URL{Loc: result})
	}

	c.XML(200, si)
}

func (w *WebUI) SitemapBlocks(c *gin.Context) {
	index := c.Param("index")

	parts := strings.Split(index, ".")

	if len(parts) != 2 || parts[1] != "xml" {
		c.AbortWithStatus(404)
		return
	}

	blockIndex, err := strconv.Atoi(parts[0])
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	if blockIndex < 0 || blockIndex > 199 {
		c.AbortWithStatus(400)
		return
	}

	si := sitemap.New()
	startBlock := blockIndex * 10000
	for i := startBlock; i < startBlock+10000; i++ {
		result := "https://burned.money/block/" + strconv.Itoa(i)
		si.Add(&sitemap.URL{Loc: result})
	}

	c.XML(200, si)
}
