package pdfparser

import (
	"encoding/json"
	"io"
	"mime/multipart"

	"github.com/ledongthuc/pdf"
)

func ParsePdf(file *multipart.FileHeader) []byte {
	f, err := file.Open()
	catch(err)

	defer f.Close()
	r, err := pdf.NewReader(io.ReaderAt(f), file.Size)
	catch(err)

	return readFromReader(r)
}

func ParsePdfFile(filename string) []byte {
	f, r, err := pdf.Open(filename)
	catch(err)

	defer f.Close()

	return readFromReader(r)
}

func readFromReader(r *pdf.Reader) []byte {
	data := getPDF(r)
	titleFonts := findTitles(data)
	parent := hierarchizeText(data, titleFonts)
	file, err := json.Marshal(parent.Children)
	catch(err)

	return file
}
