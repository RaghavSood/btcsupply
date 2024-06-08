package sqlite

import (
	"encoding/csv"
	"fmt"
)

func (d *SqliteBackend) seedBurnScriptsFromCSV() error {
	csvFile, err := embeddedMigrations.Open("migrations/burn_scripts.csv")
	if err != nil {
		return fmt.Errorf("failed to open burn scripts CSV: %w", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read burn scripts CSV: %w", err)
	}

	tx, err := d.db.Begin()
	for _, record := range records {
		_, err = tx.Exec("INSERT INTO burn_scripts (script, confidence_level, provenance, script_group) VALUES (?, ?, ?, ?) ON CONFLICT (script) DO UPDATE SET confidence_level=excluded.confidence_level, provenance=excluded.provenance, script_group=excluded.script_group", record[0], record[1], record[2], record[3])
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert burn script: %w", err)
		}
	}

	var nonExistentGroups []string
	rows, err := tx.Query("SELECT DISTINCT script_group FROM burn_scripts WHERE script_group NOT IN (SELECT DISTINCT script_group FROM script_groups)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get non-existent script groups: %w", err)
	}

	for rows.Next() {
		var group string
		err = rows.Scan(&group)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to scan non-existent script group: %w", err)
		}

		nonExistentGroups = append(nonExistentGroups, group)
	}

	for _, group := range nonExistentGroups {
		_, err = tx.Exec("INSERT INTO script_groups (script_group) VALUES (?)", group)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert non-existent script group: %w", err)
		}
	}

	return tx.Commit()
}
