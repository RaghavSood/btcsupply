package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetOpReturnSummaries(limit int) ([]types.OpReturnSummary, error) {
	query := `SELECT
	  					burn_script,
							COUNT(DISTINCT(tx_id)) as transactions,
							SUM(amount) as total_loss
						FROM losses WHERE burn_script LIKE '6a%'
						GROUP BY burn_script ORDER BY total_loss DESC LIMIT ?`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	summaries, err := scanOpReturnSummaries(rows)
	return summaries, err
}

func (d *SqliteBackend) GetOpReturnSummary(script string) (types.OpReturnSummary, error) {
	query := `SELECT
							burn_script,
							COUNT(DISTINCT(tx_id)) as transactions,
							SUM(amount) as total_loss
						FROM losses WHERE burn_script = ?`

	var summary types.OpReturnSummary
	err := d.db.QueryRow(query, script).Scan(&summary.Script, &summary.Transactions, &summary.TotalLoss)
	if err != nil {
		return types.OpReturnSummary{}, err
	}

	return summary, nil
}

func scanOpReturnSummaries(rows *sql.Rows) ([]types.OpReturnSummary, error) {
	var summaries []types.OpReturnSummary
	for rows.Next() {
		var summary types.OpReturnSummary
		err := rows.Scan(&summary.Script, &summary.Transactions, &summary.TotalLoss)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
