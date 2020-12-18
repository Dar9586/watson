package any

import (
	"testing"

	"github.com/genkami/watson/pkg/vm"

	"github.com/google/go-cmp/cmp"
)

func TestFromValueConvertsInt(t *testing.T) {
	val := vm.NewIntValue(123)
	var want interface{} = int64(123)
	got := FromValue(val)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestFromValueConvertsFloat(t *testing.T) {
	val := vm.NewFloatValue(1.23)
	var want interface{} = float64(1.23)
	got := FromValue(val)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestFromValueConvertsString(t *testing.T) {
	val := vm.NewStringValue([]byte("hey"))
	var want interface{} = "hey"
	got := FromValue(val)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestFromValueConvertsObject(t *testing.T) {
	val := vm.NewObjectValue(map[string]*vm.Value{
		"name": vm.NewStringValue([]byte("Taro")),
		"age":  vm.NewIntValue(25),
	})
	var want interface{} = map[string]interface{}{
		"name": "Taro",
		"age":  int64(25),
	}
	got := FromValue(val)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestFromValueConvertsArray(t *testing.T) {
	val := vm.NewArrayValue([]*vm.Value{
		vm.NewStringValue([]byte("Yo")),
		vm.NewIntValue(123),
	})
	var want interface{} = []interface{}{"Yo", int64(123)}
	got := FromValue(val)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestFromValueConvertsBool(t *testing.T) {
	val := vm.NewBoolValue(true)
	var want interface{} = true
	got := FromValue(val)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestFromValueConvertsNil(t *testing.T) {
	val := vm.NewNilValue()
	var want interface{} = nil
	got := FromValue(val)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestToValueConvertsNil(t *testing.T) {
	want := vm.NewNilValue()
	got := ToValue(nil)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestToValueConvertsTrue(t *testing.T) {
	want := vm.NewBoolValue(true)
	got := ToValue(true)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestToValueConvertsFalse(t *testing.T) {
	want := vm.NewBoolValue(false)
	got := ToValue(false)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
