package utils

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
)

func LoadImageFromDirs(fileName, processedDir, uploadDir string) (image.Image, string, error) {
	processedPath := filepath.Join(processedDir, fileName)
	uploadPath := filepath.Join(uploadDir, fileName)
	img, format, err := LoadImage(processedPath)
	if err == nil {
		log.Printf("Файл найден в processedDir: %s, формат: %s", processedPath, format)
		return img, format, nil
	}
	img, format, err = LoadImage(uploadPath)
	if err == nil {
		log.Printf("Файл найден в uploadDir: %s, формат: %s", uploadPath, format)
		return img, format, nil
	}
	log.Printf("Файл не найден в директориях: %s, %s", processedPath, uploadPath)
	return nil, "", errors.New("file not found in both processedDir and uploadDir")
}

func LoadImage(filePath string) (image.Image, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

func ResponseImage(w http.ResponseWriter, img image.Image, format string) error {
	var err error
	switch format {
	case "jpeg", "jpg":
		w.Header().Set("Content-Type", "image/jpeg")
		err = jpeg.Encode(w, img, nil)
	case "png":
		w.Header().Set("Content-Type", "image/png")
		err = png.Encode(w, img)
	case "gif":
		w.Header().Set("Content-Type", "image/gif")
		err = gif.Encode(w, img, nil)
	case "webp":
		w.Header().Set("Content-Type", "image/webp")
		err = webp.Encode(w, img, nil)
	default:
		log.Printf("Неизвестный формат изображения: %s. Используется JPEG по умолчанию.", format)
		w.Header().Set("Content-Type", "image/jpeg")
		err = jpeg.Encode(w, img, nil)
	}
	if err != nil {
		log.Printf("Ошибка при кодировании изображения: %v. Попытка сохранить в формате JPEG.", err)
		err = jpeg.Encode(w, img, nil)
		w.Header().Set("Content-Type", "image/jpeg")
	}
	if err != nil {
		log.Printf("Ошибка при кодировании изображения в формате JPEG: %v", err)
		http.Error(w, "Не удалось закодировать изображение", http.StatusInternalServerError)
		return err
	}
	return nil
}

func SaveImageToFile(w *os.File, img image.Image, format string) error {
	var err error
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(w, img, nil)
	case "png":
		err = png.Encode(w, img)
	case "gif":
		err = gif.Encode(w, img, nil)
	case "webp":
		err = webp.Encode(w, img, nil)
	default:
		log.Printf("Неизвестный формат изображения: %s. Используется JPEG по умолчанию.", format)
		err = jpeg.Encode(w, img, nil)
	}
	if err != nil {
		log.Printf("Ошибка при кодировании изображения: %v. Попытка сохранить в формате JPEG.", err)
		err = jpeg.Encode(w, img, nil)
	}
	if err != nil {
		log.Printf("Ошибка при кодировании изображения в формате JPEG: %v", err)
		return err
	}
	return nil
}
