package migrations

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gu-io/midash/pkg/db"
)

// Migrations defines a table which will be used to run migrations for a given
// database.
type Migrations struct {
	database string
	db       db.DB
	ml       sync.Mutex
	tables   []TableMigration
}

// New returns a new instance of Migration.
func New(database string, db db.DB) *Migrations {
	return &Migrations{
		database: database,
		db:       db,
	}
}

// Migrate takes the individual query supplied and attempts to
// execute them returning any error found.
func (migration *Migrations) Migrate() error {
	migration.ml.Lock()
	defer migration.ml.Unlock()

	if migration.db == nil {
		return errors.New("No db provided available")
	}

	dbi, err := migration.db.New()
	if err != nil {
		return err
	}

	for _, table := range migration.tables {
		if _, err := dbi.Exec(table.String()); err != nil {
			return err
		}
	}

	return nil
}

// Migration adds the giving TableMigration into the giving migrations query.
func (migration *Migrations) Migration(table TableMigration) {
	migration.ml.Lock()
	migration.tables = append(migration.tables, table)
	migration.ml.Unlock()
}

// Database returns the data to use for migrations.
func (migration *Migrations) Database() string {
	return migration.database
}

// String returns the index query for the giving table migration.
func (migration *Migrations) String() string {
	var b bytes.Buffer

	if migration.database != "" {
		fmt.Fprintf(&b, "use %s;\r\n\n", migration.database)
	}

	migration.ml.Lock()
	for _, table := range migration.tables {
		fmt.Fprint(&b, table.String())
	}
	migration.ml.Unlock()

	return b.String()
}

// TableMigration defines a struct which defines a query field to be run against.
type TableMigration struct {
	TableName   string           `json:"table_name"`
	Timestamped bool             `json:"timestamped"`
	Fields      []FieldMigration `json:"fields"`
	Indexes     []IndexMigration `json:"indexes"`
	Queries     []string         `json:"queries"` // complete sql queries which will be ran.
}

// String returns the index query for the giving table migration.
func (table TableMigration) String() string {
	var b bytes.Buffer

	if table.TableName != "" {
		fmt.Fprintf(&b, "CREATE TABLE IF NOT EXISTS %s (\n", table.TableName)

		total := len(table.Fields) - 1
		hasFields := table.Fields != nil

		for index, field := range table.Fields {
			fmt.Fprintf(&b, "\t%s", field.String())

			if index < total {
				fmt.Fprintf(&b, ",\n")
			}
		}

		if table.Timestamped {
			if hasFields {
				fmt.Fprint(&b, ",\n")
			}

			fmt.Fprint(&b, "\tcreated_at timestamp NOT NULL,\n")
			fmt.Fprint(&b, "\tupdated_at timestamp NOT NULL")

			if table.Indexes != nil && table.Fields != nil {
				fmt.Fprint(&b, ",")
			}
		}

		if table.Indexes != nil && table.Fields != nil {
			if !table.Timestamped {
				fmt.Fprint(&b, ",")
			}

			if hasFields {
				fmt.Fprint(&b, "\n")
			}

			totalIndexes := len(table.Indexes) - 1

			for index, ind := range table.Indexes {
				fmt.Fprintf(&b, "\t%s", ind.String())

				if index < totalIndexes {
					fmt.Fprintf(&b, ",\n")
				}
			}
		}

		fmt.Fprint(&b, "\n);\r\n")
	}

	if table.Queries != nil {
		fmt.Fprint(&b, "\r\n")

		for _, query := range table.Queries {
			// Attempt to swap in tablename incase of format string
			query = fmt.Sprintf(query, table.TableName)

			if strings.HasSuffix(query, ";") {
				fmt.Fprintf(&b, "%s\n", query)
				continue
			}

			fmt.Fprintf(&b, "%s;\n", query)
		}

		fmt.Fprint(&b, "\r\n")
	}

	return b.String()
}

// FieldMigration defines a struct which defines the fields for a tableMigrations.
type FieldMigration struct {
	FieldName     string `json:"field_name"`
	FieldType     string `json:"field_type"`
	NotNull       bool   `json:"not_null"`
	PrimaryKey    bool   `json:"primary_key"`
	AutoIncrement bool   `json:"auto_increment"`
}

// String returns the index query for the giving field migration.
func (field FieldMigration) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "%s", field.FieldName)
	fmt.Fprintf(&b, " ")
	fmt.Fprintf(&b, "%s", field.FieldType)

	if field.NotNull {
		fmt.Fprintf(&b, " ")
		fmt.Fprintf(&b, "NOT NULL")
	}

	if field.PrimaryKey {
		fmt.Fprintf(&b, " ")
		fmt.Fprintf(&b, "PRIMARY KEY")
	}

	if field.AutoIncrement {
		fmt.Fprintf(&b, " ")
		fmt.Fprintf(&b, "AUTO_INCREMENT")
	}

	return b.String()
}

// IndexMigration defines a struct which contain index instruction for creating index fields
// for tables
type IndexMigration struct {
	Field     string `json:"field"`
	IndexName string `json:"index_name"`
}

// String returns the index query for the giving index migration.
func (index IndexMigration) String() string {
	return fmt.Sprintf("INDEX %s (%s)", index.Field, index.IndexName)
}

// Register adds the giving table and query into the tables list.
func (m *Migrations) Register(tableName string, query string) {

}
