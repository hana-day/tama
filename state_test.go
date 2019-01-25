package tama

import (
	"github.com/hyusuk/tama/types"
	"testing"
)

func TestExecuteString(t *testing.T) {
	// default state
	state := NewState(Option{})

	testcases := []struct {
		stateFactory func() *State
		source       string
		resultString string
	}{
		{
			func() *State { return state },
			" 1 ",
			"1",
		},
		{
			func() *State { s := NewState(Option{}); s.Global["test"] = types.Number(2); return s },
			"test",
			"2",
		},
		{
			func() *State { return state },
			"(+ 1 2 3 4)",
			"10",
		},
		{
			func() *State { return state },
			"(+ 1 (+ 2 3) 4)",
			"10",
		},
		{
			func() *State { return state },
			"(define a 1) (+ a 1)",
			"2",
		},
		{
			func() *State { return state },
			"((lambda (a) (+ 1 2 a)) 3)",
			"6",
		},
		{
			func() *State { return state },
			"(begin 1 2 3 4)",
			"4",
		},
		{
			func() *State { return state },
			"(((lambda (a) (lambda (b) (+ a b))) 1) 2)",
			"3",
		},
		{
			func() *State { return state },
			"((((lambda (a) (lambda (b) (lambda (c) (+ a b c)))) 1) 2) 3)",
			"6",
		},
		{
			func() *State { return state },
			"((lambda (a) (+ a 1) (+ a 2)) 1)",
			"3",
		},
		{
			func() *State { return state },
			"((lambda (a) (+ (car a) (cdr a))) (cons 1 2))",
			"3",
		},
		{
			func() *State { return state },
			"(set! a 1) a",
			"1",
		},
		{
			func() *State { return state },
			"((lambda (a) (set! a 2) (+ 1 a)) 1)",
			"3",
		},
		{
			func() *State { return state },
			"(((lambda (a) (lambda (b) (set! a 1) (+ a b))) 100) 2)",
			"3",
		},
		{
			func() *State { return state },
			"(define a 1) (set! a 2)",
			types.UndefinedObject.String(),
		},
		{
			func() *State { return state },
			"(define (test a) (+ a 1) (+ a 2)) (test 1)",
			"3",
		},
		{
			func() *State { return state },
			"(car (quote (1 2 3)))",
			"1",
		},
		{
			func() *State { return state },
			"(car '(1 2 3))",
			"1",
		},
		{
			func() *State { return state },
			"#t",
			"#t",
		},
		{
			func() *State { return state },
			"#f",
			"#f",
		},
		{
			func() *State { return state },
			"'#t",
			"#t",
		},
		{
			func() *State { return state },
			"(if #f 1 2)",
			"2",
		},
		{
			func() *State { return state },
			"(define a 1) (if #f (set! a 2) (set! a 100)) a",
			"100",
		},
		{
			func() *State { return state },
			"(define a 1) (if #t (set! a 2) (set! a 100)) a",
			"2",
		},
		{
			func() *State { return state },
			"(define a 1) (if #f (set! a 2)) a",
			"1",
		},
		{
			func() *State { return state },
			"(define a 1) (if 1 (set! a 2)) a",
			"2",
		},
		{
			func() *State { return state },
			"(if #f 1)",
			types.UndefinedObject.String(),
		},
		{
			func() *State { return state },
			"((lambda args (+ (car args) 100)) 1 2 3)",
			"101",
		},
		{
			func() *State { return state },
			"((lambda (a b . rest) (+ a b (car rest))) 1 2 3 4)",
			"6",
		},
		{
			func() *State { return state },
			"(define (a . rest) (car rest)) (a 1 2 3)",
			"1",
		},
		{
			func() *State { return state },
			"(= 1 1 1 1)",
			"#t",
		},
		{
			func() *State { return state },
			"(= 1 1 2 1)",
			"#f",
		},
		{
			func() *State { return state },
			"(define (factorial n) (if (= n 1) 1 (* n (factorial (- n 1))))) (factorial 3)",
			"6",
		},

		// if tail call doesn't work, tests below will fail due to the insufficient stack size.
		{
			func() *State { return NewState(Option{StackSize: 100}) },
			"(define (recur n) (if (= n 1) 1 (recur (- n 1)))) (recur 100)",
			"1",
		},
		{
			func() *State { return NewState(Option{StackSize: 100}) },
			"(define (recur a) (if (= a 1) 1 (begin (recur (- a 1))))) (recur 100)",
			"1",
		},
		{
			func() *State { return NewState(Option{StackSize: 100}) },
			"(define (factorial n init) (if (= n 1) init (factorial (- n 1) (* init n)))) (factorial 100 1)",
			"9.332621544394418e+157",
		},
	}
	for i, tc := range testcases {
		s := tc.stateFactory()
		if err := s.ExecString(tc.source); err != nil {
			t.Fatalf("case %d: unexpected error %v ;  source: %s", i, err, tc.source)
		}
		s.CallStack.Dump()
		top := s.CallStack.Top()
		v, ok := top.(types.Object)
		if !ok {
			t.Fatalf("case %d: unsupported object %v ; source: %s", i, top, tc.source)
		}
		if v.String() != tc.resultString {
			t.Fatalf("case %d: expected %s, but got %s ; source %s", i, tc.resultString, v.String(), tc.source)
		}
	}
}
