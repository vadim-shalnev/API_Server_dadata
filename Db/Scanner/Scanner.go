package Scanner

import (
	"reflect"
	"strings"
	"time"
)

const (
	AllFields = "all"
	Create    = "create"
	Update    = "update"
	Upsert    = "upsert"
	Conflict  = "conflict"
)

type Tabler interface {
	TableName() string
	OnCreate() []string
}

type TableUpdater interface {
	SetUpdatedAt(updatedAt time.Time) Tabler
}

type Table struct {
	Name            string
	Fields          []Field
	FieldsMap       map[string]Field
	Constraints     []Constraint
	OperationFields map[string][]string
	Entity          Tabler
}

type Scanner interface {
	RegisterTable(entities ...Tabler)
	OperationFields(tableName, operation string) []string
	Table(tableName string) Table
	Tables() map[string]Table
}

type TableScanner struct {
	tables map[string]Table
}

func NewTableScanner() Scanner {
	return &TableScanner{}
}

func (t *TableScanner) RegisterTable(entities ...Tabler) {
	tableEntities := make(map[string]Tabler, len(entities))
	t.tables = make(map[string]Table, len(entities))
	for i := range entities {
		tableEntities[entities[i].TableName()] = entities[i]
	}

	for name, entity := range tableEntities {
		table := Table{
			Name:            name,
			FieldsMap:       map[string]Field{},
			OperationFields: map[string][]string{},
			Entity:          entity,
		}
		reflected := reflect.TypeOf(entity).Elem()

		for i := 0; i < reflected.NumField(); i++ {
			// Get the field, returns https://golang.org/pkg/reflect/#StructField
			structField := reflected.Field(i)
			// Get the structField tag value
			fieldName := structField.Tag.Get("db")

			if fieldName == "" || fieldName == "-" {
				continue
			}
			table.OperationFields[AllFields] = append(table.OperationFields[AllFields], fieldName)

			field := Field{
				Name:    fieldName,
				Type:    structField.Tag.Get("db_type"),
				Default: structField.Tag.Get("db_default"),
				Table:   &table,
			}
			constraintRaw := structField.Tag.Get("db_index")
			constraintPieces := strings.Split(constraintRaw, ",")
			if len(constraintPieces) < 1 {
				field.Constraint = Constraint{}
			}
			if len(constraintPieces) > 0 {
				for i := range constraintPieces {
					switch constraintPieces[i] {
					case "index":
						field.Constraint.Index = true
					case "unique":
						field.Constraint.Unique = true
					}
				}
			}
			if field.Constraint.Index {
				field.Constraint.Field = &field
				table.Constraints = append(table.Constraints, field.Constraint)
			}
			table.Fields = append(table.Fields, field)
			table.FieldsMap[field.Name] = field

			opsRaw := structField.Tag.Get("db_ops")
			ops := strings.Split(opsRaw, ",")
			if opsRaw != "" {
				for j := range ops {
					table.OperationFields[ops[j]] = append(table.OperationFields[ops[j]], fieldName)
				}
			}
		}

		t.tables[name] = table
	}
}

func (t *TableScanner) OperationFields(tableName, operation string) []string {
	return t.tables[tableName].OperationFields[operation]
}

func (t *TableScanner) Table(tableName string) Table {
	return t.tables[tableName]
}

func (t *TableScanner) Tables() map[string]Table {
	return t.tables
}

type Field struct {
	Name       string
	Type       string
	Default    string
	Constraint Constraint
	Table      *Table
}

type Constraint struct {
	Index  bool
	Unique bool
	Field  *Field
}
