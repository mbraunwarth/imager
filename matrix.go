package main

import "strconv"

// Pixel type as placeholder for size of pixel
type Pixel uint16

// Row type alias for initialization of ImageMatrix.
type Row []Pixel

func (r Row) Sum() uint64 {
	var sum uint64
	for _, p := range r {
		sum = sum + uint64(p)
	}
	return sum
}

// Matrix for Netpbm images.
// Netpbm only supports color values up to 2^16 (0-65535).
type Matrix []Row

func (m *Matrix) Sum() uint64 {
	var sum uint64
	for _, r := range *m {
		sum = sum + r.Sum()
	}
	return sum
}

func (m *Matrix) Get(x, y uint32) Pixel {
	return (*m)[x][y]
}

func (m *Matrix) Set(x, y uint32, p Pixel) {
	(*m)[x][y] = p
}

func fillMatrix(w, h uint32, toks <-chan string) Matrix {
	im := make(Matrix, h)
	var i, j uint32
	for i = 0; i < h; i++ {
		im[i] = make(Row, w)
		for j = 0; j < w; j++ {
			p, _ := strconv.ParseUint(<-toks, ConvUintBase, ConvUintBitSize)
			im.Set(i, j, Pixel(p))
		}
	}

	return im
}
