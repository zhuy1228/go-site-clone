package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Bucket 名称常量
var (
	BucketSites    = []byte("sites")
	BucketDownload = []byte("download_records")
	BucketLogs     = []byte("access_logs")
	BucketSettings = []byte("settings")
)

// Store BBolt 存储服务
type Store struct {
	db *bolt.DB
}

// NewStore 创建新的存储服务
func NewStore(dbPath string) (*Store, error) {
	// 确保数据库目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %v", err)
	}

	// 打开数据库
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %v", err)
	}

	// 创建必要的 buckets
	err = db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{BucketSites, BucketDownload, BucketLogs, BucketSettings}
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return fmt.Errorf("创建 bucket %s 失败: %v", string(bucket), err)
			}
		}
		return nil
	})

	if err != nil {
		db.Close()
		return nil, err
	}

	log.Printf("数据库初始化成功: %s", dbPath)

	return &Store{db: db}, nil
}

// Close 关闭数据库连接
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Set 存储键值对到指定 bucket
func (s *Store) Set(bucket []byte, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %v", err)
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s 不存在", string(bucket))
		}
		return b.Put([]byte(key), data)
	})
}

// Get 从指定 bucket 获取值
func (s *Store) Get(bucket []byte, key string, dest interface{}) error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s 不存在", string(bucket))
		}

		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("键 %s 不存在", key)
		}

		return json.Unmarshal(data, dest)
	})
}

// Delete 删除指定键
func (s *Store) Delete(bucket []byte, key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s 不存在", string(bucket))
		}
		return b.Delete([]byte(key))
	})
}

// GetAll 获取指定 bucket 的所有数据
func (s *Store) GetAll(bucket []byte, destSlice interface{}) error {
	var items []json.RawMessage

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s 不存在", string(bucket))
		}

		return b.ForEach(func(k, v []byte) error {
			items = append(items, json.RawMessage(v))
			return nil
		})
	})

	if err != nil {
		return err
	}

	// 将 items 转换为目标切片类型
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, destSlice)
}

// Exists 检查键是否存在
func (s *Store) Exists(bucket []byte, key string) bool {
	exists := false
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		exists = b.Get([]byte(key)) != nil
		return nil
	})
	return exists
}

// Count 统计指定 bucket 的键数量
func (s *Store) Count(bucket []byte) int {
	count := 0
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		b.ForEach(func(k, v []byte) error {
			count++
			return nil
		})
		return nil
	})
	return count
}

// Backup 备份数据库到指定路径
func (s *Store) Backup(backupPath string) error {
	// 确保备份目录存在
	dir := filepath.Dir(backupPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建备份目录失败: %v", err)
	}

	return s.db.View(func(tx *bolt.Tx) error {
		return tx.CopyFile(backupPath, 0600)
	})
}
