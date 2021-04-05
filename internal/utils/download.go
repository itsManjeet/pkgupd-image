package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/dustin/go-humanize"
)

// https://golangcode.com/download-a-file-with-progress/

type WriteCounter struct {
	Total uint64
}

func (wc WriteCounter) printProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 32))
	fmt.Printf("\rDownloading.... %s complete", humanize.Bytes(wc.Total))
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.printProgress()
	return n, nil
}

func DownloadFile(filepath, url string) error {
	fmt.Println("downloading ", url, filepath)
	if _, err := os.Stat(filepath); err == nil {
		log.Println("Skipping, already in cache")
		return nil
	}
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	counter := &WriteCounter{}
	if _, err := io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	fmt.Print("\n")

	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}

	return nil
}

func DownloadSources(urls []string, outpath string) error {
	for _, url := range urls {
		log.Println("Downloading source from", url)
		filename := path.Base(url)
		if err := DownloadFile(outpath+"/"+filename, url); err != nil {
			log.Println("Failed to download", url, err)
			return err
		}
	}

	return nil
}
