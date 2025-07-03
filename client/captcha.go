package client

import (
	"errors"
	captcha "github.com/alibabacloud-go/captcha-20230305/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/weitienwong/aliyun-captcha/config"
	"log/slog"
	"os"
)

type Client struct {
	cfg    *config.CaptchaConfig
	client *captcha.Client
}

func NewClient(cfg *config.CaptchaConfig) (*Client, error) {
	if cfg.AccessKeyId == "" {
		return nil, errors.New("accessKeyId is required")
	}
	if cfg.AccessKeySecret == "" {
		return nil, errors.New("accessKeySecret is required")
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = os.Getenv("ALIYUN_CAPTCHA_ENDPOINT")
		if cfg.Endpoint == "" {
			cfg.Endpoint = "captcha.cn-shanghai.aliyuncs.com"
		}
	}

	credential, err := credentials.NewCredential(&credentials.Config{
		Type:            tea.String("access_key"),
		AccessKeyId:     tea.String(cfg.AccessKeyId),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
		ConnectTimeout:  tea.Int(5000),
		Timeout:         tea.Int(5000),
	})

	if err != nil {
		slog.Error("Create credential error", "err", err)
		return nil, err
	}

	captchaClient, err := captcha.NewClient(&openapi.Config{
		Credential: credential,
		Endpoint:   tea.String(cfg.Endpoint),
	})

	if err != nil {
		slog.Error("Create captcha client error", "err", err)
		return nil, err
	}

	return &Client{
		cfg:    cfg,
		client: captchaClient,
	}, nil
}

func (c *Client) Verify(param string) bool {
	request := &captcha.VerifyIntelligentCaptchaRequest{
		CaptchaVerifyParam: tea.String(param),
		SceneId:            tea.String(c.cfg.SceneId),
	}

	var captchaVerifyResult *bool
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		resp, _err := c.client.VerifyIntelligentCaptcha(request)
		if _err != nil {
			return _err
		}
		slog.Info("Aliyun captcha verification success", "response", resp)
		captchaVerifyResult = resp.Body.Result.VerifyResult

		return nil
	}()
	if tryErr != nil {
		var err = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			err = _t
		} else {
			err.Message = tea.String(tryErr.Error())
		}
		slog.Error("Aliyun captcha verification error", "error", err, "request", request)
		// 出现异常认为验证通过，优先保证业务可用
		captchaVerifyResult = tea.Bool(true)

	}

	return *captchaVerifyResult
}
