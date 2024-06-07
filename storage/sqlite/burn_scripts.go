package sqlite

import (
	"database/sql"
	"strings"

	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetOnlyBurnScripts() ([]string, error) {
	var scripts []string
	rows, err := d.db.Query("SELECT script FROM burn_scripts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var script string
		err = rows.Scan(&script)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, script)
	}

	return scripts, nil
}

func (d *SqliteBackend) GetBurnScripts() ([]types.BurnScript, error) {
	rows, err := d.db.Query("SELECT * FROM burn_scripts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanBurnScripts(rows)
}

func (d *SqliteBackend) GetBurnScriptsByScripts(scripts []string) ([]types.BurnScript, error) {
	anyScripts := make([]interface{}, len(scripts))
	for i, script := range scripts {
		anyScripts[i] = script
	}

	rows, err := d.db.Query("SELECT * FROM burn_scripts WHERE script IN (?"+strings.Repeat(",?", len(scripts)-1)+")", anyScripts...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanBurnScripts(rows)
}

func (d *SqliteBackend) BurnScriptExists(script string) (bool, error) {
	var exists bool
	err := d.db.QueryRow("SELECT EXISTS(SELECT 1 FROM burn_scripts WHERE script = ?)", script).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func scanBurnScripts(rows *sql.Rows) ([]types.BurnScript, error) {
	var scripts []types.BurnScript
	for rows.Next() {
		var script types.BurnScript
		err := rows.Scan(&script.ID, &script.Script, &script.ConfidenceLevel, &script.Provenance, &script.CreatedAt, &script.ScriptGroup)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, script)
	}
	return scripts, nil
}
