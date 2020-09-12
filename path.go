package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/cekys/gopkg"
)

//扫描路径下文件,并将文件的文件信息存储返回一个切片
func pathToStruct(root string, mode string, filter []string) (interface{}, error) {
	var fileList []string
	var fileInfo FileInfo
	var fileSlice []FileInfo
	var hash Hash
	var hashSlice []Hash

	//获取root的绝对路径
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	root = strings.Replace(root, "\\", "/", -1)

	//获取文件列表
	err = mypkg.PathWalk(root, &fileList, filter, true)
	if err != nil {
		return nil, err
	}

	for _, file := range fileList {
		//所有 "\" 替换为 "/"
		file = strings.Replace(file, "\\", "/", -1)

		//获取每个文件的信息切片
		f, err := os.Stat(file)
		if err != nil {
			return nil, err
		}

		//文件路径(使用相对路径)
		filePath, err := filepath.Rel(root, file)
		if err != nil {
			return nil, err
		}
		filePath = strings.Replace(filePath, "\\", "/", -1)

		switch mode {
		case "fileInfo":
			//填充节点数据
			fileInfo.Path = filePath
			fileInfo.Size = f.Size()
			fileInfo.Time = f.ModTime()
			//节点连接到切片
			fileSlice = append(fileSlice, fileInfo)
		case "crc32", "md5", "sha1", "sha256":
			//填充节点数据
			hash.Path = filePath
			hash.Value = mypkg.Checksum(mode, file)
			//节点连接到切片
			hashSlice = append(hashSlice, hash)
		}
	}

	//返回对应的切片
	switch mode {
	case "fileInfo":
		return fileSlice, nil
	case "crc32", "md5", "sha1", "sha256":
		return hashSlice, nil
	}

	//默认返回空
	return nil, errors.New("default pathToStruct error")
}
