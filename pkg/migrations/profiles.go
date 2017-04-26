package migrations

import "github.com/gu-io/midash/pkg/db/sql/tables"

// Profiles defines the migration table for creating the profiles's table.
var Profiles = tables.TableMigration{
	TableName:   "profiles",
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
			FieldName: "address",
			FieldType: "text",
			NotNull:   true,
		},
		{
			FieldName:  "public_id",
			FieldType:  "VARCHAR(255)",
			PrimaryKey: true,
			NotNull:    true,
		},
		{
			FieldName: "first_name",
			FieldType: "VARCHAR(255)",
			NotNull:   true,
		},
		{
			FieldName: "last_name",
			FieldType: "VARCHAR(255)",
			NotNull:   true,
		},
	},
}
