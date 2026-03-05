package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Структура для представления информации о файле/директории
type FileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"is_dir"`
	Size    int64  `json:"size,omitempty"` // Будет опущено для директорий
	ModTime string `json:"mod_time"`
}

func main() {
	// Получаем имя директории из аргументов командной строки
	if len(os.Args) < 2 {
		fmt.Println("Использование: program <directory_path>")
		os.Exit(1)
	}

	dirPath := os.Args[1]

	// Проверяем существование директории
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Ошибка: директория '%s' не существует\n", dirPath)
		} else {
			fmt.Printf("Ошибка при проверке директории: %v\n", err)
		}
		os.Exit(1)
	}

	if !info.IsDir() {
		fmt.Printf("Ошибка: '%s' — это не директория\n", dirPath)
		os.Exit(1)
	}

	// Собираем информацию о файлах и директориях
	filesInfo, err := getFilesInfo(dirPath)
	if err != nil {
		fmt.Printf("Ошибка при чтении директории: %v\n", err)
		os.Exit(1)
	}

	// Конвертируем в JSON
	jsonData, err := json.MarshalIndent(filesInfo, "", "  ")
	if err != nil {
		fmt.Printf("Ошибка при создании JSON: %v\n", err)
		os.Exit(1)
	}

	// Выводим JSON
	fmt.Println(string(jsonData))
}

// Функция для получения информации о файлах в директории
func getFilesInfo(dirPath string) ([]FileInfo, error) {
	var result []FileInfo

	// Читаем содержимое директории
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		fileInfo := FileInfo{
			Name:    entry.Name(),
			Path:    entryPath,
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		}

		// Размер только для файлов
		if !entry.IsDir() {
			fileInfo.Size = info.Size()
		}

		result = append(result, fileInfo)
	}

	return result, nil
}
