package it

import (
	"context"
	"reflect"
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
	xk6_thrift "github.com/lavenderses/xk6-thrift"
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
		t.Logf("closing tranport. %v", trans)
		trans.Close()
	})

	return &trans, nil
}

func setupClient(t *testing.T) (*thrift.TStandardClient, error) {
	var trans *thrift.TTransport
	var err error
	if trans, err = setupTransport(t); err != nil {
		t.Fatalf("error opening transport. %v", err)
	}

	if err := (*trans).Open(); err != nil {
		return nil, err
	}

	conf := thrift.TConfiguration{}
	pf := thrift.NewTBinaryProtocolFactoryConf(&conf)
	iprot := pf.GetProtocol(*trans)
	oprot := pf.GetProtocol(*trans)
	client := thrift.NewTStandardClient(iprot, oprot)
	return client, nil
}

func TestSimpleCall(t *testing.T) {
	// prepare
	var client *thrift.TStandardClient
	var err error
	if client, err = setupClient(t); err != nil {
		t.Fatalf("error creating client. %v", err)
	}

	cxt := context.Background()
	method := "simpleCall"
	values := make(map[int16]xk6_thrift.TValue)
	values[1] = xk6_thrift.NewTstring("ID")
	arg := xk6_thrift.NewTRequestWithValue(&values)
	expect := xk6_thrift.NewTResponse()
	expect.Add(0, xk6_thrift.NewTstring("Success: ID"))
	actual := xk6_thrift.NewTResponse()

	// do & verify
	if _, err = (*client).Call(cxt, method, arg, actual); err != nil {
		t.Fatalf("error calling RPC. %v", err)
	}

	if !reflect.DeepEqual(actual.Values(), expect.Values()) {
		t.Errorf("expected %v, but was %v", expect, actual)
	}
}

func TestBoolCall(t *testing.T) {
	// prepare
	var client *thrift.TStandardClient
	var err error
	if client, err = setupClient(t); err != nil {
		t.Fatalf("error creating client. %v", err)
	}

	cxt := context.Background()
	method := "boolCall"
	values := make(map[int16]xk6_thrift.TValue)
	values[1] = xk6_thrift.NewTBool(true)
	arg := xk6_thrift.NewTRequestWithValue(&values)
	expect := xk6_thrift.NewTResponse()
	expect.Add(0, xk6_thrift.NewTBool(true))
	actual := xk6_thrift.NewTResponse()

	// do & verify
	if _, err = (*client).Call(cxt, method, arg, actual); err != nil {
		t.Fatalf("error calling RPC. %v", err)
	}

	if !reflect.DeepEqual(actual.Values(), expect.Values()) {
		t.Errorf("expected %v, but was %v", expect, actual)
	}
}
