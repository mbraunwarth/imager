package main

import (
	"fmt"
	"io"
	"os"
	"time"
	"unicode"
)

// Scanner for Netpbm image file formats ASCII only.

// Netpbm data structure holding necessary information for the scanned image.
type Netpbm struct {
	matrix   ImageMatrix
	height   uint64
	width    uint64
	maxColor uint8
	format   NetpbmFormat
}

// ImageMatrix
type ImageMatrix [][]int64

type NetpbmFormat uint8

// Note: No support for binary data format P4 - P6
const (
	PBM NetpbmFormat = iota // Portable BitMap - P1
	PGM                     // Portable GrayMap - P2
	PPM                     // Portable PixMap - P3
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

func main() {
	imgPath := "./assets/lena.ascii.pgm"

	toks := make(chan string)

	sc := NewScanner(imgPath)

	go sc.ScanImage(toks)
	time.Sleep(time.Millisecond * 100)

	//magicNumber := <-toks
	//fmt.Println("Magic Number", magicNumber)

	for t := range toks {
		fmt.Printf("-%s-\n", t)
	}
}
