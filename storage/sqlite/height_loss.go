package sqlite

import (
	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetHeightLossSummary() ([]types.HeightLossSummary, error) {
	query := `SELECT 
  					    block_height,
  					    SUM(amount) OVER (ORDER BY block_height) AS cumulative_loss
  					FROM 
  					    (SELECT DISTINCT block_height, SUM(amount) AS amount FROM losses GROUP BY block_height) AS block_sums
  					ORDER BY 
  					    block_height;`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []types.HeightLossSummary
	for rows.Next() {
		var summary types.HeightLossSummary
		err = rows.Scan(&summary.BlockHeight, &summary.TotalLoss)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}
