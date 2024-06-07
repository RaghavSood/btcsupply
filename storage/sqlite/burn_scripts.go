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

func (d *SqliteBackend) GetBurnScript(script string) (types.BurnScript, error) {
	var burnScript types.BurnScript
	err := d.db.QueryRow("SELECT * FROM burn_scripts WHERE script = ?", script).Scan(&burnScript.ID, &burnScript.Script, &burnScript.ConfidenceLevel, &burnScript.Provenance, &burnScript.CreatedAt, &burnScript.ScriptGroup)

	return burnScript, err
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

func (d *SqliteBackend) GetBurnScriptSummary(script string) (types.BurnScriptSummary, error) {
	query := `SELECT 
						    bs.script,
								bs.confidence_level,
								bs.provenance,
								bs.script_group,
							  COUNT(DISTINCT(l.tx_id)) AS transactions,
						    SUM(l.amount) AS total_loss
						FROM 
						    burn_scripts bs
						JOIN 
						    losses l ON bs.script = l.burn_script
						WHERE l.burn_script = ?
						GROUP BY 
						    bs.script`

	var summary types.BurnScriptSummary
	err := d.db.QueryRow(query, script).Scan(&summary.Script, &summary.ConfidenceLevel, &summary.Provenance, &summary.Group, &summary.Transactions, &summary.TotalLoss)

	return summary, err
}

func (d *SqliteBackend) GetBurnScriptSummariesForGroup(group string) ([]types.BurnScriptSummary, error) {
	query := `SELECT
						  bs.script,
							bs.confidence_level,
							bs.provenance,
							bs.script_group,
							COUNT(DISTINCT(l.tx_id)) AS transactions,
							SUM(l.amount) AS total_loss
						FROM
						  burn_scripts bs
						JOIN
						  losses l ON bs.script = l.burn_script
						WHERE
							bs.script_group = ?
						GROUP BY
						bs.script;`

	rows, err := d.db.Query(query, group)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	summaries, err := scanBurnScriptSummaries(rows)

	return summaries, err
}

func (d *SqliteBackend) GetBurnScriptSummaries(limit int) ([]types.BurnScriptSummary, error) {
	query := `SELECT 
						    bs.script,
								bs.confidence_level,
								bs.provenance,
								bs.script_group,
							  COUNT(DISTINCT(l.tx_id)) AS transactions,
						    SUM(l.amount) AS total_loss
						FROM 
						    burn_scripts bs
						JOIN 
						    losses l ON bs.script = l.burn_script
						GROUP BY 
						    bs.script
						ORDER BY 
						    total_loss DESC
						LIMIT ?;`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries, err := scanBurnScriptSummaries(rows)

	return summaries, err
}

func (d *SqliteBackend) BurnScriptExists(script string) (bool, error) {
	var exists bool
	err := d.db.QueryRow("SELECT EXISTS(SELECT 1 FROM burn_scripts WHERE script = ?)", script).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func scanBurnScriptSummaries(rows *sql.Rows) ([]types.BurnScriptSummary, error) {
	var summaries []types.BurnScriptSummary
	for rows.Next() {
		var summary types.BurnScriptSummary
		err := rows.Scan(&summary.Script, &summary.ConfidenceLevel, &summary.Provenance, &summary.Group, &summary.Transactions, &summary.TotalLoss)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
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
