package thrift

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

type TString struct {
	value string
}

func NewTstring(v string) TString {
	return TString{value: v}
}

func (p TString) Equals(other *TValue) bool {
	o, ok := (*other).(TString)
	if !ok {
		return false
	}
	return p.value == o.value
}

func (p TString) WriteField(cxt context.Context, oprot thrift.TProtocol, fid int16, fname string) (err error) {
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
