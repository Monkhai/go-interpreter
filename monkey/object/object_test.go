package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() == hello1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	eq1 := &Integer{Value: 1}
	eq2 := &Integer{Value: 1}
	diff1 := &Integer{Value: 2}
	diff2 := &Integer{Value: 2}

	if eq1.HashKey() != eq2.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}
	if diff1.HashKey() == eq1.HashKey() {
		t.Errorf("integers with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	eq1 := &Boolean{Value: true}
	eq2 := &Boolean{Value: true}
	diff1 := &Boolean{Value: false}
	diff2 := &Boolean{Value: false}

	if eq1.HashKey() != eq2.HashKey() {
		t.Errorf("booleans with same content have different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("booleans with same content have different hash keys")
	}
	if diff1.HashKey() == eq1.HashKey() {
		t.Errorf("booleans with different content have same hash keys")
	}
}
