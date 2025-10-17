package libs

import (
	"sync"

	"github.com/go-rod/rod"
)

type Browser struct{}

type BrowserInfo struct {
	Id      string
	Browser *rod.Browser
	Url     string
	Page    *rod.Page
}

// 浏览器指纹数据
type Fingerprint struct {
	Viewport []struct {
		Height int
		Width  int
	}
	Proxy string
	Lang  string
}

var (
	BrowserPool = make(map[string]*BrowserInfo) // 初始化map
	poolMutex   sync.RWMutex                    // 添加读写锁保证并发安全
)

var ChromeExample Chrome

func (*Browser) GetBrowser(id string, params *Fingerprint) *BrowserInfo {
	poolMutex.RLock()
	value, ok := BrowserPool[id]
	poolMutex.RUnlock()

	if ok {
		return value
	}

	// 未找到时获取写锁（互斥）
	poolMutex.Lock()
	defer poolMutex.Unlock()

	// 双检查：避免在等待锁期间其他goroutine已创建
	if value, ok := BrowserPool[id]; ok {
		return value
	}

	// 创建新实例（补充CancelFunc）
	var ChromeExample Chrome
	browser, url, page := ChromeExample.Create(id, params)
	browserInfo := &BrowserInfo{
		Id:      id,
		Browser: browser,
		Url:     url,
		Page:    page,
	}
	BrowserPool[id] = browserInfo
	return browserInfo
}

func (*Browser) GetAllBrowser() map[string]*BrowserInfo {
	return BrowserPool
}

func (*Browser) CancelAllBrowser() {
	for _, v := range BrowserPool {
		v.Browser.MustClose()
	}
}

func (*Browser) CancelBrowser(id string) {
	browserInfo := BrowserPool[id]
	if browserInfo == nil {
		return
	}
	browserInfo.Browser.MustClose()
	delete(BrowserPool, id)
}
