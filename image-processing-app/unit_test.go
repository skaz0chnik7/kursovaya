package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	handlers "image-processing-app/api"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const (
	testUploadDir    = "./assets/uploads"
	testProcessedDir = "./assets/processed"
	testImageFile    = "./api/test/test_image.jpg"
	testResizeWidth  = 100
	testResizeHeight = 100
)

func TestUploadHandler(t *testing.T) {
	file, err := os.Open(testImageFile)
	if err != nil {
		t.Fatalf("Не удалось открыть тестовое изображение: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(testImageFile))
	if err != nil {
		t.Fatalf("Не удалось создать форму: %v", err)
	}
	io.Copy(part, file)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.UploadHandler).ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", resp.Code)
	}

	expectedPath := filepath.Join(testUploadDir, filepath.Base(testImageFile))
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Файл не был создан в директории uploads: %v", expectedPath)
	}
}

func TestResizeHandler(t *testing.T) {
	uploadedFilePath := filepath.Join(testUploadDir, filepath.Base(testImageFile))

	req := httptest.NewRequest("GET", fmt.Sprintf("/resize?file=%s&width=%d&height=%d",
		filepath.Base(uploadedFilePath), testResizeWidth, testResizeHeight), nil)
	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.ResizeHandler).ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", resp.Code)
	}

	originalProcessedPath := filepath.Join(testProcessedDir, filepath.Base(testImageFile))
	resizedProcessedPath := filepath.Join(testProcessedDir, strings.TrimSuffix(filepath.Base(testImageFile), filepath.Ext(testImageFile))+"_resized.jpg")

	if err := os.Rename(originalProcessedPath, resizedProcessedPath); err != nil {
		t.Fatalf("Не удалось переименовать обработанный файл: %v", err)
	}

	if _, err := os.Stat(resizedProcessedPath); os.IsNotExist(err) {
		t.Errorf("Файл не был переименован в директории processed: %v", resizedProcessedPath)
	}

	file, err := os.Open(resizedProcessedPath)
	if err != nil {
		t.Fatalf("Не удалось открыть измененное изображение: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("Не удалось декодировать измененное изображение: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != testResizeWidth || bounds.Dy() != testResizeHeight {
		t.Errorf("Размеры изображения не совпадают с ожидаемыми: %dx%d, получили %dx%d",
			testResizeWidth, testResizeHeight, bounds.Dx(), bounds.Dy())
	}
}

func TestRotateHandler(t *testing.T) {
	angles := []int{90, 180, 270}
	for _, angle := range angles {
		t.Run(fmt.Sprintf("Rotate%dDegrees", angle), func(t *testing.T) {
			uploadedFilePath := filepath.Join(testUploadDir, filepath.Base(testImageFile))

			req := httptest.NewRequest("GET", fmt.Sprintf("/rotate?file=%s&rotate=%d",
				filepath.Base(uploadedFilePath), angle), nil)
			resp := httptest.NewRecorder()

			http.HandlerFunc(handlers.RotateHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("Ожидался код 200, получен %d", resp.Code)
			}

			originalRotatedPath := filepath.Join(testProcessedDir, filepath.Base(testImageFile))
			rotatedFilePath := filepath.Join(testProcessedDir, strings.TrimSuffix(filepath.Base(testImageFile), filepath.Ext(testImageFile))+"_rotated_"+strconv.Itoa(angle)+".jpg")

			if err := os.Rename(originalRotatedPath, rotatedFilePath); err != nil {
				t.Fatalf("Не удалось переименовать обработанный файл: %v", err)
			}

			if _, err := os.Stat(rotatedFilePath); os.IsNotExist(err) {
				t.Errorf("Повернутое изображение не найдено: %v", rotatedFilePath)
			}

			file, err := os.Open(rotatedFilePath)
			if err != nil {
				t.Fatalf("Не удалось открыть повернутое изображение: %v", err)
			}
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				t.Fatalf("Не удалось декодировать повернутое изображение: %v", err)
			}

			if img.Bounds().Empty() {
				t.Errorf("Повернутое изображение имеет пустые границы")
			}
		})
	}
}

func TestCropHandler(t *testing.T) {
	uploadedFilePath := filepath.Join(testUploadDir, filepath.Base(testImageFile))

	x, y := 10, 10
	width, height := 50, 50

	req := httptest.NewRequest("GET",
		fmt.Sprintf("/crop?file=%s&x=%d&y=%d&width=%d&height=%d",
			filepath.Base(uploadedFilePath), x, y, width, height), nil)
	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.CropHandler).ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", resp.Code)
	}

	originalProcessedPath := filepath.Join(testProcessedDir, filepath.Base(testImageFile))
	croppedProcessedPath := filepath.Join(testProcessedDir, strings.TrimSuffix(filepath.Base(testImageFile), filepath.Ext(testImageFile))+"_cropped.jpg")

	if err := os.Rename(originalProcessedPath, croppedProcessedPath); err != nil {
		t.Fatalf("Не удалось переименовать обработанный файл: %v", err)
	}

	if _, err := os.Stat(croppedProcessedPath); os.IsNotExist(err) {
		t.Errorf("Обрезанное изображение не найдено: %v", croppedProcessedPath)
	}

	file, err := os.Open(croppedProcessedPath)
	if err != nil {
		t.Fatalf("Не удалось открыть обрезанное изображение: %v", err)
	}
	defer file.Close()
}

func TestFilterHandler(t *testing.T) {
	filters := []string{"grayscale", "blur", "invert"} // Список фильтров для тестирования
	for _, filter := range filters {
		t.Run(fmt.Sprintf("Apply%sFilter", strings.Title(filter)), func(t *testing.T) {
			uploadedFilePath := filepath.Join(testUploadDir, filepath.Base(testImageFile))

			req := httptest.NewRequest("GET", fmt.Sprintf("/filter?file=%s&filter=%s",
				filepath.Base(uploadedFilePath), filter), nil)
			resp := httptest.NewRecorder()

			http.HandlerFunc(handlers.FilterHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("Ожидался код 200, получен %d", resp.Code)
			}

			originalFilteredPath := filepath.Join(testProcessedDir, filepath.Base(testImageFile))
			filteredFilePath := filepath.Join(testProcessedDir, strings.TrimSuffix(filepath.Base(testImageFile), filepath.Ext(testImageFile))+"_"+filter+".jpg")

			if err := os.Rename(originalFilteredPath, filteredFilePath); err != nil {
				t.Fatalf("Не удалось переименовать обработанный файл: %v", err)
			}

			if _, err := os.Stat(filteredFilePath); os.IsNotExist(err) {
				t.Errorf("Фильтрованное изображение не найдено: %v", filteredFilePath)
			}

			file, err := os.Open(filteredFilePath)
			if err != nil {
				t.Fatalf("Не удалось открыть фильтрованное изображение: %v", err)
			}
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				t.Fatalf("Не удалось декодировать фильтрованное изображение: %v", err)
			}

			// Дополнительная проверка: убедиться, что файл декодируется корректно
			if img.Bounds().Empty() {
				t.Errorf("Фильтрованное изображение имеет пустые границы")
			}
		})
	}
}

func TestInfoHandler(t *testing.T) {
	uploadedFilePath := filepath.Join(testUploadDir, filepath.Base(testImageFile))

	req := httptest.NewRequest("GET", fmt.Sprintf("/info?file=%s",
		filepath.Base(uploadedFilePath)), nil)
	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.InfoHandler).ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Ожидался код 200, получен %d", resp.Code)
	}

	var info map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &info); err != nil {
		t.Fatalf("Не удалось декодировать JSON ответ: %v", err)
	}

	if info["file_name"] != filepath.Base(testImageFile) {
		t.Errorf("Неверное имя файла: ожидается %v, получено %v", filepath.Base(testImageFile), info["file_name"])
	}

	if _, ok := info["size_bytes"].(float64); !ok {
		t.Errorf("Размер файла отсутствует или имеет неверный формат")
	}

	if _, ok := info["width"].(float64); !ok {
		t.Errorf("Ширина изображения отсутствует или имеет неверный формат")
	}

	if _, ok := info["height"].(float64); !ok {
		t.Errorf("Высота изображения отсутствует или имеет неверный формат")
	}

	if _, ok := info["format"].(string); !ok {
		t.Errorf("Формат изображения отсутствует или имеет неверный формат")
	}
}

func TestConvertHandler(t *testing.T) {
	formats := []string{"png", "gif", "webp"}
	for _, format := range formats {
		t.Run(fmt.Sprintf("ConvertTo%s", strings.ToUpper(format)), func(t *testing.T) {
			uploadedFilePath := filepath.Join(testUploadDir, filepath.Base(testImageFile))

			req := httptest.NewRequest("GET", fmt.Sprintf("/convert?file=%s&format=%s",
				filepath.Base(uploadedFilePath), format), nil)
			resp := httptest.NewRecorder()

			http.HandlerFunc(handlers.ConvertHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("Ожидался код 200, получен %d", resp.Code)
			}
		})
	}
}

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
