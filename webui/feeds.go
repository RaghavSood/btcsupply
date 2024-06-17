package webui

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/RaghavSood/btcsupply/templates"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/rs/zerolog/log"
)

func (w *WebUI) FeedIndex(c *gin.Context) {
	tmpl := templates.New()
	err := tmpl.Render(c.Writer, "feeds.tmpl", map[string]interface{}{
		"Title": "Feeds",
		"Desc":  "Received updates on BTC burn transactions and blocks via RSS/ATOM feeds.",
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (w *WebUI) FeedBlocks(c *gin.Context) {
	minLoss := c.DefaultQuery("min_loss", "100000")

	intMinLoss, err := strconv.Atoi(minLoss)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse min_loss")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	blocks, err := w.db.GetLossyBlocksWithMinimumLoss(100, int64(intMinLoss))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get lossy blocks")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "burned.money - Burn Blocks",
		Link:        &feeds.Link{Href: "https://burned.money"},
		Description: "Blocks with a minimum loss of " + minLoss + " satoshis",
		Author:      &feeds.Author{Name: "burned.money", Email: "hello@burned.money"},
		Created:     now,
	}

	for _, block := range blocks {
		item := &feeds.Item{
			Title:       fmt.Sprintf("Block #%d", block.BlockHeight),
			Link:        &feeds.Link{Href: fmt.Sprintf("https://burned.money/block/%d", block.BlockHeight)},
			Description: fmt.Sprintf("Block #%d burned %s BTC", block.BlockHeight, block.TotalLost.SatoshisToBTC(true)),
			Author:      &feeds.Author{Name: "burned.money", Email: "hello@burned.money"},
			Created:     block.BlockTimestamp,
		}

		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToAtom()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate feed")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, "application/atom+xml", []byte(rss))
}

func (w *WebUI) FeedTransactions(c *gin.Context) {
	minLoss := c.DefaultQuery("min_loss", "100000")

	intMinLoss, err := strconv.Atoi(minLoss)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse min_loss")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	transactions, err := w.db.GetTransactionSummary(100, int64(intMinLoss))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get lossy transactions")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "burned.money - Burn Transactions",
		Link:        &feeds.Link{Href: "https://burned.money"},
		Description: "Transactions with a minimum loss of " + minLoss + " satoshis",
		Author:      &feeds.Author{Name: "burned.money", Email: "hello@burned.money"},
		Created:     now,
	}

	for _, tx := range transactions {
		item := &feeds.Item{
			Title:       fmt.Sprintf("Transaction %s", tx.Txid),
			Link:        &feeds.Link{Href: fmt.Sprintf("https://burned.money/transaction/%s", tx.Txid)},
			Description: fmt.Sprintf("%s BTC burned in Transaction %s", tx.TotalLoss.SatoshisToBTC(true), tx.Txid),
			Author:      &feeds.Author{Name: "burned.money", Email: "hello@burned.money"},
			Created:     tx.BlockTimestamp,
		}

		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToAtom()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate feed")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, "application/atom+xml", []byte(rss))
}
