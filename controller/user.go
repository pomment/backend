package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"pomment-go/common"
	"pomment-go/dao"
	"pomment-go/model"
	"pomment-go/utils"
)

// ThreadMeta 获取评论串元数据
func ThreadMeta(c *gin.Context) {
	res, err := model.GetThreadMeta(c.Param("id"))
	if err != nil {
		utils.AjaxSuccess(c, common.Thread{
			ID: c.Param("id"),
		})
		return
	}
	utils.AjaxSuccess(c, res)
}

type ThreadMetaByUrlRequest struct {
	URL string `json:"url"`
}

// ThreadMetaByURL 获取评论串元数据（通过隶属博文地址）
func ThreadMetaByURL(c *gin.Context) {
	req := ThreadMetaByUrlRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}

	res, err := model.GetThreadMetaByURL(req.URL)
	if err != nil {
		utils.AjaxSuccess(c, common.Thread{
			URL: req.URL,
		})
		return
	}
	utils.AjaxSuccess(c, res)
}

// ThreadMetaByURLs 获取评论串元数据（通过多个隶属博文地址）
func ThreadMetaByURLs(c *gin.Context) {
	var req []string
	err := c.BindJSON(&req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}

	var resMap = make(map[string]common.Thread)

	for _, e := range req {
		res, err := model.GetThreadMetaByURL(e)
		if err != nil {
			resMap[e] = common.Thread{
				URL: e,
			}
			continue
		}
		resMap[e] = *res
	}

	utils.AjaxSuccess(c, resMap)
}

// PostsByID 获取评论列表
func PostsByID(c *gin.Context) {
	cacheKey := common.CachePostIDKeyPrefix + c.Param("id")
	cache, err := dao.GetCache(cacheKey)
	if err != nil {
		utils.AjaxError(c, 500, err)
		return
	}
	//fmt.Printf("Get cache key %s, result: %s\n", cacheKey, cache)
	if cache != "" {
		utils.AjaxJSONString(c, cache)
		return
	}

	resPost, err := model.GetPostsSimpleByID(c.Param("id"))
	resMeta, err := model.GetThreadMeta(c.Param("id"))
	if err != nil {
		utils.AjaxSuccess(c, gin.H{
			"meta": common.Thread{
				ID: c.Param("id"),
			},
			"post": make([]string, 0),
		})
		return
	}

	res, err := utils.AjaxSuccessWithReturns(c, gin.H{
		"meta": resMeta,
		"post": resPost,
	})
	if err != nil {
		fmt.Printf("Error occured while processing cache content generation: %s\n", err)
		return
	}
	err = dao.SetCache(cacheKey, string(res))
	if err != nil {
		fmt.Printf("Error occured while processing cache saving: %s\n", err)
		return
	}
}

// PostsByURL 获取评论列表（通过隶属博文地址）
func PostsByURL(c *gin.Context) {
	req := ThreadMetaByUrlRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}

	cacheKey := common.CachePostURLKeyPrefix + req.URL
	cache, err := dao.GetCache(cacheKey)
	if err != nil {
		utils.AjaxError(c, 500, err)
		return
	}
	//fmt.Printf("Get cache key %s, result: %s\n", cacheKey, cache)
	if cache != "" {
		utils.AjaxJSONString(c, cache)
		return
	}

	resPost, err := model.GetPostsSimpleByURL(req.URL)
	resMeta, err := model.GetThreadMetaByURL(req.URL)
	if err != nil {
		utils.AjaxSuccess(c, gin.H{
			"meta": common.Thread{
				ID: c.Param("id"),
			},
			"post": make([]string, 0),
		})
		return
	}

	res, err := utils.AjaxSuccessWithReturns(c, gin.H{
		"meta": resMeta,
		"post": resPost,
	})
	if err != nil {
		fmt.Printf("Error occured while processing cache content generation: %s\n", err)
		return
	}
	err = dao.SetCache(cacheKey, string(res))
	if err != nil {
		fmt.Printf("Error occured while processing cache saving: %s\n", err)
		return
	}
}

// PostsAdd 新增评论
func PostsAdd(c *gin.Context) {
	// 数据解析
	req := model.AppendPostUserArgs{}
	err := c.BindJSON(&req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}

	// 数据校验
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}

	// 数据操作
	res, err := model.AppendPostUser(req)
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}

	utils.AjaxSuccess(c, res)
}
