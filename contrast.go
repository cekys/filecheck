package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"../mypkg"
)

func compare(configFile string) error {
	//读取设置文件
	var setting Settings
	err := readJSON(&setting, configFile)
	if err != nil {
		return err
	}

	//获取记录文件存储路径
	logPath, err := filepath.Abs(setting.Log.LogLocation)
	if err != nil {
		return err
	}

	//对每一个配置文件中target进行操作
	for _, target := range setting.Targets {
		var NewFiles []string
		var filter []string

		//获取存储数据表的json文件名称
		jsonName := filepath.Join(logPath, target.Name+".json")
		jsonName = strings.Replace(jsonName, "\\", "/", -1)

		//获取存储数据表zip文件的名称
		zipName := mypkg.StringTrimSuffix(jsonName) + ".zip"

		//获取存储数据表的root路径
		root := target.Location
		root, err := filepath.Abs(root)
		if err != nil {
			return err
		}
		root = strings.Replace(root, "\\", "/", -1)

		//获取zip数据压缩包以及json数据文件的存在情况
		zipExist := mypkg.PathExist(zipName)
		jsonExist := mypkg.PathExist(jsonName)

		//如果数据json与数据zip文件均不存在,返回错误
		if !(zipExist && jsonExist) {
			return fmt.Errorf("can't find data file %s or %s", jsonName, zipName)
		}

		//如果存在zip文件就解压
		if zipExist {
			err := mypkg.Unzip(zipName, logPath)
			if err != nil {
				return err
			}
		}

		//读取新的文件列表
		err = mypkg.PathWalk(root, NewFiles, filter, true)
		if err != nil {
			return err
		}

		//判断当前target的mode,使用对应的模式来读取json文件到结构体
		switch target.Mode {
		case "fileInfo":
			var fileList FileList
			readJSON(&fileList, jsonName)
		case "crc32", "md5", "sha1", "sha256":
			var hashList HashList
			readJSON(&hashList, jsonName)
		default:
			return fmt.Errorf("unknown mode: %s", target.Mode)
		}

	}

	return nil
}
