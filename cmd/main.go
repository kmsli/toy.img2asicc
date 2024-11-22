package main

import (
	img2ascii "toy.img2ascii/internal"
	"toy.img2ascii/internal/gray"
)

func main() {
	img2ascii.Webp2Ascii("test.webp", gray.Sobel)
	// img2ascii.Webp2Ascii("test.webp", gray.Canny)
}
