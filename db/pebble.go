package db

import (
	"log"

	"github.com/cockroachdb/pebble"
)

func GetDB(path string) *pebble.DB {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CloseDB(db *pebble.DB) {
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func Get(k []byte, db *pebble.DB) []byte {
	value, closer, err := db.Get(k)
	if err != nil {
		return nil
	}
	if err := closer.Close(); err != nil {
		log.Fatal(err)
	}
	return value
}

func Set(k1, v1 []byte, db *pebble.DB) {
	if err := db.Set(k1, v1, pebble.Sync); err != nil {
		log.Fatal(err)
	}
}
