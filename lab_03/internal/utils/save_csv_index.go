package utils

import (
	"encoding/csv"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
)

func ReadCsv(csvFileName string, index bleve.Index) error {
	file, err := openCSVFile(csvFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	r := csv.NewReader(file)

	//headers
	headers, err := r.Read()
	if err != nil {
		return err
	}

	//content
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	for idx, record := range records {
		if len(record) != len(headers) {
			continue
		}
		if idx%500 == 0 {
			log.Println(idx)
		}
		if err = index.Index(uuid.NewString(), record[1]); err != nil {
			log.Printf("Ошибка индексации документа: %v", err)
		}
	}
	return nil
}

func openCSVFile(csvFileName string) (*os.File, error) {
	// 1. Получаем текущую рабочую директорию
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("не удалось определить рабочую директорию: %v", err)
	}

	// 2. Строим относительный путь
	relPath := filepath.Join("csv", csvFileName)

	// 3. Преобразуем в абсолютный путь
	absPath := filepath.Join(wd, relPath)

	// 4. Нормализуем путь (убираем ../ и ./)
	cleanPath := filepath.Clean(absPath)

	log.Printf("Открываем файл: %s", cleanPath)

	// 5. Открываем файл
	return os.Open(cleanPath)
}
