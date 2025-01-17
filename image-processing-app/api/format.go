package api

import (
	"image"
	"image-processing-app/api/utils"
	"image/color"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на применение фильтра к изображению")
	fileName := r.URL.Query().Get("file")
	filter := r.URL.Query().Get("filter")
	if fileName == "" {
		http.Error(w, "Missing file parameter", http.StatusBadRequest)
		return
	}
	img, format, err := utils.LoadImageFromDirs(fileName, processedDir, uploadDir)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	var filteredImg image.Image
	switch strings.ToLower(filter) {
	case "grayscale":
		filteredImg = ConvertToGrayscale(img)
	case "blur":
		filteredImg = imaging.Blur(img, 5)
	case "invert":
		filteredImg = imaging.Invert(img)
	default:
		filteredImg = ConvertToGrayscale(img)
	}
	outputPath := filepath.Join(processedDir, fileName)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Ошибка при сохранении форматированного изображения: %v", err)
		http.Error(w, "Failed to save formated image", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()
	if err := utils.SaveImageToFile(outputFile, filteredImg, format); err != nil {
		log.Printf("Ошибка при сохранении изображения: %v", err)
		http.Error(w, "Failed to save formated image", http.StatusInternalServerError)
		return
	}
	log.Printf("Форматированное изображение сохранено: %s", outputPath)
	utils.ResponseImage(w, filteredImg, format)
}

// ConvertToGrayscale преобразует изображение в указанный тип оттенков серого.
func ConvertToGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, _ := originalColor.RGBA()
			gray := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			grayImg.Set(x, y, color.Gray{Y: gray})
		}
	}
	return grayImg
}
