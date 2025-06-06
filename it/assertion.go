package it

import (
	"fmt"
	"testing"

	xk6_thrift "github.com/lavenderses/xk6-thrift"
)

func assertEquals(t *testing.T, actual, expect xk6_thrift.TResponse) {
	t.Logf("Got: %v, expected: %v", actual, expect)
	avs := *actual.Values()
	evs := *expect.Values()
	if len(avs) != len(evs) {
		t.Errorf("expected size %d, but was %d", len(evs), len(avs))
		return
	}

	failure := false
	msg := "[FAILED]"
	for ek, ev := range evs {
		av := avs[ek]
		if !av.Equals(&ev) {
			msg = fmt.Sprintf("%s\nexpected %v, but was %v", msg, ev, av)
			failure = true
		}
	}
	for ak, av := range avs {
		ev := evs[ak]
		if ev == nil {
			msg = fmt.Sprintf("%s\ngo unexpected value %v", msg, av)
			failure = true
		}
	}

	if failure {
		t.Error(msg)
	}
}
