
package main

import (
	"fmt"
	"os"
	"image"
	"image/color"
	"image/png"
	_ "image/jpeg"
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

	c := new(color.RGBA)
	c.A = 255

	dev := uint16(36<<8)

	for i := 0; i < 1; i++ {
		newm := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		    for x := bounds.Min.X; x < bounds.Max.X; x++ {

		    	var r, g, b, ct uint32

				_r, _g, _b, _ := m.At(x,y).RGBA()

		    	for i :=-1; i < 2; i++ {
		    		for j := -1; j < 2; j++ {
		    			rt, gt, bt, _ := m.At(x+i,y+j).RGBA()

		    			if uint16(rt-_r) > dev  || uint16(gt-_g) > dev || uint16(bt-_b) > dev {
		    				continue
		    			}

		    			r += rt
		    			g += gt
		    			b += bt
		    			ct++
		    		}
		    	}
		    	c.R = uint8((r/ct)>>8)
		    	c.G = uint8((g/ct)>>8)
		    	c.B = uint8((b/ct)>>8)

		    	// r, g, b, a := m.At(x, y).RGBA()

		    	newm.Set(x,y, c)

		    }
		}

		m = newm
	}

	outfile, _ := os.Create("output.png")
	defer outfile.Close()

	err = png.Encode(outfile, m)
	if err != nil {
		log.Fatal(err)
	}


	fmt.Println("done")
}











