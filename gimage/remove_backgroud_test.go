package gimage

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestRemoveApi(t *testing.T) {
	file, err := os.Open("/Users/me/Downloads/01-道场小恶魔-3D版/主图/01.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a new multipart form
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add the image file to the form
	part, err := writer.CreateFormFile("image_file", "input.jpg")
	if err != nil {
		panic(err)
	}
	io.Copy(part, file)

	// Add the API key to the form
	part, err = writer.CreateFormField("size")
	if err != nil {
		panic(err)
	}
	part.Write([]byte("auto"))

	// Close the form
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	// Send the HTTP request to the API endpoint
	req, err := http.NewRequest("POST", "https://api.remove.bg/v1.0/removebg", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Api-Key", "e2EmV9bShqJ41us1a2PX6aAK")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Save the output image to a file
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("/Users/me/Downloads/01-道场小恶魔-3D版/主图/01-remove.png", data, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Background removed successfully!")
}
