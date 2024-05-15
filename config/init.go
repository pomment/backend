package config

import (
	"github.com/hoisie/mustache"
	"github.com/pelletier/go-toml/v2"
	"github.com/pomment/pomment/common"
	"github.com/pomment/pomment/dao"
	"github.com/pomment/pomment/utils/recaptcha"
	"io"
	"os"
	"path/filepath"
	"time"
)

var Content common.PommentConfig

func InitConfig(basePath string) (err error) {
	// 打开文件
	absPath := filepath.Join(basePath, "config.toml")
	jsonFile, err := os.Open(absPath)
	defer jsonFile.Close()
	if err != nil {
		return err
	}

	// 读取数据
	data, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	// 解析 TOML
	err = toml.Unmarshal(data, &Content)
	if err != nil {
		return err
	}

	// 处理默认值
	if Content.Push.GravatarServer == "" {
		Content.Push.GravatarServer = "https://secure.gravatar.com/avatar/"
	}

	// 邮件模板读取
	{
		collection := Content.WebTemplate
		bodyRaw, err := dao.ReadTextFile("template/email.html")
		if err != nil {
			return err
		}

		unsubscribeConfirmRaw, err := dao.ReadTextFile("template/unsubscribe.html")
		if err != nil {
			return err
		}

		unsubscribeSuccessRaw, err := dao.ReadTextFile("template/unsubscribe_success.html")
		if err != nil {
			return err
		}

		unsubscribeErrorRaw, err := dao.ReadTextFile("template/unsubscribe_error.html")
		if err != nil {
			return err
		}

		title, err := mustache.ParseString(Content.Email.Title)
		body, err := mustache.ParseString(bodyRaw)
		unsubscribeConfirm, err := mustache.ParseString(unsubscribeConfirmRaw)
		unsubscribeSuccess, err := mustache.ParseString(unsubscribeSuccessRaw)
		unsubscribeError, err := mustache.ParseString(unsubscribeErrorRaw)
		if err != nil {
			return err
		}

		collection.EmailTitle = *title
		collection.EmailBody = *body
		collection.UnsubscribeConfirm = *unsubscribeConfirm
		collection.UnsubscribeSuccess = *unsubscribeSuccess
		collection.UnsubscribeError = *unsubscribeError
		Content.WebTemplate = collection
	}

	// ReCAPTCHA 初始化
	if Content.ReCAPTCHA.Enabled {
		reCAPTCHA, err := recaptcha.NewReCAPTCHA(Content.ReCAPTCHA.SecretKey, recaptcha.V3, time.Second*10, "https://recaptcha.net/recaptcha/api/siteverify")
		if err != nil {
			return err
		}
		Content.ReCAPTCHA.Object = reCAPTCHA
	}
	return err
}
