// Package store 负责数据的本地持久化存储
package store

import (
	"encoding/json"

	"github.com/RiverMint78/pone-quest/internal/pone"
	"go.etcd.io/bbolt"
)

const bucketName = "images"

type KVStore struct {
	db *bbolt.DB
}

// New 创建或打开一个 bbolt 数据库
func New(path string) (*KVStore, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	// 确保 bucket 存在
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})

	return &KVStore{db: db}, err
}

func (s *KVStore) Close() error {
	return s.db.Close()
}

// SaveImageItem 存入一条图片数据
func (s *KVStore) SaveImageItem(item pone.ImageItem) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		buf, err := json.Marshal(item)
		if err != nil {
			return err
		}
		return b.Put([]byte(item.ID), buf)
	})
}

// LoadAll 加载库中所有数据
func (s *KVStore) LoadAll() ([]pone.ImageItem, error) {
	var items []pone.ImageItem
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.ForEach(func(k, v []byte) error {
			var item pone.ImageItem
			if err := json.Unmarshal(v, &item); err == nil {
				items = append(items, item)
			}
			return nil
		})
	})
	return items, err
}
