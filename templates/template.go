package templates

import (
	"embed"
	"html/template"
	"io"
	"time"

	"github.com/RaghavSood/btcsupply/types"
	"github.com/RaghavSood/btcsupply/util"
)

//go:embed *
var Templates embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	funcMap := template.FuncMap{
		"now":                   func() interface{} { return time.Now() },
		"Int64ToBTC":            util.Int64ToBTC,
		"NoEscape":              util.NoEscapeHTML,
		"PrettyPrintJSON":       util.PrettyPrintJSON,
		"ScriptPubKeyDisplay":   types.ScriptPubKeyDisplay,
		"IsScriptInBurnScripts": util.IsScriptInBurnScripts,
		"IsScriptInNotes":       util.IsScriptInNotes,
		"ValueToBigInt":         types.ValueToBigInt,
		"PrettyDuration":        util.PrettyDuration,
		"GitCommit":             util.GitCommit,
		"ParseDecodeScript":     types.ParseDecodeScriptJSON,
		"HexToASCII":            util.HexToASCII,
		"BTCValueToUSD":         util.BTCValueToUSD,
		"FormatNumber":          util.FormatNumber,
	}

	templates := template.Must(template.New("").Funcs(funcMap).ParseFS(Templates, "footer.tmpl", "header.tmpl", "base.tmpl", "note.tmpl", "btc_fiat.tmpl"))
	return &Template{
		templates: templates,
	}
}

func (t *Template) Render(w io.Writer, contentTemplate string, data interface{}) error {
	tmpl, err := t.templates.Clone()
	if err != nil {
		return err
	}

	// Parse the specific content template
	_, err = tmpl.ParseFS(Templates, contentTemplate)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "base.tmpl", data)
}
