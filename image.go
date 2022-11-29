package main

import (
	"fmt"
	"os"
	"strconv"
)

// TODO rething typing of Pixel as int

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

func Save(img *Image, fname string) {
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// write meta data
	metadata := fmt.Sprintf("%s\n%d %d\n%d\n", img.MagicNumber, img.Width, img.Height, img.MaxColor)
	_, err = f.WriteString(metadata)
	if err != nil {
		panic(err)
	}

	// write image matrix
	var i, j uint32
	for i = 0; i < img.Height; i++ {
		for j = 0; j < img.Width; j++ {
			_, err := f.WriteString(fmt.Sprintf("%d ", img.IM[i][j]))
			if err != nil {
				panic(err)
			}
		}
		_, err := f.WriteString("\n")
		if err != nil {
			panic(err)
		}
	}
}

func ShowMetadata(i *Image) {
	fmt.Println(i.Name)
	fmt.Println(i.MagicNumber)
	fmt.Println(i.Width)
	fmt.Println(i.Height)
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

func (img *Image) Copy() *Image {
	outMatrix := make(ImageMatrix, img.Height)

	var i, j uint32
	for i = 0; i < img.Height; i++ {
		outMatrix[i] = make(ImageRow, img.Width)
		for j = 0; j < img.Width; j++ {
			outMatrix[i][j] = img.IM[i][j]
		}
	}

	return &Image{
		Name:        "out-" + img.Name,
		MagicNumber: img.MagicNumber,
		Width:       img.Width,
		Height:      img.Height,
		MaxColor:    img.MaxColor,
		IM:          outMatrix,
	}
}

func (img *Image) Blur() *Image {
	gaussian := ImageMatrix{
		ImageRow{1, 4, 7, 4, 1},
		ImageRow{4, 16, 26, 16, 4},
		ImageRow{7, 26, 41, 26, 4},
		ImageRow{4, 16, 26, 16, 4},
		ImageRow{1, 4, 7, 4, 1},
	}
	//gaussian := ImageMatrix{
	//	ImageRow{1, 2, 1},
	//	ImageRow{2, 4, 2},
	//	ImageRow{1, 2, 1},
	//}

	return img.Conv(gaussian)
}

// TODO rewrite convolution and ImageMatrix for use as kernel
func (img *Image) Conv(kernel ImageMatrix) *Image {
	out := img.Copy()

	kSize := len(kernel)

	for i := kSize / 2; i < int(img.Height-uint32(kSize/2)); i++ {
		for j := kSize / 2; j < int(img.Width-uint32(kSize/2)); j++ {
			var p uint16
			for m := -kSize / 2; m <= kSize/2; m++ {
				for n := -kSize / 2; n <= kSize/2; n++ {
					p = p + img.IM[uint32(i+m)][uint32(j+n)]*kernel[m+kSize/2][n+kSize/2]
				}
			}

			out.IM[i][j] = p / 273
		}
	}

	return out
}
