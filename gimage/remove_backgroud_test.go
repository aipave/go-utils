package gimage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestRemoveApi(t *testing.T) {
	outDirPath := filepath.Join("/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/", "removed")
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outDirPath, 0755); err != nil {
		panic(err)
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outDirPath, 0755); err != nil {
		panic(err)
	}
	file, err := os.Open("/Users/ricco/Downloads/01-道场小恶魔-3D版/主图/01.png")
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
	req.Header.Set("X-Api-Key", "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	//ext := filepath.Ext(file.Name())
	//filename := strings.TrimSuffix(filepath.Base(file.Name()), ext)
	outFile, err := os.Create(filepath.Join(outDirPath, file.Name()))
	if err != nil {
		logrus.Errorf("Error:%v", err)
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		logrus.Errorf("Error:%v", err)
	}

	// Save the resulting image to a new file
	//outFile, err := os.Create(filepath.Join("/Users/me/Downloads/01-道场小恶魔-3D版/主图/remove/", file.Name()))
	//if err != nil {
	//	fmt.Printf("Error creating output file: %v\n", err)
	//}
	//defer outFile.Close()

	//if _, err := io.Copy(outFile, resp.Body); err != nil {
	//	fmt.Printf("Error saving output file: %v\n", err)
	//}

	// Save the output image to a file
	//data, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//err = ioutil.WriteFile("/Users/me/Downloads/01-道场小恶魔-3D版/主图/remove/01-remove.png", data, 0644)
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println("Background removed successfully!")
}
