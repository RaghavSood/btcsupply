package webui

import (
	"fmt"
	"net/http"

	"github.com/RaghavSood/btcsupply/static"
	"github.com/RaghavSood/btcsupply/storage"
	"github.com/RaghavSood/btcsupply/templates"
	"github.com/gin-gonic/gin"
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
	router.StaticFS("/static", http.FS(static.Static))

	router.Run(":8080")
}

func (w *WebUI) Index(c *gin.Context) {
	tmpl := templates.New()
	err := tmpl.Render(c.Writer, "index.tmpl", map[string]interface{}{
		"Title": "Home",
	})
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

}
