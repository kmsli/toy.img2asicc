package gray

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

type sobelOperator struct{}

// Sobel 核
var sobelX = [][]int{
	{-1, 0, 1},
	{-2, 0, 2},
	{-1, 0, 1},
}

var sobelY = [][]int{
	{-1, -2, -1},
	{0, 0, 0},
	{1, 2, 1},
}

// 灰度图转换
func (op *sobelOperator) Apply(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, img, bounds.Min, draw.Src)
	return sobelEdgeDetection(gray)
}

// 应用 Sobel 算子
func sobelEdgeDetection(gray *image.Gray) *image.Gray {
	bounds := gray.Bounds()
	edgeImg := image.NewGray(bounds)

	for y := 1; y < bounds.Max.Y-1; y++ {
		for x := 1; x < bounds.Max.X-1; x++ {
			var gx, gy int
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					pixel := int(gray.GrayAt(x+j, y+i).Y)
					gx += pixel * sobelX[i+1][j+1]
					gy += pixel * sobelY[i+1][j+1]
				}
			}
			// 计算梯度幅值
			g := math.Sqrt(float64(gx*gx + gy*gy))
			if g > 255 {
				g = 255
			}
			edgeImg.SetGray(x, y, color.Gray{Y: uint8(g)})
		}
	}

	return edgeImg
}
