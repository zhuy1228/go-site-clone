package services

import (
	"go-site-clone/libs"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type SiteService struct {
	BaseDomain string
	Port       string
}

type RequestParams struct {
	Type proto.NetworkResourceType `json:"type"`
	URL  string                    `json:"url"`
}

type ResourcesList struct {
	Script []string `json:"script"`
	CSS    []string `json:"css"`
	Image  []string `json:"image"`
	Dom    []string `json:"dom"`
	Video  []string `json:"video"`
}

// 获取当前页面的所有本站资源 js、css、img、链接
func (s SiteService) GetAllResources(rawURL string) ([]RequestParams, []string) {
	var RequestParamsAll []RequestParams
	var RouterAll []string
	var PastHref map[string]struct{} = make(map[string]struct{})
	var ActiveHref map[string]struct{} = make(map[string]struct{})
	// 解析域名
	parsed, err := url.Parse(rawURL)
	if err != nil {
		log.Println(err)
	}

	var Browser libs.Browser
	host := parsed.Hostname()
	if parsed.Port() != "" {
		host += parsed.Port()
		s.Port = parsed.Port()
	}
	obj := Browser.GetBrowser(host, &libs.Fingerprint{})
	s.BaseDomain = parsed.Hostname()
	ActiveHref["https://"+host] = struct{}{}
	// 进入主站
	browser := obj.Browser
	page := obj.Page
	s.LoopGet(page, PastHref, ActiveHref, &RequestParamsAll)
	for k := range PastHref {
		RouterAll = append(RouterAll, k)
	}
	list := deduplicationRequest(RequestParamsAll)
	browser.Close()
	return list, RouterAll
}

// 循环获取当前的资源及a标签中的链接
func (s SiteService) LoopGet(page *rod.Page, pastHref map[string]struct{}, activeHref map[string]struct{}, requestParamsAll *[]RequestParams) {
	// 取出
	for len(activeHref) > 0 {
		var url string
		for k := range activeHref {
			url = k
			pastHref[k] = struct{}{}
			delete(activeHref, k)
			break
		}
		if url == "" {
			return
		}
		log.Println("当前页面：", url)
		page.MustNavigate(url)
		// 提取当前页面的所有资源
		rs := getApiResponse(page)
		// 将当前资源加到全局资源中
		*requestParamsAll = append(*requestParamsAll, rs...)
		// 获取页面 HTML 提取当前页面的所有a标签
		html, err := page.HTML()
		if err != nil {
			panic(err)
		}
		aList := s.GetSrcByHtml(html)
		for _, v := range aList {
			if _, seen := pastHref[v]; !seen {
				activeHref[v] = struct{}{}
			}
		}
	}

}

func getApiResponse(page *rod.Page) []RequestParams {
	// 定义一个通道接收消息
	var RequestParamsList []RequestParams
	wait := page.EachEvent(func(e *proto.PageLoadEventFired) {
		RequestParamsList = []RequestParams{}
	}, func(e *proto.NetworkResponseReceived) {
		RequestParamsList = append(RequestParamsList, RequestParams{
			Type: e.Type,
			URL:  e.Response.URL,
		})
	})
	go wait()
	page.MustWaitLoad()
	return RequestParamsList
}

func (s SiteService) GetSrcByHtml(html string) []string {
	urlList := []string{}
	seen := make(map[string]struct{}) // 用 map 去重

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists || href == "" {
			return
		}

		// 过滤掉锚点和伪链接
		if strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript:") || href == "/" {
			return
		}

		// 解析 URL
		u, err := url.Parse(href)
		if err != nil {
			return
		}

		// 相对路径 → 拼接为绝对路径
		if !u.IsAbs() {
			if strings.HasPrefix(href, "/") {
				href = "https://" + s.BaseDomain + u.Path
			} else if strings.HasPrefix(href, "./") {
				result := strings.TrimPrefix(href, ".")
				href = "https://" + s.BaseDomain + result
			} else {
				href = "https://" + s.BaseDomain + "/" + u.Path
			}

		} else {
			// 只保留站内链接
			if !strings.Contains(u.Host, s.BaseDomain) {
				return
			}
		}

		// 去重
		if _, ok := seen[href]; !ok {
			seen[href] = struct{}{}
			urlList = append(urlList, href)
		}
	})

	return urlList
}

// 去重资源
func deduplicationRequest(list []RequestParams) []RequestParams {
	var newList []RequestParams
	for _, v := range list {
		is := true
		for _, l := range newList {
			if l.URL == v.URL && l.Type == v.Type {
				is = false
			}
		}
		if is {
			newList = append(newList, v)
		}
	}

	return newList
}

func (s SiteService) DeduplicationRequestByUrl(list []RequestParams) []RequestParams {
	var newList []RequestParams
	for _, v := range list {
		is := true
		u, _ := url.Parse(v.URL)
		u.RawQuery = ""
		for _, l := range newList {
			ur, _ := url.Parse(l.URL)
			ur.RawQuery = ""
			if u.String() == ur.String() && l.Type == v.Type {
				is = false
			}
		}
		if is {
			newList = append(newList, v)
		}
	}

	return newList
}
