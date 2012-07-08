# Skeeter for Go

This is a port of [skeeter](http://github.com/blakesmith/skeeter) to Go.

## Why?

The original skeeter implementation was built with ZeroMQ in Ruby and had seperate processes for each piece. Ensuring all the pieces were behaving correctly became a hassle, so I ported it to Go to get familiar with the net/http library in Go.

## Dependencies

See [Native Dependencies](https://github.com/blakesmith/skeeter#native-dependencies) from [skeeter](http://github.com/blakesmith/skeeter)

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



