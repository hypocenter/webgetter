package main

import (
	"fmt"
	"github.com/opesun/goquery"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

type Parser struct {
	URL         *URL
	FileContent *FileContent
}

func NewParser(url *URL, fileContent *FileContent) *Parser {
	ps := &Parser{url, fileContent}
	return ps
}

// 分析文件内容
func (ps *Parser) Do() []string {

	list := []string{}

	switch strings.ToLower(ps.URL.File.Extname) {
	case ".css":
		list = ps.parseCss()
	case ".html", ".php":
		list = ps.parseHtml()
	}

	res := make([]string, 0, len(list))
	for _, v := range list {
		if strings.Replace(v, " ", "", -1) != "" {
			res = append(res, ps.ToAbs(v))
		}
	}

	return res
}

func (ps *Parser) parseCss() (links []string) {
	links = make([]string, 5)

	//匹配背景图url
	exp, _ := regexp.Compile(`url\s*\(\s*['\"]?\s*([^'\"\)]+)\s*\s*['\"]?\s*\)`)
	matchs := exp.FindAllSubmatch([]byte(ps.FileContent.Content), 1024)
	if matchs == nil {
		return
	}

	for _, match := range matchs {
		links = append(links, string(match[1]))
	}

	return
}

func (ps *Parser) parseHtml() []string {
	q, _ := goquery.Parse(strings.NewReader(ps.FileContent.Content))

	ls := []string{}

	var nodes goquery.Nodes
	// 链接
	nodes = q.Find("a")
	ls = append(ls, ps.getAttr(nodes, "href", nil)...)

	//css
	nodes = q.Find("link")
	ls = append(ls, ps.getAttr(nodes, "href", map[string]string{"type": "text/css"})...)

	//js
	nodes = q.Find("script")
	ls = append(ls, ps.getAttr(nodes, "src", nil)...)

	//图片
	nodes = q.Find("img")
	ls = append(ls, ps.getAttr(nodes, "src", nil)...)

	return ls
}

func (ps *Parser) getAttr(nodes goquery.Nodes, key string, condition map[string]string) []string {
	sl := []string{}
	var (
		flag int    //记录条件是否都成立
		val  string //暂存当前获取的key的值
	)

	for _, j := range nodes {
		//初始化
		flag = len(condition)
		val = ""

		for _, v := range j.Attr {
			if _, ok := condition[v.Key]; ok && strings.ToLower(v.Val) == condition[v.Key] {
				flag-- //一个条件成立
			}

			if v.Key == key {
				val = v.Val
			}
		}

		//判断是否找到了值并且所有条件都满足
		if flag == 0 && val != "" {
			sl = append(sl, val)
		}
	}

	return sl
}

// 转换成绝对链接
func (ps *Parser) ToAbs(val string) (link string) {
	val = strings.Trim(val, " ")

	//排除空连接，#，javascript开头的链接
	if val != "" && val != "#" && !strings.HasPrefix(val, "javascript") {
		urlInfo, err := url.Parse(val)

		if err != nil {
			fmt.Println(err)
			return link
		}

		if urlInfo.IsAbs() {
			//TODO 对于本身就是绝对路径的url，判断是否是非当前主机下的文件
			//如果不是讲做链接替换，把HTML的链接替换成本地路径
			//目前暂时跳过非当前主机连接，不下载，不替换，返回空字符串
			if urlInfo.Host == ps.URL.Host {
				link = val
			}
		} else {
			//处理相对链接
			if strings.HasPrefix(val, "/") {
				// 站点绝对路径,拼接当前文档URL的主机地址
				link = ps.URL.HostAddr + "/" + val
			} else {
				// 完全的相对路径,拼接当前文档URL的目录地址
				link = ps.URL.DirAddr + "/" + val
			}
		}

		//清除./ // ../ 这样的路径
		link = filepath.Clean(link)
		link = strings.Replace(link, `\`, "/", -1)
		link = strings.Replace(link, `http:/`, "http://", 1)
	}

	return
}
