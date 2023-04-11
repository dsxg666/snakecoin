package db

import (
	"log"

	"github.com/cockroachdb/pebble"
)

func GetDB(path string) *pebble.DB {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		log.Panic("Failed to Open:", err)
	}
	return db
}

func CloseDB(db *pebble.DB) {
	if err := db.Close(); err != nil {
		log.Panic("Failed to Close:", err)
	}
}

func Get(k []byte, db *pebble.DB) []byte {
	value, closer, err := db.Get(k)
	if err != nil {
		return nil
	}
	if err := closer.Close(); err != nil {
		log.Panic("Failed to Close:", err)
	}
	return value
}

func Set(k, v []byte, db *pebble.DB) {
	if err := db.Set(k, v, pebble.Sync); err != nil {
		log.Panic("Failed to Set:", err)
	}
}
