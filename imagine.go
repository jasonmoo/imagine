
package main

import (
	"fmt"
	"os"
	"image"
	"image/color"
	"image/png"
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

	newm := image.NewRGBA(bounds)
	c := new(color.RGBA)
	c.A = 255

	dev := uint32(30)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	    for x := bounds.Min.X; x < bounds.Max.X; x++ {

	    	var r, g, b, ct uint32

			_r, _g, _b, _ := m.At(x,y).RGBA()

	    	for i := -6; i < 7; i++ {
	    		for j := -6; j < 7; j++ {
	    			rt, gt, bt, _ := m.At(x+i,y+j).RGBA()

	    			if uint32(rt-_r) > dev  || uint32(gt-_g) > dev || uint32(bt-_b) > dev {
	    				continue
	    			}

	    			rt >>= 8
	    			gt >>= 8
	    			bt >>= 8

	    			r += rt
	    			g += gt
	    			b += bt
	    			ct++
	    		}
	    	}
	    	c.R = uint8(r/ct)
	    	c.G = uint8(g/ct)
	    	c.B = uint8(b/ct)

	    	// r, g, b, a := m.At(x, y).RGBA()

	    	newm.Set(x,y, c)

	    }
	}

	toimg, _ := os.Create("output.png")
	defer toimg.Close()

	err = png.Encode(toimg, newm)
	if err != nil {
		log.Fatal(err)
	}


	fmt.Println("done")
}











