package thrift

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
)

type TValue interface {
	Equals(other *TValue) bool
	// WriteField outputs TValue to Thrift protocol with field ID `fid` and field name `fname`.
	// This must start with `WriteFieldBegin` method call, and end with `WriteFieldEnd` method call.
	WriteField(cxt context.Context, oprot thrift.TProtocol, fid int16, fname string) error
}
