package gimage

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func removeApi(apiKey string, path string, p string) error {
	file, err := os.Open(p)
	if err != nil {
		logrus.Errorf("Error:%v", err)
		return err
	}
	defer file.Close()

	// Create a new multipart form
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add the image file to the form
	part, err := writer.CreateFormFile("image_file", p)
	if err != nil {
		logrus.Errorf("Error:%v", err)
		return err
	}
	io.Copy(part, file)

	// Add the API key to the form
	part, err = writer.CreateFormField("size")
	if err != nil {
		logrus.Errorf("Error:%v", err)
		return err
	}
	part.Write([]byte("auto"))

	// Close the form
	err = writer.Close()
	if err != nil {
		logrus.Errorf("Error:%v", err)
		return err
	}

	// Send the HTTP request to the API endpoint
	req, err := http.NewRequest("POST", "https://api.remove.bg/v1.0/removebg", body)
	if err != nil {
		logrus.Errorf("Error:%v", err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error:%v", err)
		return err
	}
	defer resp.Body.Close()

	// Save the output image to a file
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Error:%v", err)
		return err
	}

	outDirPath := filepath.Join(path, "removed")
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outDirPath, 0755); err != nil {
		panic(err)
	}
	splitArr := strings.Split(file.Name(), "/")
	err = ioutil.WriteFile(fmt.Sprintf("%v/%v",
		outDirPath, splitArr[len(splitArr)-1]), data, 0644)
	if err != nil {
		logrus.Errorf("Error:%v", err)
		return err
	}

	logrus.Infof("Background removed successfully!")

	return nil
}

func RemoveBackground(apiKey string, inDirPath string) error {
	return filepath.Walk(inDirPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(p, ".png") || strings.HasSuffix(p, ".jpg") {
			err := removeApi(apiKey, inDirPath, p)
			if err != nil {
				return err
			}
		}
		return nil
	})

}
