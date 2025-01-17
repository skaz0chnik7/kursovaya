package api

import (
	"encoding/json"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на получение информации об изображении")
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing file parameter", http.StatusBadRequest)
		return
	}
	inputPath := filepath.Join("./assets/uploads", fileName)
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Printf("Файл не найден: %s", inputPath)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	file, err := os.Open(inputPath)
	if err != nil {
		log.Printf("Ошибка при открытии файла: %s", inputPath)
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	img, format, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("Ошибка при декодировании изображения: %v", err)
		http.Error(w, "Failed to decode image", http.StatusInternalServerError)
		return
	}
	info := map[string]interface{}{
		"file_name":     fileName,
		"size_bytes":    fileInfo.Size(),
		"last_modified": fileInfo.ModTime().Format(time.RFC3339),
		"width":         img.Width,
		"height":        img.Height,
		"format":        format,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
