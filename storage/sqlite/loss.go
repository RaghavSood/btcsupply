package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetRecentLosses(limit int) ([]types.Loss, error) {
	rows, err := d.db.Query("SELECT * FROM losses ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}

	losses, err := scanLosses(rows)
	if err != nil {
		return nil, err
	}

	return losses, nil
}

func (d *SqliteBackend) GetBlockLosses(hash string) ([]types.Loss, error) {
	rows, err := d.db.Query("SELECT * FROM losses WHERE block_hash = ?", hash)
	if err != nil {
		return nil, err
	}

	losses, err := scanLosses(rows)
	if err != nil {
		return nil, err
	}

	return losses, nil
}

func (d *SqliteBackend) GetTransactionLosses(hash string) ([]types.Loss, error) {
	rows, err := d.db.Query("SELECT * FROM losses WHERE tx_id = ?", hash)
	if err != nil {
		return nil, err
	}

	losses, err := scanLosses(rows)
	if err != nil {
		return nil, err
	}

	return losses, nil
}

func scanLosses(rows *sql.Rows) ([]types.Loss, error) {
	var losses []types.Loss
	for rows.Next() {
		var loss types.Loss
		err := rows.Scan(&loss.ID, &loss.TxID, &loss.BlockHash, &loss.BlockHeight, &loss.Vout, &loss.Amount, &loss.CreatedAt, &loss.BurnScript)
		if err != nil {
			return nil, err
		}
		losses = append(losses, loss)
	}

	return losses, nil
}
