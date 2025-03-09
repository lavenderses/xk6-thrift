package thrift

import (
	"context"
	"testing"
)

func TestEquals_TBool_equals(t *testing.T) {
	// prepare
	a := NewTBool(true)
	var b TValue = NewTBool(true)
	expected := true

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_TBool_not_equals(t *testing.T) {
	// prepare
	a := NewTBool(true)
	var b TValue = NewTBool(false)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_TBool_other_type(t *testing.T) {
	// prepare
	a := NewTBool(true)
	var b TValue = NewTstring("other")
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestWriteFieldData(t *testing.T) {
	// prepare
	oprot := setupProtocol(t)
	cxt := context.Background()
	value := NewTBool(true)
	expected := true

	// do
	err := value.WriteFieldData(cxt, oprot)
	checkError(t, err)

	// verify
	oprot.Flush(cxt)
	actual, err := oprot.ReadBool(cxt)
	checkError(t, err)
	assert(t, "", actual, expected)
}

func TestReadBool_true(t *testing.T) {
	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	checkError(
		t,
		iprot.WriteBool(cxt, true),
	)
	checkError(
		t,
		iprot.Flush(cxt),
	)

	// do
	actual, err := ReadBool(cxt, iprot)
	checkError(t, err)

	// verfiy
	e := NewTBool(true)
	a, ok := actual.(TBool)
	assertTrue(t, "cast to TBool", ok)
	assert(t, "result", a.value, e.value)
}

func TestReadBool_false(t *testing.T) {
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
	actual, err := ReadBool(cxt, iprot)
	checkError(t, err)

	// verfiy
	e := NewTBool(false)
	a, ok := actual.(TBool)
	assertTrue(t, "cast to TBool", ok)
	assert(t, "result", a.value, e.value)
}

func TestReadBool_invalid_type(t *testing.T) {
	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	checkError(
		t,
		iprot.WriteString(cxt, "FOO"),
	)
	checkError(
		t,
		iprot.Flush(cxt),
	)

	// do
	_, err := ReadBool(cxt, iprot)

	// verfiy
	assertTrue(t, "error expected", err != nil)
}
