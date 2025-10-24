package main

import (
	"go-site-clone/services"
	"go-site-clone/utils"

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
		if v.Type == proto.NetworkResourceTypeMedia {
			resources.Video = append(resources.Video, v.URL)
		}
	}

	return resources
}

func (a *App) DownloadSite(obj services.ResourcesList) bool {
	// 将页面及资源一起返回

	var File utils.File
	if len(obj.CSS) > 0 {
		for k, v := range obj.CSS {
			File.Download(v)
			a.app.Event.Emit("download:css", k)
		}
	}
	if len(obj.Script) > 0 {
		for k, v := range obj.Script {
			File.Download(v)
			a.app.Event.Emit("download:script", k)
		}
	}
	if len(obj.Image) > 0 {
		for k, v := range obj.Image {
			File.Download(v)
			a.app.Event.Emit("download:script", k)
		}
	}
	if len(obj.Video) > 0 {
		for k, v := range obj.Video {
			File.Download(v)
			a.app.Event.Emit("download:video", k)
		}
	}
	if len(obj.Dom) > 0 {
		for k, v := range obj.Dom {
			File.HTMLDownload(v)
			a.app.Event.Emit("download:dom", k)
		}
	}
	return true
}
