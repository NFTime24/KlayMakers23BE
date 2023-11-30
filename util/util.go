package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"time"
	"unicode"
)

func CreateUnixFilename(initialFileName string) string {
	timestamp := time.Now().Unix()
	timestampString := strconv.Itoa(int(timestamp))
	fileName := "%s_%s.png"
	fileName = fmt.Sprintf(fileName, initialFileName, timestampString)
	return fileName
}

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// PKCS#7 Padding
	paddingSize := aes.BlockSize - len(plaintext)%aes.BlockSize
	padding := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	plaintext = append(plaintext, padding...)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// Decrypt takes a ciphertext, decrypts it using AES with a key, and returns the plaintext.
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove PKCS#7 padding
	paddingSize := int(ciphertext[len(ciphertext)-1])
	if paddingSize < 1 || paddingSize > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding size")
	}

	return ciphertext[:len(ciphertext)-paddingSize], nil
}

func RemoveLetters(input string) string {
	var result []rune

	for _, r := range input {
		if unicode.IsDigit(r) {
			result = append(result, r)
		}
	}

	return string(result)
}
