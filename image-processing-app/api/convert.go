package api

import (
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"image-processing-app/api/utils"
)

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на конвертацию изображения")
	fileName, format := r.URL.Query().Get("file"), strings.ToLower(r.URL.Query().Get("format"))
	if fileName == "" || format == "" {
		http.Error(w, "Отсутствуют параметры file или format", http.StatusBadRequest)
		return
	}
	supportedFormats := map[string]bool{"jpeg": true, "jpg": true, "png": true, "gif": true, "webp": true}
	if !supportedFormats[format] {
		http.Error(w, "Неподдерживаемый формат", http.StatusBadRequest)
		return
	}
	img, _, err := utils.LoadImageFromDirs(fileName, processedDir, uploadDir)
	if err != nil {
		http.Error(w, "Не удалось загрузить изображение", http.StatusInternalServerError)
		return
	}
	ext := filepath.Ext(fileName)
	outputFileName := strings.TrimSuffix(fileName, ext) + "." + format
	outputPath := filepath.Join(processedDir, outputFileName)
	if err := saveImageToFile(outputPath, img, format); err != nil {
		http.Error(w, "Не удалось сохранить конвертированное изображение", http.StatusInternalServerError)
		return
	}
	oldFilePath := filepath.Join(processedDir, fileName)
	if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
		log.Printf("Ошибка при удалении старого файла: %v", err)
	} else {
		log.Printf("Старое изображение удалено: %s", oldFilePath)
	}
	log.Printf("Изображение сконвертировано и сохранено: %s", outputPath)
	if err := utils.ResponseImage(w, img, format); err != nil {
		log.Printf("Ошибка при отправке изображения клиенту: %v", err)
	}
}

func saveImageToFile(path string, img image.Image, format string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Printf("Ошибка при создании файла: %v", err)
		return err
	}
	defer file.Close()

	if err := utils.SaveImageToFile(file, img, format); err != nil {
		log.Printf("Ошибка при сохранении изображения в файл: %v", err)
		return err
	}

	return nil
}
