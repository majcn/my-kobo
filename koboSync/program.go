package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// curl -L -o ca-bundle.crt https://curl.haxx.se/ca/cacert.pem
	os.Setenv("SSL_CERT_FILE", "/mnt/onboard/.adds/certs/ca-bundle.crt")

	project := "my-kobo"
	branch := "master"
	url := fmt.Sprintf("https://github.com/majcn/%s/archive/refs/heads/%s.zip", project, branch)

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	zipReader, _ := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	alterFPath := func(p string)(r string) {
		rootPathInZip := fmt.Sprintf("%s-%s", project, branch)
		re := regexp.MustCompile(fmt.Sprintf(`^(%s/?)(sdcard/?)?`, rootPathInZip))
		return re.ReplaceAllString(p, "/")
	}

	xx, err := Unzip(zipReader, "", alterFPath)
	if err != nil {
		println("ERROR: " + err.Error())
	}
	println("majcn")
	for _, x := range xx {
		println(x)
	}
}

// Unzip https://golangcode.com/unzip-files-in-go/
func Unzip(zipReader *zip.Reader, dest string, alterFPath func(string)string) ([]string, error) {

	var filenames []string

	for _, f := range zipReader.File {

		// Store filename/path for returning and using later on
		fpath := alterFPath(filepath.Join(dest, f.Name))

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		// if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			// return filenames, fmt.Errorf("%s: illegal file path", fpath)
		// }

		if !strings.HasPrefix(fpath, `/mnt/onboard/.adds/`) {
			continue
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
