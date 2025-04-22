package sqlc

//go:generate go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"iter"
	"path"
	"slices"
	"unsafe"
)

var (
	//go:embed init.sql
	initSQL string

	//go:embed schema/*.sql
	migrationsFS embed.FS
)

type migration struct {
	name string
	sql  string
}

func migrations(skip []string) iter.Seq2[migration, error] {
	return func(yield func(migration, error) bool) {
		files, err := fs.ReadDir(migrationsFS, "schema")
		if err != nil {
			yield(migration{}, err)
			return
		}

		for _, file := range files {
			name := file.Name()
			if slices.Contains(skip, name) {
				continue
			}

			sql, err := fs.ReadFile(migrationsFS, path.Join("schema", name))
			if err != nil {
				yield(migration{}, err)
				return
			}

			m := migration{
				name: name,
				sql:  unsafe.String(unsafe.SliceData(sql), len(sql)),
			}
			if !yield(m, nil) {
				return
			}
		}
	}
}

func (m migration) run(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction for migration: %w", err)
	}

	_, err = tx.ExecContext(ctx, m.sql)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("execute migration %q: %w", m.name, err)
	}

	err = New(tx).addMigration(ctx, m.name)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("add migration %q: %w", m.name, err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("commit transaction for migration %q: %w", m.name, err)
	}

	return nil
}

func schema(ctx context.Context, db *sql.DB) ([]string, error) {
	_, err := db.ExecContext(ctx, initSQL)
	if err != nil {
		return nil, fmt.Errorf("create migrations table: %w", err)
	}

	return New(db).listMigrations(ctx)
}

func Migrate(ctx context.Context, db *sql.DB) error {
	skip, err := schema(ctx, db)
	if err != nil {
		return fmt.Errorf("get existing migrations: %w", err)
	}

	for m, err := range migrations(skip) {
		if err != nil {
			return err
		}

		err := m.run(ctx, db)
		if err != nil {
			return err
		}
	}

	return nil
}
