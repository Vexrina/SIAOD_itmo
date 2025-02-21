package extendible_hashing

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (b *Bucket) writeToFile(key, value string) error {
	// Проверяем, не пустой ли путь
	if b.fileName == "" {
		return fmt.Errorf("fileName is empty")
	}
	dir := filepath.Dir(b.fileName)

	// Проверяем, не является ли путь корневым
	if dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create dir %s: %w", dir, err)
		}
	}

	file, err := os.OpenFile(b.fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", b.fileName, err)
	}
	defer file.Close()

	if key != "" && value != "" {
		_, err = file.WriteString(fmt.Sprintf("%s:%s\n", key, value))
	}
	return err
}

func (b *Bucket) readFromFile() (map[string]string, error) {
	data := make(map[string]string)

	file, err := os.Open(b.fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (eh *ExtendableHash) insertFile(key, value string) {
	index := eh.getIndex(key)
	bucket := eh.buckets[index]

	data, err := bucket.readFromFile()
	if err != nil {
		panic(err)
	}

	if len(data) < bucket.maxSize {
		err := bucket.writeToFile(key, value)
		if err != nil {
			panic(err)
		}
		return
	}

	// Если ведро заполнено, происходит разделение
	eh.splitBucketNotFileSystem(index)
	eh.Insert(key, value)
}

func (eh *ExtendableHash) getByKeyFileSystem(key string) (any, bool) {
	index := eh.getIndex(key)
	bucket := eh.buckets[index]

	data, err := bucket.readFromFile()
	if err != nil {
		panic(err)
	}

	val, exists := data[key]
	return val, exists
}

func (eh *ExtendableHash) getAllKeysFileSystem() []string {
	keys := make([]string, 0, len(eh.buckets)*eh.buckets[0].maxSize)
	for _, bucket := range eh.buckets {
		data, err := bucket.readFromFile()
		if err != nil {
			panic(err)
		}
		for k := range data {
			keys = append(keys, k)
		}
	}
	return keys
}

func (eh *ExtendableHash) splitBucket(index int) {
	bucket := eh.buckets[index]
	bucket.localDepth++

	// Если хотим поделить бакет на количество, превышающее кол-во директорий
	if bucket.localDepth > eh.globalDepth {
		eh.doubleDirectory()
	}

	newBucket := newBucket(bucket.localDepth, bucket.maxSize, eh.GetNumDirs(), eh.fileSystem)

	// Читаем данные из старого бакета
	data, err := bucket.readFromFile()
	if err != nil {
		panic(err)
	}

	// Перераспределяем ключи между старым и новым бакетом
	for k, v := range data {
		if eh.getIndex(k) != index {
			err := newBucket.writeToFile(k, v)
			if err != nil {
				panic(err)
			}
			delete(data, k)
		}
	}

	// Записываем оставшиеся данные обратно в старый бакет
	err = bucket.rewriteFile(data)
	if err != nil {
		panic(err)
	}

	// Удаляем файл старого бакета, если он больше не нужен
	err = bucket.deleteFile()
	if err != nil {
		panic(err)
	}

	// Обновляем ссылки на бакеты
	for i := 0; i < len(eh.buckets); i++ {
		if eh.buckets[i] == bucket && (i>>uint(bucket.localDepth-1))&1 == 1 {
			eh.buckets[i] = newBucket
		}
	}
}

func (b *Bucket) rewriteFile(data map[string]string) error {
	file, err := os.Create(b.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for k, v := range data {
		_, err := file.WriteString(fmt.Sprintf("%s:%s\n", k, v))
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bucket) deleteFile() error {
	return os.Remove(b.fileName)
}

func (eh *ExtendableHash) getByKeyFile(key string) (any, bool) {
	index := eh.getIndex(key)
	bucket := eh.buckets[index]

	data, err := bucket.readFromFile()
	if err != nil {
		panic(err)
	}

	val, exists := data[key]
	return val, exists
}

func (eh *ExtendableHash) getAllKeysFile() []string {
	keys := make([]string, 0, len(eh.buckets)*eh.buckets[0].maxSize)
	for _, bucket := range eh.buckets {
		data, err := bucket.readFromFile()
		if err != nil {
			panic(err)
		}
		for k := range data {
			keys = append(keys, k)
		}
	}
	return keys
}
