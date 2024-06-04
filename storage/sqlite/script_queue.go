package sqlite

import (
	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetScriptQueue() ([]types.ScriptQueue, error) {
	rows, err := d.db.Query("SELECT * FROM script_queue")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scripts []types.ScriptQueue
	for rows.Next() {
		var script types.ScriptQueue
		err := rows.Scan(&script.ID, &script.Script, &script.TryCount, &script.CreatedAt, &script.UpdatedAt)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, script)
	}

	return scripts, nil
}

func (d *SqliteBackend) IncrementScriptQueueTryCount(id int) error {
	_, err := d.db.Exec("UPDATE script_queue SET try_count = try_count + 1 WHERE id = ?", id)
	return err
}

func (d *SqliteBackend) RecordScriptUnspents(script types.ScriptQueue, unspentTxids []string, unspentHeights []int64) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	for i, txid := range unspentTxids {
		_, err = tx.Exec("INSERT INTO transaction_queue (txid, block_height) VALUES (?, ?)", txid, unspentHeights[i])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec("DELETE FROM script_queue WHERE id = ?", script.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
