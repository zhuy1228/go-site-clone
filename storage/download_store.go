package storage

import (
	"fmt"
	"time"
)

// DownloadRecord 下载记录
type DownloadRecord struct {
	ID          string    `json:"id"`           // 唯一标识
	URL         string    `json:"url"`          // 下载的 URL
	SiteName    string    `json:"site_name"`    // 站点名称
	Status      string    `json:"status"`       // 状态：success, failed, processing
	TotalFiles  int       `json:"total_files"`  // 总文件数
	Downloaded  int       `json:"downloaded"`   // 已下载数
	CSSCount    int       `json:"css_count"`    // CSS 文件数
	ScriptCount int       `json:"script_count"` // JS 文件数
	ImageCount  int       `json:"image_count"`  // 图片文件数
	VideoCount  int       `json:"video_count"`  // 视频文件数
	ErrorMsg    string    `json:"error_msg"`    // 错误信息
	StartTime   time.Time `json:"start_time"`   // 开始时间
	EndTime     time.Time `json:"end_time"`     // 结束时间
	Duration    int64     `json:"duration"`     // 耗时（秒）
}

// AddDownloadRecord 添加下载记录
func (s *Store) AddDownloadRecord(record DownloadRecord) error {
	if record.ID == "" {
		record.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if record.StartTime.IsZero() {
		record.StartTime = time.Now()
	}
	return s.Set(BucketDownload, record.ID, record)
}

// GetDownloadRecord 获取下载记录
func (s *Store) GetDownloadRecord(id string) (*DownloadRecord, error) {
	var record DownloadRecord
	if err := s.Get(BucketDownload, id, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

// UpdateDownloadRecord 更新下载记录
func (s *Store) UpdateDownloadRecord(record DownloadRecord) error {
	if record.ID == "" {
		return fmt.Errorf("记录 ID 不能为空")
	}
	return s.Set(BucketDownload, record.ID, record)
}

// DeleteDownloadRecord 删除下载记录
func (s *Store) DeleteDownloadRecord(id string) error {
	return s.Delete(BucketDownload, id)
}

// GetAllDownloadRecords 获取所有下载记录
func (s *Store) GetAllDownloadRecords() ([]DownloadRecord, error) {
	var records []DownloadRecord
	if err := s.GetAll(BucketDownload, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// GetDownloadRecordsBySite 根据站点名称获取下载记录
func (s *Store) GetDownloadRecordsBySite(siteName string) ([]DownloadRecord, error) {
	allRecords, err := s.GetAllDownloadRecords()
	if err != nil {
		return nil, err
	}

	var filtered []DownloadRecord
	for _, record := range allRecords {
		if record.SiteName == siteName {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

// GetDownloadRecordsByStatus 根据状态获取下载记录
func (s *Store) GetDownloadRecordsByStatus(status string) ([]DownloadRecord, error) {
	allRecords, err := s.GetAllDownloadRecords()
	if err != nil {
		return nil, err
	}

	var filtered []DownloadRecord
	for _, record := range allRecords {
		if record.Status == status {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

// GetRecentDownloadRecords 获取最近的 N 条下载记录
func (s *Store) GetRecentDownloadRecords(limit int) ([]DownloadRecord, error) {
	allRecords, err := s.GetAllDownloadRecords()
	if err != nil {
		return nil, err
	}

	// 按时间倒序排序
	for i := 0; i < len(allRecords)-1; i++ {
		for j := i + 1; j < len(allRecords); j++ {
			if allRecords[i].StartTime.Before(allRecords[j].StartTime) {
				allRecords[i], allRecords[j] = allRecords[j], allRecords[i]
			}
		}
	}

	if limit > 0 && limit < len(allRecords) {
		return allRecords[:limit], nil
	}
	return allRecords, nil
}

// GetDownloadStats 获取下载统计信息
func (s *Store) GetDownloadStats() (map[string]interface{}, error) {
	records, err := s.GetAllDownloadRecords()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total":       len(records),
		"success":     0,
		"failed":      0,
		"processing":  0,
		"total_files": 0,
	}

	for _, record := range records {
		switch record.Status {
		case "success":
			stats["success"] = stats["success"].(int) + 1
		case "failed":
			stats["failed"] = stats["failed"].(int) + 1
		case "processing":
			stats["processing"] = stats["processing"].(int) + 1
		}
		stats["total_files"] = stats["total_files"].(int) + record.TotalFiles
	}

	return stats, nil
}

// ClearOldDownloadRecords 清理指定天数之前的下载记录
func (s *Store) ClearOldDownloadRecords(days int) (int, error) {
	records, err := s.GetAllDownloadRecords()
	if err != nil {
		return 0, err
	}

	cutoffTime := time.Now().AddDate(0, 0, -days)
	deleted := 0

	for _, record := range records {
		if record.StartTime.Before(cutoffTime) {
			if err := s.DeleteDownloadRecord(record.ID); err == nil {
				deleted++
			}
		}
	}

	return deleted, nil
}
