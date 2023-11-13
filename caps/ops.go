package caps

type CaptchaCfg struct {
	// 水印大小
	FontSize *int
	// 水印内容
	Text *string
	// 文件位置微调，由于文字宽度，需要根据具体水印进行微调
	Offset *int
}

type CaptchaOp func(ops *CaptchaCfg)

func WaterMark(text string) CaptchaOp {
	return func(ops *CaptchaCfg) {
		ops.Text = &text
	}
}

func WaterMarkSize(fontSize int) CaptchaOp {
	return func(ops *CaptchaCfg) {
		ops.FontSize = &fontSize
	}
}

func WaterMarkOffset(offset int) CaptchaOp {
	return func(ops *CaptchaCfg) {
		ops.Offset = &offset
	}
}
