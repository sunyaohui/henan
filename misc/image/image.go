package image

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"maoguo/henan/misc/utils/parse"
	"os"

	"github.com/nfnt/resize"
	"github.com/wonderivan/logger"
)

/*
* 缩略图生成
* 入参:
* 规则: 当宽或者高为空时，等比例缩放
* 矩形坐标系起点是左上
* 返回:error
 */
func Scale(src, dst string, width, height int) error {
	in, _ := os.Open(src)
	defer in.Close()
	out, _ := os.Create(dst)
	defer out.Close()
	origin, fm, err := image.Decode(in)
	if err != nil {
		logger.Error("images Scale failed", err)
		return err
	}
	w := origin.Bounds().Max.X
	h := origin.Bounds().Max.Y

	if width == 0 && height == 0 {
		height = origin.Bounds().Max.Y
		width = origin.Bounds().Max.X
	}
	if width == 0 && height != 0 {
		in_h := parse.IntToFloat64(h)
		bili := in_h / (parse.IntToFloat64(height))
		out_w := (parse.IntToFloat64(w)) / bili
		width = parse.Float64ToInt(out_w)
	}
	if height == 0 && width != 0 {
		in_w := parse.IntToFloat64(w)
		bili := in_w / (parse.IntToFloat64(width))
		int_h := parse.IntToFloat64(h)
		out_h := int_h / bili
		height = parse.Float64ToInt(out_h)
	}

	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)
	switch fm {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{100})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
	case "jpg":
		return jpeg.Encode(out, canvas, &jpeg.Options{100})
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}
