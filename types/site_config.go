package types

// NginxSiteConfig 站点配置结构
type NginxSiteConfig struct {
	ID      string   `json:"id"`      // 站点唯一ID
	Name    string   `json:"name"`    // 站点名称
	Domains []string `json:"domains"` // 域名列表
	Port    int      `json:"port"`    // 监听端口
	Path    string   `json:"path"`    // 网站根目录路径
	Index   string   `json:"index"`   // 默认首页文件
	Enabled bool     `json:"enabled"` // 是否启用
}
