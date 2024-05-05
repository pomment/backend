package model

import (
	"encoding/json"
	"github.com/pomment/pomment/dao"
	"path/filepath"
)

func GetThreadList() (files []string, err error) {
	files, err = filepath.Glob(filepath.Join("threads/*.meta.json"))
	return files, err
}

func LoadFCMTokenList() (list []string, err error) {
	res, err := dao.ReadJSON("fcm-tokens.json")
	if err != nil {
		return list, err
	}

	// 读取 JSON，追加列表
	err = json.Unmarshal([]byte(res), &list)
	return list, err
}

func SaveFCMTokenList(list []string) error {
	jsonStr, err := json.Marshal(list)
	if err != nil {
		return err
	}

	err = dao.WriteJSON("fcm-tokens.json", string(jsonStr))
	return err
}

func ReplaceFCMToken(oldToken string, newToken string) (err error) {
	list, err := LoadFCMTokenList()
	if err != nil {
		return err
	}
	var newList []string

	// 逐步复制到新的数组，如果有则替换，没有则添加
	found := false
	for _, e := range list {
		if e == oldToken {
			found = true
			newList = append(newList, newToken)
			continue
		}
		if e == newToken {
			found = true
			newList = append(newList, e)
			continue
		}
		newList = append(newList, e)
	}
	if !found {
		newList = append(newList, newToken)
	}

	return SaveFCMTokenList(newList)
}

func DeleteFCMToken(token string) (err error) {
	list, err := LoadFCMTokenList()
	if err != nil {
		return err
	}
	var newList []string

	// 逐步复制到新的数组
	for _, e := range list {
		if e == token {
			continue
		}
		newList = append(newList, e)
	}

	return SaveFCMTokenList(newList)
}
