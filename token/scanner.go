package token

import (
	"log"
)

type Scanner struct {
	src     []byte
	ch      byte
	forward int
	offset  int
}

func (s *Scanner) next() {
	if s.forward < len(s.src) {
		s.offset = s.forward
		s.ch = s.src[s.forward]
		s.forward += 1
	} else {
		s.offset = len(s.src)
		s.ch = 0
	}
}

func (s *Scanner) Init(src []byte) {
	s.src = src
	s.ch = ' '
	s.forward = 0
	s.offset = 0
	s.next()
}

func (s *Scanner) skipSpaces() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) scanNumber() (Token, string) {
	off := s.offset
	for s.ch >= '0' && s.ch <= '9' {
		s.next()
	}
	return INT, string(s.src[off:s.offset])
}

func (s *Scanner) Scan() (tok Token, lit string) {
	s.skipSpaces()
	ch := s.ch
	if ch >= '0' && ch <= '9' {
		return s.scanNumber()
	}
	if ch == 0 {
		tok = EOF
		return
	}
	log.Fatalf("Unknown character %c", ch)
	return
}
