package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("Decryption Key")
	keyEntry := widget.NewPasswordEntry()
	keyEntry.PlaceHolder = "Enter Decryption Key"

	submitButton := widget.NewButton("Submit", func() {
		key := keyEntry.Text
		if len(key) > 0 {
			// Get the path for the folders containing the files to be decrypted:
			usr, _ := user.Current()
			downloads := filepath.Join(usr.HomeDir, "Downloads")
			documents := filepath.Join(usr.HomeDir, "Documents")
			foldersPath := []string{downloads, documents}

			// Decrypt every file in each folder
			for _, folderPath := range foldersPath {
				files, _ := ioutil.ReadDir(folderPath)
				for _, file := range files {
					filePath := filepath.Join(folderPath, file.Name())
					if strings.HasSuffix(file.Name(), ".aes") {
						ciphertext, _ := ioutil.ReadFile(filePath)
						keyDecoded, err := base64.StdEncoding.DecodeString(key)
						block, err := aes.NewCipher(keyDecoded)
						if err != nil {
							dialog.ShowError(err, w)
							return
						}
						if len(ciphertext) < aes.BlockSize {
							dialog.ShowError(fmt.Errorf("ciphertext too short"), w)
							return
						}
						iv := ciphertext[:aes.BlockSize]
						ciphertext = ciphertext[aes.BlockSize:]
						stream := cipher.NewCFBDecrypter(block, iv)
						stream.XORKeyStream(ciphertext, ciphertext)

						// Move the decrypted file
						destinationPath := filepath.Join(folderPath, "decrypted_"+file.Name()[:len(file.Name())-4])
						err = ioutil.WriteFile(destinationPath, ciphertext, 0644)
						if err != nil {
							dialog.ShowError(err, w)
							return
						}

						// Delete the encrypted file
						os.Remove(filePath)
					}
				}
			}
			fmt.Println("Decryption complete")
			w.Hide()
		}
	})

	w.SetContent(widget.NewVBox(
		keyEntry,
		submitButton,
	))
	w.ShowAndRun()
}
