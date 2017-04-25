package migrations_test

import (
	"testing"

	"github.com/gu-io/midash/pkg/db/migrations"
	"github.com/influx6/faux/tests"
)

var expected = "use buba;\r\n\nCREATE TABLE IF NOT EXISTS users (\n\tid INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,\n\tname VARCHAR(255) NOT NULL,\n\tcreated_at timestamp NOT NULL,\n\tupdated_at timestamp NOT NULL,\n\tINDEX name (name_idx)\n);\r\n\r\nALTER TABLE users ADD bottle VARCHAR(255) NOT NULL;\nALTER TABLE users ADD wine_id VARCHAR(255) NOT NULL;\n\r\n"

func TestMigration(t *testing.T) {
	migration := migrations.New("buba", nil)

	migration.Migration(migrations.TableMigration{
		TableName:   "users",
		Timestamped: true,
		Indexes: []migrations.IndexMigration{
			{
				IndexName: "name_idx",
				Field:     "name",
			},
		},
		Fields: []migrations.FieldMigration{
			{
				FieldName:     "id",
				FieldType:     "INTEGER",
				NotNull:       true,
				PrimaryKey:    true,
				AutoIncrement: true,
			},
			{
				FieldName: "name",
				FieldType: "VARCHAR(255)",
				NotNull:   true,
			},
		},
		Queries: []string{
			"ALTER TABLE %s ADD bottle VARCHAR(255) NOT NULL;",
			"ALTER TABLE %s ADD wine_id VARCHAR(255) NOT NULL",
		},
	})

	if migration.String() != expected {
		tests.Info("Expected: %q", expected)
		tests.Info("Recieved: %q", migration.String())
		tests.Failed("Should have successfully matched expected query")
	}
	tests.Passed("Should have successfully matched expected query")
}
