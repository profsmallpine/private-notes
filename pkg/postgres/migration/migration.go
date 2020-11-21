package migration

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Migration struct {
	DB       *gorm.DB
	Executor func(*gorm.DB) error
	Key      string
}

type MigrationResult struct {
	Error error
}

func (m Migration) execute() MigrationResult {
	var err error

	// Start transaction
	tx := m.DB.Begin()
	if tx.Error != nil {
		return MigrationResult{Error: err}
	}

	// Run migration logic
	err = m.Executor(tx)
	if err != nil {
		tx.Rollback()
		return MigrationResult{Error: err}
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return MigrationResult{Error: err}
	}

	// Return result
	return MigrationResult{Error: err}
}

func Up(db *gorm.DB) error {
	// Ensure schema
	ensurePublicSchema(db)

	// Ensure migrations table exists
	ensureMigrationsTable(db)

	// Fetch correct list of migrations for each DB
	migrations := migrationsList(db)

	// Run migrations
	migrationsToRun := determineMigrationsToRun(db, migrations)
	for _, m := range migrationsToRun {
		if result := m.execute(); result.Error != nil {
			panic(result.Error)
		}

		// There was no error, so create a record for the migration
		createMigrationRecord(db, m.Key)
	}

	return nil
}

func ensurePublicSchema(db *gorm.DB) {
	err := db.Exec("CREATE SCHEMA IF NOT EXISTS public;").Error
	if err != nil {
		panic(fmt.Sprintf("Error creating public schema. Cannot continue: %s", err))
	}
}

func ensureMigrationsTable(db *gorm.DB) {
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			ran_at bigint,
			key text,
			CONSTRAINT migrations_key UNIQUE (key)
		)
	`).Error
	if err != nil {
		panic(fmt.Sprintf("Error creating migrations table. Cannot continue: %s", err))
	}
}

func createMigrationRecord(db *gorm.DB, key string) {
	err := db.Exec(`INSERT INTO migrations (key, ran_at) VALUES (?, ?)`, key, time.Now().Unix()).Error
	if err != nil {
		panic(fmt.Sprintf("Error creating migration. Cannot continue: %s", err))
	}
}

type migrationKeyCol struct {
	Key string
}

func determineMigrationsToRun(db *gorm.DB, allMigrations []Migration) []Migration {
	ranMigrations := []migrationKeyCol{}
	r := db.Raw("SELECT key FROM migrations;")
	if r.Error != nil {
		panic(fmt.Sprintf("Error fetching ran migrations. Cannot continue: %s", r.Error))
	}
	r.Scan(&ranMigrations)

	// If key is empty, we haven't ran *any* migrations
	if len(ranMigrations) == 0 {
		return allMigrations
	}

	// Compare ran migration keys to all migration keys to determine which need to run
	migrationsToRun := []Migration{}
	for _, migrationToCheck := range allMigrations {
		itsBeenRun := false
		// If this migration has already ran, continue to next iteration (don't add to list to run)
		for _, ranMigration := range ranMigrations {
			if migrationToCheck.Key == ranMigration.Key {
				itsBeenRun = true
				continue
			}
		}

		if !itsBeenRun {
			migrationsToRun = append(migrationsToRun, migrationToCheck)
		}
	}

	return migrationsToRun
}

func migrationsList(db *gorm.DB) []Migration {
	// List migrations in order
	return []Migration{
		{DB: db, Executor: CreateGroupsTable, Key: "20201109_create_groups"},
		{DB: db, Executor: CreateUsersTable, Key: "20201109_create_users"},
		{DB: db, Executor: CreateNotesTable, Key: "20201110_create_notes"},
		{DB: db, Executor: CreateCommentsTable, Key: "20201113_create_comments"},
		// -- END LIST -- (do not edit, used for autogeneration)
	}
}
