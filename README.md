# File Encryption and Decryption

This code encrypts files in the Downloads and Documents folders on the user's local machine and sends the decryption key via email. It uses AES encryption with a randomly generated key.

## Requirements

- Go version 1.14 or higher

## Libraries used
- crypto/aes
- crypto/cipher
- crypto/rand
- encoding/base64
- encoding/json
- io
- io/ioutil
- net/http
- os
- os/user
- path/filepath
- strings

## Functionality
- The code gets the path for the Downloads and Documents folders on the user's local machine.
- It generates a random AES key of 32 bytes.
- The key is then sent via email to the address "paulkiragu621@gmail.com" with the subject "Decryption Key for [USERNAME]" using the SendGrid API.
- The code then encrypts every file in the Downloads and Documents Folders using the AES encryption and saves the encrypted files with a ".aes" file extension.

## Usage
- Replace the X-RapidAPI-Key and X-RapidAPI-Host in the headers variable with your own API key and host obtained from SendGrid.
- Replace the 'to' email address in the payload variable with the desired email address to receive the decryption key.
- Replace the 'from' email address in the payload variable with the desired email address to send the decryption key.
- Run the code by executing go run main.go in the terminal.
- The encrypted files can be decrypted using the key sent via email and the AES decryption process.

## Disclaimer: 

This code is for demonstration purposes only and should not be used in a production environment without modification and proper security measures.