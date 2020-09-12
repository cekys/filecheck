package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cekys/gopkg"
)

func compare(configFile string) ([]Diff, []string, error) {
	var diffs []Diff
	newFiles := []string{""}
	filter := []string{"$Recycle.Bin", "System Volume Information"}

	//读取设置文件
	var setting Settings
	err := readJSON(&setting, configFile)
	if err != nil {
		return diffs, newFiles, err
	}

	//获取记录文件存储路径
	logPath, err := filepath.Abs(setting.Log.LogLocation)
	if err != nil {
		return diffs, newFiles, err
	}

	//对每一个配置文件中target进行操作
	for _, target := range setting.Targets {
		var diff Diff

		//获取存储数据表的json文件名称
		jsonName := filepath.Join(logPath, target.Name+".json")
		jsonName = strings.Replace(jsonName, "\\", "/", -1)

		//获取存储数据表zip文件的名称
		zipName := mypkg.StringTrimSuffix(jsonName) + ".zip"

		//获取存储数据表的root路径
		root := target.Location
		root, err := filepath.Abs(root)
		if err != nil {
			return diffs, newFiles, err
		}
		root = strings.Replace(root, "\\", "/", -1)

		//获取zip数据压缩包以及json数据文件的存在情况
		zipExist := mypkg.PathExist(zipName)
		jsonExist := mypkg.PathExist(jsonName)

		//如果数据json与数据zip文件均不存在,返回错误
		if !(zipExist || jsonExist) {
			return diffs, newFiles, fmt.Errorf("can't find data file %s or %s", jsonName, zipName)
		}

		//如果存在zip文件就解压
		if zipExist {
			err := mypkg.Unzip(zipName, logPath)
			if err != nil {
				return diffs, newFiles, err
			}
		}

		//读取新的文件列表
		err = mypkg.PathWalk(root, &newFiles, filter, true)
		if err != nil {
			return diffs, newFiles, err
		}

		//判断当前target的mode,使用对应的模式来读取json文件到结构体
		switch target.Mode {
		case "fileInfo":
			//读取记录文件
			var oldFileList FileList
			readJSON(&oldFileList, jsonName)
			//处理每一个记录
			for _, ofile := range oldFileList.FileInfo {
				//拼接记录文件完整地址
				oldFile := filepath.Join(oldFileList.Info.Root, ofile.Path)
				oldFile = strings.Replace(oldFile, "\\", "/", -1)

				//在新文件列表中找到记录文件中的数据
				position := mypkg.StringFindInSlice(newFiles, oldFile)

				if position != -1 {
					//如果匹配,删除新文件列表里的对应元素
					mypkg.SliceDelete(&newFiles, position)

					//获取文件的信息切片
					f, err := os.Stat(oldFile)
					if err != nil {
						return diffs, newFiles, err
					}
					//计算特征值是否符合
					timeMatch := f.ModTime() == ofile.Time
					sizeMatch := f.Size() == ofile.Size

					//填写diffs表
					if !timeMatch || !sizeMatch {
						diff.Path = oldFile
						if timeMatch && !sizeMatch {
							diff.Reason = "time"
						}
						if !timeMatch && sizeMatch {
							diff.Reason = "size"
						}
						if !timeMatch && !sizeMatch {
							diff.Reason = "both"
						}
						diffs = append(diffs, diff)
					}
				}
			}
		case "crc32", "md5", "sha1", "sha256":
			//读取记录文件
			var oldHashList HashList
			readJSON(&oldHashList, jsonName)
			//处理每一个记录
			for _, oHash := range oldHashList.Hash {
				//拼接记录文件完整地址
				oldHash := filepath.Join(oldHashList.Info.Root, oHash.Path)
				oldHash = strings.Replace(oldHash, "\\", "/", -1)

				//在新文件列表中找到记录文件中的数据
				position := mypkg.StringFindInSlice(newFiles, oldHash)

				if position != -1 {
					//如果匹配,删除新文件列表里的对应元素
					mypkg.SliceDelete(&newFiles, position)

					//计算特征值是否符合
					hashMatch := mypkg.Checksum(target.Mode, oldHash) == oHash.Value

					//填写diffs表
					if !hashMatch {
						diff.Path = oldHash
						diffs = append(diffs, diff)
					}
				}
			}
		default:
			return diffs, newFiles, fmt.Errorf("unknown mode: %s", target.Mode)
		}
	}
	return diffs, newFiles, nil
}
