package main

import (
	"errors"
	"net/url"
	"path"
	"strings"
)

type URL struct {
	Addr     string //原始地址，如：http://xxx.com/dir/index.html 必须是带主机名的绝对链接
	HostAddr string //域名地址，如：http://xxx.com
	DirAddr  string //目录地址，如：http://xxx.com/dir
	*url.URL        //组合进url.URL
	File     *File  //文件信息
}

func NewURL(rawurl string) (*URL, error) {
	//*url.URL
	u, err := url.ParseRequestURI(rawurl) //此处保证非相对链接
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		// 没有主机名，一样视为非法链接
		return nil, errors.New("invalid URI for request")
	}

	//截取url的域名部分如：http://xxx.com
	hostAddr := rawurl[:strings.Index(rawurl, u.Path)]

	//目录地址
	dirAddr := path.Dir(rawurl)

	//*File
	f := &File{
		path.Join(u.Host, u.Path),
		path.Base(u.Path),
		path.Ext(u.Path),
	}

	return &URL{rawurl, hostAddr, dirAddr, u, f}, nil
}
