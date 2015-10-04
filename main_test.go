package main

import (
	//	"fmt"
	"strings"
	"testing"
)

func TestGetHtml(t *testing.T) {
	//创建URL
	url, _ := NewURL("http://themepixels.com/themes/demo/webpage/starlight/gray/index.html")

	//获取URL内容
	getter := NewGetter(1)
	content, err := getter.LoadUrl(url)
	if err != nil {
		t.Error("Error:", err)
		return
	}

	//创建文件内容
	fileContent := NewFileContent(string(content), "")

	//分析文件内容
	parser := NewParser(url, fileContent)
	links := parser.Do()

	//判断分析出的链接数量
	if len(links) != 5 {
		t.Errorf("Number of links is %d, except 5", len(links))
	}

	//判断分析出的链接是为绝对路径
	for _, l := range links {
		if !strings.HasPrefix(l, "http") {
			t.Errorf("URL:%s is not absolute url", l)
		}
	}

	//保存文件
	fpath, bs, err := fileContent.Save("./data/test", url)
	if err != nil {
		t.Error(err)
		return
	}
	if fpath != `data\test\themepixels.com\themes\demo\webpage\starlight\gray\index.html` {
		t.Errorf(`Path Error. Get %s, except "%s"`,
			fpath,
			`data\test\themepixels.com\themes\demo\webpage\starlight\gray\index.html`)
	}
	if bs != len(fileContent.RawContent) {
		t.Error("No data write to file " + fpath)
	}
}
