package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestImageToByte(t *testing.T) {
	imageUrl := "https://d3c9sn363qnamy.cloudfront.net/images/nftime-nft-vip_1696552003.png" // Replace with the URL of the image you want to fetch

	// Make an HTTP GET request to the image URL
	response, err := http.Get(imageUrl)
	if err != nil {
		t.Errorf("err: %v\n", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("HTTP request failed with status code: %d", response.StatusCode)
		return
	}

	imageBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		// Handle the error
		return
	}

	// Get the image byte data from the buffer
	fmt.Println(imageBytes)
}
