package config

type CaptchaConfig struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Endpoint        string `json:"endpoint"`
	SceneId         string `json:"scene_id"`
}
