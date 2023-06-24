// Code generated by ent, DO NOT EDIT.

package template

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the template type in the database.
	Label = "template"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldData holds the string denoting the data field in the database.
	FieldData = "data"
	// FieldAvgDistance holds the string denoting the avg_distance field in the database.
	FieldAvgDistance = "avg_distance"
	// FieldDiffDistance holds the string denoting the diff_distance field in the database.
	FieldDiffDistance = "diff_distance"
	// FieldPerceptionDistance holds the string denoting the perception_distance field in the database.
	FieldPerceptionDistance = "perception_distance"
	// Table holds the table name of the template in the database.
	Table = "templates"
)

// Columns holds all SQL columns for template fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldData,
	FieldAvgDistance,
	FieldDiffDistance,
	FieldPerceptionDistance,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
)

// OrderOption defines the ordering options for the Template queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByAvgDistance orders the results by the avg_distance field.
func ByAvgDistance(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAvgDistance, opts...).ToFunc()
}

// ByDiffDistance orders the results by the diff_distance field.
func ByDiffDistance(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDiffDistance, opts...).ToFunc()
}

// ByPerceptionDistance orders the results by the perception_distance field.
func ByPerceptionDistance(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPerceptionDistance, opts...).ToFunc()
}
