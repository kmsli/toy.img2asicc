package img2ascii

import (
	"testing"

	"toy.img2ascii/internal/gray"
)

func TestWebp2Ascii(t *testing.T) {
	type args struct {
		imgPath string
		op      gray.Operator
	}
	tests := []struct {
		name string
		args args
	}{
		{"canny", args{"test.webp", gray.Canny}},
		{"sobel", args{"test.webp", gray.Sobel}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Webp2Ascii(tt.args.imgPath, tt.args.op)
		})
	}
}
