package main

import (
	"downloader/helpers"
	"fmt"
	"os"
)

func main() {
	// creating output file
	file, err := os.OpenFile("output/file.zip", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	url := "https://github.com/teddy-codes/url-shortner/archive/refs/heads/master.zip"
	download := helpers.NewDownload(url, 10)
	err = download.Do(file)
	fmt.Println(err)
}
