package it

import (
	"context"
	"log/slog"
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/lavenderses/xk6-thrift/pkg/gen-go/idl"
)

var url = "http://127.0.0.1:8080/thrift"

func setupTransport(t *testing.T) (*thrift.TTransport, error) {
	var trans thrift.TTransport
	var err error
	if trans, err = thrift.NewTHttpClient(url); err != nil {
		return nil, err
	}
	// alaways close transport every test execution
	t.Cleanup(func() {
		slog.Info("closing tranport.", slog.Any("transport", trans))
		trans.Close()
	})

	return &trans, nil
}

func setupClient(trans *thrift.TTransport) (*idl.TestServiceClient, error) {
	if err := (*trans).Open(); err != nil {
		return nil, err
	}

	conf := thrift.TConfiguration{}
	pf := thrift.NewTBinaryProtocolFactoryConf(&conf)
	iprot := pf.GetProtocol(*trans)
	oprot := pf.GetProtocol(*trans)
	client := idl.NewTestServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return client, nil
}

func TestSimpleCall(t *testing.T) {
	// prepare
	var trans *thrift.TTransport
	var err error

	if trans, err = setupTransport(t); err != nil {
		t.Fatalf("error opening transport. %v", err)
	}

	var client *idl.TestServiceClient
	if client, err = setupClient(trans); err != nil {
		t.Fatalf("error creating client. %v", err)
	}

	cxt := context.Background()
	arg := "ID"
	expect := "Success: ID"
	var actual string

	// do & verify
	if actual, err = (*client).SimpleCall(cxt, arg); err != nil {
		t.Fatalf("error calling RPC. %v", err)
	}

	if actual != expect {
		t.Errorf("expected %v, but was %v", expect, actual)
	}
}
