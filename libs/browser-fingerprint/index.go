package browserfingerprint

import (
	paramsTypes "go-site-clone/types"
	"math/rand"

	"github.com/go-rod/rod"
)

type BrowserFingerprint struct {
}

func (b *BrowserFingerprint) SetBrowserFingerprint(page *rod.Page, params *paramsTypes.BrowserFingerprintParams) *rod.Page {
	var p *rod.Page = page
	if params.Canvas {
		p = setCanvas(p)
	}
	if params.TimeZone != "" {
		p = setTimezoneAndLangAndGeo(p, params)
	}

	return p
}

func setCanvas(page *rod.Page) *rod.Page {
	// 注入噪声脚本
	// 生成随机 RGBA 偏移值 (Go 端控制)
	rgba := [4]int{
		rand.Intn(10) - 5, // R: -5 ~ +4
		rand.Intn(10) - 5, // G
		rand.Intn(10) - 5, // B
		rand.Intn(10) - 5, // A
	}
	// 使用更隐蔽的 Canvas 指纹修改方法
	page.MustEvalOnNewDocument(GetChangeCanvasJavaScript(rgba))
	return page
}

func setTimezoneAndLangAndGeo(page *rod.Page, params *paramsTypes.BrowserFingerprintParams) *rod.Page {
	page.MustEvalOnNewDocument(GetChangeTimezoneJavaScript(params.TimeZone, params.Language, params.GeoLocation))
	return page
}
