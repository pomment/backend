package message

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pomment/pomment/common"
	"github.com/pomment/pomment/config"
	"net/http"
	//"github.com/pomment/pomment/model"
)

// PushToClient 将评论提醒推送到手机客户端
func PushToClient(post common.Post, tokenList []string) error {
	if !config.Content.Push.Enabled {
		return nil
	}

	prefix := config.Content.Push.Gateway

	// 遍历 token 列表，推送
	for _, e := range tokenList {
		// 发送请求到网关
		client := resty.New()
		data := map[string]string{
			"title":     post.Name,
			"body":      post.Content,
			"imageURL":  config.Content.Push.GravatarServer + post.EmailHashed,
			"userToken": e,
		}
		jsonStr, err := json.Marshal(data)
		if err != nil {
			return err
		}

		resp, err := client.
			R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Pomment-Site-Name", config.Content.Push.SiteName).
			SetHeader("Pomment-Site-Key", config.Content.Push.SiteKey).
			SetHeader("User-Agent", "PommentServer").
			SetBody(string(jsonStr)).
			Post(prefix + "/user/message")

		if err != nil {
			fmt.Printf("Warning: unable to send push request: %e\n", err)
			continue
		}

		if resp.RawResponse.StatusCode != http.StatusOK {
			fmt.Printf("Warning: Push gateway server returned %d\n", resp.RawResponse.StatusCode)
		}
	}
	return nil
}
