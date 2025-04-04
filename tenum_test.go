package thrift

import (
	"context"
	"testing"
)

func TestEquals_TEnum_Equals(t *testing.T) {
	// prepare
	a := NewTEnum(3)
	var b TValue = NewTEnum(3)
	expected := true

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_TEnum_NotEquals(t *testing.T) {
	// prepare
	a := NewTEnum(3)
	var b TValue = NewTEnum(1)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_TEnum_OtherType(t *testing.T) {
	// prepare
	a := NewTEnum(3)
	var b TValue = NewTBool(false)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestWriteFieldData_Enum(t *testing.T) {
	// prepare
	oprot := setupProtocol(t)
	cxt := context.Background()
	value := NewTEnum(3)
	var expected int32 = 3

	// do
	err := value.WriteFieldData(cxt, oprot)
	checkError(t, err)

	// verify
	oprot.Flush(cxt)
	actual, err := oprot.ReadI32(cxt)
	checkError(t, err)
	assert(t, "", actual, expected)
}

func TestReadEnum(t *testing.T) {
	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	checkError(
		t,
		iprot.WriteI32(cxt, 3),
	)
	checkError(
		t,
		iprot.Flush(cxt),
	)

	// do
	actual, err := ReadEnum(cxt, iprot)
	checkError(t, err)

	// verfiy
	e := NewTEnum(3)
	a, ok := actual.(TEnum)
	assertTrue(t, "cast to TEnum", ok)
	assert(t, "result", a.value, e.value)
}

func TestReadEnum_InvalidType(t *testing.T) {
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
	_, err := ReadEnum(cxt, iprot)

	// verfiy
	assertTrue(t, "error expected", err != nil)
}
