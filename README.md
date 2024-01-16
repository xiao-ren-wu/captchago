package captchago_test

import (
	"encoding/base64"
	"github.com/xiao-ren-wu/captchago"
	"github.com/xiao-ren-wu/captchago/caps"
	"os"
	"testing"
)

func TestGenSliderCaptcha(t *testing.T) {
	captcha := captchago.NewSliderCaptcha(caps.WaterMark("hello"), caps.WaterMarkOffset(-50))
	captchaData, err := captcha.Get()
	if err != nil {
		t.Fatal(err)
	}
	decodeString, err := base64.StdEncoding.DecodeString(captchaData.BackgroundImageBase64)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile("bk.png", decodeString, os.ModePerm)
	bytes, err := base64.StdEncoding.DecodeString(captchaData.SliderImageBase64)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile("sli.png", bytes, os.ModePerm)
	t.Logf("Percent: %d", captchaData.ResultPercent)
}
