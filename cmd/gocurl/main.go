package main

/* gocurl - A tiny curl replacement BECAUSE :
   How often did you have to download curl or unzip during Docker builds?
   This won't become larger. Peace. jan@hacker.ch, 2023 */

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	client          http.Client
	unzipTo         = flag.String("unzipTo", "", "directory to unzip a .zip download to")
	followRedirects = flag.Bool("follow", false, "follow HTTP redirects")
	remoteURL       = flag.String("remote", "", "URL of file to download")
	localFilename   = flag.String("local", "", "local filename")
)

func fatal(err error) {
	// stuck in old go time. feel free to improve.
	if err != nil {
		fmt.Println("# gocurl 1.0 - usage:")
		fmt.Println("# gocurl -remote URL -local FILENAME [-follow -unzipTo DIRECTORY]")
		fmt.Fprintf(os.Stderr, "FATAL download error: %s\n", err.Error())
		os.Exit(1)
	}
}

func empty(g string /* g like generic soon */) bool {
	return g == ""
}

func main() {
	flag.Parse()

	if empty(*remoteURL) {
		fatal(errors.New("missing -remote URL to download file from"))
	}

	if empty(*localFilename) {
		fatal(errors.New("missing -local FILENAME for downloaded file"))
	}

	if _, err := os.Stat(*localFilename); err == nil {
		fatal(errors.New("local output file exists, cowardly aborting to not overwrite it"))
	}

	if _, err := os.Stat(*unzipTo); err == nil {
		fatal(errors.New("directory to -unzipTo exits, aborting like a gentleman"))
	}

	if *remoteURL == *localFilename {
		fatal(errors.New("i don't write tests, but this happend to me: -remote = -local"))
	}

	out, err := os.Create(*localFilename)
	fatal(err)
	defer out.Close()

	if *followRedirects {
		// thanks - https://dev.to/fuchiao/http-redirect-for-golang-http-client-2i35 🍝
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	fmt.Printf("Downloading %s to %s\n", *remoteURL, *localFilename)
	resp, err := client.Get(*remoteURL)
	fatal(err)
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	fatal(err)
	fmt.Printf("Download OK, %d bytes ..in n secs tbd.. curl awe logs\n", n)

	fatal(unzip(*localFilename, *unzipTo))
	fmt.Println("OK. Seems that download worked. Enjoy!")
}

func unzip(zipFile, targetDirectory string) error {
	if empty(targetDirectory) {
		return nil
	}
	fmt.Printf("Unzipping %s to %s\n", zipFile, targetDirectory)
	// thanks! blunt copy from: https://golang.cafe/blog/golang-unzip-file-example.html
	archive, err := zip.OpenReader(zipFile)
	fatal(err)
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(targetDirectory, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(targetDirectory)+string(os.PathSeparator)) {
			fatal(errors.New("invalid file path"))
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		fatal(os.MkdirAll(filepath.Dir(filePath), os.ModePerm))

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		fatal(err)

		fileInArchive, err := f.Open()
		fatal(err)

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			fatal(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}
