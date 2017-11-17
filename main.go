package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"sort"
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func getPixels(file io.Reader) (map[string]int, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()

	width, height := bounds.Max.X, bounds.Max.Y

	var pixels = make(map[string]int)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			hex := rgbaToPixel(img.At(j, i).RGBA())
			total, ok := pixels[hex]

			if ok {
				pixels[hex] = total + 1
			} else {
				pixels[hex] = 1
			}
		}
	}

	return pixels, nil
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) string {
	return fmt.Sprintf("%02x%02x%02x", int(r/257), int(g/257), int(b/257))
}

func getHank(words map[string]int) PairList {
	pl := make(PairList, len(words))
	i := 0

	for k, v := range words {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func getPredominant(colors PairList, amount int) PairList {
	pl := make(PairList, 10)
	counter := 0

	for i, v := range colors {
		if i < 10 {
			pl[i] = Pair{v.Key, v.Value}
		}
		counter++
	}

	return pl
}

func main() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open("./golang_image.png")

	if err != nil {
		fmt.Println("Error opening image")
		os.Exit(1)
	}

	defer file.Close()

	pixels, err := getPixels(file)

	if err != nil {
		panic(err)
	}

	hank := getHank(pixels)

	predominant := getPredominant(hank, 10)
	fmt.Println(predominant)
}
