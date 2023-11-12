package captchago

import (
	"fmt"
	"github.com/xiao-ren-wu/captchago/captcha_embed"
	"github.com/xiao-ren-wu/captchago/util"
	"golang.org/x/image/colornames"
	"image/color"
	"strconv"
)

type SliderCaptchaData struct {
	// 滑块验证码背景图片
	BackgroundImageBase64 string `json:"background_image_base64"`
	// 滑块
	SliderImageBase64 string `json:"slider_image_base64"`
	// 滑块所在背景位置百分比
	ResultPercent int `json:"-"`
}

type SliderCaptcha struct {
	// 图片资源库
	fr  *captcha_embed.FileResources
	wmc *waterMarkCfg
}

func NewSliderCaptcha() *SliderCaptcha {
	return &SliderCaptcha{
		fr: captcha_embed.DefaultFileResource,
		wmc: &waterMarkCfg{
			fontSize: 20,
			color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
	}
}

func (sc *SliderCaptcha) Get() (*SliderCaptchaData, error) {
	sliderResource, err := sc.fr.RandSliderResource()
	if err != nil {
		return nil, err
	}
	// 写水印
	if len(sc.wmc.text) > 0 {
		wm := waterMark{
			imgResource:  sliderResource.BKImg,
			fontRes:      sliderResource.FontRes,
			waterMarkCfg: sc.wmc,
		}
		if err := wm.WriteWaterMark(); err != nil {
			return nil, err
		}
	}
	x := sc.pictureTemplatesCut(sliderResource)
	bkImgBase64, err := sliderResource.BKImg.Base64()
	if err != nil {
		return nil, err
	}
	sliBase64, err := sliderResource.SliImg.Base64()
	if err != nil {
		return nil, err
	}
	percent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(x)/float64(sliderResource.BKImg.Width)), 64)
	return &SliderCaptchaData{
		BackgroundImageBase64: bkImgBase64,
		SliderImageBase64:     sliBase64,
		ResultPercent:         int(percent * 100),
	}, nil
}

func (sc *SliderCaptcha) pictureTemplatesCut(sr *captcha_embed.SliderResource) int {
	// 生成拼图坐标点
	x := sc.randPoint(sr.BKImg, sr.SliImg)
	// 裁剪模板图
	sc.cutByTemplate(sr.BKImg, sr.SliImg, x)

	// 插入干扰图
	var offsetX int
	for {
		offsetX = util.RandomInt(sr.NoiseImg.Width, sr.BKImg.Width-sr.NoiseImg.Width)
		if offsetX > x && x+sr.SliImg.Width < offsetX {
			break
		} else if sr.NoiseImg.Width+offsetX < x {
			break
		}
	}
	sc.cutByTemplate(sr.BKImg, sr.NoiseImg, offsetX)
	return x
}

func (sc *SliderCaptcha) cutByTemplate(backgroundImage, templateImage *captcha_embed.ImgResource, x1 int) {
	xLength := templateImage.Width
	yLength := templateImage.Height
	for x := 0; x < xLength; x++ {
		for y := 0; y < yLength; y++ {
			// 如果模板图像当前像素点不是透明色 copy源文件信息到目标图片中
			isOpacity := templateImage.IsOpacity(x, y)

			// 当前模板像素在背景图中的位置
			backgroundX := x + x1
			backgroundY := y

			// 当不为透明时
			if !isOpacity {
				// 获取原图像素
				backgroundRgba := backgroundImage.RgbaImage.RGBAAt(backgroundX, backgroundY)
				// 将原图的像素扣到模板图上
				templateImage.SetPixel(backgroundRgba, x, y)
				// 背景图区域模糊
				backgroundImage.VagueImage(backgroundX, backgroundY)
			}

			//防止数组越界判断
			if x == (xLength-1) || y == (yLength-1) {
				continue
			}

			rightOpacity := templateImage.IsOpacity(x+1, y)
			downOpacity := templateImage.IsOpacity(x, y+1)

			//描边处理，,取带像素和无像素的界点，判断该点是不是临界轮廓点,如果是设置该坐标像素是白色
			if (isOpacity && !rightOpacity) || (!isOpacity && rightOpacity) || (isOpacity && !downOpacity) || (!isOpacity && downOpacity) {
				templateImage.RgbaImage.SetRGBA(x, y, colornames.White)
				backgroundImage.RgbaImage.SetRGBA(backgroundX, backgroundY, colornames.White)
			}
		}
	}
}

// 生成模板图在背景图中的随机坐标点
func (sc *SliderCaptcha) randPoint(bkImg, sliImg *captcha_embed.ImgResource) (x int) {
	x = util.RandomInt(sliImg.Width+10, bkImg.Width-sliImg.Width)
	return x
}
