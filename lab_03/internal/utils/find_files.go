package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FindMatchingFiles ищет файлы в dirName, чьи имена начинаются с префикса папки
func FindMatchingFiles(csvDir, bleveDir string) ([]string, []string, error) {
	var indexed, notIndexed []string

	// Получаем список CSV файлов
	csvFiles, err := os.ReadDir(csvDir)
	if err != nil {
		log.Println(err)
		return nil, nil, fmt.Errorf("error reading CSV directory: %v", err)
	}

	// Получаем список .bleve папок
	bleveItems, err := os.ReadDir(bleveDir)
	if err != nil {
		log.Println(err)
		return nil, nil, fmt.Errorf("error reading Bleve directory: %v", err)
	}

	// Создаем множество проиндексированных файлов
	indexedSet := make(map[string]bool)
	for _, item := range bleveItems {
		if item.IsDir() && filepath.Ext(item.Name()) == ".bleve" {
			// Извлекаем имя CSV файла из названия .bleve папки
			csvName := strings.TrimSuffix(item.Name(), ".bleve")
			indexedSet[csvName] = true
		}
	}

	// Классифицируем CSV файлы
	for _, file := range csvFiles {
		if file.IsDir() {
			continue
		}

		fullPath := filepath.Join(file.Name())
		if indexedSet[file.Name()] {
			indexed = append(indexed, fullPath)
		} else {
			notIndexed = append(notIndexed, fullPath)
		}
	}

	return indexed, notIndexed, nil
}
