package main

import (
	ascii "ascii-go/internal"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

// TODO
// error handling
// tests/int tests

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5656"
	}

	baseDir := os.Getenv("BASEDIR")
	if baseDir == "" {
		baseDir = "."
	}

	// relative to root/where you call `go run` from
	fs := http.FileServer(http.Dir(path.Join(baseDir, "static")))

	http.Handle("GET /static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("GET /", serveTemplate)

	http.HandleFunc("POST /api/convert", convertImage)

	log.Print("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var allowedTypes = map[string]string{
	"image/jpg":  ".jpg",
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

func readImage(contentType string, reader *bufio.Reader) (image.Image, error) {
	switch contentType {
	case ".png":
		return png.Decode(reader)
	case ".jpg":
		return jpeg.Decode(reader)
	}

	return nil, errors.New("Invalid file type")
}

func convertImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10mb
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	file, fileHeader, err := r.FormFile("img")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Missing image field", http.StatusBadRequest)
		return
	}
	fmt.Println(fileHeader.Size, fileHeader.Filename, fileHeader.Header)

	defer file.Close()

	// check actual file type
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType((buf[:n]))

	ext, ok := allowedTypes[contentType]
	if !ok {
		http.Error(w, fmt.Sprintf("Invalid file type %s", contentType), http.StatusBadRequest)
		return
	}

	// rewind file to actually read it
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, "Couldn't process file", http.StatusInternalServerError)
		return
	}

	reader := bufio.NewReader(file)
	img, err := readImage(ext, reader)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	converted := ascii.Convert(img)
	writeJson(w, 200, converted)
}

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// https://www.alexedwards.net/blog/serving-static-sites-with-go
// TODO parse on startup
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	baseDir := os.Getenv("BASEDIR")
	if baseDir == "" {
		baseDir = "."
	}

	files := []string{
		path.Join(baseDir, "templates/layout.html"),
		path.Join(baseDir, "templates/index.html"),
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
