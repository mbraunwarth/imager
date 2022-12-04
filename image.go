package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	ConvUintBase    = 10
	ConvUintBitSize = 16
)

// Image
type Image struct {
	Name        string
	MagicNumber string
	Width       uint32
	Height      uint32
	MaxColor    uint16
	Mtx         Matrix
}

func (img *Image) Get(x, y uint32) Pixel {
	return img.Mtx.Get(x, y)
}

func (img *Image) Set(x, y uint32, p Pixel) {
	img.Mtx.Set(x, y, p)
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
		Mtx:         im,
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
			_, err := f.WriteString(fmt.Sprintf("%d ", img.Get(i, j)))
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

func ShowMetadata(img *Image) {
	fmt.Println(img.Name)
	fmt.Println(img.MagicNumber)
	fmt.Println(img.Width)
	fmt.Println(img.Height)
	fmt.Println(img.MaxColor)
}

func ShowMatrix(img *Image) {
	var i, j uint32
	for i = 0; i < img.Height; i++ {
		for j = 0; j < img.Width; j++ {
			fmt.Print(img.Get(i, j), " ")
		}
		fmt.Println()
	}
}

// Copy returns a deep copy of the image applied to.
func (img *Image) Copy() *Image {
	outMatrix := make(Matrix, img.Height)

	var i, j uint32
	for i = 0; i < img.Height; i++ {
		outMatrix[i] = make(Row, img.Width)
		for j = 0; j < img.Width; j++ {
			outMatrix.Set(i, j, img.Get(i, j))
		}
	}

	return &Image{
		Name:        "out-" + img.Name,
		MagicNumber: img.MagicNumber,
		Width:       img.Width,
		Height:      img.Height,
		MaxColor:    img.MaxColor,
		Mtx:         outMatrix,
	}
}

func (img *Image) Blur() *Image {
	gaussian := Matrix{
		Row{1, 4, 7, 4, 1},
		Row{4, 16, 26, 16, 4},
		Row{7, 26, 41, 26, 4},
		Row{4, 16, 26, 16, 4},
		Row{1, 4, 7, 4, 1},
	}
	//gaussian := ImageMatrix{
	//	Row{1, 2, 1},
	//	Row{2, 4, 2},
	//	Row{1, 2, 1},
	//}

	return img.Conv(gaussian)
}

// TODO rewrite convolution and ImageMatrix for use as kernel
// Conv convolves the image `img` with the given `kernel`.
// The `kernel` must be of odd, quadratic shape e.g. 3x3, 5x5, 17x17 etc.
func (img *Image) Conv(kernel Matrix) *Image {
	if len(kernel)%2 == 0 || len(kernel) != len(kernel[0]) {
		panic("kernel must be of odd, quadratic shape like 3x3, 5x5, 17x17 etc.")
	}

	out := img.Copy()

	kSize := len(kernel)

	// sum kernel for later scaling since the kernel may only consist of integer values
	// less than fractions
	kSum := kernel.Sum()

	for i := kSize / 2; i < img.Height-(kSize/2); i++ {
		for j := uint32(kSize / 2); j < img.Width-uint32(kSize/2); j++ {
			var p Pixel
			for m := -kSize / 2; m <= kSize/2; m++ {
				for n := -kSize / 2; n <= kSize/2; n++ {
					p = p + img.Get(i+m, j+n)*kernel.Get(m+kSize/2, n+kSize/2)
					fmt.Println(p)
				}
			}

			out.Set(i, j, p/Pixel(kSum))
		}
	}

	return out
}
