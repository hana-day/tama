package types

import "testing"

func TestPairString(t *testing.T) {
	testcases := []struct {
		object Object
		str    string
	}{
		{Cons(Number(1), Number(2)), "(1 . 2)"},
		{Cons(Number(1), Cons(Number(2), NilObject)), "(1 . (2 . ()))"},
	}

	for i, tc := range testcases {
		actual := tc.object.String()
		if actual != tc.str {
			t.Fatalf("case %d: expected %s, but got %s", i, tc.str, actual)
		}
	}

}

func TestList(t *testing.T) {
	if List().Type() != TyNil {
		t.Fatalf("expected nil")
	}

	l, ok := List(Number(1), Number(2)).(*Pair)
	if !ok {
		t.Fatalf("expected pair")
	}
	cdar, err := l.Cdar()
	if err != nil {
		t.Fatal(err)
	}
	if cdar.String() != "2" {
		t.Fatalf("expected %v, but got %v", "2", cdar.String())
	}
	cddr, err := l.Cddr()
	if err != nil {
		t.Fatal(err)
	}
	if cddr.Type() != TyNil {
		t.Fatalf("expected nil")
	}

}

func TestLen(t *testing.T) {
	list, _ := List(Number(1), Number(2)).(*Pair)
	if l := list.Len(); l != 2 {
		t.Fatalf("expected %d, but got %d", l, 2)
	}
}
