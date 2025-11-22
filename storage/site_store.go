package storage

import (
	"fmt"
	"go-site-clone/types"
	"time"
)

// SiteRecord 站点数据记录（带时间戳）
type SiteRecord struct {
	types.NginxSiteConfig
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddSite 添加站点配置
func (s *Store) AddSite(site types.NginxSiteConfig) error {
	record := SiteRecord{
		NginxSiteConfig: site,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	return s.Set(BucketSites, site.Name, record)
}

// GetSite 获取站点配置
func (s *Store) GetSite(siteName string) (*types.NginxSiteConfig, error) {
	var record SiteRecord
	if err := s.Get(BucketSites, siteName, &record); err != nil {
		return nil, err
	}
	return &record.NginxSiteConfig, nil
}

// UpdateSite 更新站点配置
func (s *Store) UpdateSite(site types.NginxSiteConfig) error {
	var record SiteRecord

	// 尝试获取现有记录以保留创建时间
	if err := s.Get(BucketSites, site.Name, &record); err == nil {
		record.NginxSiteConfig = site
		record.UpdatedAt = time.Now()
	} else {
		// 如果不存在，创建新记录
		record = SiteRecord{
			NginxSiteConfig: site,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
	}

	return s.Set(BucketSites, site.Name, record)
}

// DeleteSite 删除站点配置
func (s *Store) DeleteSite(siteName string) error {
	return s.Delete(BucketSites, siteName)
}

// GetAllSites 获取所有站点配置
func (s *Store) GetAllSites() ([]types.NginxSiteConfig, error) {
	var records []SiteRecord
	if err := s.GetAll(BucketSites, &records); err != nil {
		return nil, err
	}

	sites := make([]types.NginxSiteConfig, len(records))
	for i, record := range records {
		sites[i] = record.NginxSiteConfig
	}
	return sites, nil
}

// SiteExists 检查站点是否存在
func (s *Store) SiteExists(siteName string) bool {
	return s.Exists(BucketSites, siteName)
}

// GetSitesByStatus 根据启用状态获取站点
func (s *Store) GetSitesByStatus(enabled bool) ([]types.NginxSiteConfig, error) {
	allSites, err := s.GetAllSites()
	if err != nil {
		return nil, err
	}

	var filtered []types.NginxSiteConfig
	for _, site := range allSites {
		if site.Enabled == enabled {
			filtered = append(filtered, site)
		}
	}
	return filtered, nil
} // UpdateSiteStatus 更新站点启用状态
func (s *Store) UpdateSiteStatus(siteName string, enabled bool) error {
	site, err := s.GetSite(siteName)
	if err != nil {
		return fmt.Errorf("站点不存在: %v", err)
	}

	site.Enabled = enabled
	return s.UpdateSite(*site)
}

// GetSitesCount 获取站点总数
func (s *Store) GetSitesCount() int {
	return s.Count(BucketSites)
}
