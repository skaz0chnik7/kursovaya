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

func crop(img image.Image, rect image.Rectangle) image.Image {
	if !rect.Overlaps(img.Bounds()) {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}
	intersect := rect.Intersect(img.Bounds())
	dst := image.NewRGBA(image.Rect(0, 0, intersect.Dx(), intersect.Dy()))
	offsetX, offsetY := intersect.Min.X, intersect.Min.Y
	for y := 0; y < intersect.Dy(); y++ {
		for x := 0; x < intersect.Dx(); x++ {
			color := img.At(x+offsetX, y+offsetY)
			dst.Set(x, y, color)
		}
	}
	return dst
}

func CropHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос на обрезку изображения")
	fileName := r.URL.Query().Get("file")
	xStr := r.URL.Query().Get("x")
	yStr := r.URL.Query().Get("y")
	widthStr := r.URL.Query().Get("width")
	heightStr := r.URL.Query().Get("height")

	if fileName == "" || xStr == "" || yStr == "" || widthStr == "" || heightStr == "" {
		log.Println("Недостаточно параметров для обрезки изображения")
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}
	x, err1 := strconv.Atoi(xStr)
	y, err2 := strconv.Atoi(yStr)
	width, err3 := strconv.Atoi(widthStr)
	height, err4 := strconv.Atoi(heightStr)
	if (err1 != nil) || (err2 != nil) || (err3 != nil) || (err4 != nil) || width <= 0 || height <= 0 {
		log.Println("Неверные параметры обрезки")
		http.Error(w, "Invalid crop parameters", http.StatusBadRequest)
		return
	}
	img, format, err := utils.LoadImageFromDirs(fileName, processedDir, uploadDir)
	if err != nil {
		log.Printf("Ошибка при загрузке изображения: %v", err)
		http.Error(w, "Failed to load image", http.StatusInternalServerError)
		return
	}
	rect := image.Rect(x, y, x+width, y+height)
	cropped := crop(img, rect)
	outputPath := filepath.Join(processedDir, fileName)
	outputFile, err := os.Create(outputPath)
	defer outputFile.Close()
	log.Printf("Обрезанное изображение сохранено: %s", outputPath)
	utils.ResponseImage(w, cropped, format)
}
