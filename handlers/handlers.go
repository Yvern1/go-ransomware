package handlers

import (
	"fmt"
	"log"
	"io"
	"time"
	"bytes"
	"math/rand"
	"net/http"
	"encoding/json"
)

var KeyStore = make(map[string]string)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}
	return string(b)
}

func StoreKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var keyData map[string]string

		err = json.Unmarshal(bodyBytes, &keyData)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		for k, v := range keyData {
			KeyStore[k] = v
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Keys stored successfully.")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func RetrieveKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && len(r.URL.Query()) >  0 {
		key := r.URL.Query().Get("key")
		value, ok := KeyStore[key]
		if !ok {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, value)
	} else {
		http.Error(w, "Invalid request method or missing key query parameter", http.StatusBadRequest)
	}
}

func SendKey(id string, key string) {
	// Define the key and value you want to send to the server
	keyData := map[string]string{
        id: key,
    }

	// Convert the key data to JSON
	jsonData, err := json.Marshal(keyData)
	if err != nil {
		log.Fatalf("Error marshaling key data: %v", err)
	}

	// Create a new HTTP request with the JSON data in the body
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/store", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	log.Println("Key stored successfully.")
}

func GetKey(id string) ([]byte, error){
	
	url := fmt.Sprintf("http://localhost:8080/retrieve?key=%s", id)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	data := []byte(fmt.Sprintf("Value for key '%s': %s\n", id, string(bodyBytes)))

	return data, nil
}