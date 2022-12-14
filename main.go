package main

// Scanner for Netpbm image file formats ASCII only.

// Netpbm data structure holding necessary information for the scanned image.
type Netpbm struct {
	matrix   ImageMatrix
	height   uint64
	width    uint64
	maxColor uint8
	format   NetpbmFormat
}

type NetpbmFormat uint8

// Note: No support for binary data format P4 - P6
const (
	PBM NetpbmFormat = iota // Portable BitMap - P1
	PGM                     // Portable GrayMap - P2
	PPM                     // Portable PixMap - P3
)

func main() {
	imgPath := "./assets/lena.ascii.pgm"

	// run scanner on image and fill image matrix
	//img := Load(imgPath)
	img := Load(imgPath)
	ShowMetadata(img)
	ShowMatrix(img)
}
