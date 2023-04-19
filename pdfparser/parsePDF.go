package pdfparser

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"pdfparser/utils"

	"github.com/ledongthuc/pdf"
)

func ParsePdf(file *multipart.FileHeader) []byte {
	f, err := file.Open()
	utils.Catch(err)

	defer f.Close()
	r, err := pdf.NewReader(io.ReaderAt(f), file.Size)
	utils.Catch(err)

	return readFromReader(r)
}

func ParsePdfFile(filename string) []byte {
	f, r, err := pdf.Open(filename)
	utils.Catch(err)

	defer f.Close()

	return readFromReader(r)
}

func readFromReader(r *pdf.Reader) []byte {
	data := getPDF(r)
	titleFonts := findTitles(data)
	parent := hierarchizeText(data, titleFonts)
	file, err := json.Marshal(parent.Children)
	utils.Catch(err)

	return file
}
