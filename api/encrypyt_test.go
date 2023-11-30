package api

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/nftime/logger"
	"github.com/nftime/util"
	"io"
	"strings"
	"testing"
)

func TestDecrypt(t *testing.T) {
	prefix := "nftime"

	//redisCli := redis.RedisManager()
	//keyString := fmt.Sprintf("%s:%s", "qr", "aes_key")
	//aesHexKey, err := redisCli.Get(context.Background(), keyString).Result()
	//if err != nil {
	//	logger.Error.Printf("err:%v\n", err)
	//	aesHexKey = "d4e9a43db1f266980f3cca85f00ca47606b9d92b82e8fac632c5b45090eeb5b4"
	//}
	aesHexKey := "d4e9a43db1f266980f3cca85f00ca47606b9d92b82e8fac632c5b45090eeb5b4"

	aesKey, err := hex.DecodeString(aesHexKey)
	if err != nil {
		logger.Error.Printf("Error decoding key from hexadecimal: %v", err)
		return
	}

	decodedData, err := base64.RawURLEncoding.DecodeString("UtgKQHKuav7V3+y39xuGGpeFuJ8UeKmdvtq9dlZCtic")
	if err != nil {
		fmt.Printf("Error decoding data: %v", err)
		fmt.Printf("decodedData: %v\n", decodedData)
		return
	}
	logger.Info.Printf("Decoded data: %v\n", string(decodedData))

	decryptedData, err := util.Decrypt(decodedData, aesKey)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	decryptedStr := string(decryptedData)
	if strings.HasPrefix(decryptedStr, prefix) {
		fmt.Printf("%v\n", util.RemoveLetters(decryptedStr))
	} else {
		return
	}
}
func TestEncrypt(t *testing.T) {
	aesKeyHex := "d4e9a43db1f266980f3cca85f00ca47606b9d92b82e8fac632c5b45090eeb5b4"

	// Convert the hexadecimal AES key to bytes
	aesKey, err := hexStringToBytes(aesKeyHex)
	if err != nil {
		fmt.Println("Error converting AES key:", err)
		return
	}

	fmt.Printf("aesKey: %s\n", aesKey)
	// Data to be encrypted
	dataToEncrypt := []byte("nftime1") // Replace this with your data

	// Encrypt the data
	encryptedData, err := util.Encrypt(dataToEncrypt, aesKey)
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return
	}

	fmt.Printf("%v\n", string(encryptedData))
	// Base64 encode the encrypted data
	base64EncodedData := base64.StdEncoding.EncodeToString(encryptedData)

	fmt.Printf("Base64 Encoded and AES Encrypted Data: %v\n", base64EncodedData)

	prefix := "nftime"

	decodedData, err := base64.RawURLEncoding.DecodeString(base64EncodedData)
	if err != nil {
		fmt.Printf("Error decoding data: %v", err)
		fmt.Printf("decodedData: %v\n", decodedData)
		return
	}
	logger.Info.Printf("Decoded data: %v\n", string(decodedData))

	decryptedData, err := util.Decrypt(decodedData, aesKey)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	decryptedStr := string(decryptedData)
	if strings.HasPrefix(decryptedStr, prefix) {
		fmt.Printf("%v\n", util.RemoveLetters(decryptedStr))
	} else {
		return
	}
}

func hexStringToBytes(hexString string) ([]byte, error) {
	return hex.DecodeString(hexString)
}

func TestBase64(t *testing.T) {
	data := []byte("Hello, World!")

	// Encode the data to base64
	encodedData := base64.StdEncoding.EncodeToString(data)

	fmt.Println("Base64 Encoded Data:", encodedData)
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}

	fmt.Println("Decoded Data:", string(decodedData))
}

func TestEncode(t *testing.T) {
	aesKeyHex := "d4e9a43db1f266980f3cca85f00ca47606b9d92b82e8fac632c5b45090eeb5b4"

	// Convert the hexadecimal AES key to bytes
	aesKey, err := hexStringToBytes(aesKeyHex)
	if err != nil {
		fmt.Println("Error converting AES key:", err)
		return
	}

	// Data to be encrypted
	dataToEncrypt := "nftime4" // Replace this with your data

	// Encrypt the data
	encryptedData, err := encryptAES([]byte(dataToEncrypt), aesKey)
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return
	}

	// Base64 encode the encrypted data
	base64EncodedData := base64.RawURLEncoding.EncodeToString(encryptedData)

	fmt.Println("Base64 Encoded and AES Encrypted Data:", base64EncodedData)

	// Decrypt the data
	decodedData, err := base64.RawURLEncoding.DecodeString(base64EncodedData)
	if err != nil {
		fmt.Println("Error decoding data:", err)
		return
	}

	decryptedData, err := decryptAES(decodedData, aesKey)
	if err != nil {
		fmt.Println("Error decrypting data:", err)
		return
	}

	fmt.Println("Decrypted Data:", string(decryptedData))
}

// Encrypt data using AES with a given key
func encryptAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Pad the data to the block size
	padding := aes.BlockSize - (len(data) % aes.BlockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	data = append(data, padText...)

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// Decrypt data using AES with a given key
func decryptAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	// Unpad the data
	padding := int(data[len(data)-1])
	if padding < 1 || padding > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding")
	}

	return data[:len(data)-padding], nil
}
