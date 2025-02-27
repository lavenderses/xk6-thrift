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
		slog.Error(fmt.Sprintf("ERROR while getting transport: %v", err))
		return NewTCallResult(nil, err)
	}
	pf := thrift.NewTBinaryProtocolFactoryConf(&cfg)
	iprot := pf.GetProtocol(transport)
	oprot := pf.GetProtocol(transport)
	tclient := thrift.NewTStandardClient(iprot, oprot)
	defer transport.Close()

	err = transport.Open()
	if err != nil {
		slog.Error(fmt.Sprintf("ERROR while opening transport: %v", err))
		return NewTCallResult(nil, err)
	}

	values := make(map[int16]TValue)
	values[1] = NewTstring("ID")
	req := NewTRequestWithValue(&values)
	res := NewTResponse()

	cxt := context.Background()
	_, err = tclient.Call(cxt, method, req, res)
	if err != nil {
		slog.Error(fmt.Sprintf("ERROR calling RPC: %v", err))
		return NewTCallResult(nil, err)
	}

	slog.Info(fmt.Sprintf("Response: %v", res))

	body := res.values[0]
	if body == nil {
		return NewTCallResult(nil, fmt.Errorf("Empty body"))
	}
	
	return NewTCallResult(&body, nil)
}

func (m *TModule) Call(method string, req *TRequest) *TCallResult {
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
		slog.Error(fmt.Sprintf("ERROR: %v", err))
		return NewTCallResult(nil, err)
	}
	pf := thrift.NewTBinaryProtocolFactoryConf(&cfg)
	iprot := pf.GetProtocol(transport)
	oprot := pf.GetProtocol(transport)
	tclient := thrift.NewTStandardClient(iprot, oprot)
	defer transport.Close()

	err = transport.Open()
	if err != nil {
		slog.Error(fmt.Sprintf("ERROR: %v", err))
		return NewTCallResult(nil, err)
	}

	res := NewTResponse()

	cxt := context.Background()
	_, err = tclient.Call(cxt, method, req, res)
	if err != nil {
		slog.Error(fmt.Sprintf("ERROR calling RPC: %v", err))
		return NewTCallResult(nil, err)
	}
	body := res.values[0]
	if body == nil {
		return NewTCallResult(nil, fmt.Errorf("Empty body"))
	}

	slog.Info("Response:", res.values)
	
	return NewTCallResult(&body, nil)
}
