package util

import (
	"context"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/nftime/logger"
	"github.com/nftime/redis"
	"image"
	"image/jpeg"
	"log"
	"os"
	"testing"
	"time"
)

func TestThumbnail(t *testing.T) {
	file, err := os.Open("../mos-sukjaroenkraisri-dxWlUPIsnJk-unsplash.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileName := CreateUnixFilename(file.Name())

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// Resize the image to the desired dimensions (e.g., 300 pixels wide)
	resizedImg := resize.Resize(300, 0, img, resize.Lanczos3)

	// Create the output file
	outfile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	defer os.Remove(fileName)
	// Encode the resized image and save it to the output file
	jpeg.Encode(outfile, resizedImg, nil)

	fmt.Println("Image resized and saved successfully.")
}

func TestAes(t *testing.T) {
	//key := make([]byte, 32) // 32 bytes for AES-256
	//if _, err := rand.Read(key); err != nil {
	//	fmt.Println("Error generating key:", err)
	//	return
	//}
	//keyHex := hex.EncodeToString(key)
	//fmt.Println(keyHex)
	redis.Init()
	redisCli := redis.RedisManager()

	storedHexKey := "d4e9a43db1f266980f3cca85f00ca47606b9d92b82e8fac632c5b45090eeb5b4"
	err := redisCli.Set(context.Background(), "qr:aes_key", storedHexKey, time.Hour*60*24).Err()

	if err != nil {
		logger.Error.Printf("err: %v\n", err)
	}

	//file, err := os.Create("./output.csv")
	//if err != nil {
	//	panic(err)
	//}
	//wr := csv.NewWriter(bufio.NewWriter(file))
	//key, err := hex.DecodeString(storedHexKey)
	//if err != nil {
	//	fmt.Println("Error decoding key from hexadecimal:", err)
	//	return
	//}
	//fmt.Println(key)
	//
	//for i := 6; i < 250; i++ {
	//	countStr := strconv.Itoa(i)
	//	workPlainText := fmt.Sprintf("nftime%s", countStr)
	//	data, err := Encrypt([]byte(workPlainText), key)
	//	if err != nil {
	//		t.Errorf("error: %v\n", err)
	//	}
	//	fmt.Printf("here:%v\n", string(data))
	//
	//	encodedData := base64.RawURLEncoding.EncodeToString(data)
	//	fmt.Printf("Base64 Encoded data:%v\n", encodedData)
	//
	//	decodedData, err := base64.RawURLEncoding.DecodeString(encodedData)
	//	if err != nil {
	//		fmt.Println("Error decoding data:", err)
	//		return
	//	}
	//	fmt.Println("Decoded data:", string(decodedData))
	//
	//	//data2, err := Decrypt(data, key)
	//	//if err != nil {
	//	//	fmt.Println(err)
	//	//}
	//	wr.Write([]string{countStr, "https://nftime.store/klip/mint/swf/work/" + string(encodedData)})
	//}
	//wr.Flush()
}
