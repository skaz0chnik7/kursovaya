package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const uploadDir = "./assets/uploads/"
const processedDir = "./assets/processed/"

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на загрузку файла")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileName := strings.ReplaceAll(header.Filename, " ", "_")
	filePath := filepath.Join(uploadDir, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Ошибка при сохранении файла: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Printf("Ошибка при записи файла: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	log.Printf("Файл успешно загружен: %s", filePath)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"fileName": fileName})
}
