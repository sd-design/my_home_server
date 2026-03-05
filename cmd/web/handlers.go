package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type FileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"is_dir"`
	Size    int64  `json:"size,omitempty"` // Будет опущено для директорий
	ModTime string `json:"mod_time"`
}

const RootDir = "D:\\DEV\\my_home_server\\"

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}
func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func getClientIp(w http.ResponseWriter, r *http.Request) {
	ipClient := r.RemoteAddr
	fmt.Fprintf(w, "Your Ip-adress %s", ipClient)
}

func readFolder(w http.ResponseWriter, r *http.Request) {
	folderName := r.URL.Query().Get("folder")

	dirPath := RootDir + folderName

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
	fmt.Fprintf(w, string(jsonData))
}

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
