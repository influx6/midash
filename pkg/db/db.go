package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
	"github.com/jmoiron/sqlx"
)

// contains templates of sql statement for use in operations.
const (
	countTemplate         = "SELECT %s FROM %s"
	selectAllTemplate     = "SELECT * FROM %s ORDER BY %s %s"
	selectLimitedTemplate = "SELECT * FROM %s ORDER BY %s %s LIMIT %d %d"
	selectItemTemplate    = "SELECT * FROM %s WHERE %s=?"
	insertTemplate        = "INSERT INTO %s %s VALUES %s"
	updateTemplate        = "UPDATE %s SET %s WHERE %s=%s"
	deleteTemplate        = "DELETE FROM %s WHERE %s=?"
)

// DB defines an interface which exposes a method to return a new
// underline sql.Database.
type DB interface {
	New() (*sqlx.DB, error)
}

// TableIdentity defines an interface which exposes a method returning table name
// associated with the giving implementing structure.
type TableIdentity interface {
	Table() string
}

// TableFields defines an interface which exposes method to return a map of all data
// associated with the defined structure.
type TableFields interface {
	TableIdentity
	Fields() map[string]interface{}
}

// TableConsumer defines an interface which accepts a map of data which will be loaded
// into the giving implementing structure.
type TableConsumer interface {
	WithFields(map[string]interface{}) error
}

// Save takes the giving table name with the giving fields and attempts to save this giving
// data appropriately into the giving db.
func Save(log sink.Sink, db *sqlx.DB, table TableFields) error {
	defer log.Emit(sinks.Info("Save to DB").With("table", table.Table()).Trace("db.Save").End())

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	fields := table.Fields()
	fieldNames := FieldNames(fields)
	values := FieldValues(fieldNames, fields)

	fieldNames = append(fieldNames, "created_at")
	fieldNames = append(fieldNames, "updated_at")

	values = append(values, time.Now())
	values = append(values, time.Now())

	query := fmt.Sprintf(insertTemplate, table.Table(), FieldNameMarkers(fieldNames), FieldMarkers(len(fieldNames)))
	log.Emit(sinks.Info("DB:Query").With("query", query))

	if _, err := db.Exec(query, values...); err != nil {
		log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return err
	}

	return tx.Commit()
}

// Update takes the giving table name with the giving fields and attempts to update this giving
// data appropriately into the giving db.
// index - defines the string which should identify the key to be retrieved from the fields to target the
// data to be updated in the db.
func Update(log sink.Sink, db *sqlx.DB, table TableFields, index string) error {
	defer log.Emit(sinks.Info("Update to DB").With("table", table.Table()).Trace("db.Update").End())

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	tableFields := table.Fields()
	tableFields["updated_at"] = time.Now()

	// Given index was not found, return error.
	indexValue, ok := tableFields[index]
	if !ok {
		return errors.New("Index key not found in fields")
	}

	// Delete given index from fieldNames
	delete(tableFields, index)

	fieldNames := FieldNamesFromMap(tableFields)

	query := fmt.Sprintf(updateTemplate, table.Table(), FieldMarkers(len(fieldNames)), index, indexValue)
	log.Emit(sinks.Info("DB:Query").With("query", query))

	if _, err := db.Exec(query, FieldValues(fieldNames, tableFields)...); err != nil {
		log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return err
	}

	return tx.Commit()
}

// GetAllPerPage retrieves the giving data from the specific db with the specific index and value.
func GetAllPerPage(log sink.Sink, db *sqlx.DB, table TableIdentity, order string, orderBy string, page int, responserPerPage int) ([]map[string]interface{}, int, error) {
	defer log.Emit(sinks.Info("Retrieve all records from DB").With("table", table.Table()).WithFields(sink.Fields{
		"order":            order,
		"page":             page,
		"responserPerPage": responserPerPage,
	}).Trace("db.GetAll").End())

	if page < 0 && responserPerPage < 0 {
		records, err := GetAll(log, db, table, order, orderBy)
		return records, len(records), err
	}

	// Get total number of records.
	totalRecords, err := Count(log, db, table, "public_id")
	if err != nil {
		return nil, 0, err
	}

	switch strings.ToLower(order) {
	case "asc":
		order = "ASC"
	case "dsc", "desc":
		order = "DESC"
	default:
		order = "ASC"
	}

	var fields []map[string]interface{}

	var totalWanted, indexToStart int

	if page < 0 && responserPerPage > 0 {
		totalWanted = responserPerPage
		indexToStart = 0
	} else {
		totalWanted = responserPerPage * page
		indexToStart = totalWanted / 2

		if page > 1 {
			indexToStart++
		}
	}

	// If we are passed the total, just return nil records and total without error.
	if indexToStart > totalRecords {
		return nil, totalRecords, nil
	}

	query := fmt.Sprint(selectAllTemplate, table, orderBy, order, indexToStart, totalWanted)
	log.Emit(sinks.Info("DB:Query").With("query", query))

	if err := db.Select(&fields, query); err != nil {
		log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))

		return nil, totalRecords, err
	}

	return fields, totalRecords, nil
}

// GetAll retrieves the giving data from the specific db with the specific index and value.
func GetAll(log sink.Sink, db *sqlx.DB, table TableIdentity, order string, orderBy string) ([]map[string]interface{}, error) {
	defer log.Emit(sinks.Info("Retrieve all records from DB").With("table", table.Table()).Trace("db.GetAll").End())

	switch strings.ToLower(order) {
	case "asc":
		order = "ASC"
	case "dsc", "desc":
		order = "DESC"
	default:
		order = "ASC"
	}

	var fields []map[string]interface{}

	query := fmt.Sprint(selectAllTemplate, table, orderBy, order)
	log.Emit(sinks.Info("DB:Query").With("query", query))

	if err := db.Select(&fields, query); err != nil {
		log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return nil, err
	}

	return fields, nil
}

// Get retrieves the giving data from the specific db with the specific index and value.
func Get(log sink.Sink, db *sqlx.DB, table TableIdentity, consumer TableConsumer, index string, indexValue interface{}) error {
	defer log.Emit(sinks.Info("Get record from DB").WithFields(sink.Fields{
		"table":      table.Table(),
		"index":      index,
		"indexValue": indexValue,
	}).Trace("db.Get").End())

	var fields map[string]interface{}

	query := fmt.Sprintf(selectItemTemplate, table.Table(), index)
	log.Emit(sinks.Info("DB:Query").With("query", query))

	if err := db.Get(&fields, query, indexValue); err != nil {
		log.Emit(sinks.Error("DB:Query").WithFields(sink.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return err
	}

	return nil
}

// Count retrieves the total number of records from the specific table from the db.
func Count(log sink.Sink, db *sqlx.DB, table TableIdentity, index string) (int, error) {
	defer log.Emit(sinks.Info("Count record from DB").WithFields(sink.Fields{
		"table": table.Table(),
	}).Trace("db.Get").End())

	var records []int

	query := fmt.Sprintf(countTemplate, index, table.Table())
	log.Emit(sinks.Info("DB:Query").With("query", query))

	if err := db.Get(&records, query); err != nil {
		log.Emit(sinks.Error("DB:Query").WithFields(sink.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
		return 0, err
	}

	return len(records), nil
}

// Delete removes the giving data from the specific db with the specific index and value.
func Delete(log sink.Sink, db *sqlx.DB, table TableIdentity, index string, indexValue interface{}) error {
	defer log.Emit(sinks.Info("Delete record from DB").WithFields(sink.Fields{
		"table":      table.Table(),
		"index":      index,
		"indexValue": indexValue,
	}).Trace("db.GetAll").End())

	var fields map[string]interface{}

	query := fmt.Sprintf(deleteTemplate, table.Table(), index)
	log.Emit(sinks.Info("DB:Query").With("query", query))

	if err := db.Get(&fields, query, indexValue); err != nil {
		log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"err":   err,
			"query": query,
			"table": table.Table(),
		}))
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
func FieldValues(names []string, fields map[string]interface{}) []interface{} {
	var vals []interface{}

	for _, name := range names {
		vals = append(vals, fields[name])
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
func FieldNames(fields map[string]interface{}) []string {
	var names []string

	for key := range fields {
		names = append(names, key)
	}

	return names
}
