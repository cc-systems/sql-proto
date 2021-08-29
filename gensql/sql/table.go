package sql

import (
	"fmt"
	"strings"
)

type Schemer interface {
	Schema() string
}

type Table struct {
	Name        string
	Columns     []*Column
	Constraints []Schemer
}

type Column struct {
	Type  FieldType
	Name  string
	Null  bool
	Array bool
}

type PrimaryKeyConstraint struct {
	Columns []*Column
}

type ForeignKeyConstraint struct {
	ForeignTableName string
	Columns          []*Column
}

func (t *Table) Schema() string {
	columnSchema := make([]string, len(t.Columns))
	for idx, c := range t.Columns {
		columnSchema[idx] = c.Schema()
	}

	constraintSchema := make([]string, len(t.Constraints))
	for idx, c := range t.Constraints {
		constraintSchema[idx] = c.Schema()
	}

	return fmt.Sprintf(`
		CREATE TABLE %s IF NOT EXISTS (
			%s,
			%s
		);
	`, t.Name, strings.Join(columnSchema, ",\n"), strings.Join(constraintSchema, ",\n"))
}

func (c *Column) Schema() string {
	arrayModifier := ""
	if c.Array {
		arrayModifier = "[]"
	}

	nullModifier := "null"
	if !c.Null {
		nullModifier = "not null"
	}

	return fmt.Sprintf("%s %s%s %s", c.Name, c.Type, arrayModifier, nullModifier)
}

func (c *Column) PrefixName(prefix string) *Column {
	return &Column{
		Name: prefix + "_" + c.Name,
		Type: c.Type,
		Null: c.Null,
	}
}

func (p *PrimaryKeyConstraint) Schema() string {

	columnNames := make([]string, len(p.Columns))
	for i, col := range p.Columns {
		columnNames[i] = col.Name
	}

	primaryKeySchema := ""
	if len(p.Columns) > 0 {
		primaryKeySchema = fmt.Sprintf("\nPRIMARY KEY (%s)", strings.Join(columnNames, ", "))
	}
	return primaryKeySchema
}

type FieldType string

// See https://www.postgresql.org/docs/current/datatype.html
const (
	SQLUnknown FieldType = ""

	SQLBoolean FieldType = "boolean"

	SQLSmallInt FieldType = "smallint"
	SQLInteger  FieldType = "integer"
	SQLBigInt   FieldType = "bigint"
	SQLReal     FieldType = "real"
	SQLDouble   FieldType = "double precision"

	SQLText FieldType = "text"

	SQLBytea FieldType = "bytea"

	SQLTimestampTimezone FieldType = "timestamp with time zone"
	SQLTimestamp         FieldType = "timestamp"
	SQLDate              FieldType = "date"
)
