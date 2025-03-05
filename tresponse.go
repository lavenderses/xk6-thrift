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
			switch fieldTypeId {
			case thrift.STRING:
				v, err = ReadString(cxt, iprot)
			case thrift.BOOL:
				v, err = ReadBool(cxt, iprot)
			case thrift.LIST:
				v, err = p.ReadList(cxt, iprot, fieldId)
			case thrift.MAP:
				v, err = ReadMap(cxt, iprot)
			case thrift.STRUCT:
				v, err = ReadStruct(cxt, iprot)
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

func (p *TResponse) ReadList(cxt context.Context, iproto thrift.TProtocol, fieldId int16) (*TList, error) {
	valueType, size, err := iproto.ReadListBegin(cxt)
	if err != nil {
		return nil, thrift.PrependError(fmt.Sprintf("error reading list field %d: ", fieldId), err)
	}

	var tlist []TValue
	for i := 0; i < size; i++ {
		var tv TValue
		tv, err = p.readListField(cxt, iproto, valueType)
		if err != nil {
			return nil, thrift.PrependError(fmt.Sprintf("error reading list %d: ", fieldId), err)
		}
		tlist = append(tlist, tv)
	}

	res := NewTList(&tlist, valueType)
	return res, nil
}

func (p *TResponse) readListField(cxt context.Context, iprot thrift.TProtocol, ttype thrift.TType) (tv TValue, err error) {
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
		return thrift.PrependError(fmt.Sprintf("%T write struct end error: ", p), err)
	}
	return nil
}

func (p *TResponse) Add(key int16, value TValue) {
	p.values[key] = value
}
