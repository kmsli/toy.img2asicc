package gray

import (
	"fmt"
	"image/jpeg"
	"os"
	"testing"

	"github.com/nfnt/resize"
	"golang.org/x/image/webp"
)

func Test_toGray(t *testing.T) {
	type args struct {
		filepath string
		op       Operator
	}
	tests := []struct {
		name string
		args args
	}{
		{"sobel", args{filepath: "../test.webp", op: Sobel}},
		{"canny", args{filepath: "../test.webp", op: Canny}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.OpenFile(tt.args.filepath, os.O_RDONLY, 0666)
			if err != nil {
				t.Errorf("toGray() error = %v", err)
				return
			}
			img, err := webp.Decode(file)
			if err != nil {
				t.Errorf("toGray() error = %v", err)
				return
			}
			img = resize.Resize(uint(img.Bounds().Dx()/5), 0, img, resize.Lanczos3)
			gray := tt.args.op.Apply(img)
			if gray == nil {
				t.Errorf("toGray() got nil")
				return
			}

			// Save the gray image
			grayFile, err := os.Create(fmt.Sprintf("%s.jpg", tt.name))
			if err != nil {
				t.Errorf("toGray() error = %v", err)
				return
			}

			if err := jpeg.Encode(grayFile, gray, nil); err != nil {
				t.Errorf("toGray() error = %v", err)
			}

		})
	}
}
