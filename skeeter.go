package main

import (
	"fmt"
	"net/http"
	"image"
	"os"
	"strconv"
)

import _ "image/png"
import _ "image/jpeg"
import _ "image/gif"

func getImage(url string) image.Image {
	fmt.Printf("Ignoring url %s\n", url)

	file, err := os.Open("/Users/blake/projects/skeeter-go/images/moose-ascii.jpg")
	if err != nil {
		fmt.Println("Failed to open file!")
	}

	im, _, err := image.Decode(file)
	
	if err != nil {
		fmt.Println("Failed to decode image!")
	}

	return im
}

func toAscii (img image.Image, width int) string {
	return "00000000000"
}

func makeImage(url string, width string) string {
	img := getImage(url)
	w, err := strconv.Atoi(width)
	if err != nil {
		fmt.Printf("Failed to convert width to an int!")
	}
	out := toAscii(img, w)

	return string(out)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var url, width string

	if val, ok := r.Form["image_url"]; ok {
		url = val[0]
	} else {
		// 400 Bad Request: [:error, "Image_url query param is required"]
		http.Error(w, fmt.Sprintf("[:error, \"Image_url query param is required\"]"), 400)

		return
	}

	if val, ok := r.Form["width"]; ok {
		width = val[0]
	} else {
		width = "80"
	}

	fmt.Printf("Processing request with url %s, and width %s\n", url, width)
	out := makeImage(url, width)

	fmt.Fprintf(w, out)
}

func main() {
	fmt.Println("Running server...")
	http.HandleFunc("/", imageHandler)
	http.ListenAndServe(":9001", nil)
}
