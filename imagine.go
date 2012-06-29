
package main

import (
	"fmt"
	"os"
	"image"
	_ "image/png"
	"log"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Please supply an image file to process.")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
	    log.Fatal(err)
	}
	defer file.Close()

	// Decode the image.
	m, _, err := image.Decode(file)
	if err != nil {
	    log.Fatal(err)
	}
	bounds := m.Bounds()

	m := image.NewRGBA(bounds)
	color := new(color.RGBA)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	    for x := bounds.Min.X; x < bounds.Max.X; x++ {
	        r, g, b, a := m.At(x, y).RGBA()
	        // A color's RGBA method returns values in the range [0, 65535].
	        // Shifting by 12 reduces this to the range [0, 15].
	        histogram[r>>12][0]++
	        histogram[g>>12][1]++
	        histogram[b>>12][2]++
	        histogram[a>>12][3]++
	    }
	}

	fmt.Printf("%#v\n",bounds)
}