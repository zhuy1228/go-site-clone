package types

// DownloadMode 下载模式
type DownloadMode string

const (
	// DownloadModeSameDomain 只下载同域名资源
	DownloadModeSameDomain DownloadMode = "same-domain"
	// DownloadModeAllResources 下载所有资源
	DownloadModeAllResources DownloadMode = "all-resources"
	// DownloadModeCustom 自定义下载规则
	DownloadModeCustom DownloadMode = "custom"
)

// DownloadOptions 下载配置选项
type DownloadOptions struct {
	// 下载模式
	Mode DownloadMode `json:"mode"`

	// 自定义域名列表（当 Mode = custom 时使用）
	CustomDomains []string `json:"customDomains"`

	// 是否跳过超大文件
	SkipLargeFiles bool `json:"skipLargeFiles"`

	// 最大文件大小（MB）
	MaxFileSize int `json:"maxFileSize"`

	// 是否下载外部CSS
	DownloadExternalCSS bool `json:"downloadExternalCSS"`

	// 是否下载外部JS
	DownloadExternalJS bool `json:"downloadExternalJS"`

	// 是否下载外部图片
	DownloadExternalImages bool `json:"downloadExternalImages"`

	// 是否下载外部视频
	DownloadExternalVideos bool `json:"downloadExternalVideos"`
}

// DefaultDownloadOptions 返回默认下载配置
func DefaultDownloadOptions() DownloadOptions {
	return DownloadOptions{
		Mode:                   DownloadModeSameDomain,
		CustomDomains:          []string{},
		SkipLargeFiles:         true,
		MaxFileSize:            10, // 10MB
		DownloadExternalCSS:    false,
		DownloadExternalJS:     false,
		DownloadExternalImages: false,
		DownloadExternalVideos: false,
	}
}

// ShouldDownloadExternal 判断是否应该下载外部资源
func (opt DownloadOptions) ShouldDownloadExternal(resourceHost, baseHost, resourceType string) bool {
	// 如果是同域名，始终下载
	if resourceHost == baseHost {
		return true
	}

	// 根据模式判断
	switch opt.Mode {
	case DownloadModeSameDomain:
		// 只下载同域名，外部资源不下载
		return false

	case DownloadModeAllResources:
		// 下载所有资源，但要检查类型过滤
		return opt.isResourceTypeEnabled(resourceType)

	case DownloadModeCustom:
		// 检查是否在自定义域名列表中
		for _, domain := range opt.CustomDomains {
			if resourceHost == domain {
				return opt.isResourceTypeEnabled(resourceType)
			}
		}
		return false

	default:
		return false
	}
}

// isResourceTypeEnabled 检查资源类型是否启用下载
func (opt DownloadOptions) isResourceTypeEnabled(resourceType string) bool {
	switch resourceType {
	case "css":
		return opt.DownloadExternalCSS
	case "script", "js":
		return opt.DownloadExternalJS
	case "image", "img":
		return opt.DownloadExternalImages
	case "video":
		return opt.DownloadExternalVideos
	default:
		return true
	}
}
