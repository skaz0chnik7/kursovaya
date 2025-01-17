package api

import (
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"image-processing-app/api/utils"
)

func ResizeImage(img image.Image, newWidth, newHeight int) image.Image {
	srcBounds := img.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	xScale := float64(srcW) / float64(newWidth)
	yScale := float64(srcH) / float64(newHeight)

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) * xScale)
			srcY := int(float64(y) * yScale)
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
}

func ResizeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на изменение размера изображения")
	fileName := r.URL.Query().Get("file")
	widthStr := r.URL.Query().Get("width")
	heightStr := r.URL.Query().Get("height")

	if fileName == "" || widthStr == "" || heightStr == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}
	width, err1 := strconv.Atoi(widthStr)
	height, err2 := strconv.Atoi(heightStr)
	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid width or height", http.StatusBadRequest)
		return
	}
	img, format, err := utils.LoadImageFromDirs(fileName, processedDir, uploadDir)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	resizedImg := ResizeImage(img, width, height)
	outputPath := filepath.Join(processedDir, fileName)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Ошибка при сохранении форматированного изображения: %v", err)
		http.Error(w, "Failed to save formated image", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()
	if err := utils.SaveImageToFile(outputFile, resizedImg, format); err != nil {
		log.Printf("Ошибка при сохранении изображения: %v", err)
		http.Error(w, "Failed to save resized image", http.StatusInternalServerError)
		return
	}
	log.Printf("Изображение с измененным размером сохранено: %s", outputPath)
	utils.ResponseImage(w, resizedImg, format)
}
