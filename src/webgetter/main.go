package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"webgetter/ptr"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	ptr.Show()

	var (
		//要获取的地址
		urls = []string{
			"http://themepixels.com/themes/demo/webpage/starlight/gray/index.html",
			"http://themepixels.com/themes/demo/webpage/starlight/gray/dashboard.html",
			"http://themepixels.com/themes/demo/webpage/starlight/darkblue/index.html",
			"http://themepixels.com/themes/demo/webpage/starlight/darkblue/dashboard.html",
			"http://themepixels.com/themes/demo/webpage/starlight/index.html",
			"http://themepixels.com/themes/demo/webpage/starlight/dashboard.html",
			// "http://psd.dev/080_sparklingpixels/Sparkpsd2html01.html",
			// "http://psd.dev/080_sparklingpixels/Sparkpsd2html02.html",
		}
		//获取链接的并发数
		getterNum = 10
		//项目名
		projectName = "test"
	)

	//获取器
	var getter = NewGetter(getterNum)

	for _, v := range urls {
		//创建URL
		url, err := NewURL(v)
		if err == nil {
			getter.AddUrl(url)
		} else {
			fmt.Fprintln(os.Stderr, "Error to parse URL: ", url, err.Error())
		}
	}

	var timeout = make(chan bool)
	var counter = 0
	go func(t chan bool) {
		for {
			time.Sleep(time.Second)
			counter++
			if counter > 30 { //30秒超时
				t <- true
			}
		}
	}(timeout)

	var result *Result
	for {
		select {
		case result = <-getter.Results:
			//处理获得的内容
			go func(res *Result) {
				var fileContent = NewFileContent(string(res.Content), "")

				//分析文件内容
				var parser = NewParser(res.URL, fileContent)
				var links = parser.Do()
				if len(links) > 0 {
					fmt.Println("From url:", res.URL.String(), "get links:")
					for _, v := range links {
						fmt.Println(" -", v)
					}
				}
				//添加结果中的新链接到获取器中
				for _, v := range links {
					u, err := NewURL(v)
					if err == nil {
						getter.AddUrl(u)
					}
				}
				//保存
				fn, bs, err := fileContent.Save("./data/"+projectName, res.URL)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				} else {
					fmt.Printf("Saved file '%s' with %d bytes.\n", fn, bs)
				}

			}(result)

			//重置超时检测的状态
			counter = 0

		case <-timeout:
			//超时退出
			os.Exit(0)
		}
	}

}

func checkError(err error) {
	fmt.Println(err)
}
