package api

import (
	"image-processing-app/api/utils"
	"image/color"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
)

func RotateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на поворот изображения")
	fileName := r.URL.Query().Get("file")
	rotateStr := r.URL.Query().Get("rotate")
	if fileName == "" || rotateStr == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}
	rotate, err := strconv.Atoi(rotateStr)
	if err != nil {
		http.Error(w, "Invalid rotation angle", http.StatusBadRequest)
		return
	}
	img, format, err := utils.LoadImageFromDirs(fileName, processedDir, uploadDir)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	rotatedImg := imaging.Rotate(img, float64(rotate), color.Transparent)
	outputPath := filepath.Join(processedDir, fileName)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Ошибка при сохранении форматированного изображения: %v", err)
		http.Error(w, "Failed to save formated image", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()
	if err := utils.SaveImageToFile(outputFile, rotatedImg, format); err != nil {
		log.Printf("Ошибка при сохранении изображения: %v", err)
		http.Error(w, "Failed to save rotated image", http.StatusInternalServerError)
		return
	}
	log.Printf("Повернутое изображение сохранено: %s", outputPath)
	utils.ResponseImage(w, rotatedImg, format)
}
