package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	REDWEIGHT    = 0.2989
	GREENWEIGHT  = 0.5866
	BLUEWEIGHT   = 0.1145
	ASCIIPALETTE = "   ...',;:clodxkO0KXNWM"
)

var (
	port = flag.String("port", "9001", "Port to run the server on")
)

func getImage(url string) (image.Image, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, errors.New("Failed to open file!")
	}
	defer res.Body.Close()

	im, _, err := image.Decode(res.Body)

	if err != nil {
		return nil, errors.New("Failed to decode image!")
	}

	return im, nil
}

func getAsciiChar(red uint32, green uint32, blue uint32) string {
	v := (float64(red>>8) * REDWEIGHT / 255.0) +
		(float64(green>>8) * GREENWEIGHT / 255.0) +
		(float64(blue>>8) * BLUEWEIGHT / 255.0)
	idx := int(v * float64(len(ASCIIPALETTE)-1))
	char := ASCIIPALETTE[idx]

	return string(char)
}

func asciiDimensions(b image.Rectangle, width int) (w int, h int) {
	ratio := float64(b.Max.Y-b.Min.Y) / float64(b.Max.X-b.Min.X)

	return width, int(ratio * float64(width))
}

func printAscii(img image.Image) string {
	bounds := img.Bounds()
	out := bytes.NewBufferString("")

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			fmt.Fprint(out, getAsciiChar(r, g, b))
		}
		fmt.Fprint(out, "\n")
	}

	return out.String()
}

func toAscii(img image.Image, width int) string {
	bounds := img.Bounds()

	log.Printf("Bounds of image: %d, %d, %d, %d\n", bounds.Min.X, bounds.Max.X, bounds.Min.Y, bounds.Max.Y)

	w, h := asciiDimensions(bounds, width)
	scaledImage := Resize(img, bounds, w, h)

	log.Printf("Ascii dimensions: %d, %d\n", w, h)

	out := printAscii(scaledImage)

	return out
}

func makeImage(url string, width string) (string, error) {
	w, err := strconv.Atoi(width)

	if err != nil {
		return "", errors.New("Please enter a valid integer 'width' param")
	}

	img, err := getImage(url)

	if err != nil {
		return "", err
	}

	out := toAscii(img, w)

	return string(out), nil
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	r.ParseForm()
	var url, width string

	if val, ok := r.Form["image_url"]; ok {
		url = val[0]
	} else {
		http.Error(w, fmt.Sprintf("[:error, \"Image_url query param is required\"]"), 400)

		return
	}

	if val, ok := r.Form["width"]; ok {
		width = val[0]
	} else {
		width = "80"
	}

	log.Println("---")
	log.Printf("Processing request with url %s, and width %s\n", url, width)

	out, err := makeImage(url, width)

	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		fmt.Fprintf(w, out)
	}
	log.Printf("Ascii generated in: %d ms", time.Now().Sub(start).Nanoseconds()/1000000)
	log.Println("---")
}

func main() {
	flag.Parse()

	log.Printf("Running skeeter on port %s\n", *port)
	http.HandleFunc("/", imageHandler)
	http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
}
