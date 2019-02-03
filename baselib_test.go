package tama

import (
	"github.com/hyusuk/tama/types"

	"testing"
)

type tcase struct {
	src       string
	expect    string
	option    Option
	expectErr bool
}

func (tc *tcase) fail(t *testing.T, caseNo int, actual string) {
	t.Fatalf("case %d: expected %s, but got %s\nsrc: %s", caseNo, tc.expect, actual, tc.src)
}

func (tc *tcase) error(t *testing.T, caseNo int, err error) {
	t.Fatalf("case %d: unexpected error %v\nsrc: %s", caseNo, err, tc.src)
}

func (tc *tcase) noerror(t *testing.T, caseNo int) {
	t.Fatalf("case %d: expected error, but got no error\nsrc: %s", caseNo, tc.src)
}

func testTcases(t *testing.T, tcases []*tcase) {
	for i, tc := range tcases {
		s := NewState(tc.option)
		err := s.ExecString(tc.src)
		if tc.expectErr {
			if err == nil {
				tc.noerror(t, i)
			}
			continue
		}
		if err != nil {
			tc.error(t, i, err)
		}
		top := s.CallStack.Top()
		v := top.(types.Object)
		if v.String() != tc.expect {
			tc.fail(t, i, v.String())
		}
	}
}

func TestDefine(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(define x 1) (define x (+ 2 x)) x", expect: "3"},
		&tcase{src: "(define (a . rest) (car rest)) (a 1 2 3)", expect: "1"},
	}
	testTcases(t, tcases)
}

func TestLambda(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "((lambda (x) (car x)) '(5 6 7))", expect: "5"},
		&tcase{src: "((lambda (x y) ((lambda (z) (* (car z) (cdr z))) (cons x y))) 3 4)", expect: "12"},
		&tcase{src: "((lambda (a) (set! a 2) (+ 1 a)) 1)", expect: "3"},
		&tcase{src: "(((lambda (a) (lambda (b) (set! a 1) (+ a b))) 100) 2)", expect: "3"},
		&tcase{src: "((lambda args (+ (car args) 100)) 1 2 3)", expect: "101"},
		&tcase{src: "((lambda (a b . rest) (+ a b (car rest))) 1 2 3 4)", expect: "6"},
	}
	testTcases(t, tcases)
}

func TestBegin(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(begin 1 2 3 4)", expect: "4"},
	}
	testTcases(t, tcases)
}

func TestSet(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(set! a 1)", expect: types.UndefinedObject.String()},
		&tcase{src: "(set! a 2) a", expect: "2"},
		&tcase{src: "(define a 1) (set! a 2) a", expect: "2"},
		&tcase{src: "((lambda (a) (set! a 2) (+ 1 a)) 1)", expect: "3"},
		&tcase{src: "(((lambda (a) (lambda (b) (set! a 1) (+ a b))) 100) 2)", expect: "3"},
	}
	testTcases(t, tcases)
}

func TestQuote(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(car (quote (1 2 3)))", expect: "1"},
		&tcase{src: "(car '(1 2 3))", expect: "1"},
		&tcase{src: "'1", expect: "1"},
		&tcase{src: "'#t", expect: "#t"},
	}
	testTcases(t, tcases)
}

func TestIf(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(if #f 1 2)", expect: "2"},
		&tcase{src: "(define a 1) (if #f (set! a 2) (set! a 100)) a", expect: "100"},
		&tcase{src: "(define a 1) (if 0 (set! a 2)) a", expect: "2"},
		&tcase{src: "(if #f 1)", expect: types.UndefinedObject.String()},
	}
	testTcases(t, tcases)
}

func TestCallCC(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(call/cc (lambda (cc) (cc 3) 5))", expect: "3"},
		&tcase{src: "((lambda (x) (call/cc (lambda (cc) (cc x)))) 3)", expect: "3"},
		&tcase{src: "((lambda (x) (call/cc (lambda (cc) (cc x) 5))) 3)", expect: "3"},
		&tcase{src: "((lambda (z) ((lambda (x y) (call/cc (lambda (cc) (cc x)))) 3 4)) 100)", expect: "3"},
		&tcase{src: "((lambda (z) ((lambda (x y) (call/cc (lambda (cc) (cc x) 99))) 3 4)) 100)", expect: "3"},
		&tcase{src: "((lambda (z a) ((lambda (x y) (call/cc (lambda (cc) (cc x) 99))) 3 4)) 100 99)", expect: "3"},
	}
	testTcases(t, tcases)
}

func TestFnCar(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(car '(a b c))", expect: "a"},
	}
	testTcases(t, tcases)
}

func TestFnAdd(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(+)", expect: "0"},
		&tcase{src: "(+ 1)", expect: "1"},
		&tcase{src: "(+ 1 2 3 4)", expect: "10"},
		&tcase{src: "(+ 1 'a)", expectErr: true},
	}
	testTcases(t, tcases)
}

func TestFnSub(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(- 1)", expect: "-1"},
		&tcase{src: "(- 3 4)", expect: "-1"},
		&tcase{src: "(- 3 4 5)", expect: "-6"},
		&tcase{src: "(- 1 'a)", expectErr: true},
		&tcase{src: "(-)", expectErr: true},
	}
	testTcases(t, tcases)
}

func TestFnMul(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(*)", expect: "1"},
		&tcase{src: "(* 4)", expect: "4"},
		&tcase{src: "(* 4 5)", expect: "20"},
		&tcase{src: "(* 1 'a)", expectErr: true},
	}
	testTcases(t, tcases)
}

func TestFnDiv(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(/ 3)", expect: "0.3333333333333333"},
		&tcase{src: "(/ 4 2 2)", expect: "1"},
		&tcase{src: "(/ 1 'a)", expectErr: true},
	}
	testTcases(t, tcases)
}

func TestFnNumEq(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(= 3 3)", expect: "#t"},
		&tcase{src: "(= 3 3 3)", expect: "#t"},
		&tcase{src: "(= 3 3 4)", expect: "#f"},
		&tcase{src: "(= 1 'a)", expectErr: true},
	}
	testTcases(t, tcases)
}

func TestFnComp(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(< 2 3)", expect: "#t"},
		&tcase{src: "(< 2 3 4)", expect: "#t"},
		&tcase{src: "(< 2 3 3)", expect: "#f"},
		&tcase{src: "(> 3 2)", expect: "#t"},
		&tcase{src: "(> 3 2 1)", expect: "#t"},
		&tcase{src: "(> 3 2 2)", expect: "#f"},
		&tcase{src: "(<= 3 4)", expect: "#t"},
		&tcase{src: "(<= 3 3 3)", expect: "#t"},
		&tcase{src: "(<= 3 3 2)", expect: "#f"},
		&tcase{src: "(>= 3 2)", expect: "#t"},
		&tcase{src: "(>= 3 2 2)", expect: "#t"},
		&tcase{src: "(>= 3 2 3)", expect: "#f"},
		&tcase{src: "(< 3 2 'a)", expectErr: true},
	}
	testTcases(t, tcases)
}

// 6.3.5. Strings
func TestFnStrLen(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(string-length \"test\")", expect: "4"},
		&tcase{src: "(string-length \"\")", expect: "0"},
		&tcase{src: "(string-length 1)", expectErr: true},
	}
	testTcases(t, tcases)
}

// 6.3.6 Vectors
func TestFnVecRef(t *testing.T) {
	tcases := []*tcase{
		&tcase{src: "(vector-ref #(1 2 3 4 5) 1)", expect: "2"},
		&tcase{src: "(vector-ref '(1 2 3 4 5) 1)", expectErr: true},
		&tcase{src: "(vector-ref #(1 2 3 4 5) 5)", expectErr: true},
	}
	testTcases(t, tcases)
}
