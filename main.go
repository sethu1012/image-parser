package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"

	"gopkg.in/alecthomas/kingpin.v2"

	"image/jpeg"
	"image/png"
)

var (
	debug        = kingpin.Flag("debug", "Enable debug mode.").Bool()
	imgPath      = kingpin.Flag("path", "Path to image").Required().String()
	outputFormat = kingpin.Flag("output", "Output format").Required().String()
	rzImage      = kingpin.Flag("resize", "Enable resize image").Bool()
	thumbnail    = kingpin.Flag("thumbnail", "Generate thumbnail").Bool()
	quality      = kingpin.Flag("q", "Image Quality").Int()
	height       = kingpin.Flag("h", "Height").Int()
	width        = kingpin.Flag("w", "Width").Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	img, _, err := GetImage(*imgPath)
	if err != nil {
		log.Panicln(err.Error())
	}
	if *rzImage {
		ConvertImage(ResizeImage(img), *outputFormat)
	} else {
		ConvertImage(img, *outputFormat)
	}
}

// GetImage - Reads image and returns image object
func GetImage(path string) (image.Image, string, error) {
	reader, err := os.Open(path)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer reader.Close()
	return image.Decode(reader)
}

// ResizeImage - Resize the provided image
func ResizeImage(img image.Image) image.Image {
	if *thumbnail {
		return resize.Thumbnail(uint(*width), uint(*height), img, resize.NearestNeighbor)
	}
	return resize.Resize(uint(*width), uint(*height), img, resize.NearestNeighbor)
}

// ConvertImage - Write image to file system
func ConvertImage(img image.Image, format string) {
	f, err := os.Create(fmt.Sprintf("rome.%s", format))
	if err != nil {
		log.Panicln("Could not create file")
	}
	defer f.Close()
	switch format {
	case "png":
		png.Encode(f, img)
	case "jpg":
		jpeg.Encode(f, img, &jpeg.Options{Quality: *quality})
	case "webp":
		webp.Encode(f, img, &webp.Options{Lossless: false, Quality: float32(*quality)})
	default:
		log.Panicln("Format not supported")
	}
}
