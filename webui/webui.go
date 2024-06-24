package webui

import (
	"net/http"

	"github.com/RaghavSood/btcsupply/btclogger"
	"github.com/RaghavSood/btcsupply/middleware"
	"github.com/RaghavSood/btcsupply/static"
	"github.com/RaghavSood/btcsupply/storage"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type WebUI struct {
	db       storage.Storage
	readonly bool
}

func NewWebUI(db storage.Storage, noindex bool) *WebUI {
	return &WebUI{
		db:       db,
		readonly: noindex,
	}
}

func (w *WebUI) Serve() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.StructuredLogger(btclogger.NewLogger("gin-webui")))
	router.Use(gin.Recovery())

	router.GET("/", w.Index)

	router.GET("/blocks", w.Blocks)
	router.GET("/block/:identifier", w.Block)

	router.GET("/transactions", w.Transactions)
	router.GET("/transactions/coinbase", w.CoinbaseTransactions)
	router.GET("/transaction/:hash", w.Transaction)

	router.GET("/scripts", w.Scripts)
	router.GET("/script/:script", w.Script)

	router.GET("/scriptgroups", w.ScriptGroups)
	router.GET("/scriptgroup/:scriptgroup", w.ScriptGroup)

	router.GET("/opreturns", w.OpReturns)
	router.GET("/opreturn/:opreturn", w.OpReturn)

	router.GET("/tips", w.Tips)
	router.GET("/why", w.Why)
	router.GET("/methodology", w.Methodology)
	router.GET("/schedule", w.HalvingSchedule)

	router.GET("/stats", w.Stats)

	router.GET("/search", w.Search)

	router.GET("/ogimage/:slug", w.OGImage)

	router.StaticFS("/static", http.FS(static.Static))
	// Serve /favicon.ico and /robots.txt from the root
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", http.FS(static.Static))
	})

	router.GET("/robots.txt", func(c *gin.Context) {
		c.FileFromFS("robots.txt", http.FS(static.Static))
	})

	sitemap := router.Group("/sitemap")
	{
		sitemap.GET("/index/blocks", w.SitemapIndexBlocks)
		sitemap.GET("/blocks/:index", w.SitemapBlocks)

		sitemap.GET("/index/transactions", w.SitemapIndexTransactions)
		sitemap.GET("/transactions/:index", w.SitemapTransactions)

		sitemap.GET("/index/scripts", w.SitemapIndexScripts)
		sitemap.GET("/scripts/:index", w.SitemapScripts)
	}

	feeds := router.Group("/feeds")
	{
		feeds.GET("/", w.FeedIndex)
		feeds.GET("/blocks", w.FeedBlocks)
		feeds.GET("/transactions", w.FeedTransactions)
	}

	router.Run(":8080")
}

func (w *WebUI) Index(c *gin.Context) {
	losses, err := w.db.GetTransactionSummary(50, 1, false)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recent losses")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	indexStats, err := w.statsForHeight(-1)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get index statistics")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	w.renderTemplate(c, "index.tmpl", map[string]interface{}{
		"Title":  "Track Bitcoin Supply",
		"Losses": losses,
		"Stats":  indexStats,
	})
}

func (w *WebUI) Tips(c *gin.Context) {
	w.renderTemplate(c, "tips.tmpl", map[string]interface{}{
		"Title": "Tips - Submit a BTC loss",
		"Desc":  "Submit a tip about lost, missing, or burned BTC and help us track the total supply of Bitcoin",
	})
}

func (w *WebUI) Why(c *gin.Context) {
	w.renderTemplate(c, "why.tmpl", map[string]interface{}{
		"Title": "Why?",
	})
}

func (w *WebUI) Methodology(c *gin.Context) {
	w.renderTemplate(c, "methodology.tmpl", map[string]interface{}{
		"Title": "Methodology",
		"Desc":  "How we calculate the total supply of Bitcoin",
	})
}

func (w *WebUI) renderTemplate(c *gin.Context, template string, params map[string]interface{}) {
	tmpl := templates.New()
	params["Readonly"] = w.readonly
	err := tmpl.Render(c.Writer, template, params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
