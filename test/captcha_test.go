package test

import (
	"testing"

	"github.com/weitienwong/aliyun-captcha/client"
	"github.com/weitienwong/aliyun-captcha/config"
)

const (
	AccessKeyID     = "xxx"
	AccessKeySecret = "xxx"
	SceneId         = "xxx"
	Param           = `{"sceneId":"177jxdo1","certifyId":"xxxxxx","deviceToken":"xxxxxxx==","data":"xxxxxx==","..."}`
)

func TestVerify(t *testing.T) {
	cfg := &config.CaptchaConfig{
		AccessKeyId:     AccessKeyID,
		AccessKeySecret: AccessKeySecret,
		SceneId:         SceneId,
	}
	c, err := client.NewClient(cfg)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(c.Verify(Param))
}
