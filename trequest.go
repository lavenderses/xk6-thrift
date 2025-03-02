package thrift

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/apache/thrift/lib/go/thrift"
)

type TRequest struct {
	values map[int16]TValue
}

func NewTRequest() *TRequest {
	return &TRequest{values: make(map[int16]TValue)}
}

func NewTRequestWithValue(v *map[int16]TValue) *TRequest {
	return &TRequest{values: *v}
}

func (p *TRequest) Read(cxt context.Context, iprot thrift.TProtocol) (err error) {
	slog.Error("*Trequest.Read is not expected to be called.")
	return
}

func (p *TRequest) Write(cxt context.Context, oprot thrift.TProtocol) (err error) {
	if err = oprot.WriteStructBegin(cxt, "simple_args"); err != nil {
		err = thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
		return
	}

	if err = p.writeFields(cxt, oprot); err != nil {
		return
	}

	if err = oprot.WriteFieldStop(cxt); err != nil {
		err = thrift.PrependError("write field stop error: ", err)
		return
	}
	if err = oprot.WriteStructEnd(cxt); err != nil {
		err = thrift.PrependError("write struct stop error: ", err)
		return
	}
	return
}

func (p *TRequest) writeFields(cxt context.Context, oprot thrift.TProtocol) (err error) {
	if p == nil {
		return
	}

	for fid, v := range p.values {
		ttype := v.TType()
		fname := "dummy"

		if err = oprot.WriteFieldBegin(cxt, fname, ttype, fid); err != nil {
			err = thrift.PrependError(fmt.Sprintf("%T write field begin error %d:%s: ", p, fid, fname), err)
			return
		}
		if err = v.WriteFieldData(cxt, oprot); err != nil {
			return
		}
		if err = oprot.WriteFieldEnd(cxt); err != nil {
			err = thrift.PrependError(fmt.Sprintf("%T write field end error %d:%s: ", p, fid, fname), err)
			return
		}
	}
	return
}
