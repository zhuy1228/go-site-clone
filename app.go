package main

import (
	"fmt"
	"go-site-clone/config"
	"go-site-clone/services"
	"go-site-clone/utils"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-rod/rod/lib/proto"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	app *application.App
}

var siteService services.SiteService

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
// func (a *App) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {

// 	return nil
// }

// func (a *App) ServiceShutdown(ctx context.Context) error {

// 	return nil
// }

func (a *App) GetResources(rawURL string) *services.ResourcesList {
	RequestParamsAll, RouterAll := siteService.GetAllResources(rawURL)
	// 再次去重 根据url链接去除重复项
	dataList := siteService.DeduplicationRequestByUrl(RequestParamsAll)
	resources := &services.ResourcesList{
		Dom: RouterAll,
	}
	// 将页面及资源一起返回
	for _, v := range dataList {
		if v.Type == proto.NetworkResourceTypeImage {
			resources.Image = append(resources.Image, v.URL)
		}
		if v.Type == proto.NetworkResourceTypeStylesheet {
			resources.CSS = append(resources.CSS, v.URL)
		}
		if v.Type == proto.NetworkResourceTypeScript {
			resources.Script = append(resources.Script, v.URL)
		}
		if v.Type == proto.NetworkResourceTypeMedia {
			resources.Video = append(resources.Video, v.URL)
		}
	}
	return resources
}

func (a *App) DownloadSite(uri string, obj services.ResourcesList) bool {
	// 将页面及资源一起返回
	parsed, _ := url.Parse(uri)
	var File utils.File
	if len(obj.CSS) > 0 {
		for k, v := range obj.CSS {
			u, _ := url.Parse(v)
			if parsed.Hostname() == u.Hostname() {
				File.Download(v)
			}
			a.app.Event.Emit("download:css", k)
		}
	}
	if len(obj.Script) > 0 {
		for k, v := range obj.Script {
			u, _ := url.Parse(v)
			if parsed.Hostname() == u.Hostname() {
				File.Download(v)
			}
			a.app.Event.Emit("download:script", k)
		}
	}
	if len(obj.Image) > 0 {
		for k, v := range obj.Image {
			u, _ := url.Parse(v)
			if parsed.Hostname() == u.Hostname() {
				File.Download(v)
			}
			a.app.Event.Emit("download:image", k)
		}
	}
	if len(obj.Video) > 0 {
		for k, v := range obj.Video {
			u, _ := url.Parse(v)
			if parsed.Hostname() == u.Hostname() {
				File.Download(v)
			}
			a.app.Event.Emit("download:video", k)
		}
	}
	if len(obj.Dom) > 0 {
		for k, v := range obj.Dom {
			u, _ := url.Parse(v)
			if parsed.Hostname() == u.Hostname() {
				File.HTMLDownload(v)
			}
			a.app.Event.Emit("download:dom", k)
		}
	}
	return true
}

func (a *App) GetDownloadList() []utils.FileDir {
	return siteService.GetLocalSiteList()
}

// 打开网站文件夹
func (a *App) OpenSiteFileDir(pathDir string) bool {
	appConfig, _ := config.LoadConfig()
	newPath := filepath.Join(appConfig.SiteFileDir, pathDir)
	cmd := exec.Command("explorer", newPath)
	err := cmd.Start()
	if err != nil {
		fmt.Println("打开文件夹失败:", err)
		return false
	}
	return true
}

func (a *App) DeleteSiteFileDir(pathDir string) bool {
	appConfig, _ := config.LoadConfig()
	newPath := filepath.Join(appConfig.SiteFileDir, pathDir)
	err := os.RemoveAll(newPath)
	if err != nil {
		fmt.Println("删除文件夹失败:", err)
		return false
	}
	return true
}
