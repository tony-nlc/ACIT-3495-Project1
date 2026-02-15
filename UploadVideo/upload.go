// upload-service/main.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse the form
    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("video")
    if err != nil {
        http.Error(w, "Missing 'video' field in form", http.StatusBadRequest)
        return
    }
    defer file.Close()

    fmt.Printf("Uploading: %s\n", handler.Filename)
    
    os.MkdirAll("./temp", os.ModePerm)
    
    dst, err := os.Create("./temp/" + handler.Filename)
    if err != nil {
        http.Error(w, "Save error", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    io.Copy(dst, file)
    w.Write([]byte("Upload successful"))
}

func main() {
    http.HandleFunc("/upload", uploadHandler)
    http.ListenAndServe(":8080", nil)
}