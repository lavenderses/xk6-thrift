package thrift

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

type TResponse struct {
	values map[int16]TValue
}

func NewTResponse() *TResponse {
	return &TResponse{values: make(map[int16]TValue)}
}

func (p TResponse) Values() *map[int16]TValue {
	return &p.values
}

func (p *TResponse) Read(cxt context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(cxt); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		var err error
		var fieldTypeId thrift.TType
		var fieldId int16
		_, fieldTypeId, fieldId, err = iprot.ReadFieldBegin(cxt)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
		case 0:
			var v TValue
			v, err = ReadContainerData(fieldTypeId, cxt, iprot)
			if err != nil {
				return thrift.PrependError(fmt.Sprintf("%T read field (%d, %v) error: ", p, fieldId, fieldTypeId), err)
			}

			p.values[fieldId] = v
		default:
			err = iprot.Skip(cxt, fieldTypeId)
		}

		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T read field (%d, %v) error: ", p, fieldId, fieldTypeId), err)
		}
	}

	return nil
}

// dummy.
func (p *TResponse) Write(cxt context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(cxt, "dummy"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := oprot.WriteStructEnd(cxt); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
	}
	return nil
}

func (p *TResponse) Add(key int16, value TValue) {
	p.values[key] = value
}
