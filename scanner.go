package main

import (
	"io"
	"os"
	"unicode"
)

type Scanner struct {
	source string
	r      rune
	offset int
	line   int
	column int
}

func NewScanner(imagePath string) *Scanner {
	f, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	src := string(bs)

	return &Scanner{
		source: src,
		r:      rune(src[0]),
		offset: 0,
		line:   1,
		column: 1,
	}
}

func (sc *Scanner) ScanImage(ch chan<- string) {
	for {
		switch {
		case unicode.IsLetter(sc.r) || unicode.IsNumber(sc.r):
			// in number/string to parse
			start := sc.offset

			for next := sc.peek(); ; next = sc.peek() {
				if !unicode.IsNumber(next) && !unicode.IsLetter(next) {
					break
				}

				if !sc.advance() {
					break
				}
			}

			ch <- sc.source[start : sc.offset+1]
			// ignore
		case sc.r == '\n':
			sc.line++
		case sc.r == '#':
			// rest of line will be comment
			for sc.peek() != '\n' {
				if !sc.advance() {
					break
				}
			}
		}

		if !sc.advance() {
			break
		}
	}
	close(ch)
}

func (sc *Scanner) advance() bool {
	if sc.offset+1 >= len(sc.source) {
		return false
	}

	sc.offset++
	sc.column++
	sc.r = rune(sc.source[sc.offset])
	return true
}

func (sc *Scanner) peek() rune {
	if sc.offset+1 >= len(sc.source) {
		var r rune
		return r
	}

	return rune(sc.source[sc.offset+1])
}
