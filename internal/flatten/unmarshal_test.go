package flatten

import (
	"reflect"
	"testing"
)

type demo struct {
	Foo    string
	Bar    int64
	Nested nested
}

type nested struct {
	Baz bool
}

func TestUnmarshal(t *testing.T) {
	in := demo{
		Foo:    "Hello World",
		Bar:    10,
		Nested: nested{Baz: true},
	}

	props, err := Marshal(in)
	if err != nil {
		t.Fatal(err)
	}

	var out demo
	if err = Unmarshal(props, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fail()
	}
}
