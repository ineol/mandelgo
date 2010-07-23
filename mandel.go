package main

import (
	"image"
        "os"
	"image/png"
	"fmt"
	"runtime"
	"flag"
)

const max_iteration = 100

func getPixelAt(x0, y0 float64) image.RGBAColor  {
	iteration := 0; 
	var x float64 =  0.0
	var y float64 = 0.0
	var xtemp float64 = 0.0
	for x*x + y*y <= 4.0 && iteration < max_iteration {
		xtemp = x*x - y*y + x0
		y = 2*x*y + y0
		x = xtemp
		iteration += 1
	}
	if iteration == max_iteration {
		return image.RGBAColor{0, 0, 0, 255}
	} else {
		color := uint8(255 - 255 * iteration / max_iteration)
		return image.RGBAColor{color, color, color, 255}
	}
	fmt.Printf("This is impossible\n")
	panic(1)
}

func Mandelbrot(im *image.RGBA, lineMin, lineMax int, vpx, vpy, d float64) {
	width := float64(im.Width())
	height := float64(im.Height())
	for i := lineMin; i < lineMax; i++ {
		for j := 0; j < im.Height(); j++ {
			x0 := float64(i) / (width / d) - d / 2.0 + vpx
			y0 := float64(j) / (height / d) - d / 2.0 + vpy
			im.Pixel[i][j] = getPixelAt(y0, x0)
		}
	}
}

func Start(im *image.RGBA, num int, vpx, vpy, d float64) {
	share := im.Height() / num
	for i := 0; i < num - 1; i += 1 {
		go Mandelbrot(im, i * share, (i+1) * share, vpx, vpy, d)
	}
	Mandelbrot(im, (num - 1) * share, num * share, vpx, vpy, d)
}

func main() {
	size := flag.Int("psize", 500, "physical size of the square image")
	vpx := flag.Float64("x", 0, "x coordinate of the center of the image")
	vpy := flag.Float64("y", 0, "y coordinate of the center of the image")
	d := flag.Float64("size", 2, "size of the represented part of the plane")
	filename := flag.String("name", "image", "name of the image file produced (w/o extension")
	numberOfProcs := flag.Int("procs", 2, "number of procs to use")

	flag.Parse()

	runtime.GOMAXPROCS(*numberOfProcs)

	file, err := os.Open(*filename + ".png", os.O_RDWR | os.O_CREAT, 0666)
	if err != nil {
		panic("error with opening file \"" + *filename + "\"")
	}
	

	im := image.NewRGBA(*size, *size)	
	Start(im, 2, *vpx, *vpy, *d)
	png.Encode(file, im)
}
