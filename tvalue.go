package thrift

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
)

type TValue interface {
	Equals(other *TValue) bool
	// WriteFieldData outputs TValue to Thrift protocol with field ID `fid` and field name `fname`.
	// This is equivalent to `<field-data>` in [Thrift protocol spec].
	//
	//	<field> ::= <field-begin> <field-data> <field-end>
	//	<field-data> ::= I8 | I16 | I32 | I64 | DOUBLE | STRING | BINARY
	//			<struct> | <map> | <list> | <set>
	//
	// [Thrift protocol spec]: https://github.com/apache/thrift/blob/1a31d9051d35b732a5fce258955ef95f576694ba/doc/specs/thrift-protocol-spec.md (v0.21.0)
	WriteFieldData(cxt context.Context, oprot thrift.TProtocol) error
	// TType returns type in Thrift.
	TType() thrift.TType
}
