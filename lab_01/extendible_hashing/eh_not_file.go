package extendible_hashing

func (eh *ExtendableHash) insertNotFileSystem(key, value string) {
	index := eh.getIndex(key)
	bucket := eh.buckets[index]

	if len(bucket.data) < bucket.maxSize {
		bucket.data[key] = value
		return
	}

	// Если ведро заполнено, происходит разделение
	eh.splitBucketNotFileSystem(index)
	eh.Insert(key, value)
}

// делим бакет на два
func (eh *ExtendableHash) splitBucketNotFileSystem(index int) {
	bucket := eh.buckets[index]
	bucket.localDepth++

	// если хотим поделить бакет на количество, превышающее кол-во директорий
	if bucket.localDepth > eh.globalDepth {
		eh.doubleDirectory()
	}

	newBucket := newBucket(bucket.localDepth, bucket.maxSize, eh.GetNumDirs(), eh.fileSystem)

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

// Поиск элемента
func (eh *ExtendableHash) getByKeyNotFile(key string) (any, bool) {
	index := eh.getIndex(key)
	bucket := eh.buckets[index]
	val, exists := bucket.data[key]
	return val, exists
}

func (eh *ExtendableHash) getAllKeysNotFile() []string {
	keys := make([]string, 0, len(eh.buckets)*eh.buckets[0].maxSize)
	for _, bucket := range eh.buckets {
		for k := range bucket.data {
			keys = append(keys, k)
		}
	}
	return keys
}
