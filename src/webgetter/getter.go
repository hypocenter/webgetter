package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type Result struct {
	URL     *URL
	Content []byte
}

type Getter struct {
	loadedList map[string]*URL
	client     *http.Client
	Results    chan *Result
	mutex      sync.RWMutex
}

func NewGetter(n int) *Getter {
	getter := &Getter{
		loadedList: make(map[string]*URL), //连接用
		client:     &http.Client{},        //已加载的文件
		Results:    make(chan *Result, n), //结果队列
	}

	return getter
}

// 添加url到处理队列中
func (g *Getter) AddUrl(url *URL) {
	fmt.Println("Add URL:", url.String())

	//处理,获取
	_, err := g.LoadUrl(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		return
	}
}

// 加载一个URL
func (g *Getter) LoadUrl(url *URL) (content []byte, err error) {
	//先检查是否已经加载过了
	if g.loaded(url) {
		return
	}

	content, err = g.load(url)

	if err == nil {
		g.mutex.RLock()
		defer g.mutex.RUnlock()

		//添加当前URL到已加载的列表中去
		g.loadedList[url.Addr] = url
		//添加结果到结果队列
		g.Results <- &Result{url, content}
	}

	return
}

// 创建请求头
func (g *Getter) newRequest(url string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	//req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.2; WOW64; rv:16.0) Gecko/20100101 Firefox/16.0 FirePHP/0.7.1")

	return
}

// 判断链接是否被加载过
func (g *Getter) loaded(url *URL) bool {
	_, ok := g.loadedList[url.Addr]
	return ok
}

// 加载一个url
func (g *Getter) load(url *URL) (content []byte, err error) {

	//创建连接请求
	req, err := g.newRequest(url.Addr)
	if err != nil {
		return
	}

	//连接
	resp, err := g.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	//获取内容
	content, err = ioutil.ReadAll(resp.Body)

	return
}
