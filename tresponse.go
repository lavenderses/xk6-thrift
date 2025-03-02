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
				v, err = p.ReadMap(cxt, iprot, fieldId)
			case thrift.STRUCT:
				v, err = p.ReadStruct(cxt, iprot, fieldId)
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

func (p *TResponse) ReadMap(cxt context.Context, iproto thrift.TProtocol, fieldId int16) (*TMap, error) {
	keyType, valueType, size, err := iproto.ReadMapBegin(cxt)
	if err != nil {
		return nil, thrift.PrependError(fmt.Sprintf("error reading map field %d: ", fieldId), err)
	}

	tmap := make(map[TValue]TValue)
	for i := 0; i < size; i++ {
		if err = p.readFeidlDataList(cxt, iproto, &tmap, keyType, valueType); err != nil {
			return nil, thrift.PrependError(fmt.Sprintf("error reading map %d: ", fieldId), err)
		}
	}

	res := NewTMap(&tmap)
	return res, nil
}

func (p *TResponse) readFeidlDataList(cxt context.Context, iprot thrift.TProtocol, tmap *map[TValue]TValue, ktype, vtype thrift.TType) error {
	var key, value TValue
	var err error
	if key, err = p.readField(cxt, iprot, ktype); err != nil {
		return err
	}
	if value, err = p.readField(cxt, iprot, vtype); err != nil {
		return err
	}

	(*tmap)[key] = value
	return nil
}

func (p *TResponse) readField(cxt context.Context, iprot thrift.TProtocol, ttype thrift.TType) (tv TValue, err error) {
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

func (p *TResponse) ReadStruct(cxt context.Context, iprot thrift.TProtocol, fieldId int16) (*TStruct, error) {
	fieldName, err := iprot.ReadStructBegin(cxt)
	if err != nil {
		return nil, thrift.PrependError(fmt.Sprintf("%T read struct begin (%d, %s) error: ", p, fieldId, fieldName), err)
	}

	tvalue := make(map[TStructField]TValue)
	for {
		fname, ftype, fid, err := iprot.ReadFieldBegin(cxt)
		if err != nil {
			return nil, thrift.PrependError(fmt.Sprintf("%T field read field begin (%d, %s) error: ", p, fid, fname), err)
		}
		if ftype == thrift.STOP {
			break
		}

		var tv TValue
		// slog.Info(fmt.Sprintf("GOT %s, %v, %d", fieldName, fieldTypeId, fieldId))
		switch ftype {
		case thrift.STRING:
			tv, err = ReadString(cxt, iprot)
		case thrift.BOOL:
			tv, err = ReadBool(cxt, iprot)
		case thrift.MAP:
			tv, err = p.ReadMap(cxt, iprot, fid)
		case thrift.STRUCT:
			tv, err = p.ReadStruct(cxt, iprot, fid)
		default:
			err = iprot.Skip(cxt, ftype)
		}
		if err != nil {
			return nil, err
		}

		err = iprot.ReadFieldEnd(cxt)
		if err != nil {
			return nil, err
		}

		// TODO: somehow fname gecomes an empty string
		tvalue[*NewTStructField(fid, fname)] = tv
	}

	err = iprot.ReadStructEnd(cxt)
	if err != nil {
		return nil, thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}

	res := NewTStruct(&tvalue)
	return res, nil
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
