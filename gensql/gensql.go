package gensql

import (
	"fmt"

	gensql "github.com/cc-systems/sql-proto/gensql/proto"
	"github.com/cc-systems/sql-proto/gensql/sql"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func GenerateFile(gen *protogen.Plugin, f *protogen.File) error {
	generatedFile := gen.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.sql", f.GoImportPath)

	tableLookup := make(map[string]*sql.Table)

	for _, m := range f.Messages {
		primaryKey, err := getPrimary(m.Desc)
		if err != nil {
			return err
		}
		var cols []*sql.Column
		m.Desc.Fields()
		for _, field := range m.Fields {
			fieldName := string(field.Desc.Name())

			// nested message -> needs join
			if joinedMessage := field.Desc.Message(); joinedMessage != nil {
				// 1:n reference
				if field.Desc.IsList() {
					if foreignTable, ok := tableLookup[string(joinedMessage.Name())]; ok {
						for _, primaryKeyPart := range primaryKey {
							foreignTable.Columns = append(foreignTable.Columns, primaryKeyPart.PrefixName(string(m.Desc.Name())))
						}
					} else {
						return fmt.Errorf("could not resolve foreign table for message %s", joinedMessage.Name())
					}
					continue
				}

				// 1:1 reference
				// build foreign keys
				foreignKey, err := getPrimary(joinedMessage)
				if err != nil {
					return err
				}

				for _, foreignKeyPart := range foreignKey {
					cols = append(cols, foreignKeyPart.PrefixName(string(joinedMessage.Name())))
				}

				continue
			}

			fieldType, err := sql.ProtoTypeToSQL(field.Desc.Kind())
			if err != nil {
				return err
			}

			cols = append(cols, &sql.Column{
				Name:  fieldName,
				Type:  fieldType,
				Null:  field.Desc.HasOptionalKeyword(),
				Array: field.Desc.IsList(),
			})
		}

		table := &sql.Table{
			Name:    string(m.Desc.Name()),
			Columns: cols,
			Constraints: []sql.Schemer{
				&sql.PrimaryKeyConstraint{
					Columns: primaryKey,
				},
			},
		}

		tableLookup[table.Name] = table
	}

	for _, table := range tableLookup {
		generatedFile.P(table.Schema())
	}

	return nil
}

func getPrimary(m protoreflect.MessageDescriptor) ([]*sql.Column, error) {
	var primaryKey []*sql.Column

	fields := m.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if proto.GetExtension(field.Options(), gensql.E_Primary).(bool) {
			col, err := fieldToColumn(field)
			if err != nil {
				return nil, err
			}

			if col.Null {
				return nil, fmt.Errorf("primary key '%s' can not be null", col.Name)
			}

			primaryKey = append(primaryKey, col)
		}
	}

	return primaryKey, nil
}

func fieldToColumn(field protoreflect.FieldDescriptor) (*sql.Column, error) {
	fieldName := string(field.Name())
	fieldType, err := sql.ProtoTypeToSQL(field.Kind())
	if err != nil {
		return nil, err
	}

	return &sql.Column{
		Name: fieldName,
		Type: fieldType,
		Null: field.HasOptionalKeyword(),
	}, nil
}
