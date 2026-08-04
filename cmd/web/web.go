package main

import (
	ascii "ascii-go/internal"
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"image/png"
	"log"
	"net/http"
)

// TODO
// validation
// error handling
// styling
// tests/int tests

func main() {
	// relative to root/where you call `go run` from
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("GET /", serveTemplate)

	http.HandleFunc("POST /api/convert", convertImage)

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func convertImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10mb
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	file, fileHeader, err := r.FormFile("img")
	fmt.Println(fileHeader.Size, fileHeader.Filename, fileHeader.Header)

	reader := bufio.NewReader(file)
	img, err := png.Decode(reader)
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
	files := []string{
		"./templates/layout.html",
		"./templates/index.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
