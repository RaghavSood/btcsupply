package sqlite

import (
	"fmt"

	btypes "github.com/RaghavSood/btcsupply/bitcoinrpc/types"
	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) RecordBlockIndexResults(block types.Block, txoutset types.TxOutSetInfo, blockstats btypes.BlockStats, losses []types.Loss, transactions []types.Transaction, spentTxids []string, spentVouts []int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Insert the block record into the block table
	_, err = tx.Exec("INSERT INTO blocks (block_height, block_hash, block_timestamp, parent_block_hash, num_transactions) VALUES (?, ?, ?, ?, ?) ON CONFLICT (block_hash) DO NOTHING", block.BlockHeight, block.BlockHash, block.BlockTimestamp, block.ParentBlockHash, block.NumTransactions)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert block record: %v", err)
	}

	// Insert the txoutset record into the txoutset table
	_, err = tx.Exec("INSERT INTO tx_out_set_info (height, bestblock, txouts, bogosize, muhash, total_amount, total_unspendable_amount, prevout_spent, coinbase, new_outputs_ex_coinbase, unspendable, genesis_block, bip30, scripts, unclaimed_rewards) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (bestblock) DO NOTHING", txoutset.Height, txoutset.Bestblock, txoutset.Txouts, txoutset.Bogosize, txoutset.Muhash, txoutset.TotalAmount, txoutset.TotalUnspendableAmount, txoutset.PrevoutSpent, txoutset.Coinbase, txoutset.NewOutputsExCoinbase, txoutset.Unspendable, txoutset.GenesisBlock, txoutset.Bip30, txoutset.Scripts, txoutset.UnclaimedRewards)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert txoutset record: %v", err)
	}

	// Insert the blockstats record into the blockstats table
	_, err = tx.Exec("INSERT INTO block_stats (avgfee, avgfeerate, avgtxsize, blockhash, height, ins, maxfee, maxfeerate, maxtxsize, medianfee, mediantime, mediantxsize, minfee, minfeerate, mintxsize, outs, subsidy, swtotal_size, swtotal_weight, swtxs, time, total_out, total_size, total_weight, totalfee, txs, utxo_increase, utxo_size_inc, utxo_increase_actual, utxo_size_inc_actual) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (blockhash) DO NOTHING", blockstats.Avgfee, blockstats.Avgfeerate, blockstats.Avgtxsize, blockstats.Blockhash, blockstats.Height, blockstats.Ins, blockstats.Maxfee, blockstats.Maxfeerate, blockstats.Maxtxsize, blockstats.Medianfee, blockstats.Mediantime, blockstats.Mediantxsize, blockstats.Minfee, blockstats.Minfeerate, blockstats.Mintxsize, blockstats.Outs, blockstats.Subsidy, blockstats.SwtotalSize, blockstats.SwtotalWeight, blockstats.Swtxs, blockstats.Time, blockstats.TotalOut, blockstats.TotalSize, blockstats.TotalWeight, blockstats.Totalfee, blockstats.Txs, blockstats.UtxoIncrease, blockstats.UtxoSizeInc, blockstats.UtxoIncreaseActual, blockstats.UtxoSizeIncActual)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert blockstats record: %v", err)
	}

	// Insert the losses records into the losses table
	for _, loss := range losses {
		_, err = tx.Exec("INSERT INTO losses (tx_id, block_hash, vout, amount) VALUES (?, ?, ?, ?)", loss.TxID, loss.BlockHash, loss.Vout, loss.Amount)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert loss record: %v", err)
		}
	}

	// Insert the transactions records into the transactions table
	for _, transaction := range transactions {
		_, err = tx.Exec("INSERT INTO transactions (tx_id, transaction_details) VALUES (?, ?)", transaction.TxID, transaction.TransactionDetails)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert transaction record: %v", err)
		}
	}

	// Remove any previously recorded outputs that were spent in this block
	for i := range spentTxids {
		_, err = tx.Exec("DELETE FROM losses WHERE tx_id = ? AND vout = ?", spentTxids[i], spentVouts[i])
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete spent output: %v", err)
		}
	}

	err = tx.Commit()

	return err
}

func (d *SqliteBackend) RecordTransactionIndexResults(losses []types.Loss, transactions []types.Transaction, spentTxids []string, spentVouts []int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Insert the losses records into the losses table
	for _, loss := range losses {
		_, err = tx.Exec("INSERT INTO losses (tx_id, block_hash, vout, amount) VALUES (?, ?, ?, ?) ON CONFLICT (tx_id, vout) DO NOTHING", loss.TxID, loss.BlockHash, loss.Vout, loss.Amount)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert loss record: %v", err)
		}
	}

	// Insert the transactions records into the transactions table
	for _, transaction := range transactions {
		_, err = tx.Exec("INSERT INTO transactions (tx_id, transaction_details) VALUES (?, ?) ON CONFLICT (tx_id) DO UPDATE SET transaction_details=excluded.transaction_details", transaction.TxID, transaction.TransactionDetails)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert transaction record: %v", err)
		}
	}

	// Remove any previously recorded outputs that were spent by this transaction
	for i := range spentTxids {
		_, err = tx.Exec("DELETE FROM losses WHERE tx_id = ? AND vout = ?", spentTxids[i], spentVouts[i])
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete spent output: %v", err)
		}
	}

	// Delete the transaction from the transaction queue
	_, err = tx.Exec("DELETE FROM transaction_queue WHERE txid = ?", transactions[0].TxID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete transaction from queue: %v", err)
	}

	return tx.Commit()
}
