# GOMB - Go Mandelbrot

Gomb is a simple tool that renders the Mandelbrot fractal set in ASCII

## Building

go get .
go build .

## Usage

Gomb will detect the dimensions of your terminal and use the entire space to
render the image. The following switches are supported:

| Argument        | Description                         |
|-----------------|-------------------------------------|
| -invert         | Invert palette                      |
| -iterations int | Number of iterations (default 128)  |
| -x float        | Position x                          |
| -y float        | Position y                          |
| -zoom float     | Zoom factor (default 1)             |
