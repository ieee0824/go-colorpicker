package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/ieee0824/go-colorpicker"
)

func main() {
	imageDir := "image/"
	infos, err := ioutil.ReadDir(imageDir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, info := range infos {
		file, err := os.Open(imageDir + info.Name())
		if err != nil {
			log.Fatalln(err)
		}
		img, _, _ := image.Decode(file)
		colorClass := colorpicker.ExtractTypicalColors(img, 5)
		orgName := info.Name()
		fmt.Println(orgName)
		fmt.Println(len(colorClass))
		for i, hsv := range colorClass {
			fmt.Println(hsv)
			r, g, b, _ := hsv.RGBA()
			rgba := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: 0xff,
			}
			m := image.NewRGBA(image.Rect(0, 0, 100, 100))
			draw.Draw(m, m.Bounds(), &image.Uniform{rgba}, image.ZP, draw.Src)
			toimg, err := os.Create("result/" + orgName + "_" + fmt.Sprintf("%04d", i) + ".jpg")
			defer toimg.Close()
			if err != nil {
				log.Fatalln(err)
			}
			jpeg.Encode(toimg, m, &jpeg.Options{jpeg.DefaultQuality})
		}
		fmt.Println("=====")
	}
}
