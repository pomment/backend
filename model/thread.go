package model

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/pomment/pomment/common"
	"github.com/pomment/pomment/dao"
	"log"
)

// GetThreadRelationList 获取评论串关系列表
func GetThreadRelationList() (list []common.ThreadMapItem, err error) {
	res, err := dao.ReadJSON("index.json")
	if err != nil {
		log.Printf("Unable to read index.json: %s. Loading initial list instead.", err)
		return make([]common.ThreadMapItem, 0), nil
	}

	err = json.Unmarshal([]byte(res), &list)
	return list, err
}

// GetThreadRelationByURL 获取评论串关系（URL 换 ID）
func GetThreadRelationByURL(url string) (item *common.ThreadMapItem, err error) {
	list, err := GetThreadRelationList()
	if err != nil {
		return nil, err
	}

	for _, e := range list {
		if e.URL == url {
			return &e, nil
		}
	}
	return nil, errors.New("thread not found")
}

// GetThreadRelationByID 获取评论串关系（ID 换 URL）
func GetThreadRelationByID(id string) (item *common.ThreadMapItem, err error) {
	list, err := GetThreadRelationList()
	if err != nil {
		return nil, err
	}

	for _, e := range list {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, errors.New("thread not found")
}

// AddThreadRelation 增加评论串关系记录
func AddThreadRelation(id string, url string) (err error) {
	list, err := GetThreadRelationList()
	if err != nil {
		return err
	}

	list = append(list, common.ThreadMapItem{
		ID:  id,
		URL: url,
	})
	jsonStr, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = dao.WriteJSON("index.json", string(jsonStr))
	return err
}

// GetThreadMeta 获取评论串元数据
func GetThreadMeta(id string) (item *common.Thread, err error) {
	res, err := dao.ReadJSON(dao.GetThreadMetaPath(id))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(res), &item)
	return item, nil
}

// GetThreadMetaByURL 获取评论串元数据（通过 URL）
func GetThreadMetaByURL(url string) (item *common.Thread, err error) {
	relation, err := GetThreadRelationByURL(url)
	if err != nil {
		return nil, err
	}
	return GetThreadMeta(relation.ID)
}

// SetThreadMeta 单独更新评论串元数据
func SetThreadMeta(item common.Thread) (err error) {
	res, err := json.Marshal(item)
	if err != nil {
		return err
	}

	// 元数据写入
	err = dao.WriteJSON(dao.GetThreadMetaPath(item.ID), string(res))
	return err
}

// RegisterThread 手动注册评论串元数据
func RegisterThread(item common.Thread) (err error) {
	res, err := json.Marshal(item)
	if err != nil {
		return err
	}

	// 元数据写入
	err = dao.WriteJSON(dao.GetThreadMetaPath(item.ID), string(res))
	err = dao.WriteJSON(dao.GetThreadPath(item.ID), "[]")

	// 注册评论串
	err = AddThreadRelation(item.ID, item.URL)
	return err
}

// GetThreadMetaForSubmit 获取评论串（如果不存在则增加一条）
func GetThreadMetaForSubmit(url string, title string, id string) (item *common.Thread, err error) {
	resRelation, err := GetThreadRelationByURL(url)
	// 没有评论串，需要新建元数据
	if err != nil {
		var item common.Thread

		if id == "" {
			item.ID = uuid.New().String()
		} else {
			item.ID = id
		}
		item.URL = url
		item.Title = title
		err = RegisterThread(item)
		if err != nil {
			return nil, err
		}
		return &item, nil
	}

	resMeta, err := GetThreadMeta(resRelation.ID)
	return resMeta, err
}

func UpdateThreadMeta(id string) (meta *common.Thread, err error) {
	post, err := GetPostsSimpleByID(id)
	meta, err = GetThreadMeta(id)
	meta.Amount = 0

	if err != nil {
		return nil, err
	}

	if len(*post) >= 1 {
		meta.FirstPostAt = (*post)[0].CreatedAt
		meta.LatestPostAt = (*post)[0].CreatedAt
	}

	// 更新最晚评论时间和评论计数
	for _, e := range *post {
		if meta.FirstPostAt > e.CreatedAt {
			meta.FirstPostAt = e.CreatedAt
		}
		if meta.LatestPostAt < e.CreatedAt {
			meta.LatestPostAt = e.CreatedAt
		}
		if !e.Hidden {
			meta.Amount++
		}
	}

	// 应用元数据更新
	err = SetThreadMeta(*meta)
	return meta, err
}

func UpdateAllThreadMeta() (err error) {
	list, err := GetThreadRelationList()
	if err != nil {
		return err
	}
	for _, e := range list {
		_, err = UpdateThreadMeta(e.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
