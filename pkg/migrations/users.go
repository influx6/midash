package migrations

import (
	"github.com/gu-io/midash/pkg/db/sql/tables"
)

// Users defines the migration table for creating the user's table.
var Users = tables.TableMigration{
	TableName:   "users",
	Timestamped: true,
	Indexes:     []tables.IndexMigration{},
	Fields: []tables.FieldMigration{
		{
			FieldName: "email",
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
			FieldName: "private_id",
			FieldType: "VARCHAR(255)",
			NotNull:   true,
		},
		{
			FieldName: "hash",
			FieldType: "VARCHAR(255)",
			NotNull:   true,
		},
	},
}
