package main

import "time"

const (
	//DefaultConfig Default Config file path
	DefaultConfig = "config.json"
	//DefaultFileInfoMode Default hash mode
	DefaultFileInfoMode = "fileInfo"
	//DefaultHashMode Default hash mode
	DefaultHashMode = "crc32"
)

//Settings 设置结构体第一层
type Settings struct {
	Log     Log       `json:"log"`
	Targets []Targets `json:"targets"`
}

//Log 设置结构体第二层
type Log struct {
	LogLocation    string `json:"logLocation"`
	ReportLocation string `json:"reportLocation"`
}

//Targets 设置结构体第二层
type Targets struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Location    string   `json:"location"`
	Mode        string   `json:"mode"`
	Filter      []string `json:"filter,omitempty"`
}

//FileList 文件信息保存表第一层
type FileList struct {
	Info     Info       `json:"info"`
	FileInfo []FileInfo `json:"fileInfo"`
}

//FileInfo 文件信息保存表第二层
type FileInfo struct {
	Path string    `json:"path"`
	Size int64     `json:"size"`
	Time time.Time `json:"time"`
}

//HashList 哈希信息保存表第一层
type HashList struct {
	Info Info   `json:"info"`
	Hash []Hash `json:"hash"`
}

//Info 信息表第二层(多个表使用)
type Info struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"`
	Root string    `json:"root"`
	Mode string    `json:"mode"`
}

//Hash 哈希信息保存表第二层
type Hash struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

//DiffList 差异保存表第一层
type DiffList struct {
	Info Info   `json:"info"`
	Diff []Diff `json:"diff"`
}

//Diff 差异保存表第二层
type Diff struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}
