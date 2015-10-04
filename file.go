package main

import (
	//	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Path     string //路径，如：xxx.com/xx/index.html
	Basename string //文件名，如：index.html
	Extname  string //扩展名，如：.html
}

type FileContent struct {
	RawContent string //原始字符串
	Content    string //转换成UTF8后的字符串
	Charset    string //原始编码
}

func NewFileContent(rawcontent string, charset string) *FileContent {
	f := &FileContent{
		RawContent: rawcontent,
		Charset:    charset,
	}

	//编码转换
	if f.Charset != "" && f.Charset != "utf8" {
		reader := transform.NewReader(strings.NewReader(rawcontent), simplifiedchinese.GBK.NewDecoder())
		data, _ := ioutil.ReadAll(reader)
		f.Content = string(data)
	} else {
		f.Content = rawcontent
	}

	return f
}

func (f *FileContent) Save(basedir string, url *URL) (file string, bs int, err error) {
	file = filepath.Clean(filepath.Join(basedir, url.File.Path))
	dir := filepath.Dir(file)

	//创建目录
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	//创建文件
	fileobj, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer fileobj.Close()

	//写入原始文件内容
	bs, err = fileobj.WriteString(f.RawContent)
	if err != nil {
		return
	}

	return
}
