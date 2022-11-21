package main

import "fmt"

type ImageMatrix [][]uint32

type Netpbm struct {
	Extension   string
	MagicNumber uint8
	Width       uint32
	Height      uint32
	Matrix      ImageMatrix
}

// PBM portable bitmap format. Only containing black or white pixel. Zero represents
// white, otherwise black.
type PBM struct {
	Netpbm
}

type PPM struct{}
type PGM struct{}
type PNM struct{}

func main() {
	var m ImageMatrix = make()
	m[0] = make([]uint32, 0)
	fmt.Println(m)
}
