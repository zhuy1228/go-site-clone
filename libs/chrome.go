package libs

import (
	"log"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

type Chrome struct{}

func (c *Chrome) Create(accountId string, params *Fingerprint) (*rod.Browser, string, *rod.Page) {
	userDataDir := "user-data/" + accountId
	absPath, _ := filepath.Abs(userDataDir)
	log.Printf("用户数据目录: %s", absPath)
	path, _ := launcher.LookPath() // 获取当前系统的浏览器目录
	// 判断是否传递代理
	// var proxyHandle utils.ProxyHandle
	// protocol, username, password, ip, port, errProxy := proxyHandle.ParseProxy(params.Proxy)

	// ipInfo := external.GetIpInfo(strings.TrimSpace(string(params.Proxy)))

	l := launcher.New().
		Bin(path).
		// Proxy("210.51.27.121:50003"). // 直接传入完整认证URL
		Delete("use-mock-keychain").                           // delete flag "--use-mock-keychain"
		Set("disable-blink-features", "AutomationControlled"). // 绕过自动化检测
		// Set("incognito").                                      // 无痕模式
		Set("user-data-dir", absPath).   // 数据持久化目录
		Set("window-size", "1920,1480"). // 窗口尺寸
		Set("disable-infobars", "true"). // 隐藏自动化提示栏
		Set("no-sandbox", "true").
		Set("excludeSwitches", "enable-automation").
		Set("enable-gpu").                                  // 启用 GPU 加速
		Set("ignore-certificate-errors").                   // 忽略证书错误
		Set("use-fake-ui-for-media-stream").                // 允许媒体流
		Set("autoplay-policy", "no-user-gesture-required"). // 自动播放
		Set("ignore-certificate-errors").
		Set("disable-application-cache").
		Set("disable-dev-shm-usage").
		// 禁用webrtc
		Set("disable-webrtc", "true").                               // 核心禁用参数
		Set("disable-features", "WebRtcHideLocalIpsWithMdns").       // 隐藏本地IP
		Set("webrtc-ip-handling-policy", "disable_non_proxied_udp"). // 禁用非代理UDP
		Set("force-webrtc-ip-handling-policy").
		Headless(false) // 关闭无头模式

	// 修改语言
	l.Set("lang", "en-US")
	l.Set("accept-lang", "en-US")

	// // 修改时区
	// l.Set("timezone", ipInfo.Timezone)
	// l.Set("disable-geolocation", "true")  // 启用地理定位API
	// l.Set("disable-web-security", "true") // 允许跨域（某些网站需要）
	// l.Set("disable-notifications", "true")

	// // 设置环境变量 - 核心时区设置
	// l.Env(append(os.Environ(), "TZ="+ipInfo.Timezone)...)
	// l.
	// 	Set("timezone", ipInfo.Timezone).
	// 	Set("geolocation", ipInfo.Lat+","+ipInfo.Lon)

	// if errProxy == nil {
	// 	if strings.Contains(protocol, "http") {
	// 		log.Println(protocol + "://" + ip + ":" + port)
	// 		l = l.Proxy(ip + ":" + port)
	// 	}
	// }

	uri := l.MustLaunch()

	browser := rod.New().
		ControlURL(uri).
		MustConnect().NoDefaultDevice()
	browser.MustIgnoreCertErrors(true)
	// if errProxy == nil {
	// 	if strings.Contains(protocol, "http") {
	// 		log.Println(username, password)
	// 		go browser.MustHandleAuth(username, password)()
	// 	}
	// }
	page := stealth.MustPage(browser)
	// page := browser.MustPage()

	go browser.EachEvent(func(e *proto.TargetTargetCreated) {
		if e.TargetInfo.Type != proto.TargetTargetInfoTypePage {
			return
		}
		browser.MustPageFromTargetID(e.TargetInfo.TargetID).MustEvalOnNewDocument(stealth.JS)
	})()

	// // 在页面加载前注入时区覆盖代码
	// var browserFingerprint browserfingerprint.BrowserFingerprint
	// page = browserFingerprint.SetBrowserFingerprint(page, &paramsTypes.BrowserFingerprintParams{
	// 	Canvas:      true,
	// 	TimeZone:    ipInfo.Timezone, // 太平洋时间 (UTC-8)
	// 	Language:    "en-US",
	// 	GeoLocation: ipInfo.Lat + "," + ipInfo.Lon, // 洛杉矶坐标
	// })
	page.MustSetViewport(1920, 1480, 1.0, false)

	return browser, uri, page
}
