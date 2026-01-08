package migrations

import (
	"embed"
)

// FS contains all SQL migration files embedded at compile time.
// The path is relative to the module root since we use go:generate to copy files.
//
//go:embed sql/*.sql
var FS embed.FS

// Dir is the directory name within the embedded filesystem.
const Dir = "sql"
