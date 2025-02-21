package perfect_hash

import (
	"fmt"
	"hash/fnv"
	"vexrina/siaod_itmo/lab_01/utils"
)

type PerfectHash struct {
	table     []string
	keyValues []*utils.KeyValue
	size      int
}

/*
	hashFunc получает хэш от ключа

берет побайтово символы ключа, выполняет операции:

hash = hash_const (какая-то стандартная большая константа)

hash = hash XOR char_byte

hash *= prime_const (какая-то стандартная большая константа)

после чего sum32 суммирует все полученные таким образом hash'и
*/
func hashFunc(key string, size int) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32()) % size
}


func cleanCollision(size, outerKeyIdx int, keys []string, values []any)([]string, []*utils.KeyValue){
	newTable := make([]string, size)
	newKeyValues := make([]*utils.KeyValue, size)
	for innerKeyIdx, innerKey := range keys[:outerKeyIdx+1] {
		idx := hashFunc(innerKey, size)
		if newTable[idx] != "" { // в теории, не достижимо
			newTable, newKeyValues = cleanCollision(size*2, innerKeyIdx, keys, values)
		}
		newTable[idx] = innerKey
		newKeyValues[idx] = utils.NewKeyValue(innerKey, values[innerKeyIdx])
	}
	return newTable, newKeyValues
}
// NewPerfectHash создает хеш-таблицу без коллизий для фиксированного набора ключей
func NewPerfectHash(keys []string, values []any) *PerfectHash {
	size := len(keys) * 2 // Начальный размер таблицы
	table := make([]string, size)
	keyValues := make([]*utils.KeyValue, size)

	for outerKeyIdx, outerKey := range keys {
		index := hashFunc(outerKey, size)
		wasCollision := false
		// коллизия
		for table[index] != "" {
			// увеличиваем размер таблицы вдвое
			size *= 2
			table, keyValues= cleanCollision(size, outerKeyIdx, keys, values)
			wasCollision = true
		}

		if !wasCollision {
			table[index] = outerKey
			keyValues[index] = utils.NewKeyValue(outerKey, values[outerKeyIdx])
		}
	}

	return &PerfectHash{table: table, size: size, keyValues: keyValues}
}
// findIdx находит необходимый индекс по ключу
func (ph *PerfectHash) findIdx(key string) (bool, int) {
	index := hashFunc(key, ph.size)
	return ph.table[index] == key, index
}

// Lookup проверяет, есть ли необходимый ключ в таблице
func (ph *PerfectHash) Lookup(key string) bool {
	exist, _ := ph.findIdx(key)
	return exist
}

// GetValueByKey получает значение по ключу в таблице. Если ключа нет, то получаем ошибку
func (ph *PerfectHash) GetValueByKey(key string) (any, error) {
	var (
		idx   int
		exist bool
	)

	if exist, idx = ph.findIdx(key); !exist {
		return nil, fmt.Errorf("don't have a key \"%s\" in hashTable", key)
	}

	return ph.keyValues[idx].Value, nil
}

// GetAllKeys получает все ключи из таблицы
func (ph *PerfectHash) GetAllKeys() []string {
	var keys []string
	for _, key := range ph.table {
		if key != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

// GetAllValues получает все значения
func (ph *PerfectHash) GetAllValues() []any {
	var values []any
	for idx := range ph.table {
		if ph.table[idx] != "" {
			values = append(values, ph.keyValues[idx].Value)
		}
	}
	return values
}

// GetAllKeysValues получает все ключ-значения
func (ph *PerfectHash) GetAllKeysValues() []*utils.KeyValue {
	var values []*utils.KeyValue
	for idx := range ph.table {
		if ph.table[idx] != "" {
			values = append(values, ph.keyValues[idx])
		}
	}
	return values
}

// GetAllIndexes получает все индексы. Задумывалась как дебаг функция
func (ph *PerfectHash) GetAllIndexes() []int {
	var indexes []int
	for idx, key := range ph.table {
		if key != "" {
			indexes = append(indexes, idx)
		}
	}
	return indexes
}

// PutNewKeyValue кладет новую пару ключ-значение в таблицу
func (ph *PerfectHash) PutNewKeyValue(newKey string, newValue any) *PerfectHash {
	keys := append(ph.GetAllKeys(), newKey)
	values := append(ph.GetAllValues(), newValue)
	return NewPerfectHash(keys, values)
}
