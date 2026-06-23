package storage

type MemTable struct {
	data map[string]string
}

func NewMemTable() *MemTable {
	return &MemTable{
		data: make(map[string]string),
	}
}

func (m *MemTable) Set(key, value string) {
	m.data[key] = value
}

func (m *MemTable) Get(key string) (string, bool) {
	value, ok := m.data[key]
	return value, ok
}

func (m *MemTable) Delete(key string) {
	delete(m.data, key)
}