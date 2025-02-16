package thrift

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

type TString struct {
	value string
}

func NewTstring(v string) *TString {
	return &TString{value: v}
}

func (p *TString) Equals(other *TValue) bool {
	o, ok := (*other).(*TString)
	if !ok {
		return false
	}
	return p.value == o.value
}

func (p *TString) WriteField(cxt context.Context, oprot thrift.TProtocol, fid int16, fname string) (err error) {
	if err = oprot.WriteFieldBegin(cxt, fname, thrift.STRING, fid); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write field begin error %d:%s: ", p, fid, fname), err)
		return
	}
	if err = oprot.WriteString(cxt, string(p.value)); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T.id (1) field write error: ", p), err)
		return
	}
	if err = oprot.WriteFieldEnd(cxt); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write field end error %d:%s: ", p, fid, fname), err)
		return
	}
	return
}

func (p *TString) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *TString) ReadField0(ctx context.Context, iprot thrift.TProtocol) error {
	v, err := iprot.ReadString(ctx)
	fmt.Println("READ 0", v, err)
	if err != nil {
		return thrift.PrependError("error reading field 0: ", err)
	} else {
		p.value = v
	}
	return nil
}

func (p *TString) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.value = v
	}
	return nil
}

func (p *TString) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "simple_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *TString) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "id", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:id: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.value)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.id (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:id: ", p), err)
	}
	return err
}
