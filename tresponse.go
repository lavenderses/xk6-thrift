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
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(cxt)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}

		switch fieldId {
			case 0:
				switch fieldTypeId {
					case thrift.STRING:
						p.ReadString(cxt, iprot, fieldId)
					case thrift.BOOL:
						p.ReadBool(cxt, iprot, fieldId)
				}
			default:
				if err := iprot.Skip(cxt, fieldTypeId); err != nil {
					return err
				}
		}
	}

	return nil
}

func (p *TResponse) ReadString(cxt context.Context, iproto thrift.TProtocol, fieldId int16) error {
	v, err := iproto.ReadString(cxt)
	if err != nil {
		return thrift.PrependError(fmt.Sprintf("error reading string field %d: ", fieldId), err)
	}
	
	tv := NewTstring(v)
	p.values[fieldId] = tv
	return nil
}

func (p *TResponse) ReadBool(cxt context.Context, iproto thrift.TProtocol, fieldId int16) error {
	v, err := iproto.ReadBool(cxt)
	if err != nil {
		return thrift.PrependError(fmt.Sprintf("error reading boolean field %d: ", fieldId), err)
	}
	
	tv := NewTBool(v)
	p.values[fieldId] = tv
	return nil
}

// dummy.
func (p *TResponse) Write(cxt context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(cxt, "dummy"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := oprot.WriteStructEnd(cxt); err != nil {
		return thrift.PrependError(fmt.Sprintf("write struct end error: ", p), err)
	}
	return nil
}

func (p *TResponse) Add(key int16, value TValue) {
	p.values[key] = value
}
