package thrift

import (
	"context"
	"fmt"

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
	// dummy
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
		if tv, ok := v.(*TString); ok {
			if err = tv.WriteField(cxt, oprot, fid, "dummy"); err != nil {
				return
			}
		}
		if tv, ok := v.(*TBool); ok {
			if err = tv.WriteField(cxt, oprot, fid, "dummy"); err != nil {
				return
			}
		}
	}
	return
}
