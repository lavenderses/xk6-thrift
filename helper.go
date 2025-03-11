package thrift

import (
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
)

func setupProtocol(t *testing.T) thrift.TProtocol {
	tf := thrift.NewTMemoryBufferTransportFactory(1_024)
	trans, err := tf.GetTransport(nil)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	pf := thrift.NewTJSONProtocolFactory()
	debugPf := thrift.TDebugProtocolFactory{
		Underlying: pf,
		LogPrefix:  "xk6-thrif5",
	}
	return debugPf.GetProtocol(trans)
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func assertTrue(t *testing.T, title string, result bool) {
	assert(t, title, result, true)
}

func assert[T string | bool | int | int16 | thrift.TType](t *testing.T, title string, actual, expected T) {
	if actual != expected {
		t.Fatalf("[%v] Expected %v but was %v", title, expected, actual)
	}
}
