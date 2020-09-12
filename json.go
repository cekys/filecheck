package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/cekys/gopkg"
)

//输入结构体指针与json文件路径,将json内部数据存储到结构体中
func readJSON(stPointer interface{}, jsonFile string) error {
	//打开json文件
	fileData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return err
	}

	//json数据写入结构体
	err = json.Unmarshal(fileData, stPointer)
	if err != nil {
		return err
	}

	return nil
}

//根据现有的文件信息结构体生成json文件
func createJSONTable(configFile string, zipbool bool) error {
	//读取设置文件
	var setting Settings
	err := readJSON(&setting, configFile)
	if err != nil {
		return err
	}

	//创建存放json数据表的目录
	logPath, err := filepath.Abs(setting.Log.LogLocation)
	if err != nil {
		return err
	}
	mypkg.PathCreate(logPath)

	//对每一个配置文件中target进行操作
	for _, target := range setting.Targets {
		//获取target的绝对路径,并将"\"转换为统一的"/"
		root := target.Location
		root, err := filepath.Abs(root)
		if err != nil {
			return err
		}
		root = strings.Replace(root, "\\", "/", -1)

		//查看此目录是否确实存在
		if mypkg.PathExist(root) {

			//获取存储数据表的json文件名称
			jsonName := filepath.Join(logPath, target.Name+".json")
			jsonName = strings.Replace(jsonName, "\\", "/", -1)

			//获取存储数据表zip文件的名称
			zipName := mypkg.StringTrimSuffix(jsonName) + ".zip"

			switch target.Mode {
			case "fileInfo":
				var container []FileInfo
				//将目录下文件根据所选方式扫描后连接进存储信息表的切片
				st, err := pathToStruct(root, target.Mode, target.Filter)
				if err != nil {
					return err
				}
				container = append(container, st.([]FileInfo)...)
				//生成需要写入文件信息表的结构体
				data, err := structSetup(target.Name, root, target.Mode, container)
				if err != nil {
					return err
				}
				//data写入json
				err = writeJSON(data, jsonName)
				if err != nil {
					return err
				}
			case "crc32", "md5", "sha1", "sha256":
				var container []Hash
				//将目录下文件根据所选方式扫描后连接进存储信息表的切片
				st, err := pathToStruct(root, target.Mode, target.Filter)
				if err != nil {
					return err
				}
				container = append(container, st.([]Hash)...)
				//生成需要写入文件信息表的结构体
				data, err := structSetup(target.Name, root, target.Mode, container)
				if err != nil {
					return err
				}
				//data写入json
				err = writeJSON(data, jsonName)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown mode: %s", target.Mode)
			}

			//zip压缩,并删除原有json文件
			if zipbool {
				err := mypkg.Zipit(jsonName, zipName)
				if err != nil {
					return err
				}
				//当成功压缩后删除json文件
				os.RemoveAll(jsonName)
			}
		}
	}
	return nil
}

//生成需要写入文件信息表的结构体
func structSetup(name string, root string, mode string, slice interface{}) (interface{}, error) {
	var fileInfo FileList
	var hash HashList
	var f []FileInfo
	var h []Hash

	//根据mode不同决定生成哪一种结构体
	switch mode {
	case "fileInfo":
		//获取切片类型,当输入切片类型不是指定类型时返回空
		v := reflect.ValueOf(slice)
		if v.Type() != reflect.TypeOf(f) {
			return nil, fmt.Errorf("wrong slice type: %s", v.Type())
		}
		//填充数据
		fileInfo.Info.Name = name
		fileInfo.Info.Root = root
		fileInfo.Info.Time = time.Now()
		fileInfo.Info.Mode = mode
		fileInfo.FileInfo = slice.([]FileInfo)
		return fileInfo, nil

	case "crc32", "md5", "sha1", "sha256":
		//获取切片类型,当输入切片类型不是指定类型时返回空
		v := reflect.ValueOf(slice)
		if v.Type() != reflect.TypeOf(h) {
			return nil, fmt.Errorf("wrong slice type: %s", v.Type())
		}
		//填充数据
		hash.Info.Name = name
		hash.Info.Root = root
		hash.Info.Time = time.Now()
		hash.Info.Mode = mode
		hash.Hash = slice.([]Hash)
		return hash, nil
	}
	//默认返回空
	return nil, errors.New("default structSetup error")
}

//写入json文件
func writeJSON(st interface{}, fileName string) error {
	//从结构体生成json的byte数据
	jsonData, err := json.MarshalIndent(st, "", "    ")
	if err != nil {
		return err
	}

	//json的byte数据写入json文件
	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
