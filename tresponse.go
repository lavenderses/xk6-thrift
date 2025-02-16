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
			switch fieldTypeId {
			case thrift.STRING:
				err = p.ReadString(cxt, iprot, fieldId)
			case thrift.BOOL:
				err = p.ReadBool(cxt, iprot, fieldId)
			case thrift.MAP:
				err = p.ReadMap(cxt, iprot, fieldId)
			}
		default:
			err = iprot.Skip(cxt, fieldTypeId)
		}

		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T read field (%d, %v) error: ", p, fieldId, fieldTypeId), err)
		}
	}

	return nil
}

func (p *TResponse) ReadString(cxt context.Context, iproto thrift.TProtocol, fieldId int16) error {
	v, err := iproto.ReadString(cxt)
	if err != nil {
		return thrift.PrependError(fmt.Sprintf("error reading string field %d: ", fieldId), err)
	}

	p.values[fieldId] = NewTstring(v)
	return nil
}

func (p *TResponse) ReadBool(cxt context.Context, iproto thrift.TProtocol, fieldId int16) error {
	v, err := iproto.ReadBool(cxt)
	if err != nil {
		return thrift.PrependError(fmt.Sprintf("error reading boolean field %d: ", fieldId), err)
	}

	p.values[fieldId] = NewTBool(v)
	return nil
}

func (p *TResponse) ReadMap(cxt context.Context, iproto thrift.TProtocol, fieldId int16) error {
	keyType, valueType, size, err := iproto.ReadMapBegin(cxt)
	if err != nil {
		return thrift.PrependError(fmt.Sprintf("error reading map field %d: ", fieldId), err)
	}

	tmap := make(map[TValue]TValue)
	for i := 0; i < size; i++ {
		if err = p.readPair(cxt, iproto, &tmap, keyType, valueType); err != nil {
			return thrift.PrependError(fmt.Sprintf("error reading map %d: ", fieldId), err)
		}
	}

	p.values[fieldId] = NewTMap(&tmap)
	return nil
}

func (p *TResponse) readPair(cxt context.Context, iprot thrift.TProtocol, tmap *map[TValue]TValue, ktype, vtype thrift.TType) error {
	var key, value TValue
	var err error
	if key, err = p.readEntry(cxt, iprot, ktype); err != nil {
		return err
	}
	if value, err = p.readEntry(cxt, iprot, vtype); err != nil {
		return err
	}

	(*tmap)[key] = value
	return nil
}

func (p *TResponse) readEntry(cxt context.Context, iprot thrift.TProtocol, ttype thrift.TType) (tv TValue, err error) {
	switch ttype {
	case thrift.STRING:
		if v, err := iprot.ReadString(cxt); err != nil {
			return nil, err
		} else {
			tv = NewTstring(v)
		}
	case thrift.BOOL:
		if v, err := iprot.ReadBool(cxt); err != nil {
			return nil, err
		} else {
			tv = NewTBool(v)
		}
	}

	return
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
