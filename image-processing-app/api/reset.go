package api

import (
	"fmt"
	"image-processing-app/api/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на сброс изображения до изначального состояния")
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
		return
	}
	oldFilePath := filepath.Join(processedDir, fileName)
	if _, err := os.Stat(oldFilePath); err == nil {
		if err := os.Remove(oldFilePath); err != nil {
			log.Printf("Ошибка при удалении файла из processed: %v", err)
			http.Error(w, "Failed to delete file", http.StatusInternalServerError)
			return
		}
		log.Printf("Файл удален из processed: %s", oldFilePath)
	}
	uploadFilePath, format, err := findOriginalFile(fileName, uploadDir)
	if err != nil {
		log.Printf("Ошибка при поиске исходного файла: %v", err)
		http.Error(w, "Original file not found", http.StatusNotFound)
		return
	}
	newProcessedFilePath := filepath.Join(processedDir, fileName)
	if err := copyFile(uploadFilePath, newProcessedFilePath); err != nil {
		log.Printf("Ошибка при копировании файла: %v", err)
		http.Error(w, "Failed to reset file", http.StatusInternalServerError)
		return
	}
	img, _, err := utils.LoadImage(uploadFilePath)
	if err != nil {
		log.Printf("Ошибка при загрузке изображения: %v", err)
		http.Error(w, "Failed to load image", http.StatusInternalServerError)
		return
	}
	utils.ResponseImage(w, img, format)
	log.Printf("Файл успешно сброшен: %s", newProcessedFilePath)
}

// findOriginalFile ищет файл в указанной директории по имени
func findOriginalFile(baseName, dir string) (string, string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", "", err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), strings.TrimSuffix(baseName, filepath.Ext(baseName))) {
			fullPath := filepath.Join(dir, file.Name())
			ext := filepath.Ext(file.Name())[1:] // Расширение без точки
			return fullPath, ext, nil
		}
	}
	return "", "", fmt.Errorf("file not found")
}

// copyFile копирует файл из src в dst
func copyFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}
