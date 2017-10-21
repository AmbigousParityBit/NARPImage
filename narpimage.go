// notaregularpixel package
package NARPImage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
)

type NARPImage struct {
	NARPixels []NotARegularPixel
	Size      struct{ X, Y uint16 }
	Version   string
}

func (narpimage *NARPImage) deconstructToImage() (img *image.RGBA, err error) {
	img = image.NewRGBA(image.Rect(0, 0, int(narpimage.Size.X), int(narpimage.Size.Y)))
	var visited [][]bool
	initVisitedArray(&visited, int(narpimage.Size.X), int(narpimage.Size.Y))
	x, y := uint16(0), uint16(0)

	for _, v := range narpimage.NARPixels {
		v.drawNARP(img, int(x), int(y))
		v.markVisited(int(x), int(y), &visited, int(narpimage.Size.X), int(narpimage.Size.Y))
		end := false

		for !end && visited[x][y] {
			x++
			if x >= narpimage.Size.X {
				x = 0
				y++
				if y >= narpimage.Size.Y {
					end = true
				}
			}
		}
	}
	/*
		for _, narpixel := range narpimage.NARPixels {
			color := color.RGBA{narpixel.Color.R, narpixel.Color.G, narpixel.Color.B, 255}

			if !(visited[x][y]) {
				drawAndMark(img, x, y, color, &visited)
				for h := uint8(0); h < narpixel.HSize; h++ {
					xH := x + uint16(h)
					drawAndMark(img, xH, y, color, &visited)
					if narpixel.VSize != nil && len(narpixel.VSize) > 0 {
						vsize := putBytesToUint16(narpixel.VSize[h])
						for v := uint16(0); v < vsize; v++ {
							yV := y + uint16(v)
							drawAndMark(img, xH, yV, color, &visited)
						}
					}
				}
			}
			for visited[x][y] {
				x++
				if x >= narpimage.Size.X {
					y++
					x = 0
				}
			}
		}
	*/
	return img, nil
}

func (narpimage *NARPImage) DeconstructToPngFile(s string) error {
	f, err := os.OpenFile(s, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := narpimage.deconstructToImage()
	if err != nil {
		return err
	}
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}

func (narpimage *NARPImage) DeconstructToJpgFile(s string) error {
	f, err := os.OpenFile(s, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := narpimage.deconstructToImage()
	if err != nil {
		return err
	}

	opt := jpeg.Options{Quality: 100}
	err = jpeg.Encode(f, img, &opt)
	if err != nil {
		return err
	}

	return nil
}

func (narpimage *NARPImage) ConstructFromPngFile(s string, showprogress bool) error {
	reader, err := os.Open(s)
	if err != nil {
		log.Println(err)
		return err
	}
	defer reader.Close()

	img, err := png.Decode(reader)
	if err != nil {
		log.Println(err)
		return err
	}

	narpimage.initNARPImage()
	narpimage.putToNarpImage(img, showprogress)

	return nil
}

func (narpimage *NARPImage) ConstructFromJpgFile(s string, showprogress bool) error {
	reader, err := os.Open(s)
	if err != nil {
		log.Println(err)
		return err
	}
	defer reader.Close()

	img, err := jpeg.Decode(reader)
	if err != nil {
		log.Println(err)
		return err
	}

	narpimage.initNARPImage()
	narpimage.putToNarpImage(img, showprogress)

	return nil
}

func (narpimage *NARPImage) Load(s string) error {
	file, err := os.Open(s)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	narpimage.initNARPImage()

	err = gob.NewDecoder(file).Decode(narpimage)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

func (narpimage *NARPImage) Save(s string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(s); !os.IsNotExist(err) {
			log.Fatalf("Save: error, file <%s> already exists", s)
			return err
		}
	}

	file, err := os.Create(s)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	b := new(bytes.Buffer)
	err = gob.NewEncoder(b).Encode(narpimage)
	if err != nil {
		log.Println(err)
		return err
	}
	file.Write(b.Bytes())

	return err
}

func (narpimage *NARPImage) Print() {
	s := ""
	var visited [][]bool
	if narpimage.NARPixels == nil {
		s = "nil"
	}
	s = strconv.Itoa(len(narpimage.NARPixels))

	log.Printf("===========================================================================================")
	log.Printf("(NotARegularPixelImage):: size: %v, codec version: %v, pixels (#=%v):", narpimage.Size, narpimage.Version, s)
	log.Printf("===========================================================================================")

	if narpimage.Size.X == 0 || narpimage.Size.Y == 0 {
		return
	}

	x, y := 0, 0
	for _, v := range narpimage.NARPixels {
		v.Print(fmt.Sprintf("(x:%v,y:%v) ", x+1, y+1))
		v.markVisited(x, y, &visited, int(narpimage.Size.X), int(narpimage.Size.Y))
		end := false

		for !end && visited[x][y] {
			x++
			if x >= int(narpimage.Size.X) {
				x = 0
				y++
				if y >= int(narpimage.Size.Y) {
					end = true
				}
			}
		}
	}
}
