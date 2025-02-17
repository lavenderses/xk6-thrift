package it

import (
	"context"
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
	assertEquals(t, *actual, *expect)
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
	assertEquals(t, *actual, *expect)
}

func TestMapCall(t *testing.T) {
	// prepare
	var client *thrift.TStandardClient
	var err error
	if client, err = setupClient(t); err != nil {
		t.Fatalf("error creating client. %v", err)
	}

	cxt := context.Background()
	method := "mapCall"
	value := map[xk6_thrift.TValue]xk6_thrift.TValue{
		xk6_thrift.NewTstring("key 1"): xk6_thrift.NewTBool(true),
		xk6_thrift.NewTstring("key 2"): xk6_thrift.NewTBool(false),
	}
	tvalue := map[int16]xk6_thrift.TValue{
		1: xk6_thrift.NewTMap(&value),
	}
	arg := xk6_thrift.NewTRequestWithValue(&tvalue)
	expectValue := map[xk6_thrift.TValue]xk6_thrift.TValue{
		xk6_thrift.NewTstring("NEW: key 1"): xk6_thrift.NewTBool(true),
		xk6_thrift.NewTstring("NEW: key 2"): xk6_thrift.NewTBool(false),
	}
	expectTValue := xk6_thrift.NewTMap(&expectValue)
	expect := xk6_thrift.NewTResponse()
	expect.Add(0, expectTValue)
	actual := xk6_thrift.NewTResponse()

	// do & verify
	if _, err = (*client).Call(cxt, method, arg, actual); err != nil {
		t.Fatalf("error calling RPC. %v", err)
	}

	assertEquals(t, *actual, *expect)
}

func TestMessageCall(t *testing.T) {
	// prepare
	var client *thrift.TStandardClient
	var err error
	if client, err = setupClient(t); err != nil {
		t.Fatalf("error creating client. %v", err)
	}

	cxt := context.Background()
	method := "messageCall"
	value := map[xk6_thrift.TStructField]xk6_thrift.TValue{
		*xk6_thrift.NewTStructField(1, "content"): xk6_thrift.NewTstring("this is a content"),
		*xk6_thrift.NewTStructField(2, "tags"): xk6_thrift.NewTMap(
			&map[xk6_thrift.TValue]xk6_thrift.TValue{
				xk6_thrift.NewTstring("bool true"): xk6_thrift.NewTBool(true),
				xk6_thrift.NewTstring("bool false"): xk6_thrift.NewTBool(false),
			},
		),
		*xk6_thrift.NewTStructField(3, "nested"): xk6_thrift.NewTStruct(
			&map[xk6_thrift.TStructField]xk6_thrift.TValue{
				*xk6_thrift.NewTStructField(1, "inner"): xk6_thrift.NewTstring("this is an inner content"),
			},
		),
	}
	tvalue := map[int16]xk6_thrift.TValue{
		1: xk6_thrift.NewTStruct(&value),
	}
	arg := xk6_thrift.NewTRequestWithValue(&tvalue)
	expectValue := map[xk6_thrift.TStructField]xk6_thrift.TValue{
		*xk6_thrift.NewTStructField(1, "content"): xk6_thrift.NewTstring("this is a content"),
		*xk6_thrift.NewTStructField(2, "tags"): xk6_thrift.NewTMap(
			&map[xk6_thrift.TValue]xk6_thrift.TValue{
				xk6_thrift.NewTstring("bool true"): xk6_thrift.NewTBool(true),
				xk6_thrift.NewTstring("bool false"): xk6_thrift.NewTBool(false),
			},
		),
		*xk6_thrift.NewTStructField(3, "nested"): xk6_thrift.NewTStruct(
			&map[xk6_thrift.TStructField]xk6_thrift.TValue{
				*xk6_thrift.NewTStructField(1, "inner"): xk6_thrift.NewTstring("this is an inner content"),
			},
		),
	}
	expectTValue := xk6_thrift.NewTStruct(&expectValue)
	expect := xk6_thrift.NewTResponse()
	expect.Add(0, expectTValue)
	actual := xk6_thrift.NewTResponse()

	// do & verify
	if _, err = (*client).Call(cxt, method, arg, actual); err != nil {
		t.Fatalf("error calling RPC. %v", err)
	}

	assertEquals(t, *actual, *expect)
}
