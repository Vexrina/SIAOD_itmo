package extendible_hashing

import (
	"fmt"
	"hash/fnv"
)

type Bucket struct {
	data       map[string]any
	localDepth int
	maxSize    int
	fileName   string
}

type ExtendableHash struct {
	globalDepth int
	buckets     []*Bucket
	fileSystem  bool
}

func hashFunc(key string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return h.Sum32()
}

func newBucket(depth, maxSize, bucketNumber int, fileSystem bool) *Bucket {
	if fileSystem {
		fileName := fmt.Sprintf("./buckets/bucket_%d.dat", bucketNumber)
		b := &Bucket{
			localDepth: depth,
			maxSize:    maxSize,
			fileName:   fileName,
		}
		if b.writeToFile("", "") != nil {
			panic("cannot write empty values to file")
		}
		return b
	}
	return &Bucket{
		data:       make(map[string]any),
		localDepth: depth,
		maxSize:    maxSize,
	}
}

func NewExtendableHash(globalDepth, bucketSize int, fileSystem bool) *ExtendableHash {
	if globalDepth <= 0 {
		panic("globalDepth must be greater than zero")
	}
	buckets := make([]*Bucket, 1<<globalDepth)
	for i := range buckets {
		buckets[i] = newBucket(1, bucketSize, i, fileSystem)
	}
	return &ExtendableHash{
		globalDepth: globalDepth,
		buckets:     buckets,
		fileSystem:  fileSystem,
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
	if eh.fileSystem {
		eh.insertFile(key, value)
		return
	}
	eh.insertNotFileSystem(key, value)
}

// Индексация ключа в директории
func (eh *ExtendableHash) getIndex(key string) int {
	mask := (1 << eh.globalDepth) - 1
	return int(hashFunc(key)) & mask
}

// Поиск элемента
func (eh *ExtendableHash) GetByKey(key string) (any, bool) {
	if eh.fileSystem {
		return eh.getByKeyFile(key)
	}
	return eh.getByKeyNotFile(key)
}

func (eh *ExtendableHash) GetAllKeys() []string {
	if eh.fileSystem {
		return eh.getAllKeysFileSystem()
	}
	return eh.getAllKeysNotFile()
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
