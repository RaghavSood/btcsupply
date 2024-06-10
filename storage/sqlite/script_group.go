package sqlite

import "github.com/RaghavSood/btcsupply/types"

func (d *SqliteBackend) GetScriptGroupSummary(group string) (types.ScriptGroupSummary, error) {
	query := `SELECT 
					    bs.script_group,
					    SUM(l.amount) AS total_loss,
					    count(DISTINCT(bs.script)) AS scripts,
					    count(DISTINCT(l.tx_id)) AS transactions
					  FROM 
					    burn_scripts bs                                           
					  JOIN 
					    losses l ON bs.script = l.burn_script
						WHERE
						  bs.script_group = ?`

	var summary types.ScriptGroupSummary
	err := d.db.QueryRow(query, group).Scan(&summary.ScriptGroup, &summary.TotalLoss, &summary.Scripts, &summary.Transactions)

	return summary, err
}

func (d *SqliteBackend) GetScriptGroupSummaries(limit int) ([]types.ScriptGroupSummary, error) {
	query := `SELECT 
					    bs.script_group,
					    SUM(l.amount) AS total_loss,
					    count(DISTINCT(bs.script)) AS scripts,
					    count(DISTINCT(l.tx_id)) AS transactions
					  FROM 
					    burn_scripts bs                                           
					  JOIN 
					    losses l ON bs.script = l.burn_script
					  GROUP BY bs.script_group
						ORDER BY total_loss DESC
					  LIMIT ?`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []types.ScriptGroupSummary
	for rows.Next() {
		var summary types.ScriptGroupSummary
		err = rows.Scan(&summary.ScriptGroup, &summary.TotalLoss, &summary.Scripts, &summary.Transactions)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}
