package parser

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-fitz"
)

type Parser interface {
	ExtractImages() error
	Close() error
}
type FitzParser struct {
	Doc  *fitz.Document
	Name string
}

func NewFitzParser(path string) (*FitzParser, error) {
	doc, err := fitz.New(path)
	if err != nil {
		return nil, err
	}
	return &FitzParser{
		Doc:  doc,
		Name: path,
	}, nil
}

func (p *FitzParser) ExtractAndSaveImages() error {
	totalPages := p.Doc.NumPage()

	name := strings.TrimSuffix(filepath.Base(p.Name), filepath.Ext(p.Name))
	outputDir := fmt.Sprintf("/home/ahsansaif/projects/advanced-resume/resources/images/%s", name)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
	}

	for i := 0; i < totalPages; i++ {
		img, err := p.Doc.Image(i)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/page-%d.png", outputDir, i)
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}

		err = png.Encode(f, img)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *FitzParser) Close() error {
	return p.Doc.Close()
}
