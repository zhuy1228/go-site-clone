package main

import (
	"go-site-clone/services"

	"github.com/go-rod/rod/lib/proto"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	app *application.App
}

var siteService services.SiteService

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
	}

	return resources
}

func (a *App) DownloadSite() {
	// 将页面及资源一起返回
	// for _, v := range dataList {
	// 	u, _ := url.Parse(v.URL)
	// 	if v.Type == proto.NetworkResourceTypeImage && u.Hostname() == parsed.Hostname() {
	// 		var File utils.File
	// 		File.Download(v.URL)
	// 	}
	// }
	// for _, v := range RouterAll {
	// 	var File utils.File
	// 	File.HTMLDownload(v)
	// }
}
