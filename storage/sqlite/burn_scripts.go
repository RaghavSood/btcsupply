package sqlite

import (
	"github.com/RaghavSood/btcsupply/types"
)

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
