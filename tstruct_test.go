package thrift

import (
	"context"
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
)

func TestEquals_Struct_Equals(t *testing.T) {
	// prepare
	a := NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTMap(
				thrift.STRING,
				thrift.BOOL,
				&map[TValue]TValue{
					NewTstring("key 1"): NewTBool(true),
					NewTstring("key 2"): NewTBool(false),
				},
			),
		},
	)
	var b TValue = NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTMap(
				thrift.STRING,
				thrift.BOOL,
				&map[TValue]TValue{
					NewTstring("key 1"): NewTBool(true),
					NewTstring("key 2"): NewTBool(false),
				},
			),
		},
	)
	expected := true

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_Struct_DifferentValue(t *testing.T) {
	// prepare
	a := NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTMap(
				thrift.STRING,
				thrift.BOOL,
				&map[TValue]TValue{
					NewTstring("key 1"): NewTBool(true),
					NewTstring("key 2"): NewTBool(false),
				},
			),
		},
	)
	var b TValue = NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTMap(
				thrift.STRING,
				thrift.BOOL,
				&map[TValue]TValue{
					NewTstring("key 1"): NewTBool(true),
					// different here
					NewTstring("key 2"): NewTBool(true),
				},
			),
		},
	)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_Struct_DifferentValueCount(t *testing.T) {
	// prepare
	a := NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTMap(
				thrift.STRING,
				thrift.BOOL,
				&map[TValue]TValue{
					NewTstring("key 1"): NewTBool(true),
					NewTstring("key 2"): NewTBool(false),
				},
			),
		},
	)
	var b TValue = NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTMap(
				thrift.STRING,
				thrift.BOOL,
				&map[TValue]TValue{
					NewTstring("key 1"): NewTBool(true),
					// different here
					NewTstring("key 2"): NewTBool(true),
				},
			),
			*NewTStructField(3, "name 3"): NewTBool(true),
		},
	)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_Struct_OtherKey(t *testing.T) {
	// prepare
	a := NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTBool(true),
		},
	)
	var b TValue = NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "NAME"):   NewTBool(true),
		},
	)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_Struct_OtherValueType(t *testing.T) {
	// prepare
	a := NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTMap(
				thrift.STRING,
				thrift.BOOL,
				&map[TValue]TValue{
					NewTstring("key 1"): NewTBool(true),
					NewTstring("key 2"): NewTBool(false),
				},
			),
		},
	)
	var b TValue = NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTBool(true),
		},
	)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestEquals_Struct_OtherType(t *testing.T) {
	// prepare
	a := NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTMap(
				thrift.STRING,
				thrift.BOOL,
				&map[TValue]TValue{
					NewTstring("key 1"): NewTBool(true),
					NewTstring("key 2"): NewTBool(false),
				},
			),
		},
	)
	var b TValue = NewTBool(true)
	expected := false

	// do
	actual := a.Equals(&b)

	// verify
	assert(t, "", actual, expected)
}

func TestWriteFieldData_TSruct(t *testing.T) {
	// prepare
	oprot := setupProtocol(t)
	cxt := context.Background()
	value := NewTStruct(
		&map[TStructField]TValue{
			*NewTStructField(1, "name 1"): NewTstring("value 1"),
			*NewTStructField(2, "name 2"): NewTMap(
				thrift.STRING,
				thrift.BOOL,
				&map[TValue]TValue{
					NewTstring("key 1"): NewTBool(true),
					NewTstring("key 2"): NewTBool(false),
				},
			),
		},
	)

	// do
	err := value.WriteFieldData(cxt, oprot)
	checkError(t, err)

	// verify
	oprot.Flush(cxt)
	{
		name, err := oprot.ReadStructBegin(cxt)
		checkError(t, err)
		assert(t, "key type", name, "")
	}
	{
		fname, ftype, fid, err := oprot.ReadFieldBegin(cxt)
		checkError(t, err)
		assert(t, "name 1 fname", fname, "")
		assert(t, "name 1 ftype", ftype, thrift.STRING)
		assert(t, "name 1 fid", fid, 1)
	}
	{
		v, err := oprot.ReadString(cxt)
		checkError(t, err)
		assert(t, "name 1 value", v, "value 1")
	}
	{
		err := oprot.ReadFieldEnd(cxt)
		checkError(t, err)
	}
	{
		fname, ftype, fid, err := oprot.ReadFieldBegin(cxt)
		checkError(t, err)
		assert(t, "name 2 fname", fname, "")
		assert(t, "name 2 ftype", ftype, thrift.MAP)
		assert(t, "name 2 fid", fid, 2)
	}
	{
		ktype, vtype, size, err := oprot.ReadMapBegin(cxt)
		checkError(t, err)
		assert(t, "name 2 ktype", ktype, thrift.STRING)
		assert(t, "name 2 vtype", vtype, thrift.BOOL)
		assert(t, "name 2 size", size, 2)
	}
	{
		k1, err := oprot.ReadString(cxt)
		checkError(t, err)
		assert(t, "name 2 key 1 key", k1, "key 1")
	}
	{
		v1, err := oprot.ReadBool(cxt)
		checkError(t, err)
		assert(t, "name 2 key 1 value", v1, true)
	}
	{
		k2, err := oprot.ReadString(cxt)
		checkError(t, err)
		assert(t, "name 2 key 2 key", k2, "key 2")
	}
	{
		v2, err := oprot.ReadBool(cxt)
		checkError(t, err)
		assert(t, "name 2 key 2 value", v2, false)
	}
	{
		err = oprot.ReadMapEnd(cxt)
		checkError(t, err)
	}
	{
		err := oprot.ReadFieldEnd(cxt)
		checkError(t, err)
	}
	{
		err := oprot.ReadStructEnd(cxt)
		checkError(t, err)
	}
}

func TestReadStruct_StringToBool(t *testing.T) {
	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	{
		checkError(
			t,
			iprot.WriteStructBegin(cxt, "name"),
		)
		checkError(
			t,
			iprot.WriteFieldBegin(cxt, "name 1", thrift.STRING, 1),
		)
		checkError(
			t,
			iprot.WriteString(cxt, "key 1"),
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
			iprot.Flush(cxt),
		)
	}

	// do
	actual, err := ReadStruct(cxt, iprot)
	checkError(t, err)

	// verfiy
	a, ok := actual.(*TStruct)
	assertTrue(t, "cast to TStruct", ok)
	{
		assert(t, "size", len((*a).value), 1)
	}
	{
		tv := (*a).value[*NewTStructField(1, "")]
		assertTrue(t, "value is not nil", tv != nil)
		v, ok := tv.(TString)
		assertTrue(t, "cast name 1 value to TString", ok)
		assert(t, "name 1 value", v.value, "key 1")
	}
}

func TestReadStruct_ContainerData(t *testing.T) {
	t.Skip("FIXME: Somehow this test fails. will be re-enabled.")

	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	{
		checkError(
			t,
			iprot.WriteStructBegin(cxt, "name"),
		)
		checkError(
			t,
			iprot.WriteFieldBegin(cxt, "name 1", thrift.STRING, 1),
		)
		checkError(
			t,
			iprot.WriteString(cxt, "key 1"),
		)
		checkError(
			t,
			iprot.WriteFieldEnd(cxt),
		)

		checkError(
			t,
			iprot.WriteFieldBegin(cxt, "name 2", thrift.MAP, 2),
		)
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
			iprot.Flush(cxt),
		)
	}

	// do
	actual, err := ReadStruct(cxt, iprot)
	checkError(t, err)

	// verfiy
	a, ok := actual.(*TStruct)
	assertTrue(t, "cast to TStruct", ok)
	{
		assertTrue(t, "size", len((*a).value) == 2)
	}
	{
		tv := (*a).value[*NewTStructField(1, "")]
		v, ok := tv.(TString)
		assertTrue(t, "cast name 1 value to TString", ok)
		assert(t, "name 1 value", v.value, "name")
	}
	tv := (*a).value[*NewTStructField(2, "")]
	v, ok := tv.(*TMap)
	{
		assertTrue(t, "cast name 2 to TMap", ok)
		assert(t, "name 2 map size", len(v.value), 2)
	}
	{
		value := (*v).value[NewTstring("key 1")]
		v1, ok := value.(TBool)
		assertTrue(t, "cast key 1 in name 2 to TBool", ok)
		assert(t, "key 1 in name 2", v1.value, true)
	}
	{
		value := (*v).value[NewTstring("key 2")]
		v2, ok := value.(TBool)
		assertTrue(t, "cast key 2 in name 2 to TBool", ok)
		assert(t, "key 2 in name 2", v2.value, true)
	}
}

func TestReadStruct_Invalid_type(t *testing.T) {
	// prepare
	iprot := setupProtocol(t)
	cxt := context.Background()
	{
		var err error
		err = iprot.WriteString(cxt, "other")
		checkError(t, err)
	}

	// do
	_, err := ReadString(cxt, iprot)

	// verfiy
	assertTrue(t, "", err != nil)
}
