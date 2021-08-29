package sql

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func ProtoTypeToSQL(kind protoreflect.Kind) (FieldType, error) {
	switch kind {
	case protoreflect.BoolKind:
		return SQLBoolean, nil
	case protoreflect.BytesKind:
		return SQLBytea, nil
	case protoreflect.DoubleKind:
		return SQLDouble, nil
	case protoreflect.Fixed32Kind:
		return SQLReal, nil
	case protoreflect.Fixed64Kind:
		return SQLDouble, nil
	case protoreflect.FloatKind:
		return SQLReal, nil
	case protoreflect.Int32Kind:
		return SQLInteger, nil
	case protoreflect.Int64Kind:
		return SQLBigInt, nil
	case protoreflect.StringKind:
		return SQLText, nil
	}

	return SQLUnknown, fmt.Errorf("can not convert %v to sql.FieldType", kind)
}
