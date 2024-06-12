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
	db storage.Storage
}

func NewWebUI(db storage.Storage) *WebUI {
	return &WebUI{
		db: db,
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
	router.GET("/transaction/:hash", w.Transaction)

	router.GET("/scripts", w.Scripts)
	router.GET("/script/:script", w.Script)

	router.GET("/scriptgroups", w.ScriptGroups)
	router.GET("/scriptgroup/:scriptgroup", w.ScriptGroup)

	router.GET("/opreturns", w.OpReturns)
	router.GET("/opreturn/:opreturn", w.OpReturn)

	router.GET("/why", w.Why)
	router.GET("/methodology", w.Methodology)
	router.GET("/schedule", w.HalvingSchedule)

	router.GET("/stats", w.Stats)

	router.StaticFS("/static", http.FS(static.Static))
	// Serve /favicon.ico from the root
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", http.FS(static.Static))
	})

	sitemap := router.Group("/sitemap")
	{
		sitemap.GET("/index/blocks", w.SitemapIndexBlocks)
		sitemap.GET("/blocks/:index", w.SitemapBlocks)

		sitemap.GET("/index/transactions", w.SitemapIndexTransactions)
		sitemap.GET("/transactions/:index", w.SitemapTransactions)
	}

	router.Run(":8080")
}

func (w *WebUI) Index(c *gin.Context) {
	losses, err := w.db.GetTransactionSummary(50)
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

	tmpl := templates.New()
	err = tmpl.Render(c.Writer, "index.tmpl", map[string]interface{}{
		"Title":  "Track Bitcoin Supply",
		"Losses": losses,
		"Stats":  indexStats,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

}

func (w *WebUI) Why(c *gin.Context) {
	tmpl := templates.New()
	err := tmpl.Render(c.Writer, "why.tmpl", map[string]interface{}{
		"Title": "Why?",
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (w *WebUI) Methodology(c *gin.Context) {
	tmpl := templates.New()
	err := tmpl.Render(c.Writer, "methodology.tmpl", map[string]interface{}{
		"Title": "Methodology",
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
