package test

import (
	"github.com/weitienwong/aliyun-captcha"
	"testing"
)

const (
	AccessKeyID     = "xxx"
	AccessKeySecret = "xxx"
	SceneId         = "xxx"
	Param           = `{"sceneId":"177jxdo1","certifyId":"xxxxxx","deviceToken":"xxxxxxx==","data":"xxxxxx==","..."}`
)

func TestVerify(t *testing.T) {
	cfg := &captcha.Config{
		AccessKeyId:     AccessKeyID,
		AccessKeySecret: AccessKeySecret,
		SceneId:         SceneId,
	}
	c := captcha.NewClient(cfg)
	t.Log(c.Verify(Param))
}
