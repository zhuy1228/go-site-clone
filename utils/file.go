package utils

import (
	"fmt"
	"go-site-clone/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type File struct{}

func (*File) Download(uri string) string {
	// 解析文件链接及 路径
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	if u.Port() != "" {
		host += u.Port()
	}
	// 获取路径部分
	filePath := u.Path
	appConfig, _ := config.LoadConfig()
	fp := filepath.Join(appConfig.SiteFileDir, host, filePath)
	// 获取文件名
	// fileName := path.Base(filePath)
	log.Println("开始下载：", uri)
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("下载失败:", err)
		return ""
	}
	defer resp.Body.Close()
	// 创建本地文件
	outFile, err := CreateFileWithDirs(fp)
	if err != nil {
		fmt.Println("文件创建失败:", err)
		return ""
	}
	defer outFile.Close()
	// 将响应内容写入文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Println("写入失败:", err)
		return ""
	}
	fmt.Println("下载完成:", fp)
	return fp
}

func CreateFileWithDirs(filePath string) (*os.File, error) {
	// 取出目录部分
	dir := filepath.Dir(filePath)

	// 确保目录存在（递归创建）
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件（如果已存在会清空）
	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}

	return f, nil
}

func (*File) HTMLDownload(uri string) string {
	// 解析文件链接及 路径
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	if u.Port() != "" {
		host += u.Port()
	}
	// 获取路径部分
	filePath := u.Path
	appConfig, _ := config.LoadConfig()
	// 获取文件名
	fileName := path.Base(filePath)
	if fileName == "/" || filePath == "" {
		filePath = filePath + "/index.html"
	}
	fp := filepath.Join(appConfig.SiteFileDir, host, filePath)

	log.Println("开始下载：", uri)
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("下载失败:", err)
		return ""
	}
	defer resp.Body.Close()
	// 创建本地文件
	outFile, err := CreateFileWithDirs(fp)
	if err != nil {
		fmt.Println("文件创建失败:", err)
		return ""
	}
	defer outFile.Close()
	// 将响应内容写入文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Println("写入失败:", err)
		return ""
	}
	fmt.Println("下载完成:", fp)
	return fp
}
