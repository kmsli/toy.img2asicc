package gray

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

type cannyOperator struct{}

// 高斯滤波核（5x5标准高斯）
var gaussianKernel = [][]float64{
	{2, 4, 5, 4, 2},
	{4, 9, 12, 9, 4},
	{5, 12, 15, 12, 5},
	{4, 9, 12, 9, 4},
	{2, 4, 5, 4, 2},
}

// 高斯滤波
func gaussianBlur(gray *image.Gray) *image.Gray {
	bounds := gray.Bounds()
	blurred := image.NewGray(bounds)

	for y := 2; y < bounds.Max.Y-2; y++ {
		for x := 2; x < bounds.Max.X-2; x++ {
			var sum float64
			for i := -2; i <= 2; i++ {
				for j := -2; j <= 2; j++ {
					weight := gaussianKernel[i+2][j+2]
					pixel := float64(gray.GrayAt(x+j, y+i).Y)
					sum += weight * pixel
				}
			}
			blurred.SetGray(x, y, color.Gray{Y: uint8(sum / 159)}) // 标准化
		}
	}
	return blurred
}

// 梯度计算
func sobelGradient(gray *image.Gray) (*image.Gray, [][]float64) {
	bounds := gray.Bounds()
	gradient := image.NewGray(bounds)
	directions := make([][]float64, bounds.Max.Y)
	for i := range directions {
		directions[i] = make([]float64, bounds.Max.X)
	}

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
			magnitude := math.Sqrt(float64(gx*gx + gy*gy))
			if magnitude > 255 {
				magnitude = 255
			}
			gradient.SetGray(x, y, color.Gray{Y: uint8(magnitude)})

			// 方向 (atan2 返回值范围为 [-π, π])
			directions[y][x] = math.Atan2(float64(gy), float64(gx))
		}
	}
	return gradient, directions
}

// 非极大值抑制
func nonMaximumSuppression(gradient *image.Gray, directions [][]float64) *image.Gray {
	bounds := gradient.Bounds()
	nms := image.NewGray(bounds)

	for y := 1; y < bounds.Max.Y-1; y++ {
		for x := 1; x < bounds.Max.X-1; x++ {
			magnitude := gradient.GrayAt(x, y).Y
			angle := directions[y][x] * (180 / math.Pi) // 转为角度
			angle = math.Mod(angle+180, 180)            // 保证在 [0, 180) 范围

			// 确定梯度方向
			var neighbor1, neighbor2 uint8
			if (angle >= 0 && angle < 22.5) || (angle >= 157.5 && angle < 180) {
				neighbor1 = gradient.GrayAt(x-1, y).Y
				neighbor2 = gradient.GrayAt(x+1, y).Y
			} else if angle >= 22.5 && angle < 67.5 {
				neighbor1 = gradient.GrayAt(x-1, y-1).Y
				neighbor2 = gradient.GrayAt(x+1, y+1).Y
			} else if angle >= 67.5 && angle < 112.5 {
				neighbor1 = gradient.GrayAt(x, y-1).Y
				neighbor2 = gradient.GrayAt(x, y+1).Y
			} else if angle >= 112.5 && angle < 157.5 {
				neighbor1 = gradient.GrayAt(x+1, y-1).Y
				neighbor2 = gradient.GrayAt(x-1, y+1).Y
			}

			// 抑制非局部最大值
			if magnitude >= neighbor1 && magnitude >= neighbor2 {
				nms.SetGray(x, y, color.Gray{Y: magnitude})
			} else {
				nms.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return nms
}

// 双阈值检测
func doubleThreshold(nms *image.Gray, low, high uint8) *image.Gray {
	bounds := nms.Bounds()
	output := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			value := nms.GrayAt(x, y).Y
			if value >= high {
				output.SetGray(x, y, color.Gray{Y: 255}) // 强边缘
			} else if value >= low {
				output.SetGray(x, y, color.Gray{Y: 128}) // 弱边缘
			} else {
				output.SetGray(x, y, color.Gray{Y: 0}) // 非边缘
			}
		}
	}
	return output
}

func (op *cannyOperator) Apply(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, img, bounds.Min, draw.Src)

	// 高斯平滑
	blurred := gaussianBlur(gray)

	// 梯度计算
	gradient, directions := sobelGradient(blurred)

	// 非极大值抑制
	nms := nonMaximumSuppression(gradient, directions)

	// 双阈值检测
	final := doubleThreshold(nms, 50, 150)

	return final
}
