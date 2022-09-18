package utils

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
)

func ImageCompress(
	getReadSizeFile func() (io.Reader, error),
	getDecodeFile func() (*os.File, error),
	to string,
	Quality,
	base int,
	format string) error {
	/** 读取文件 */
	file_origin, err := getDecodeFile()
	defer file_origin.Close()
	if err != nil {
		return err
	}
	var origin image.Image
	var config image.Config
	var temp io.Reader
	/** 读取尺寸 */
	temp, err = getReadSizeFile()
	if err != nil {
		return err
	}
	var typeImage int64
	format = strings.ToLower(format)

	allowFormats := []string{".jpg", ".jpeg", ".png"}

	if !InSliceString(format, allowFormats) {
		return nil
	}

	/** jpg 格式 */
	if format == ".jpg" || format == ".jpeg" {
		typeImage = 1

		origin, err = jpeg.Decode(file_origin)
		if err != nil {
			return err
		}
		temp, err = getReadSizeFile()
		if err != nil {
			return err
		}

		config, err = jpeg.DecodeConfig(temp)
		if err != nil {
			return err
		}
	} else if format == ".png" {
		typeImage = 0
		origin, err = png.Decode(file_origin)
		if err != nil {
			return err
		}
		temp, err = getReadSizeFile()
		if err != nil {
			return err
		}
		config, err = png.DecodeConfig(temp)
		if err != nil {
			return err
		}
	}

	/** 做等比缩放 */
	width := uint(base) /** 基准 */
	height := uint(base * config.Height / config.Width)

	canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)
	file_out, err := os.Create(to)
	defer file_out.Close()
	if err != nil {
		return err
	}
	if typeImage == 0 {
		err = png.Encode(file_out, canvas)
		if err != nil {
			return err
		}
	} else {
		err = jpeg.Encode(file_out, canvas, &jpeg.Options{Quality})
		if err != nil {
			return err
		}
	}

	return nil
}
