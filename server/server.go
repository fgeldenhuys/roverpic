package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"roverpic/downloader"
	"roverpic/roverapi"
)

type Config struct {
	Port int
}

func (conf Config) ListenAddress() string {
	if conf.Port == 0 {
		return ":80"
	} else {
		return fmt.Sprintf(":%d", conf.Port)
	}
}

type Server struct {
	Config
	roverAPI *roverapi.RoverAPI
	dl       *downloader.Downloader
}

func Init(conf Config, roverAPI *roverapi.RoverAPI, dl *downloader.Downloader) *Server {
	return &Server{
		Config:   conf,
		roverAPI: roverAPI,
		dl:       dl,
	}
}

type downloadHttpResult struct {
	ApiSuccess bool     `json:"api_success"`
	Downloaded int      `json:"downloaded,omitempty"`
	Errors     []string `json:"errors,omitempty"`
}

// helper function to write json to the http response
func writeJson(w http.ResponseWriter, code int, out interface{}) {
	bytes, err := json.Marshal(out)
	if err != nil {
		log.Printf("Failed to write json response: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(bytes)
	if err != nil {
		// can't change the header after writing, just log error
		log.Printf("Failed to write json response: %s", err)
		return
	}
}

func (server *Server) download(w http.ResponseWriter, req *http.Request) {
	dateString := req.FormValue("date")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		writeJson(w, http.StatusBadRequest, downloadHttpResult{
			ApiSuccess: false,
			Errors:     []string{fmt.Sprintf("Unable to parse date: %s", err)},
		})
		return
	}

	// get info from the api about photos for that day
	photos, err := server.roverAPI.GetPhotos(date)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, downloadHttpResult{
			ApiSuccess: false,
			Errors:     []string{err.Error()},
		})
		return
	}

	// question: should we perhaps check for already downloaded photos?
	//           it is unspecified, so for now we just redownload everything

	log.Printf("Downloading %d photos for %s", len(photos), date)

	// download the photo files pointed to by the previous reply
	var downloads []downloader.Download
	for _, photo := range photos {
		split := strings.Split(photo.ImgSrc, "/")
		filename := filepath.Join(date.Format("2006-01-02"), split[len(split)-1])
		downloads = append(downloads, downloader.Download{
			URL:      photo.ImgSrc,
			Filename: filename,
		})
	}
	results := server.dl.Download(req.Context(), downloads)
	var errStrings []string
	for _, result := range results {
		if result.Error != nil {
			errStrings = append(errStrings, result.Error.Error())
		}
	}

	writeJson(w, http.StatusOK, downloadHttpResult{
		ApiSuccess: true,
		Downloaded: len(results) - len(errStrings),
		Errors:     errStrings,
	})
}

func (server *Server) ListenAndServe() {
	http.HandleFunc("/download", server.download)

	log.Printf("Starting server on %s", server.ListenAddress())
	http.ListenAndServe(server.ListenAddress(), nil)
}
