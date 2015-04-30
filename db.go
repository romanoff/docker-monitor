package main

import (
	"github.com/boltdb/bolt"
)

var db *bolt.DB

func WriteDockerfileSha(name string, sha string) error {
	return db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("docker-monitor"))
		if err != nil {
			return err
		}
		return b.Put([]byte(name), []byte(sha))
	})
}

func ReadDockerfileSha(name string) (string, error) {
	sha := ""
	err := db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("docker-monitor"))
		if err != nil {
			return err
		}
		value := b.Get([]byte(name))
		if value != nil {
			sha = string(value)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return sha, err
}
