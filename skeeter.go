package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"strconv"
)

const (
	REDWEIGHT    = 0.2989
	GREENWEIGHT  = 0.5866
	BLUEWEIGHT   = 0.1145
	ASCIIPALETTE = "   ...',;:clodxkO0KXNWM"
)

func getImage(url string) image.Image {
	fmt.Printf("Ignoring url %s\n", url)

	file, err := os.Open("/Users/blake/projects/skeeter-go/images/moose.png")
	if err != nil {
		fmt.Println("Failed to open file!")
	}

	im, _, err := image.Decode(file)

	if err != nil {
		fmt.Println("Failed to decode image!")
	}

	return im
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

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
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

	fmt.Printf("Bounds of image: %d, %d, %d, %d\n", bounds.Min.X, bounds.Max.X, bounds.Min.Y, bounds.Max.Y)

	w, h := asciiDimensions(bounds, width)
	scaledImage := Resize(img, bounds, w, h)

	fmt.Printf("Ascii dimensions: %d, %d\n", w, h)

	out := printAscii(scaledImage)

	return out
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
