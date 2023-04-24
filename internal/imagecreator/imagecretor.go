package imagecreator

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/karmdip-mi/go-fitz"
)

type ImageType string

const (
	Contract ImageType = "contract"
	Terms              = "terms"
)

func ImageCreator(imageType ImageType, fileName string) {
	doc, err := fitz.New(fileName + ".pdf")
	if err != nil {
		panic(err)
	}

	// Extract pages as images
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll("img", 0755)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join("img/", fmt.Sprintf("image-%s.jpg", imageType)))
		if err != nil {
			panic(err)
		}

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			panic(err)
		}

		f.Close()

	}

}

func CleanUpImages() {
	err := os.RemoveAll("img")
	if err != nil {
		fmt.Println(err)
	}
}
