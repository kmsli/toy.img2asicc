package img2ascii

import (
	"image"
	"image/color"
	"os"

	"github.com/nfnt/resize"
	"golang.org/x/image/webp"
	"toy.img2ascii/internal/gray"
)

/*
ASCII（American Standard Code for Information Interchange，美国信息交换标准代码）。
它是基于拉丁字母的一套计算机编码系统，主要用于表示现代英语和其他西欧语言中的字符，包括控制字符（如回车、换行等）和可打印字符（如字母、数字、标点符号等）。
ASCII 码用特定的二进制数字来表示每个字符，在计算机数据存储和通信中被广泛应用。
*/

func Webp2Ascii(imgPath string, op gray.Operator) {
	// 读取图片
	imgfile, err := os.OpenFile(imgPath, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}

	// 解码图片
	img, err := webp.Decode(imgfile)
	if err != nil {
		panic(err)
	}

	// 获取终端尺寸
	termWidth, termHeight := 80, 24 // 假设终端尺寸为 80x24，可以根据实际情况调整

	// 根据终端尺寸动态调整图片尺寸
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()
	scale := float64(termWidth) / float64(imgWidth)
	if imgHeight*int(scale) > termHeight {
		scale = float64(termHeight) / float64(imgHeight)
	}
	newWidth := uint(float64(imgWidth) * scale)
	newHeight := uint(float64(imgHeight) * scale)

	// 调整图片尺寸
	img = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// 转换为灰度图
	bounds := img.Bounds()
	grayImg := op.Apply(img)

	// 生成ASCII字符画
	asciiImg := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			asciiImg.Set(x, y, color.Gray{grayImg.GrayAt(x, y).Y})
		}
	}

	// 输出到终端
	chars := "@#*+=-:. "
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := asciiImg.At(x, y).(color.RGBA).R
			// 终端彩色输出
			// fmt.Printf("\x1b[48;2;%d;%d;%dm \x1b[0m ", gray, gray, gray)
			index := int(gray) * (len(chars) - 1) / 255
			// 终端彩色字符输出
			// fmt.Printf("\x1b[38;2;%d;%d;%dm%c\x1b[0m", colorful.Y, colorful.Y, colorful.Y, chars[index])
			// 终端灰度输出
			print(string(chars[index]) + " ")
		}
		println()
	}
}
