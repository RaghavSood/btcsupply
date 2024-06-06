package webui

import (
	"net/http"

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
	router := gin.Default()

	router.GET("/", w.Index)

	router.GET("/blocks", w.Blocks)
	router.GET("/block/:identifier", w.Block)

	router.GET("/transactions", w.Transactions)
	router.GET("/transaction/:hash", w.Transaction)

	router.GET("/why", w.Why)
	router.GET("/methodology", w.Methodology)
	router.GET("/schedule", w.HalvingSchedule)

	router.StaticFS("/static", http.FS(static.Static))
	// Serve /favicon.ico from the root
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", http.FS(static.Static))
	})

	router.Run(":8080")
}

func (w *WebUI) Index(c *gin.Context) {
	losses, err := w.db.GetTransactionLossSummary(50)
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
		"Title":  "Home",
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
