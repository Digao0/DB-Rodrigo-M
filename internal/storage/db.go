package storage

import "fmt"

type DB struct {
	mem *MemTable
	wal *WAL
}

func NewDB(walPath string) (*DB, error) {
	wal, err := OpenWAL(walPath)
	if err != nil {
		return nil, fmt.Errorf("open wal: %w", err)
	}

	db := &DB{
		mem: NewMemTable(),
		wal: wal,
	}

	if err := db.recover(); err != nil {
		_ = wal.Close()
		return nil, err
	}

	return db, nil
}

func (db *DB) recover() error {
	records, err := db.wal.Load()
	if err != nil {
		return fmt.Errorf("load wal: %w", err)
	}

	for _, record := range records {
		switch record.Operation {
		case OpSet:
			db.mem.Set(record.Key, record.Value)
		case OpDelete:
			db.mem.Delete(record.Key)
		default:
			return fmt.Errorf("unknown operation: %s", record.Operation)
		}
	}

	return nil
}

func (db *DB) Set(key, value string) error {
	record := Record{
		Operation: OpSet,
		Key:       key,
		Value:     value,
	}

	if err := db.wal.Append(record); err != nil {
		return err
	}

	if err := db.wal.Sync(); err != nil {
		return err
	}

	db.mem.Set(key, value)
	return nil
}

func (db *DB) Get(key string) (string, bool) {
	return db.mem.Get(key)
}

func (db *DB) Delete(key string) error {
	record := Record{
		Operation: OpDelete,
		Key:       key,
	}

	if err := db.wal.Append(record); err != nil {
		return err
	}

	if err := db.wal.Sync(); err != nil {
		return err
	}

	db.mem.Delete(key)
	return nil
}

func (db *DB) Close() error {
	if db.wal == nil {
		return nil
	}

	return db.wal.Close()
}