# Skeeter for Go

This is a port of [skeeter](http://github.com/blakesmith/skeeter) to Go.

## What is it?

Convert this:

![Original
image](https://github.com/blakesmith/skeeter-go/raw/master/images/moose.png)

Into this:

![Converted
image](https://github.com/blakesmith/skeeter-go/raw/master/images/moose-ascii.jpg)

You make a request to it like so:

  http://skeeter.blakesmith.me/?image_url=http://www.softicons.com/download/animal-icons/animal-icons-by-martin-berube/png/128/moose.png&width=100

And it spits out the ascii art! Magic!

## Why do this?

The original skeeter implementation was built with ZeroMQ in Ruby and had seperate processes for each piece. Ensuring all the pieces were behaving correctly and always up became a hassle. The Ruby version also shelled out to jp2a and imagemagick to do the actual image maniuplation. I ported everything to Go, including image fetching and borrowed the jp2a ascii algorithm to do everything in one native Go program. No external dependencies, just one binary that you drop onto the server and run!

## Building

[Install Go](http://golang.org/doc/install)

Then run ```go run skeeter.go resize.go -port=9001``` from within the skeeter directory.

To cross compile for linux 386 from OS X, I followed this wiki page to setup the compiler toolchain: http://code.google.com/p/go-wiki/wiki/WindowsCrossCompiling

Once that's setup, I ran:

```
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build
```

This generates the ```skeeter-go``` binary, which you run to start the web server.

## The Future

With the present design, ```curl``` ```convert``` and ```jp2a``` are used to process the image. I'd like to convert these to native Go code.

## Author

Skeeter is written by Blake Smith <blakesmith0@gmail.com>.



