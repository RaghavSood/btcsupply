package sqlite

import (
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
	var scripts []types.BurnScript
	rows, err := d.db.Query("SELECT * FROM burn_scripts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var script types.BurnScript
		err = rows.Scan(&script.Script, &script.ConfidenceLevel, &script.Provenance, &script.CreatedAt)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, script)
	}

	return scripts, nil
}

func (d *SqliteBackend) BurnScriptExists(script string) (bool, error) {
	var exists bool
	err := d.db.QueryRow("SELECT EXISTS(SELECT 1 FROM burn_scripts WHERE script = ?)", script).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
