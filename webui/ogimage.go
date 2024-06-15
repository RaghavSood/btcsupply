package webui

import (
	"fmt"
	"image"
	"image/color"
	"net/http"
	"strings"

	"github.com/RaghavSood/btcsupply/static"
	"github.com/RaghavSood/btcsupply/util"
	"github.com/RaghavSood/ogimage"
	"github.com/gin-gonic/gin"
)

func (w *WebUI) OGImage(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.String(http.StatusBadRequest, "missing slug")
		return
	}

	// Generate the image
	img, err := w.generateOGImage(slug)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to generate image")
		return
	}

	// Write the image to the response
	c.Header("Content-Type", "image/png")
	c.Writer.Write(img)
}

func (w *WebUI) generateOGImage(slug string) ([]byte, error) {
	parts := strings.Split(slug, "-")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid slug: %s", slug)
	}

	var title string
	var subtitle string
	var description string
	var err error

	switch parts[0] {
	case "tx":
		title, subtitle, description, err = w.getTxOGData(parts[1])
	case "block":
		title, subtitle, description, err = w.getBlockOGData(parts[1])
	case "script":
		title, subtitle, description, err = w.getScriptOGData(parts[1])
	case "scriptgroup":
		title, subtitle, description, err = w.getScriptGroupOGData(parts[1])
	case "opreturn":
		title, subtitle, description, err = w.getOpReturnOGData(parts[1])
	default:
		title = "Explore BTC Burns"
		subtitle = "Discover the history of Bitcoin burns"
	}

	templateBytes, err := static.Static.ReadFile("template.png")
	if err != nil {
		return nil, err
	}

	logoBytes, err := static.Static.ReadFile("logo.png")
	if err != nil {
		return nil, err
	}

	fontBytes, err := static.Static.ReadFile("menlo.ttf")
	if err != nil {
		return nil, err
	}

	ogImage, err := ogimage.NewOgImage(templateBytes, logoBytes)

	titleText := ogimage.Text{
		Content:  title,
		FontData: fontBytes,
		FontSize: 64,
		Color:    color.White,
		Point:    image.Point{20, 305},
	}

	subtitleText := ogimage.Text{
		Content:  subtitle,
		FontData: fontBytes,
		FontSize: 25,
		Color:    color.White,
		Point:    image.Point{20, 345},
	}

	descriptionText := ogimage.Text{
		Content:  description,
		FontData: fontBytes,
		FontSize: 30,
		Color:    color.White,
		Point:    image.Point{20, 450},
	}

	config := ogimage.Config{
		Position: ogimage.BottomRight,
		Padding:  20,
		Texts:    []ogimage.Text{titleText, subtitleText, descriptionText},
	}

	imageBytes, err := ogImage.Generate(config)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

func (w *WebUI) getTxOGData(txid string) (string, string, string, error) {
	txid = strings.ToLower(txid)

	transactionSummary, err := w.db.GetTransactionSummaryForTxid(txid)
	if err != nil {
		return "", "", "", err
	}

	title := fmt.Sprintf("%s BTC Burned", util.FormatNumber(transactionSummary.TotalLoss.SatoshisToBTC(true)))
	subtitle := fmt.Sprintf("Transaction %s", txid)
	description := fmt.Sprintf("That's worth $%s", util.FormatNumber(fmt.Sprintf("%.2f", util.BTCValueToUSD(transactionSummary.TotalLoss))))

	return title, subtitle, description, nil
}

func (w *WebUI) getBlockOGData(identifier string) (string, string, string, error) {
	identifier = strings.ToLower(identifier)

	blockSummary, err := w.db.GetBlockLossSummary(identifier)
	if err != nil {
		return "", "", "", err
	}

	title := fmt.Sprintf("%s BTC Burned", util.FormatNumber(blockSummary.TotalLost.SatoshisToBTC(true)))
	subtitle := fmt.Sprintf("Block %s", identifier)
	description := fmt.Sprintf("That's worth $%s", util.FormatNumber(fmt.Sprintf("%.2f", util.BTCValueToUSD(blockSummary.TotalLost))))

	return title, subtitle, description, nil
}

func (w *WebUI) getScriptOGData(script string) (string, string, string, error) {
	script = strings.ToLower(script)

	burnScript, err := w.db.GetBurnScriptSummary(script)
	if err != nil {
		return "", "", "", err
	}

	decodeScript, err := burnScript.ParseDecodeScript()
	if err != nil {
		return "", "", "", err
	}

	displayAddress := decodeScript.DisplayAddress(script)
	if len(displayAddress) > 64 {
		displayAddress = displayAddress[:61] + "..."
	}

	title := fmt.Sprintf("%s BTC Burned", util.FormatNumber(burnScript.TotalLoss.SatoshisToBTC(true)))
	subtitle := fmt.Sprintf("Address %s", displayAddress)
	description := fmt.Sprintf("That's worth $%s", util.FormatNumber(fmt.Sprintf("%.2f", util.BTCValueToUSD(burnScript.TotalLoss))))

	return title, subtitle, description, nil
}

func (w *WebUI) getScriptGroupOGData(group string) (string, string, string, error) {
	group = strings.ToLower(group)

	burnScriptGroup, err := w.db.GetScriptGroupSummary(group)
	if err != nil {
		return "", "", "", err
	}

	title := fmt.Sprintf("%s BTC Burned", util.FormatNumber(burnScriptGroup.TotalLoss.SatoshisToBTC(true)))
	subtitle := fmt.Sprintf("Group %s", group)
	description := fmt.Sprintf("That's worth $%s", util.FormatNumber(fmt.Sprintf("%.2f", util.BTCValueToUSD(burnScriptGroup.TotalLoss))))

	return title, subtitle, description, nil
}

func (w *WebUI) getOpReturnOGData(opreturn string) (string, string, string, error) {
	opreturn = strings.ToLower(opreturn)

	opReturnSummary, err := w.db.GetOpReturnSummary(opreturn)
	if err != nil {
		return "", "", "", err
	}

	if len(opreturn) > 64 {
		opreturn = opreturn[:61] + "..."
	}

	title := fmt.Sprintf("%s BTC Burned", util.FormatNumber(opReturnSummary.TotalLoss.SatoshisToBTC(true)))
	subtitle := fmt.Sprintf("OP_RETURN %s", opreturn)
	description := fmt.Sprintf("That's worth $%s", util.FormatNumber(fmt.Sprintf("%.2f", util.BTCValueToUSD(opReturnSummary.TotalLoss))))

	return title, subtitle, description, nil
}
