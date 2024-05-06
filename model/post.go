package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pomment/pomment/common"
	"github.com/pomment/pomment/config"
	"github.com/pomment/pomment/dao"
	"github.com/pomment/pomment/message"
	"github.com/pomment/pomment/utils"
	"github.com/pomment/pomment/utils/recaptcha"
	"path"
	"strings"
	"time"
)

func GetPostsByID(id string) (outItem *[]common.Post, err error) {
	res, err := dao.ReadJSON(dao.GetThreadPath(id))
	if err != nil {
		return nil, err
	}

	var item *[]common.Post

	err = json.Unmarshal([]byte(res), &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func GetPostsByURL(url string) (outItem *[]common.Post, err error) {
	relation, err := GetThreadRelationByURL(url)
	if err != nil {
		return nil, err
	}
	return GetPostsByID(relation.ID)
}

func GetPostByID(threadID string, postID string) (post *common.Post, err error) {
	posts, err := GetPostsByID(threadID)
	if err != nil {
		return nil, err
	}
	for _, e := range *posts {
		if e.ID == postID {
			return &e, nil
		}
	}
	return nil, errors.New("unable to find post")
}

func GetPostsSimpleByID(id string) (outItem *[]common.PostSimple, err error) {
	res, err := dao.ReadJSON(dao.GetThreadPath(id))
	if err != nil {
		return nil, err
	}

	var item *[]common.PostSimple
	var builtItem []common.PostSimple

	err = json.Unmarshal([]byte(res), &item)
	if err != nil {
		return nil, err
	}

	// 过滤隐藏的评论
	for _, e := range *item {
		if !e.Hidden {
			builtItem = append(builtItem, e)
		}
	}

	return &builtItem, nil
}

func GetPostsSimpleByURL(url string) (outItem *[]common.PostSimple, err error) {
	relation, err := GetThreadRelationByURL(url)
	if err != nil {
		return nil, err
	}
	return GetPostsSimpleByID(relation.ID)
}

func AppendPostAndSave(id string, item common.Post) (err error) {
	posts, err := GetPostsByID(id)
	if err != nil {
		return err
	}

	*posts = append(*posts, item)

	res, err := json.Marshal(*posts)
	if err != nil {
		return err
	}

	err = dao.WriteJSON(path.Join("threads", fmt.Sprintf("%s.json", id)), string(res))
	return err
}

func AppendPostBackgroundJob(post common.Post, thread common.Thread, challengeResponse string) (err error) {
	// 检查 reCAPTCHA
	if config.Content.ReCAPTCHA.Enabled {
		score, err := config.Content.ReCAPTCHA.Object.VerifyWithOptions(challengeResponse, recaptcha.VerifyOption{
			Threshold: float32(config.Content.ReCAPTCHA.MinimumScore),
		})
		if err != nil {
			fmt.Println(err)
		}

		post.Hidden = false
		post.Rating = float64(score)

		err = EditPost(thread.ID, post, false)
		if err != nil {
			fmt.Println(err)
		}
	}

	// 删除缓存
	err = dao.DeleteCacheForThread(&thread)

	// 手机推送
	tokenList, err := LoadFCMTokenList()
	if err != nil {
		return err
	}
	err = message.PushToClient(post, tokenList)
	if err != nil {
		fmt.Printf("Unable to push: %s\n", err)
	}

	// 发送邮件
	err = SendNotifyEmail(post, thread)
	return err
}

func SendNotifyEmail(post common.Post, thread common.Thread) (err error) {
	if !post.ReceiveEmail {
		fmt.Println("Parent post receiveEmail does not set, do not send notify.")
		return nil
	}
	fmt.Println("Sending email...")
	if post.Parent != "" {
		parentPost, err := GetPostByID(thread.ID, post.Parent)
		if err != nil {
			fmt.Printf("Finding parent post error, do not send notify: %s\n", err)
			return nil
		}
		err = message.SendEmail(parentPost.Email, post, *parentPost, thread)
	} else {
		fmt.Println("Parent post not set, do not send notify.")
	}
	return nil
}

type AppendPostUserArgs struct {
	ID                string `json:"id"`
	URL               string `json:"url"`
	Title             string `json:"title" validate:"required"`
	Parent            string `json:"parent"`
	Name              string `json:"name" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	Website           string `json:"website"`
	Content           string `json:"content" validate:"required"`
	ReceiveEmail      bool   `json:"receiveEmail"`
	ChallengeResponse string `json:"challengeResponse"`
}

func AppendPostUser(req AppendPostUserArgs) (post common.Post, err error) {
	var meta *common.Thread
	// 获取评论元数据
	if req.ID == "" {
		// 无 ID，使用 URL 查询 meta
		meta, err = GetThreadMetaForSubmit(req.URL, req.Title)
		if err != nil {
			return common.Post{}, err
		}
	} else {
		meta, err = GetThreadMeta(req.ID)
		if err != nil {
			return common.Post{}, err
		}
	}

	// 过滤非法 URL
	isHttp := strings.HasPrefix(req.Website, "http://")
	isHttps := strings.HasPrefix(req.Website, "https://")
	if !(isHttp || isHttps) {
		req.Website = ""
	}

	// 构建评论
	now := time.Now()
	sec := now.Unix() * 1000
	post = common.Post{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Email:        req.Email,
		EmailHashed:  utils.GetMailHash(req.Email),
		Website:      req.Website,
		Content:      req.Content,
		OrigContent:  req.Content,
		Parent:       req.Parent,
		EditKey:      utils.GetEditKey(),
		Hidden:       config.Content.ReCAPTCHA.Enabled,
		ReceiveEmail: req.ReceiveEmail,
		CreatedAt:    sec,
		UpdatedAt:    sec,
	}

	// 保存评论数据
	err = AppendPostAndSave(meta.ID, post)
	if err != nil {
		return common.Post{}, err
	}

	// 更新元数据
	meta, err = UpdateThreadMeta(meta.ID)
	if err != nil {
		return common.Post{}, err
	}

	// 删除缓存
	err = dao.DeleteCacheForThread(meta)
	if err != nil {
		return common.Post{}, err
	}

	go func() {
		err := AppendPostBackgroundJob(post, *meta, req.ChallengeResponse)
		if err != nil {
			fmt.Printf("Error when executing background tasks: %s\n", err)
		}
	}()

	return post, nil
}

type AppendPostAdminArgs struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Parent  string `json:"parent"`
	Content string `json:"content" validate:"required"`
}

func AppendPostAdmin(threadID string, req AppendPostAdminArgs) (post common.Post, err error) {
	// 获取评论元数据
	meta, err := GetThreadMeta(threadID)
	if err != nil {
		return common.Post{}, err
	}

	// 构建评论
	now := time.Now()
	sec := now.Unix() * 1000
	post = common.Post{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Email:        req.Email,
		EmailHashed:  utils.GetMailHash(req.Email),
		Website:      "",
		Content:      req.Content,
		OrigContent:  req.Content,
		Parent:       req.Parent,
		EditKey:      "",
		Hidden:       false,
		ByAdmin:      true,
		ReceiveEmail: false,
		CreatedAt:    sec,
		UpdatedAt:    sec,
	}

	// 保存评论数据
	err = AppendPostAndSave(meta.ID, post)
	if err != nil {
		return common.Post{}, err
	}

	// 更新元数据
	meta, err = UpdateThreadMeta(meta.ID)

	return post, nil
}

// EditPost 编辑评论
func EditPost(threadID string, post common.Post, alterEditTime bool) (err error) {
	posts, err := GetPostsByID(threadID)
	if err != nil {
		return err
	}

	// 当前时间
	now := time.Now()
	sec := now.Unix() * 1000

	// 重建数组
	var newPosts []common.Post
	for _, e := range *posts {
		if e.ID == post.ID {
			updatedAt := post.UpdatedAt
			if alterEditTime {
				updatedAt = sec
			}
			newPosts = append(newPosts, common.Post{
				ID:           post.ID,
				Name:         post.Name,
				Email:        post.Email,
				EmailHashed:  post.EmailHashed,
				Website:      post.Website,
				Parent:       post.Parent,
				Content:      post.Content,
				Hidden:       post.Hidden,
				ByAdmin:      post.ByAdmin,
				ReceiveEmail: post.ReceiveEmail,
				EditKey:      post.EditKey,

				CreatedAt:   post.CreatedAt,
				UpdatedAt:   updatedAt,
				OrigContent: post.OrigContent,

				Avatar: post.Avatar,
				Rating: post.Rating,
			})
		} else {
			newPosts = append(newPosts, e)
		}
	}

	// 保存
	res, err := json.Marshal(newPosts)
	if err != nil {
		return err
	}

	err = dao.WriteJSON(path.Join("threads", fmt.Sprintf("%s.json", threadID)), string(res))
	return err
}

func ConfirmUnsubscribe(postID string, threadID string, editKey string) (post *common.Post, thread *common.Thread, err error) {
	post, err = GetPostByID(threadID, postID)
	if err != nil {
		return nil, nil, err
	}

	thread, err = GetThreadMeta(threadID)
	if err != nil {
		return nil, nil, err
	}

	if post.EditKey != editKey {
		return nil, nil, errors.New("editKey is not correct")
	}
	return post, thread, nil
}

func PerformUnsubscribe(postID string, threadID string, editKey string, userConfirmed bool) (post *common.Post, thread *common.Thread, err error) {
	if !userConfirmed {
		return nil, nil, errors.New("user is not confirmed")
	}

	post, thread, err = ConfirmUnsubscribe(postID, threadID, editKey)
	post.ReceiveEmail = false
	err = EditPost(thread.ID, *post, false)
	if err != nil {
		return nil, nil, err
	}

	return post, thread, err
}
