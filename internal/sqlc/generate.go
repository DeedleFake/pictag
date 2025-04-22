//go:generate go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate
package sqlc

import (
	"context"
	_ "embed"
)

//go:embed schema.sql
var schema string

// Init initializes a database with the schema.
func Init(ctx context.Context, db DBTX) error {
	_, err := db.ExecContext(ctx, schema)
	return err
}
