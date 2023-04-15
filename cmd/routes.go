package cmd

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// set file size to 10MB max
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error parsing form")
		return
	}

	// get name of file
	name := r.Form.Get("name")
	if name == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("specify a file name please")
		return
	}

	// get file
	f, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("something went wrong")
		return
	}
	defer f.Close()

	// get file extension
	fileExtension := strings.ToLower(filepath.Ext(handler.Filename))

	// create folders
	path := filepath.Join(".", "files")
	_ = os.MkdirAll(path, os.ModePerm)
	fullPath := path + "/" + name + fileExtension

	// open and copy files
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("something went wrong")
		return
	}
	defer file.Close()

	_, err = io.Copy(file, f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("something went wrong")
		return
	}

	// send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("File uploaded successfully")

}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get the file name to download from url
	name := r.URL.Query().Get("name")

	// join to get the full file path
	directory := filepath.Join("files", name)

	// open file (check if exists)
	ff, err := os.Open(directory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		json.NewEncoder(w).Encode("Unable to open file ")
		return
	}

	// get file extension
	fileExtension := strings.ToLower(filepath.Ext(ff.Name()))
	log.Println(fileExtension)

	// force a download with the content- disposition field
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(directory))

	// serve file out.
	http.ServeFile(w, r, directory)
}
