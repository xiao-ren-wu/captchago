package captchago

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/xiao-ren-wu/captchago/capembed"
	"image"
	"image/color"
	"unicode"
)

type waterMarkCfg struct {
	// 水印大小
	fontSize int
	// 水印颜色
	color color.RGBA
	// 水印内容
	text string
	// 文件位置微调，由于文字宽度，需要根据具体水印进行微调
	offset int
}

type waterMark struct {
	imgResource *capembed.ImgResource
	fontRes     *truetype.Font
	*waterMarkCfg
}

func (c *waterMark) WriteWaterMark() error {
	x := float64(c.imgResource.Width) - float64(c.waterMarkLen())
	y := float64(c.imgResource.Height) - (25 / 2) + 7
	fc := freetype.NewContext()
	// 设置屏幕每英寸的分辨率
	//fc.SetDPI(72)
	// 设置用于绘制文本的字体
	fc.SetFont(c.fontRes)
	// 以磅为单位设置字体大小
	fc.SetFontSize(float64(c.fontSize))
	// 设置剪裁矩形以进行绘制
	fc.SetClip(c.imgResource.RgbaImage.Bounds())
	// 设置目标图像
	fc.SetDst(c.imgResource.RgbaImage)
	// 设置绘制操作的源图像，通常为 image.Uniform
	fc.SetSrc(image.NewUniform(c.color))
	// 设置水印地址
	pt := freetype.Pt(int(x), int(y))
	// 根据 Pt 的坐标值绘制给定的文本内容
	_, err := fc.DrawString(c.text, pt)
	return err
}

func (c *waterMark) waterMarkLen() int {
	enCount, zhCount := 0, 0
	for _, t := range c.text {
		if unicode.Is(unicode.Han, t) {
			zhCount++
		} else {
			enCount++
		}
	}
	chOffset := c.fontSize * zhCount
	enOffset := c.fontSize * enCount
	return chOffset + enOffset + c.offset
}
