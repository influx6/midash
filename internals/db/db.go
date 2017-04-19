package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// contains templates of sql statement for use in operations.
const (
	selectAllTemplate  = "SELECT * FROM %s"
	selectItemTemplate = "SELECT * FROM %s WHERE %s=?"
	insertTemplate     = "INSERT INTO %s %s VALUES %s"
	updateTemplate     = "UPDATE %s SET %s WHERE %s=%s"
)

// TableIdentity defines an interface which exposes a method returning table name
// associated with the giving implementing structure.
type TableIdentity interface {
	Table() string
}

// TableFields defines an interface which exposes method to return a map of all data
// associated with the defined structure.
type TableFields interface {
	Fields() map[string]interface{}
}

// TableConsumer defines an interface which accepts a map of data which will be loaded
// into the giving implementing structure.
type TableConsumer interface {
	WithFields(map[string]interface{}) error
}

// Save takes the giving table name with the giving fields and attempts to save this giving
// data appropriately into the giving db.
func Save(db *sqlx.DB, table TableIdentity, fields TableFields) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	fieldNames := FieldNames(fields)

	query := fmt.Sprintf(insertTemplate, table.Table(), FieldNameMarkers(fieldNames), FieldMarkers(len(fieldNames)))

	if _, err := db.Exec(query, FieldValues(fieldNames, fields)...); err != nil {
		return err
	}

	return tx.Commit()
}

// Update takes the giving table name with the giving fields and attempts to update this giving
// data appropriately into the giving db.
// index - defines the string which should identify the key to be retrieved from the fields to target the
// data to be updated in the db.
func Update(db *sqlx.DB, table TableIdentity, fields TableFields, index string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	tableFields := fields.Fields()
	indexValue, ok := tableFields[index]

	// Given index was not found, return error.
	if !ok {
		return errors.New("Index key not found in fields")
	}

	// Delete given index from fieldNames
	delete(tableFields, index)

	fieldNames := FieldNamesFromMap(tableFields)

	query := fmt.Sprintf(updateTemplate, table.Table(), FieldMarkers(len(fieldNames)), index, indexValue)
	if _, err := db.Exec(query, FieldValues(fieldNames, fields)...); err != nil {
		return err
	}

	return tx.Commit()
}

// GetAll retrieves the giving data from the specific db with the specific index and value.
func GetAll(db *sqlx.DB, table TableIdentity) ([]map[string]interface{}, error) {
	var fields []map[string]interface{}

	if err := db.Select(&fields, selectAllTemplate); err != nil {
		return nil, err
	}

	return fields, nil
}

// Get retrieves the giving data from the specific db with the specific index and value.
func Get(db *sqlx.DB, table TableIdentity, consumer TableConsumer, index string, indexValue interface{}) error {
	var fields map[string]interface{}

	query := fmt.Sprintf(selectItemTemplate, table.Table(), index)
	if err := db.Get(&fields, query, indexValue); err != nil {
		return err
	}

	return nil
}

// FieldMarkers returns a (?,...,>) string which represents
// all filedNames extrated from the provided TableField.
func FieldMarkers(total int) string {
	var markers []string

	for i := 0; i < total; i++ {
		markers = append(markers, "?")
	}

	return "(" + strings.Join(markers, ",") + ")"
}

// FieldNameMarkers returns a (fieldName,...,fieldName) string which represents
// all filedNames extrated from the provided TableField.
func FieldNameMarkers(fields []string) string {
	return "(" + strings.Join(fields, ", ") + ")"
}

// FieldValues returns a (fieldName,...,fieldName) string which represents
// all filedNames extrated from the provided TableField.
func FieldValues(names []string, fields TableFields) []interface{} {
	var vals []interface{}

	tableFields := fields.Fields()

	for _, name := range names {
		vals = append(vals, tableFields[name])
	}

	return vals
}

// FieldNamesFromMap returns a (fieldName,...,fieldName) string which represents
// all filedNames extrated from the provided TableField.
func FieldNamesFromMap(fields map[string]interface{}) []string {
	var names []string

	for key := range fields {
		names = append(names, key)
	}

	return names
}

// FieldNames returns a (fieldName,...,fieldName) string which represents
// all filedNames extrated from the provided TableField.
func FieldNames(fields TableFields) []string {
	var names []string

	for key := range fields.Fields() {
		names = append(names, key)
	}

	return names
}
