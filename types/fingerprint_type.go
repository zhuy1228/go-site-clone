package types

type BrowserFingerprintParams struct {
	Canvas      bool   // 是否随机指纹
	TimeZone    string // 时区名称，如 "America/New_York"
	Language    string // 语言，如 "en-US"
	UserAgent   string // 用户代理字符串
	GeoLocation string // 地理位置坐标，如 "40.7128,-74.0060"
}
