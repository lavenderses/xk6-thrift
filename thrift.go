package thrift

import (
	"context"
	"crypto/tls"
	"fmt"
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
		fmt.Println("ERROR: ", err)
		return
	}
	pf := thrift.NewTBinaryProtocolFactoryConf(&cfg)
	iprot := pf.GetProtocol(transport)
	oprot := pf.GetProtocol(transport)
	tclient := thrift.NewTStandardClient(iprot, oprot)
	defer transport.Close()

	err = transport.Open()
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	req := TString{}
	res := TString{}
	req.value = "request ID"

	cxt := context.Background()
	tclient.Call(cxt, method, &req, &res)

	fmt.Println("Response: \"", res.value, "\"")
}
