package main

import (
	"archive/zip"
	"github.com/garage44/rand/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

var config util.Config

// Compile templates on start of the application
var templates = template.Must(template.ParseFiles("public/upload.html"))

// Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	if !ok {
		log.Error().Msg("Error parsing basic auth")
		w.WriteHeader(401)
		return
	}

	if u != config.Username {
		log.Error().Msgf("Username provided is incorrect: %s", u)
		w.WriteHeader(401)
		return
	}

	if strings.TrimSpace(p) != config.Password {
		log.Error().Msgf("Password provided is incorrect: %s", p)
		w.WriteHeader(401)
		return
	}

	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("distFile")
	zipTarget := path.Join(config.UploadPath, handler.Filename)

	if err != nil {
		log.Error().Err(err).Msg("Error Retrieving the File")
		return
	}

	defer file.Close()
	log.Debug().Msgf("Uploaded File: %+v", handler.Filename)
	log.Debug().Msgf("File Size: %+v", handler.Size)
	log.Debug().Msgf("MIME Header: %+v", handler.Header)

	// Create file
	dst, err := os.Create(zipTarget)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	unzip(zipTarget, config.UploadPath)
	os.Remove(zipTarget)
	log.Info().Msgf("Uploaded zip deployed: %s", config.UploadPath)
}

func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var err error
	config, err = util.LoadConfig(".")

	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot start service")
	}

	// Upload route
	http.HandleFunc("/upload", uploadHandler)

	//Listen on port 8080
	log.Info().Msgf("Starting rand: %s", config.ServerAddress)
	http.ListenAndServe(config.ServerAddress, nil)
}
