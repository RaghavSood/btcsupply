package webui

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (w *WebUI) SitemapBlocks(c *gin.Context) {
	index := c.Param("index")

	parts := strings.Split(index, ".")

	if len(parts) != 2 || parts[1] != "txt" {
		c.AbortWithStatus(400)
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

	c.Header("Content-Type", "text/plain")
	startBlock := blockIndex * 10000
	for i := startBlock; i < startBlock+10000; i++ {
		result := "https://burned.money/block/" + strconv.Itoa(i) + "\n"
		c.Writer.WriteString(result)
	}
}
