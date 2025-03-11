package thrift

import (
	"context"
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
)

func TestEquals_Equals(t *testing.T) {
	// prepare
	a := NewTMap(
		thrift.STRING,
		thrift.BOOL,
		&map[TValue]TValue{
			NewTstring("key 1"): NewTBool(false),
			NewTstring("key 2"): NewTBool(true),
		},
	)
	var b TValue = NewTMap(
		thrift.STRING,
		thrift.BOOL,
		&map[TValue]TValue{
			NewTstring("key 2"): NewTBool(true),
			NewTstring("key 1"): NewTBool(false),
		},
	)
	expected := true

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_DifferentValue(t *testing.T) {
	// prepare
	a := NewTMap(
		thrift.STRING,
		thrift.BOOL,
		&map[TValue]TValue{
			NewTstring("key 1"): NewTBool(false),
			NewTstring("key 2"): NewTBool(true),
		},
	)
	var b TValue = NewTMap(
		thrift.STRING,
		thrift.BOOL,
		&map[TValue]TValue{
			NewTstring("key 1"): NewTBool(false),
			NewTstring("key 2"): NewTBool(false),
		},
	)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_DifferentKeyValueCount(t *testing.T) {
	// prepare
	a := NewTMap(
		thrift.STRING,
		thrift.BOOL,
		&map[TValue]TValue{
			NewTstring("key 1"): NewTBool(false),
			NewTstring("key 2"): NewTBool(false),
		},
	)
	var b TValue = NewTMap(
		thrift.STRING,
		thrift.BOOL,
		&map[TValue]TValue{
			NewTstring("key 1"): NewTBool(false),
			NewTstring("key 2"): NewTBool(false),
			NewTstring("key 3"): NewTBool(false),
		},
	)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_OtherKeyValueType(t *testing.T) {
	// prepare
	a := NewTMap(
		thrift.STRING,
		thrift.BOOL,
		&map[TValue]TValue{
			NewTstring("key 1"): NewTBool(false),
		},
	)
	var b TValue = NewTMap(
		thrift.BOOL,
		thrift.STRING,
		&map[TValue]TValue{
			NewTBool(false): NewTBool(false),
		},
	)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_OtherType(t *testing.T) {
	// prepare
	a := NewTMap(
		thrift.STRING,
		thrift.BOOL,
		&map[TValue]TValue{
			NewTstring("key 1"): NewTBool(false),
			NewTstring("key 2"): NewTBool(true),
		},
	)
	var b TValue = NewTstring("other")
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestWriteFieldData_TMap(t *testing.T) {
	// prepare
	oprot := setupProtocol(t)
	cxt := context.Background()
	value := NewTMap(
		thrift.STRING,
		thrift.BOOL,
		&map[TValue]TValue{
			NewTstring("key 1"): NewTBool(false),
			NewTstring("key 2"): NewTBool(true),
		},
	)

	// do
	err := value.WriteFieldData(cxt, oprot)
	checkError(t, err)

	// verify
	oprot.Flush(cxt)
	{
		ktype, vtype, size, err := oprot.ReadMapBegin(cxt)
		assert(t, "size", size, 2)
		checkError(t, err)
		assert(t, "key type", ktype, thrift.STRING)
		assert(t, "value type", vtype, thrift.BOOL)
	}
	{
		k, err := oprot.ReadString(cxt)
		checkError(t, err)
		assert(t, "key 1", k, "key 1")
	}
	{
		v, err := oprot.ReadBool(cxt)
		checkError(t, err)
		assert(t, "key 1 value", v, false)
	}
	{
		k, err := oprot.ReadString(cxt)
		checkError(t, err)
		assert(t, "key 2", k, "key 2")
	}
	{
		v, err := oprot.ReadBool(cxt)
		checkError(t, err)
		assert(t, "key 2 value", v, true)
	}
	{
		err := oprot.ReadMapEnd(cxt)
		checkError(t, err)
	}
}

func TestReadMap_StringToBool(t *testing.T) {
	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	{
		checkError(
			t,
			iprot.WriteMapBegin(cxt, thrift.STRING, thrift.BOOL, 2),
		)
		checkError(
			t,
			iprot.WriteString(cxt, "key 1"),
		)
		checkError(
			t,
			iprot.WriteBool(cxt, true),
		)
		checkError(
			t,
			iprot.WriteString(cxt, "key 2"),
		)
		checkError(
			t,
			iprot.WriteBool(cxt, false),
		)
		checkError(
			t,
			iprot.WriteMapEnd(cxt),
		)
		checkError(
			t,
			iprot.Flush(cxt),
		)
	}

	// do
	actual, err := ReadMap(cxt, iprot)
	checkError(t, err)

	// verfiy
	a, ok := actual.(*TMap)
	assertTrue(t, "cast to TMap", ok)
	{
		assertTrue(t, "size", len((*a).value) == 2)
	}
	{
		tv := (*a).value[NewTstring("key 1")]
		v, ok := tv.(TBool)
		assertTrue(t, "cast key 1 value to TBool", ok)
		assert(t, "key 1 value", v.value, true)
	}
	{
		tv := (*a).value[NewTstring("key 2")]
		v, ok := tv.(TBool)
		assertTrue(t, "cast key 2 to TBool", ok)
		assert(t, "key 2 value", v.value, false)
	}
}

func TestReadMap_ContainerData(t *testing.T) {
	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	{
		checkError(
			t,
			iprot.WriteMapBegin(cxt, thrift.STRING, thrift.STRUCT, 1),
		)
		// key
		checkError(
			t,
			iprot.WriteString(cxt, "key"),
		)
		// value
		checkError(
			t,
			iprot.WriteStructBegin(cxt, "struct"),
		)

		checkError(
			t,
			iprot.WriteFieldBegin(cxt, "string", thrift.STRING, 1),
		)
		checkError(
			t,
			iprot.WriteString(cxt, "string string"),
		)
		checkError(
			t,
			iprot.WriteFieldEnd(cxt),
		)

		checkError(
			t,
			iprot.WriteFieldBegin(cxt, "boolean", thrift.BOOL, 2),
		)
		checkError(
			t,
			iprot.WriteBool(cxt, true),
		)
		checkError(
			t,
			iprot.WriteFieldEnd(cxt),
		)

		checkError(
			t,
			iprot.WriteFieldStop(cxt),
		)
		checkError(
			t,
			iprot.WriteStructEnd(cxt),
		)

		checkError(
			t,
			iprot.WriteMapEnd(cxt),
		)
		checkError(
			t,
			iprot.Flush(cxt),
		)
	}

	// do
	actual, err := ReadMap(cxt, iprot)
	checkError(t, err)

	// verfiy
	a, ok := actual.(*TMap)
	assertTrue(t, "cast to TMap", ok)
	{
		assertTrue(t, "size", len((*a).value) == 1)
	}
	tv := (*a).value[NewTstring("key")]
	v, ok := tv.(*TStruct)
	assertTrue(t, "cast key 1 value to TStruct", ok)
	{
		tf := v.value[TStructField{id: 1, name: ""}]
		f, ok := tf.(TString)
		assertTrue(t, "cast field 1 to TString", ok)
		assert(t, "field string", f.value, "string string")
	}
	{
		tf := v.value[TStructField{id: 2, name: ""}]
		f, ok := tf.(TBool)
		assertTrue(t, "cast field 2 to TBool", ok)
		assert(t, "field boolean", f.value, true)
	}
}

func TestReadBool_Invalid_type(t *testing.T) {
	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	{
		err := iprot.WriteString(cxt, "other")
		checkError(t, err)
	}

	// do
	_, err := ReadMap(cxt, iprot)

	// verfiy
	assertTrue(t, "", err != nil)
}
