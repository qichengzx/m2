package rocksdb

import (
	"github.com/tecbot/gorocksdb"
)

type Storage struct {
	rdb *gorocksdb.DB
	wo  *gorocksdb.WriteOptions
	ro  *gorocksdb.ReadOptions
}

func New(dir string) *Storage {
	if dir == "" {
		dir = "data"
	}
	rdb := newRocksDB(dir)
	db := &Storage{rdb: rdb}
	db.wo = gorocksdb.NewDefaultWriteOptions()
	db.ro = gorocksdb.NewDefaultReadOptions()
	return db
}

func newRocksDB(dir string) *gorocksdb.DB {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	rdb, err := gorocksdb.OpenDb(opts, dir)
	if err != nil {
		panic(err)
	}

	return rdb
}

func (s *Storage) Set(key, value []byte) error {
	return s.rdb.Put(s.wo, key, value)
}

func (s *Storage) Get(key []byte) ([]byte, error) {
	return s.rdb.GetBytes(s.ro, key)
}

func (s *Storage) WriteBatch(batch *gorocksdb.WriteBatch) error {
	return s.rdb.Write(s.wo, batch)
}

func (s *Storage) Delete(key []byte) error {
	return s.rdb.Delete(s.wo, key)
}

func (s *Storage) Close() error {
	s.wo.Destroy()
	s.ro.Destroy()
	s.rdb.Close()

	return nil
}
