package captcha

import (
	"context"
	"fmt"
	"github.com/Jackzode/painting/commons/constants"
	"github.com/Jackzode/painting/commons/handler"
	glog "github.com/Jackzode/painting/commons/logger"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"strings"
	"time"
)

func GenerateCaptchaAndSave(ctx context.Context) (key, captchaBase64 string, err error) {
	var driverString = base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor:         &color.RGBA{R: 211, G: 211, B: 211, A: 0},
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	driver := driverString.ConvertFonts()

	id, content, answer := driver.GenerateIdQuestionAnswer()
	item, err := driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}
	err = handler.RedisClient.SetEx(ctx, id, answer, constants.CaptchaExpiration).Err()
	if err != nil {
		glog.Slog.Error(err.Error())
	}
	captchaBase64 = item.EncodeB64string()
	return id, captchaBase64, nil
}

func GetContentByCaptchaCode(ctx context.Context, key string) (captcha string, err error) {
	captcha = handler.RedisClient.Get(ctx, key).Val()
	if captcha == "" {
		return "", fmt.Errorf("captcha not exist")
	}
	return captcha, nil
}

func DelCaptcha(ctx context.Context, key string) (err error) {
	err = handler.RedisClient.Del(ctx, key).Err()
	return
}

func SetCode(ctx context.Context, code, content string, duration time.Duration) {
	err := handler.RedisClient.Set(ctx, code, content, duration).Err()
	if err != nil {
		glog.Slog.Error(err.Error())
	}
	return
}

func VerifyCaptcha(ctx context.Context, key, captcha string) (isCorrect bool, err error) {
	realCaptcha, err := GetContentByCaptchaCode(ctx, key)
	if err != nil {
		glog.Slog.Error("VerifyCaptcha GetContentByCaptchaCode Error", err.Error())
		return false, nil
	}
	err = DelCaptcha(ctx, key)
	if err != nil {
		glog.Slog.Error("VerifyCaptcha DelCaptcha Error", err.Error())
		return false, nil
	}
	return strings.TrimSpace(captcha) == realCaptcha, nil
}
