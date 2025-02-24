package thrift

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/apache/thrift/lib/go/thrift"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/thrift", new(TModule))
	modules.Register("k6/x/thrift/ttypes", new(TTypes))
}

type TModule struct{}

func (m *TModule) Echo() *TCallResult {
	host := "127.0.0.1"
	port := 8080
	path := "/thrift"
	method := "simpleCall"

	tf := thrift.NewTHttpClientTransportFactory("http://" + host + ":" + strconv.Itoa(port) + path)
	cfg := thrift.TConfiguration{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	sock := thrift.NewTSocketConf("127.0.0.1:8080", &cfg)
	transport, err := tf.GetTransport(sock)
	if err != nil {
		slog.Error("ERROR while getting transport", slog.Any("error", err))
		return NewFailureTCallResult(err)
	}
	pf := thrift.NewTBinaryProtocolFactoryConf(&cfg)
	iprot := pf.GetProtocol(transport)
	oprot := pf.GetProtocol(transport)
	tclient := thrift.NewTStandardClient(iprot, oprot)
	defer transport.Close()

	err = transport.Open()
	if err != nil {
		slog.Error("ERROR while opening transport", slog.Any("error", err))
		return NewFailureTCallResult(err)
	}

	values := make(map[int16]TValue)
	values[1] = NewTstring("ID")
	req := NewTRequestWithValue(&values)
	res := NewTResponse()

	cxt := context.Background()
	_, err = tclient.Call(cxt, method, req, res)
	if err != nil {
		slog.Error("ERROR calling RPC", slog.Any("error", err))
		return NewFailureTCallResult(err)
	}

	slog.Info(fmt.Sprintf("Response: %v", res))

	body := res.values[0]
	if body == nil {
		return NewFailureTCallResult(fmt.Errorf("Empty body"))
	}
	
	return NewTCallResult(&body)
}

func (m *TModule) Call(method string, req *TRequest) {
	host := "127.0.0.1"
	port := 8080
	path := "/thrift"

	tf := thrift.NewTHttpClientTransportFactory("http://" + host + ":" + strconv.Itoa(port) + path)
	cfg := thrift.TConfiguration{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	sock := thrift.NewTSocketConf("127.0.0.1:8080", &cfg)
	transport, err := tf.GetTransport(sock)
	if err != nil {
		slog.Error("ERROR: ", err)
		return
	}
	pf := thrift.NewTBinaryProtocolFactoryConf(&cfg)
	iprot := pf.GetProtocol(transport)
	oprot := pf.GetProtocol(transport)
	tclient := thrift.NewTStandardClient(iprot, oprot)
	defer transport.Close()

	err = transport.Open()
	if err != nil {
		slog.Error("ERROR: ", err)
		return
	}

	res := NewTResponse()

	cxt := context.Background()
	tclient.Call(cxt, method, req, res)

	slog.Info("Response:", res.values)
}
