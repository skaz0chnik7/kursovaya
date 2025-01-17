package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	handlers "image-processing-app/api"
)

const uploadDir = "./assets/uploads/"
const processedDir = "./assets/processed/"

func main() {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, 0755)
		if err != nil {
			log.Fatalf("Ошибка при создании директории загрузок: %v", err)
		}
	}

	if _, err := os.Stat(processedDir); os.IsNotExist(err) {
		err := os.Mkdir(processedDir, 0755)
		if err != nil {
			log.Fatalf("Ошибка при создании директории обработанных изображений: %v", err)
		}
	}

	http.Handle("/", http.FileServer(http.Dir("./client")))
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/reset", handlers.ResetHandler)
	http.HandleFunc("/resize", handlers.ResizeHandler)
	http.HandleFunc("/rotate", handlers.RotateHandler)
	http.HandleFunc("/crop", handlers.CropHandler)
	http.HandleFunc("/convert", handlers.ConvertHandler)
	http.HandleFunc("/filter", handlers.FilterHandler)
	http.HandleFunc("/info", handlers.InfoHandler)

	fmt.Println("Сервер запущен на http://localhost:5050")
	log.Fatal(http.ListenAndServe(":5050", nil))
}
