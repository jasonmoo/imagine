
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

	m = featurize(m)
	// m = blend(m)

	outfile, _ := os.Create("output.png")
	defer outfile.Close()

	err = png.Encode(outfile, m)
	if err != nil {
		log.Fatal(err)
	}


	fmt.Println("done")
}

// for some reason required to make anonymous
// functions able to call themselves
// type discover func(x int, y, int) (nil)

func featurize(orig image.Image) (image.Image) {

	// deviation range
	color_dev := uint16(100<<8)
	feature_dev := 10

	bounds := orig.Bounds()

	ex := [bounds.Max.X-1][bounds.Max.Y-1]bool

	// array of features
	// feature is an array of arrays [x,y]
	var features [][][]int
	feature_i := 0

	// recursively investigate all neighbors of supplied pixel
	// and build feature arrays
	// some shenanigans to make anonymous functions recurisively callable
	discover := func(x int,y int) {}
	f := func(x int,y int) {

		// add this pixel to the current feature
		features[feature_i] = append(features[feature_i], []int{x,y})
		ex[x][y] = true

		// grab the rgb for the supplied pixel
		_r, _g, _b, _ := orig.At(x,y).RGBA()

		for i := -1; i < 2; i++ {
			for j := -1; i < 2; i++ {
				rt, gt, bt, _ := orig.At(x+i,y+j).RGBA()

				if uint16(rt-_r) > color_dev || uint16(gt-_g) > color_dev || uint16(bt-_b) > color_dev {
					continue
				}
				discover(x+i, y+j)
			}
		}

	}
	discover = f

	// run through each pixel and build features
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	    for x := bounds.Min.X; x < bounds.Max.X; x++ {

	    	// skip if already processed
	    	if ex[x][y] == true {
	    		continue
	    	}

	    	discover(x,y)
	    	feature_i++

	    }
	}

	newm := image.NewRGBA(bounds)
	c := new(color.RGBA)
	c.A = 255

	for f := 0; f < len(features); f++ {

		// if the feature is large enough
		// average all pixel colors in it
		// and set all pixels to that color
		if len(features[f]) > feature_dev {
			var r, g, b, ct uint32

			for p := 0; p < len(features[f]); p++  {
				rt, gt, bt, _ := orig.At(features[f][p][0], features[f][p][1]).RGBA()

				r += rt
				g += gt
				b += bt
				ct++
			}
			c.R = uint8((r/ct)>>8)
			c.G = uint8((g/ct)>>8)
			c.B = uint8((b/ct)>>8)

			for p := 0; p < len(features[f]); p++  {
				newm.Set(features[f][p][0], features[f][p][1], c)
			}
		} else {
			// write the pixel out as-is
			for p := 0; p < len(features[f]); p++  {
				r, g, b, _ := orig.At(features[f][p][0], features[f][p][1]).RGBA()
				c.R = uint8(r>>8)
				c.G = uint8(g>>8)
				c.B = uint8(b>>8)
				newm.Set(features[f][p][0], features[f][p][1], c)
			}
		}
	}

	return newm
}


func blend(orig image.Image) (m image.Image) {
	// deviation range
	dev := uint16(36<<8)

	c := new(color.RGBA)
	c.A = 255

	bounds := orig.Bounds()

	// iterations
	for i := 0; i < 1; i++ {
		newm := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		    for x := bounds.Min.X; x < bounds.Max.X; x++ {

		    	var r, g, b, ct uint32

				_r, _g, _b, _ := orig.At(x,y).RGBA()

		    	for i :=-1; i < 2; i++ {
		    		for j := -1; j < 2; j++ {
		    			rt, gt, bt, _ := orig.At(x+i,y+j).RGBA()

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

		    	newm.Set(x,y, c)

		    }
		}

		m = newm
	}

	return
}







