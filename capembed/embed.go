package capembed

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/fs"
	"math/rand"
	"strings"
)

const (
	sliderBkImagePath = "resources/images/slider/bkimg"
	sliderBlockPath   = "resources/images/slider/block"
	fontsPath         = "resources/fonts"
)

//go:embed resources
var fileResources embed.FS

var DefaultFileResource = NewFileResources(fileResources)

type FileResources struct {
	fileResources embed.FS
}

type SliderResource struct {
	BKImg    *ImgResource
	SliImg   *ImgResource
	NoiseImg *ImgResource
	FontRes  *truetype.Font
}

type ImgResource struct {
	FilePath  string
	ImageFile image.Image
	RgbaImage *image.RGBA
	FontPath  string
	Width     int
	Height    int
	idx       int
}

func (i *ImgResource) IsOpacity(x, y int) bool {
	A := i.RgbaImage.RGBAAt(x, y).A
	if float32(A) <= 125 {
		return true
	}
	return false
}

func (i *ImgResource) SetPixel(rgba color.RGBA, x int, y int) {
	i.RgbaImage.SetRGBA(x, y, rgba)
}

func (i *ImgResource) VagueImage(x int, y int) {
	// VagueImage 模糊区域
	var red uint32
	var green uint32
	var blue uint32
	var alpha uint32

	points := [8][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

	for _, point := range points {
		pointX := x + point[0]
		pointY := y + point[1]

		if pointX < 0 || pointX >= i.Width || pointY < 0 || pointY >= i.Height {
			continue
		}

		r, g, b, a := i.RgbaImage.RGBAAt(pointX, pointY).RGBA()
		red += r >> 9
		green += g >> 9
		blue += b >> 9
		alpha += a >> 9

	}

	var avg uint32
	avg = 8

	rgba := color.RGBA{R: uint8(red / avg), G: uint8(green / avg), B: uint8(blue / avg), A: uint8(alpha / avg)}

	i.RgbaImage.SetRGBA(x, y, rgba)
}

// Base64 为像素设置颜色
func (i *ImgResource) Base64() (string, error) {
	// 开辟一个新的空buff
	var buf bytes.Buffer
	// img写入到buff
	if err := png.Encode(&buf, i.RgbaImage); err != nil {
		return "", err
	}
	//开辟存储空间
	dist := make([]byte, buf.Cap()+buf.Len())
	// buff转成base64
	base64.StdEncoding.Encode(dist, buf.Bytes())
	return strings.Trim(string(dist), "\u0000"), nil
}

type ClickResource struct {
	BKImg *ImgResource
}

func NewFileResources(fileResources embed.FS) *FileResources {
	return &FileResources{fileResources: fileResources}
}

func (fr *FileResources) RandClickResource() (resources *ClickResource, err error) {
	sliderRes, err := fr.RandImg(sliderBkImagePath, -1)
	if err != nil {
		return nil, err
	}
	return &ClickResource{BKImg: sliderRes}, nil
}

// RandSliderResource 随机返回滑块验证码背景图片和滑块
func (fr *FileResources) RandSliderResource() (resource *SliderResource, err error) {
	sliderRes, err := fr.RandImg(sliderBkImagePath, -1)
	if err != nil {
		return nil, err
	}
	blockRes, err := fr.RandImg(sliderBlockPath, -1)
	if err != nil {
		return nil, err
	}
	noiseRes, err := fr.RandImg(sliderBlockPath, blockRes.idx)
	if err != nil {
		return nil, err
	}
	fontRes, err := fr.randFont()
	if err != nil {
		return nil, err
	}
	return &SliderResource{
		BKImg:    sliderRes,
		SliImg:   blockRes,
		NoiseImg: noiseRes,
		FontRes:  fontRes,
	}, nil
}

func (fr *FileResources) RandImg(path string, excludeIdx int) (*ImgResource, error) {
	entries, err := fr.fileResources.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("path: %s, must not be empty", path)
	}
	idx := rand.Intn(len(entries))
	for excludeIdx == idx {
		idx = rand.Intn(len(entries))
	}
	fileInfo, err := entries[idx].Info()
	if err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s/%s", path, fileInfo.Name())
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("path: path must all image resource, but %s is dir", filePath)
	}

	ff, err := fr.fileResources.Open(filePath)
	if err != nil {
		return nil, err
	}
	imageFile, err := png.Decode(ff)
	if err != nil {
		return nil, err
	}
	return &ImgResource{
		FilePath:  filePath,
		ImageFile: imageFile,
		RgbaImage: fr.imageToRGBA(imageFile),
		Width:     imageFile.Bounds().Dx(),
		Height:    imageFile.Bounds().Dy(),
		idx:       idx,
	}, nil
}

// ImageToRGBA 图片转rgba
func (fr *FileResources) imageToRGBA(img image.Image) *image.RGBA {
	// No conversion needed if image is an *image.RGBA.
	if dst, ok := img.(*image.RGBA); ok {
		return dst
	}
	// Use the image/draw package to convert to *image.RGBA.
	b := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), img, b.Min, draw.Src)
	return dst
}

func (fr *FileResources) randFont() (*truetype.Font, error) {
	entries, err := fr.fileResources.ReadDir(fontsPath)
	if err != nil {
		return nil, err
	}
	var fontPaths []fs.DirEntry
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".ttf") && !entry.IsDir() {
			fontPaths = append(fontPaths, entry)
		}
	}
	if len(fontPaths) == 0 {
		return nil, fmt.Errorf("font path: %s, not found any font file", fontsPath)
	}
	idx := rand.Intn(len(fontPaths))
	fontSourceBytes, err := fr.fileResources.ReadFile(fmt.Sprintf("%s/%s", fontsPath, fontPaths[idx].Name()))
	if err != nil {
		return nil, err
	}
	return freetype.ParseFont(fontSourceBytes)
}
