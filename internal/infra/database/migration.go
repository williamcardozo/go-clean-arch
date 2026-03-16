package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations(db *sql.DB, migrationsPath string) error {
	// Cria tabela de controle de migrações
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Lê arquivos de migração
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Ordena arquivos
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Executa migrações
	for _, filename := range migrationFiles {
		// Verifica se já foi aplicada
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", filename).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", filename, err)
		}

		if count > 0 {
			log.Printf("Migration %s already applied, skipping", filename)
			continue
		}

		// Lê conteúdo da migração
		content, err := os.ReadFile(filepath.Join(migrationsPath, filename))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		// Executa migração
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Registra migração aplicada
		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", filename)
		if err != nil {
			return fmt.Errorf("failed to register migration %s: %w", filename, err)
		}

		log.Printf("Migration %s applied successfully", filename)
	}

	return nil
}
