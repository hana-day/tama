package types

import (
	"testing"
)

func TestSlice(t *testing.T) {
	arr, err := List(Number(1), Number(2), Number(3)).(*Pair).Slice()
	if err != nil {
		t.Fatal(err)
	}
	if len(arr) != 3 {
		t.Fatalf("expected %d, but got %d", 3, len(arr))
	}
	if arr[1].String() != "2" {
		t.Fatalf("expected %s, but got %s", "2", arr[1].String())
	}
	arr, err = Cons(Number(1), Number(2)).Slice()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCarCdr(t *testing.T) {
	testcases := []struct {
		pair        *Pair
		expectError bool
		expect      string
		fn          func(*Pair) (Object, error)
	}{
		{
			Cons(Number(1), Number(2)),
			false,
			"1",
			func(p *Pair) (Object, error) { return p.First(), nil },
		},
		{
			Cons(Number(1), Number(2)),
			false,
			"2",
			func(p *Pair) (Object, error) { return p.Cdr(), nil },
		},
		{
			Cons(Number(1), Number(2)),
			true,
			"",
			func(p *Pair) (Object, error) { return p.Second() },
		},
		{
			Cons(Number(1), Cons(Number(2), NilObject)),
			false,
			"2",
			func(p *Pair) (Object, error) { return p.Second() },
		},
		{
			Cons(Number(1), Cons(Number(2), NilObject)),
			true,
			"",
			func(p *Pair) (Object, error) { return p.Third() },
		},
		{
			Cons(Number(1), Cons(Number(2), Cons(Number(3), NilObject))),
			false,
			"3",
			func(p *Pair) (Object, error) { return p.Third() },
		},
		{
			Cons(Number(1), Cons(Number(2), Cons(Number(3), Number(4)))),
			false,
			"4",
			func(p *Pair) (Object, error) { return p.Cdddr() },
		},
	}
	for i, tc := range testcases {
		obj, err := tc.fn(tc.pair)
		if tc.expectError {
			if err == nil {
				t.Fatalf("case %d: expected error", i)
			}
			continue
		}
		if err != nil {
			t.Fatalf("case %d: unexpected error %v", i, err)
		}
		if obj.String() != tc.expect {
			t.Fatalf("case %d: expected %s, but got %s", i, tc.expect, obj.String())
		}
	}
}
