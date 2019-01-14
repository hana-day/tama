package types

import "testing"

func TestPairString(t *testing.T) {
	testcases := []struct {
		object Value
		str    string
	}{
		{Cons(Number(1), Number(2)), "(1 . 2)"},
		{Cons(Number(1), Cons(Number(2), Nil)), "(1 . (2 . ()))"},
	}

	for i, tc := range testcases {
		actual := tc.object.String()
		if actual != tc.str {
			t.Fatalf("case %d: expected %s, but got %s", i, tc.str, actual)
		}
	}

}
