package thrift

import (
	"context"
	"testing"
)

func TestEquals_TString_Equals(t *testing.T) {
	// prepare
	a := NewTstring("foo")
	var b TValue = NewTstring("foo")
	expected := true

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_TString_NotEquals(t *testing.T) {
	// prepare
	a := NewTstring("foo")
	var b TValue = NewTstring("bar")
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_TString_OtherType(t *testing.T) {
	// prepare
	a := NewTstring("foo")
	var b TValue = NewTBool(false)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestWriteFieldData_String(t *testing.T) {
	// prepare
	oprot := setupProtocol(t)
	cxt := context.Background()
	value := NewTstring("foo")
	expected := "foo"

	// do
	err := value.WriteFieldData(cxt, oprot)
	checkError(t, err)

	// verify
	oprot.Flush(cxt)
	actual, err := oprot.ReadString(cxt)
	checkError(t, err)
	assert(t, "", actual, expected)
}

func TestReadString(t *testing.T) {
	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	checkError(
		t,
		iprot.WriteString(cxt, "foo"),
	)
	checkError(
		t,
		iprot.Flush(cxt),
	)

	// do
	actual, err := ReadString(cxt, iprot)
	checkError(t, err)

	// verfiy
	e := NewTstring("foo")
	a, ok := actual.(TString)
	assertTrue(t, "cast to TString", ok)
	assert(t, "result", a.value, e.value)
}

func TestReadString_InvalidType(t *testing.T) {
	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	checkError(
		t,
		iprot.WriteBool(cxt, false),
	)
	checkError(
		t,
		iprot.Flush(cxt),
	)

	// do
	_, err := ReadString(cxt, iprot)

	// verfiy
	assertTrue(t, "error expected", err != nil)
}
