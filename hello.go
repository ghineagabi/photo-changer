package main

import (
	"fmt"
	_ "github.com/u2takey/ffmpeg-go"
	"image"
	"image/color"
	"image/draw"
	_ "image/draw"
	"image/jpeg"
	"io/fs"
	"math"
	"math/rand"
	"os"
	"strconv"
)

type Band struct {
	arrPoints      []image.Point
	width          int
	diagonalLength int
	offsetX        int
}

func main() {

	nbBands := 10
	width := 7
	diagLen := 5
	smoothness := 16
	offsetStart := 1

	if _, err := os.Stat("converted"); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("converted", fs.ModeDir)
			if err != nil {
				fmt.Println("could not create a folder named <converted> ...")
				return
			}
		} else {
			fmt.Println("unknown error :(")
			return
		}
	}
	filename := "download.jpg"
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Could not find <%s>. Check if file exists or change the extension", filename)
			return
		} else {
			fmt.Println("unknown error :(")
			return
		}
	}

	fmt.Println("creating photos...")
	for i := 0; i < 100; i++ {
		pathOut := "converted/" + intToCorrectString(i) + ".jpg"
		offsetStart += 10
		err := convertImage("download.jpg", pathOut, toAddGaussian,
			nbBands, width, diagLen, smoothness, offsetStart)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("done!")

	fmt.Println("Required command:", `ffmpeg -framerate 10 -i "converted/%03d.jpg" -vf "pad=ceil(iw/2)*2:ceil(ih/2)*2" output.mp4`)
}

func intToCorrectString(i int) string {
	if i/10 == 0 {
		return "00" + strconv.Itoa(i)
	} else if i/100 == 0 {
		return "0" + strconv.Itoa(i)
	}
	return strconv.Itoa(i)
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	im, _, err := image.Decode(f)
	return im, err
}

func convertImage(pathIn string, pathOut string, f func(*image.Image, ...any) *image.RGBA, args ...any) error {
	im, err := getImageFromFilePath(pathIn)
	if err != nil {
		return err
	}

	outFile, err := os.Create(pathOut)
	if err != nil {
		return err
	}

	newImg := f(&im, args...)

	defer outFile.Close()
	err = jpeg.Encode(outFile, newImg, nil)
	if err != nil {
		return err
	}

	return nil
}

func toGray(im *image.Image) *image.RGBA {
	bounds := (*im).Bounds()
	max := bounds.Max
	min := bounds.Min
	imgSet := image.NewRGBA(bounds)

	for i := min.X; i < max.X; i++ {
		for j := min.Y; j < max.Y; j++ {
			pixel := (*im).At(i, j)
			R, G, B, _ := pixel.RGBA()
			lum := (19595*R + 38470*G + 7471*B + 1<<15) >> 24
			newPixel := color.Gray{Y: uint8(lum)}
			imgSet.Set(i, j, newPixel)
		}
	}
	return imgSet
}

func toAddGaussian(im *image.Image, args ...any) *image.RGBA {
	bounds := (*im).Bounds()
	max := bounds.Max
	imgSet := image.NewRGBA(bounds)
	draw.Draw(imgSet, imgSet.Bounds(), *im, bounds.Min, draw.Src)
	nbBands := args[0].(int)
	width := args[1].(int)
	diagLen := args[2].(int)
	smoothness := args[3].(int)
	offsetStart := args[4].(int)
	var bands []Band

	offsetStart = offsetStart % max.Y / nbBands
	for offset := offsetStart; offset <= max.Y; offset = offset + max.Y/nbBands {
		for _, b := range *generateBand(width, diagLen, max.X, offset, smoothness) {
			bands = append(bands, b)
		}
	}

	for _, b := range bands {
		for _, p := range b.arrPoints {
			newPixel := color.Gray{Y: uint8(rand.Intn(255))}
			imgSet.Set(p.X, p.Y, newPixel)
		}
	}
	return imgSet

}

func generateBand(width, diagonalLength, bandLength, bandStart, smoothness int) *[]Band {
	Bands := make([]Band, 2*width-1)
	singleBand := make([]image.Point, bandLength)
	for j := 0; j < bandLength; j++ {
		singleBand[j].Y = int(float64(bandStart) +
			float64(diagonalLength)*math.Sin(math.Pi*float64(j)/float64(smoothness)))
		singleBand[j].X = j
	}

	Bands[0] = Band{arrPoints: singleBand,
		width:          width,
		diagonalLength: diagonalLength,
		offsetX:        bandStart}
	upperBand := Band{}
	lowerBand := Band{}

	for i := 1; i < width; i++ {
		upperBand = Bands[0].copy()
		lowerBand = Bands[0].copy()
		for j := 0; j < bandLength; j++ {
			upperBand.arrPoints[j] = upperBand.arrPoints[j].Add(image.Point{Y: i})
			lowerBand.arrPoints[j] = lowerBand.arrPoints[j].Add(image.Point{Y: -i})
		}
		Bands[2*i-1] = upperBand
		Bands[2*i] = lowerBand
	}
	return &Bands
}

func (b *Band) copy() Band {
	arrP := make([]image.Point, len(b.arrPoints))
	newBand := Band{arrPoints: arrP,
		width:          b.width,
		diagonalLength: b.diagonalLength,
		offsetX:        b.offsetX}
	for i, point := range b.arrPoints {
		newBand.arrPoints[i] = point
	}
	return newBand
}
