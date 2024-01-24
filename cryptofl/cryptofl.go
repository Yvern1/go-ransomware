package cryptofl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"bufio"
)

type File struct {
	os.FileInfo
	Extension string 
	Path      string 
}

func(file *File) Encrypt(key []byte, dst io.Writer) error {

	fmt.Println("encrypting", file.Path)
	
	// Open the file read only
	inFile, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer inFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCTR(block, iv)

	dst.Write(iv)

	bufWriter := bufio.NewWriter(dst)
	
	writer := &cipher.StreamWriter{S: stream, W: bufWriter}


	if _, err = io.Copy(writer, inFile); err != nil {
		return err
	}

	// Flush the buffered writer to ensure all data is written to the destination writer
	if err = bufWriter.Flush(); err != nil {
		return err
	}

	return nil
}


func(file *File) Decrypt(key []byte, dst io.Writer) error {
	// Open the encrypted file
	inFile, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer inFile.Close()
	fmt.Println("decrypting", file.Path)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	inFile.Read(iv)

	stream := cipher.NewCTR(block, iv)

	// Create a buffered reader for the input file
	reader := bufio.NewReader(inFile)

	decrptor := &cipher.StreamReader{S: stream, R: reader}

	// Copy the input file to the dst, decrypting as we go.
	if _, err = io.Copy(dst, decrptor); err != nil {
		return err
	}

	return nil
}


func ReplaceBy(path string, filename string) error {
	// Open the file
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	// Copy the reader to file
	if _, err = io.Copy(file, src); err != nil {
		return err
	}

	return nil
}
