package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pomment/pomment/auth"
	"github.com/pomment/pomment/common"
	"github.com/pomment/pomment/config"
	"github.com/pomment/pomment/dao"
	"github.com/pomment/pomment/model"
	"github.com/pomment/pomment/utils"
	"net/http"
	"strings"
)

// Auth 鉴权
func Auth(c *gin.Context) {
	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	// 用户输入内容预检
	err := c.BindJSON(&req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}

	// 密码检查
	err = auth.CheckPassword(req.Username, req.Password)
	if err != nil {
		utils.AjaxError(c, http.StatusForbidden, err)
		return
	}

	// 生成 JWT 令牌
	token, err := auth.GenerateToken(req.Username)
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}
	utils.AjaxSuccess(c, token)
}

// ThreadList 评论串元数据列表
func ThreadList(c *gin.Context) {
	list, err := model.GetThreadList()
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}
	var str = "["
	for i, s := range list {
		meta, err := dao.ReadTextFile(s, true)
		if err != nil {
			utils.AjaxError(c, http.StatusInternalServerError, err)
			return
		}
		str += meta
		if i < len(list)-1 {
			str += ","
		}
	}
	str += "]"
	utils.AjaxWrapJSONString(c, str)
}

// ThreadID 评论串单个
func ThreadID(c *gin.Context) {
	list, err := model.GetPostsByID(c.Param("id"))
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}
	utils.AjaxSuccess(c, *list)
}

// PostID 评论单个
func PostID(c *gin.Context) {
	item, err := model.GetPostByID(c.Param("threadId"), c.Param("postId"))
	if err != nil {
		utils.AjaxError(c, http.StatusNotFound, err)
		return
	}
	utils.AjaxSuccess(c, *item)
}

// ThreadMetaID 评论串元数据单个
func ThreadMetaID(c *gin.Context) {
	res, err := model.GetThreadMeta(c.Param("id"))
	if err != nil {
		// utils.AjaxError(c, http.StatusInternalServerError)
		utils.AjaxSuccess(c, common.Thread{
			ID: c.Param("id"),
		})
		return
	}
	utils.AjaxSuccess(c, res)
}

// ThreadMetaIDEdit 评论串元数据单个编辑
func ThreadMetaIDEdit(c *gin.Context) {
	req := common.Thread{}
	err := c.BindJSON(&req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}

	err = model.SetThreadMeta(req)
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}

	utils.AjaxSuccess(c, req)
}

// PostsAddAdmin 新增评论
func PostsAddAdmin(c *gin.Context) {
	// 获取当前用户资料
	reqToken := c.GetHeader("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	username, err := auth.ValidateToken(reqToken)
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}

	user, err := auth.FindUserByName(username)
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}

	// 数据解析
	req := model.AppendPostAdminArgs{}
	err = c.BindJSON(&req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}
	req.Name = user.Name
	req.Email = user.Email

	// 数据校验
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}

	// 数据操作
	res, err := model.AppendPostAdmin(c.Param("id"), req)
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}

	utils.AjaxSuccess(c, res)
}

// PostsEdit 编辑评论
func PostsEdit(c *gin.Context) {
	// 数据解析
	req := common.Post{}
	err := c.BindJSON(&req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}

	// 数据操作
	err = model.EditPost(c.Param("id"), req, true)
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}

	utils.AjaxSuccess(c, req)
}

// FCMTokenReplace 替换 FCM 令牌
func FCMTokenReplace(c *gin.Context) {
	// 未启用推送，不进行任何操作
	if !config.Content.Push.Enabled {
		utils.AjaxSuccess(c, nil)
		return
	}

	req := struct {
		OldToken string `json:"oldToken"`
		NewToken string `json:"newToken"`
	}{}

	err := c.BindJSON(&req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}
	err = model.ReplaceFCMToken(req.OldToken, req.NewToken)
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}
	utils.AjaxSuccess(c, nil)
}

func FCMTokenDelete(c *gin.Context) {
	// 未启用推送，不进行任何操作
	if !config.Content.Push.Enabled {
		utils.AjaxSuccess(c, nil)
		return
	}

	req := struct {
		Token string `json:"token"`
	}{}

	err := c.BindJSON(&req)
	if err != nil {
		utils.AjaxError(c, http.StatusBadRequest, err)
		return
	}
	err = model.DeleteFCMToken(req.Token)
	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}
	utils.AjaxSuccess(c, nil)
}

func CacheDelete(c *gin.Context) {
	err := dao.DeleteAllCache()

	if err != nil {
		utils.AjaxError(c, http.StatusInternalServerError, err)
		return
	}
	utils.AjaxSuccess(c, nil)
}
