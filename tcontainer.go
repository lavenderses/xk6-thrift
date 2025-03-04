package thrift

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
)

func ReadContainerData(ttype thrift.TType, cxt context.Context, iprot thrift.TProtocol) (TValue, error) {
	var tv TValue
	var err error

	switch ttype {
	case thrift.STRING:
		tv, err = ReadString(cxt, iprot)
	case thrift.BOOL:
		tv, err = ReadBool(cxt, iprot)
	case thrift.MAP:
		tv, err = ReadMap(cxt, iprot)
	// case thrift.STRUCT:
	// 	tv, err = ReadStruct(cxt, iprot)
	default:
		err = iprot.Skip(cxt, ttype)
	}
	if err != nil {
		return nil, err
	}

	return tv, err
}
