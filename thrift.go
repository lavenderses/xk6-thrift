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
	modules.Register("k6/x/thrift", &TModule{})
}

type TModule struct{}

func (m *TModule) Echo() {
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
		return
	}
	pf := thrift.NewTBinaryProtocolFactoryConf(&cfg)
	iprot := pf.GetProtocol(transport)
	oprot := pf.GetProtocol(transport)
	tclient := thrift.NewTStandardClient(iprot, oprot)
	defer transport.Close()

	err = transport.Open()
	if err != nil {
		slog.Error("ERROR while opening transport", slog.Any("error", err))
		return
	}

	values := make(map[int16]TValue)
	values[1] = NewTstring("ID")
	req := NewTRequestWithValue(&values)
	res := NewTResponse()

	cxt := context.Background()
	_, err = tclient.Call(cxt, method, req, res)
	if err != nil {
		slog.Error("ERROR calling RPC", slog.Any("error", err))
		return
	}

	slog.Info(fmt.Sprintf("Response: %v", res))
}
