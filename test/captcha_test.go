package test

import (
	"testing"

	"github.com/weitienwong/aliyun-captcha/client"
	"github.com/weitienwong/aliyun-captcha/config"
)

const (
	AccessKeyID     = "LTAI5t7c9pCXTnoruV2cSk9S"
	AccessKeySecret = "LlXw55fqPpL14ZNTBExw9wLwrYddPY"
	SceneId         = "177jxdo1"
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
