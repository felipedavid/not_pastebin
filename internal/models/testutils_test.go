package models

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open(
		"mysql",
		"test:test@tcp(localhost:3307)/snippetbox?collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true")
	require.NoError(t, err)

	script, err := os.ReadFile("./testdata/setup.sql")
	require.NoError(t, err)

	_, err = db.Exec(string(script))
	require.NoError(t, err)

	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		require.NoError(t, err)
		_, err = db.Exec(string(script))
		require.NoError(t, err)
		db.Close()
	})

	return db
}
