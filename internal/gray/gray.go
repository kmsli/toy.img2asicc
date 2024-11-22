package gray

import "image"

type (
	// 算子接口
	Operator interface {
		Apply(image.Image) *image.Gray
	}
)

var (
	// 索贝尔算子
	Sobel = &sobelOperator{}
	// 自定义卷积核
	Canny = &cannyOperator{}
)
