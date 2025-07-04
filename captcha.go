package captcha

import (
	captcha "github.com/alibabacloud-go/captcha-20230305/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"log/slog"
	"os"
)

type Client struct {
	cfg    *Config
	client *captcha.Client
}

func NewClient(cfg *Config) *Client {
	if cfg.AccessKeyId == "" {
		slog.Error("accessKeyId is required")
		return nil
	}
	if cfg.AccessKeySecret == "" {
		slog.Error("accessKeySecret is required")
		return nil
	}
	if cfg.Endpoint == "" {
		_endpoint := os.Getenv("ALIYUN_CAPTCHA_ENDPOINT")
		slog.Info("Get endpoint from ENV", "ALIYUN_CAPTCHA_ENDPOINT", _endpoint)
		if _endpoint == "" {
			slog.Warn("ALIYUN_CAPTCHA_ENDPOINT env was not found, use default endpoint(captcha.cn-shanghai.aliyuncs.com)")
			_endpoint = "captcha.cn-shanghai.aliyuncs.com"
		}
		cfg.Endpoint = _endpoint
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
		return nil
	}

	captchaClient, err := captcha.NewClient(&openapi.Config{
		Credential: credential,
		Endpoint:   tea.String(cfg.Endpoint),
	})

	if err != nil {
		slog.Error("Create captcha client error", "err", err)
		return nil
	}
	return &Client{
		cfg:    cfg,
		client: captchaClient,
	}
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
