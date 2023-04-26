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

func ImageCreator(imageType ImageType, fileName *string) error {
	doc, err := fitz.New(*fileName + ".pdf")
	if err != nil {
		return fmt.Errorf("failed while creating new pdf to image doc with error : %w", err)
	}

	// Extract pages as images
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			return fmt.Errorf("failed while creating image out of pdf page with error : %w", err)
		}

		err = os.MkdirAll("img", 0755)
		if err != nil {
			return fmt.Errorf("failed while creating img directory with error : %w", err)
		}

		f, err := os.Create(filepath.Join("img/", fmt.Sprintf("image-%s.jpg", imageType)))
		if err != nil {
			return fmt.Errorf("failed while creating image file with error : %w", err)
		}

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 10})
		if err != nil {
			if err != nil {
				return fmt.Errorf("failed while jpeg encoding image with error : %w", err)
			}
		}

		f.Close()

	}
	return nil

}

func CleanUpImages() error {
	err := os.RemoveAll("img")
	if err != nil {
		return fmt.Errorf("failed while trying to clean up images in /img directory with error : %w", err)
	}
	return nil
}
