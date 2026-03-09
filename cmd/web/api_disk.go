package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileInfo содержит информацию о файле или папке
type FileInfo struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	IsDir        bool   `json:"is_dir"`
	Size         int64  `json:"size"`
	SizeReadable string `json:"size_readable"`
	CreatedAt    string `json:"created_at"`
}

// formatSize преобразует размер в байтах в удобочитаемый формат (КБ, МБ, ГБ)
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGT"[exp])
}

// getFileCreationTime получает дату создания файла (кроссплатформенно)
func getFileCreationTime(info os.FileInfo) string {
	created := info.ModTime()
	return created.Format("2006-01-02 15:04:05")
}

// validateAndNormalizePath проверяет и нормализует путь
func validateAndNormalizePath(inputPath string) (string, error) {
	// Удаляем начальные и конечные пробелы
	inputPath = strings.TrimSpace(inputPath)

	// Проверяем, что путь не пустой
	if inputPath == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Нормализуем путь — убирает лишние слеши, обрабатывает .. и .
	normalizedPath := filepath.Clean(inputPath)

	// Проверяем безопасность пути — предотвращаем обход директорий через ..
	if strings.Contains(normalizedPath, "..") {
		// Разрешаем только если .. находится в начале пути (относительные пути)
		if !strings.HasPrefix(normalizedPath, "..") {
			return "", fmt.Errorf("invalid path — directory traversal detected")
		}
	}

	return normalizedPath, nil
}

// getDirectoryContents получает содержимое директории
func getDirectoryContents(directoryPath string) ([]FileInfo, error) {
	// Проверяем существование директории
	info, err := os.Stat(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory does not exist: %s", directoryPath)
		}
		return nil, fmt.Errorf("failed to access directory: %v", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", directoryPath)
	}

	// Читаем содержимое директории
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var result []FileInfo

	for _, file := range files {
		filePath := filepath.Join(directoryPath, file.Name())
		fileInfo, err := file.Info()
		if err != nil {
			continue // Пропускаем файлы, которые не удалось прочитать
		}

		result = append(result, FileInfo{
			Name:         file.Name(),
			Path:         filePath,
			IsDir:        file.IsDir(),
			Size:         fileInfo.Size(),
			SizeReadable: formatSize(fileInfo.Size()),
			CreatedAt:    getFileCreationTime(fileInfo),
		})
	}

	return result, nil
}

// handler для обработки HTTP‑запросов
func listDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр path из URL
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path parameter is required", http.StatusBadRequest)
		return
	}

	// Нормализуем и проверяем путь
	normalizedPath, err := validateAndNormalizePath(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid path: %v", err), http.StatusBadRequest)
		return
	}

	// Получаем содержимое директории
	files, err := getDirectoryContents(normalizedPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content‑Type
	w.Header().Set("Content-Type", "application/json")

	// Кодируем результат в JSON и отправляем клиенту
	jsonData, err := json.Marshal(files)
	if err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
