
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
	color_dev := uint16(50<<8)
	feature_dev := 10

	bounds := orig.Bounds()

	fmt.Println(bounds)

	// initialize the array of pixels traversed
	ex := make([][]bool, bounds.Max.X+1)
	for i := 0; i < bounds.Max.X+1; i++ {
		ex[i] = make([]bool, bounds.Max.Y+1)
	}

	// array of features
	// feature is an array of arrays [x,y]
	var features [][][2]int
	feature_i := 0

	// recursively investigate all neighbors of supplied pixel
	// and build feature arrays
	// some shenanigans to make anonymous functions recurisively callable
	discover := func(x int,y int) {}
	f := func(x int,y int) {

		fmt.Println(feature_i)

		// add this pixel to the current feature
		features[feature_i] = append(features[feature_i], [2]int{x,y})
		// fmt.Println(x,y)
		ex[x][y] = true

		// grab the rgb for the supplied pixel
		_r, _g, _b, _ := orig.At(x,y).RGBA()

		for i := -1; i < 2; i++ {
			for j := -1; i < 2; i++ {
				xx, yy := x+i, y+j
				// check if it's within our bounds and if it's been processed already
				if xx < 0 || yy < 0 || xx > bounds.Max.X || yy > bounds.Max.Y || ex[xx][yy] == true {
					continue
				}

				fmt.Println("here")
				rt, gt, bt, _ := orig.At(xx,yy).RGBA()
				// check the color range against our deviation spec
				if uint16(rt-_r) > color_dev || uint16(gt-_g) > color_dev || uint16(bt-_b) > color_dev {
					continue
				}
				discover(xx, yy)
			}
		}

	}
	discover = f

	// run through each pixel and build features
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	    for x := bounds.Min.X; x < bounds.Max.X; x++ {

	    	fmt.Println(ex[x][y])

	    	// skip if already processed
	    	if ex[x][y] == true {
	    		continue
	    	}

	    	fmt.Println(x,y)
	    	features = append(features, [][2]int{})
	    	discover(x,y)
	    	feature_i++

	    }
	}
    fmt.Println(features)

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







