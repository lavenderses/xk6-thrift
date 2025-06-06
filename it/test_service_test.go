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

func TestSimpleCallFailure(t *testing.T) {
	// prepare
	var client *thrift.TStandardClient
	var err error
	if client, err = setupClient(t); err != nil {
		t.Fatalf("error creating client. %v", err)
	}

	cxt := context.Background()
	method := "simpleCall"
	values := make(map[int16]xk6_thrift.TValue)
	values[1] = xk6_thrift.NewTstring("FAILURE")
	arg := xk6_thrift.NewTRequestWithValue(&values)
	expect := xk6_thrift.NewTResponse()
	expect.Add(0, xk6_thrift.NewTstring("Success: ID"))
	actual := xk6_thrift.NewTResponse()

	// do & verify
	if _, err = (*client).Call(cxt, method, arg, actual); err == nil {
		t.Fatalf("error NOT THROWN calling RPC. %v", err)
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
		1: xk6_thrift.NewTMap(thrift.STRING, thrift.BOOL, &value),
	}
	arg := xk6_thrift.NewTRequestWithValue(&tvalue)
	expectValue := map[xk6_thrift.TValue]xk6_thrift.TValue{
		xk6_thrift.NewTstring("NEW: key 1"): xk6_thrift.NewTBool(true),
		xk6_thrift.NewTstring("NEW: key 2"): xk6_thrift.NewTBool(false),
	}
	expectTValue := xk6_thrift.NewTMap(thrift.STRING, thrift.BOOL, &expectValue)
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
			thrift.STRING,
			thrift.BOOL,
			&map[xk6_thrift.TValue]xk6_thrift.TValue{
				xk6_thrift.NewTstring("bool true"):  xk6_thrift.NewTBool(true),
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
		*xk6_thrift.NewTStructField(1, ""): xk6_thrift.NewTstring("content: this is a content"),
		*xk6_thrift.NewTStructField(2, ""): xk6_thrift.NewTMap(
			thrift.STRING,
			thrift.BOOL,
			&map[xk6_thrift.TValue]xk6_thrift.TValue{
				xk6_thrift.NewTstring("bool true"):  xk6_thrift.NewTBool(true),
				xk6_thrift.NewTstring("bool false"): xk6_thrift.NewTBool(false),
			},
		),
		*xk6_thrift.NewTStructField(3, ""): xk6_thrift.NewTStruct(
			&map[xk6_thrift.TStructField]xk6_thrift.TValue{
				*xk6_thrift.NewTStructField(1, ""): xk6_thrift.NewTstring("this is an inner content"),
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

func TestStringCall(t *testing.T) {
	// prepare
	var client *thrift.TStandardClient
	var err error
	if client, err = setupClient(t); err != nil {
		t.Fatalf("error creating client. %v", err)
	}

	cxt := context.Background()
	method := "stringCall"
	value := []xk6_thrift.TValue{
		xk6_thrift.NewTstring("content-1"),
		xk6_thrift.NewTstring("content-2"),
	}
	tvalue := map[int16]xk6_thrift.TValue{
		1: xk6_thrift.NewTList(&value, thrift.STRING),
	}
	arg := xk6_thrift.NewTRequestWithValue(&tvalue)
	expectValue := []xk6_thrift.TValue{
		xk6_thrift.NewTstring("content-1:content-1"),
		xk6_thrift.NewTstring("content-2:content-2"),
	}
	expectTValue := xk6_thrift.NewTList(&expectValue, thrift.STRING)
	expect := xk6_thrift.NewTResponse()
	expect.Add(0, expectTValue)
	actual := xk6_thrift.NewTResponse()

	// do & verify
	if _, err = (*client).Call(cxt, method, arg, actual); err != nil {
		t.Fatalf("error calling RPC. %v", err)
	}

	assertEquals(t, *actual, *expect)
}

func TestStringsCall(t *testing.T) {
	// prepare
	var client *thrift.TStandardClient
	var err error
	if client, err = setupClient(t); err != nil {
		t.Fatalf("error creating client. %v", err)
	}

	cxt := context.Background()
	method := "stringsCall"

	msg1 := map[xk6_thrift.TStructField]xk6_thrift.TValue{
		*xk6_thrift.NewTStructField(1, "content"): xk6_thrift.NewTstring("content 1"),
		*xk6_thrift.NewTStructField(2, "tags"): xk6_thrift.NewTMap(
			thrift.STRING,
			thrift.BOOL,
			&map[xk6_thrift.TValue]xk6_thrift.TValue{
				xk6_thrift.NewTstring("bool true"): xk6_thrift.NewTBool(true),
			},
		),
		*xk6_thrift.NewTStructField(3, "nested"): xk6_thrift.NewTStruct(
			&map[xk6_thrift.TStructField]xk6_thrift.TValue{
				*xk6_thrift.NewTStructField(1, "inner"): xk6_thrift.NewTstring("inner content 1"),
			},
		),
	}
	msg2 := map[xk6_thrift.TStructField]xk6_thrift.TValue{
		*xk6_thrift.NewTStructField(1, "content"): xk6_thrift.NewTstring("content 2"),
		*xk6_thrift.NewTStructField(2, "tags"): xk6_thrift.NewTMap(
			thrift.STRING,
			thrift.BOOL,
			&map[xk6_thrift.TValue]xk6_thrift.TValue{
				xk6_thrift.NewTstring("bool false"): xk6_thrift.NewTBool(false),
			},
		),
		*xk6_thrift.NewTStructField(3, "nested"): xk6_thrift.NewTStruct(
			&map[xk6_thrift.TStructField]xk6_thrift.TValue{
				*xk6_thrift.NewTStructField(1, "inner"): xk6_thrift.NewTstring("inner content 2"),
			},
		),
	}
	tmsg := []xk6_thrift.TValue{
		xk6_thrift.NewTStruct(&msg1),
		xk6_thrift.NewTStruct(&msg2),
	}
	tvalue := map[int16]xk6_thrift.TValue{
		1: xk6_thrift.NewTList(&tmsg, thrift.STRUCT),
	}
	arg := xk6_thrift.NewTRequestWithValue(&tvalue)

	exChange := "content: "
	exMsg1 := map[xk6_thrift.TStructField]xk6_thrift.TValue{
		*xk6_thrift.NewTStructField(1, ""): xk6_thrift.NewTstring(exChange + "content 1"),
		*xk6_thrift.NewTStructField(2, ""): xk6_thrift.NewTMap(
			thrift.STRING,
			thrift.BOOL,
			&map[xk6_thrift.TValue]xk6_thrift.TValue{
				xk6_thrift.NewTstring("bool true"): xk6_thrift.NewTBool(true),
			},
		),
		*xk6_thrift.NewTStructField(3, ""): xk6_thrift.NewTStruct(
			&map[xk6_thrift.TStructField]xk6_thrift.TValue{
				*xk6_thrift.NewTStructField(1, ""): xk6_thrift.NewTstring("inner content 1"),
			},
		),
	}
	exMsg2 := map[xk6_thrift.TStructField]xk6_thrift.TValue{
		*xk6_thrift.NewTStructField(1, ""): xk6_thrift.NewTstring(exChange + "content 2"),
		*xk6_thrift.NewTStructField(2, ""): xk6_thrift.NewTMap(
			thrift.STRING,
			thrift.BOOL,
			&map[xk6_thrift.TValue]xk6_thrift.TValue{
				xk6_thrift.NewTstring("bool false"): xk6_thrift.NewTBool(false),
			},
		),
		*xk6_thrift.NewTStructField(3, ""): xk6_thrift.NewTStruct(
			&map[xk6_thrift.TStructField]xk6_thrift.TValue{
				*xk6_thrift.NewTStructField(1, ""): xk6_thrift.NewTstring("inner content 2"),
			},
		),
	}
	expectValue := []xk6_thrift.TValue{
		xk6_thrift.NewTStruct(&exMsg1),
		xk6_thrift.NewTStruct(&exMsg2),
	}
	expectTValue := xk6_thrift.NewTList(&expectValue, thrift.STRUCT)
	expect := xk6_thrift.NewTResponse()
	expect.Add(0, expectTValue)
	actual := xk6_thrift.NewTResponse()

	// do & verify
	if _, err = (*client).Call(cxt, method, arg, actual); err != nil {
		t.Fatalf("error calling RPC. %v", err)
	}

	assertEquals(t, *actual, *expect)
}

func TestEnumCall(t *testing.T) {
	// prepare
	var client *thrift.TStandardClient
	var err error
	if client, err = setupClient(t); err != nil {
		t.Fatalf("error creating client. %v", err)
	}

	cxt := context.Background()
	method := "enumCall"
	tvalue := map[int16]xk6_thrift.TValue{
		1: xk6_thrift.NewTEnum(1), // ONE
	}
	arg := xk6_thrift.NewTRequestWithValue(&tvalue)
	expectValue := []xk6_thrift.TValue{
		xk6_thrift.NewTEnum(1), // ONE
		xk6_thrift.NewTEnum(2), // TWO
		xk6_thrift.NewTEnum(3), // THREE
	}
	expectTValue := xk6_thrift.NewTList(&expectValue, thrift.I32)
	expect := xk6_thrift.NewTResponse()
	expect.Add(0, expectTValue)
	actual := xk6_thrift.NewTResponse()

	// do & verify
	if _, err = (*client).Call(cxt, method, arg, actual); err != nil {
		t.Fatalf("error calling RPC. %v", err)
	}

	assertEquals(t, *actual, *expect)
}
