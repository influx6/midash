package migrations

import "github.com/gu-io/midash/pkg/db/sql/tables"

// Sessions defines the migration table for creating the session's table.
var Sessions = tables.TableMigration{
	TableName:   "sessions",
	Timestamped: true,
	Indexes: []tables.IndexMigration{
		{
			IndexName: "user_id",
			Field:     "user_id",
		},
	},
	Fields: []tables.FieldMigration{
		{
			FieldName: "user_id",
			FieldType: "VARCHAR(255)",
			NotNull:   true,
		},
		{
			FieldName: "token",
			FieldType: "VARCHAR(255)",
			NotNull:   true,
		},
		{
			FieldName:  "public_id",
			FieldType:  "VARCHAR(255)",
			PrimaryKey: true,
			NotNull:    true,
		},
		{
			FieldName: "expires",
			FieldType: "timestamp",
			NotNull:   true,
		},
	},
}
