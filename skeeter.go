package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

// Borrowed from: https://gist.github.com/1477401
func pipe_commands(commands ...*exec.Cmd) ([]byte, error) {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		command.Start()
		commands[i+1].Stdin = out
	}
	final, err := commands[len(commands)-1].Output()
	if err != nil {
		return nil, err
	}
	return final, nil
}

func makeImage(url string, width string) string {
	curl := exec.Command("curl", "-s", url)
	convert := exec.Command("convert", "-", "jpg:-")
	jp2a := exec.Command("jp2a", "-", fmt.Sprintf("--width=%s", width))

	out, err := pipe_commands(curl, convert, jp2a)
	if err != nil {
		fmt.Printf("Commands failed to run: %s\n", err)
	}

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
	http.HandleFunc("/", imageHandler)
	http.ListenAndServe(":9001", nil)
}
