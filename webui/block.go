package webui

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/RaghavSood/btcsupply/types"
	"github.com/RaghavSood/btcsupply/util"
)

func (w *WebUI) getBlockOrFutureBlock(identifier string) (types.Block, error) {
	block, err := w.db.GetBlock(identifier)
	if err != nil && err != sql.ErrNoRows {
		return types.Block{}, fmt.Errorf("Failed to get block: %w", err)
	}

	if err == sql.ErrNoRows {
		latestBlock, err := w.db.GetLatestBlock()
		if err != nil {
			return types.Block{}, fmt.Errorf("Failed to get latest block: %w", err)
		}

		int64Identifier, err := strconv.ParseInt(identifier, 10, 64)
		if err != nil {
			return types.Block{}, fmt.Errorf("Failed to parse identifier: %w", err)
		}

		block = util.FutureBlock(int64Identifier, latestBlock.BlockHeight)
	}

	return block, nil
}
