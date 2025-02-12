package extendible_hashing

import (
	"hash/fnv"
)

type Bucket struct {
	data       map[string]any
	localDepth int
	maxSize    int
}

type ExtendableHash struct {
	globalDepth int
	buckets     []*Bucket
}

func hashFunc(key string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return h.Sum32()
}

func newBucket(depth, maxSize int) *Bucket {
	return &Bucket{
		data:       make(map[string]any),
		localDepth: depth,
		maxSize:    maxSize,
	}
}

func NewExtendableHash(globalDepth, bucketSize int) *ExtendableHash {
	if globalDepth <= 0 {
		panic("globalDepth must be greater than zero")
	}
	buckets := make([]*Bucket, 1<<globalDepth)
	for i := range buckets {
		buckets[i] = newBucket(1, bucketSize)
	}
	return &ExtendableHash{
		globalDepth: globalDepth,
		buckets:     buckets,
	}
}

func (eh *ExtendableHash) GetDepth() int {
	return eh.globalDepth
}

func (eh *ExtendableHash) GetNumDirs() int {
	return len(eh.buckets)
}

// Вставка элемента
func (eh *ExtendableHash) Insert(key, value string) {
	index := eh.getIndex(key)
	bucket := eh.buckets[index]

	if len(bucket.data) < bucket.maxSize {
		bucket.data[key] = value
		return
	}

	// Если ведро заполнено, происходит разделение
	eh.splitBucket(index)
	eh.Insert(key, value)
}

// делим бакет на два
func (eh *ExtendableHash) splitBucket(index int) {
	bucket := eh.buckets[index]
	bucket.localDepth++

	// если хотим поделить бакет на количество, превышающее кол-во директорий
	if bucket.localDepth > eh.globalDepth {
		eh.doubleDirectory()
	}

	newBucket := newBucket(bucket.localDepth, bucket.maxSize)

	// перераскидываем ключики
	for k, v := range bucket.data {
		if eh.getIndex(k) != index {
			newBucket.data[k] = v
			delete(bucket.data, k)
		}
	}

	// апдейтим ссылки на бакеты
	for i := 0; i < len(eh.buckets); i++ {
		if eh.buckets[i] == bucket && (i>>uint(bucket.localDepth-1))&1 == 1 {
			eh.buckets[i] = newBucket
		}
	}
}

// даблит директорию
func (eh *ExtendableHash) doubleDirectory() {
	// увеличиваем глубину на единицу (помним, что кол-во дир = 2^глобал_дептх
	eh.globalDepth++
	// создаем бакеты
	bucketsCount := len(eh.buckets) * 2
	newBuckets := make([]*Bucket, bucketsCount)

	for i := range eh.buckets {
		// в старой половине на том же месте
		newBuckets[2*i] = eh.buckets[i]
		// в новой половине на том же месте
		newBuckets[2*i+1] = eh.buckets[i]
	}
	eh.buckets = newBuckets
}

// Индексация ключа в директории
func (eh *ExtendableHash) getIndex(key string) int {
	mask := (1 << eh.globalDepth) - 1
	return int(hashFunc(key)) & mask
}

// Поиск элемента
func (eh *ExtendableHash) GetByKey(key string) (any, bool) {
	index := eh.getIndex(key)
	bucket := eh.buckets[index]
	val, exists := bucket.data[key]
	return val, exists
}

func (eh *ExtendableHash) GetAllKeys() []string {
	keys := make([]string, 0, len(eh.buckets)*eh.buckets[0].maxSize)
	for _, bucket := range eh.buckets {
		for k := range bucket.data {
			keys = append(keys, k)
		}
	}
	return keys
}
