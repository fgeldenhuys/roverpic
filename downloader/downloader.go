package downloader

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Download struct {
	URL      string
	Filename string
}

type DownloadResult struct {
	Download
	Error error
}

type Config struct {
	DownloadPath           string
	MaxConcurrentDownloads int
}

type Downloader struct {
	Config
	limit chan struct{}
}

func Init(conf Config) *Downloader {
	return &Downloader{
		Config: conf,
		limit:  make(chan struct{}, conf.MaxConcurrentDownloads),
	}
}

func (dl *Downloader) Download(ctx context.Context, downloads []Download) (results []DownloadResult) {
	completed := make(chan DownloadResult)
	// this outer goroutine can be blocked waiting for the download limit
	go func() {
		for _, download := range downloads {
			if err := ctx.Err(); err != nil {
				// context has been cancelled
				break
			}
			download := download
			// throttle the downloader by blocking on the limiter semaphore, or stop if cancelled
			select {
			case <-ctx.Done():
				// context has been cancelled, don't start more downloads
			case dl.limit <- struct{}{}:
				// start a separate goroutine for each download
				go func() {
					err := dl.performDownload(download)
					<-dl.limit
					completed <- DownloadResult{
						Download: download,
						Error:    err,
					}
				}()
			}
		}
	}()
	// wait for all downloads to finish, and collect errors
	for range downloads {
		select {
		case <-ctx.Done():
			log.Printf("Cancelled download of %d files", len(downloads))
			// there might be some downloads still in progress, but we don't care about the results
			// to cancel them, we would need to implement an io reader with context in performDownload
			// question: should we delete the downloaded files when request is cancelled?
			return
		case result := <-completed:
			results = append(results, result)
		}
	}
	return
}

func (dl *Downloader) performDownload(download Download) error {
	log.Printf("Downloading %s", download.Filename)
	filename := filepath.Join(dl.DownloadPath, download.Filename)

	// create the download directory if it doesn't exist
	if _, err := os.Stat(filepath.Dir(filename)); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
			return err
		}
	}

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Failed to create download %s: %w", download.URL, err)
	}
	defer out.Close()

	resp, err := http.Get(download.URL)
	if err != nil {
		return fmt.Errorf("Failed to download %s: %w", download.URL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to download %s: %s", download.URL, resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("Failed while downloading %s: %w", download.URL, err)
	}

	return nil
}
