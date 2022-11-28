package main

import (
	"fmt"
	"strconv"
)

// ImageMatrix for Netpbm images.
// Netpbm only supports color values up to 2^16 (0-65535).
type ImageMatrix []ImageRow

// ImageRow type alias for initialization of ImageMatrix.
type ImageRow []uint16

// Image
type Image struct {
	Name        string
	MagicNumber string
	Width       uint32
	Height      uint32
	MaxColor    uint16
	IM          ImageMatrix
}

func fillMatrix(w, h uint32, toks <-chan string) ImageMatrix {
	im := make(ImageMatrix, h)
	var i, j uint32
	for i = 0; i < h; i++ {
		im[i] = make(ImageRow, w)
		for j = 0; j < w; j++ {
			// TODO move magic numbers to constants
			p, _ := strconv.ParseUint(<-toks, 10, 16)
			im[i][j] = uint16(p)
		}
	}

	return im
}

// TODO Proper error handling

func Load(imgPath string) *Image {
	toks := make(chan string)
	sc := NewScanner(imgPath)

	go sc.ScanImage(toks)

	// extraxt magic number
	mn := <-toks

	// extract width, height and maximum color value
	w, _ := strconv.Atoi(<-toks)
	h, _ := strconv.Atoi(<-toks)
	max, _ := strconv.Atoi(<-toks)

	// fill matrix
	im := fillMatrix(uint32(w), uint32(h), toks)

	return &Image{
		Name:        imgPath,
		MagicNumber: mn,
		Width:       uint32(w),
		Height:      uint32(h),
		MaxColor:    uint16(max),
		IM:          im,
	}
}

func ShowMetadata(i *Image) {
	fmt.Println(i.Name)
	fmt.Println(i.MagicNumber)
	fmt.Println(i.Height)
	fmt.Println(i.Width)
	fmt.Println(i.MaxColor)
}

func ShowMatrix(img *Image) {
	var i, j uint32
	for i = 0; i < img.Height; i++ {
		for j = 0; j < img.Width; j++ {
			fmt.Print(img.IM[i][j], " ")
		}
		fmt.Println()
	}
}
