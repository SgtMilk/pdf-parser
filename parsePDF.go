package pdfparser

import (
	"encoding/json"
	"io"
	"mime/multipart"

	"github.com/ledongthuc/pdf"
)

// An agglomeration of multiple attributes of the structure and style of the text
type TextNode struct {
	Value    string     // the text value, what is written
	Font     string     // the font of the text
	Children []TextNode // the children of this piece of text (ex: this is the title of ...), or nil if there are none
	Position Rect       // The positioning of this piece of text
}

// An agglomeration of positional attributes of a rectangle
type Rect struct {
	Top    float64
	Bottom float64
	Left   float64
	Right  float64
}

// ParsePdf parses a multipart text pdf file and returns text in JSON hierarchical order.
// It takes a multipart fileheader (perfect if you are working with a web server
// and don't want to save the file), and returns the output in JSON format, in a byte stream.
// The JSON is formatted from the TextNode struct, so it follows it's format
func ParsePdf(file *multipart.FileHeader) []byte {
	tree := ParsePdfToTree(file)
	return convertToJSON(tree)
}

// ParsePdfFile parses a text pdf file and returns text JSON in hierarchical order.
// It takes a file name, and returns the output in JSON format, in a byte stream.
// The JSON is formatted from the TextNode struct, so it follows it's format
func ParsePdfFile(filename string) []byte {
	tree := ParsePdfFileToTree(filename)
	return convertToJSON(tree)
}

// ParsePdfToTree parses a text pdf file and returns text in tree hierarchical order.
// It takes a multipart fileheader (perfect if you are working with a web server
// and don't want to save the file), and returns the output in tree format, from the TextNode struct.
func ParsePdfToTree(file *multipart.FileHeader) *TextNode {
	f, err := file.Open()
	catch(err)

	defer f.Close()
	r, err := pdf.NewReader(io.ReaderAt(f), file.Size)
	catch(err)

	return readFromReader(r)
}

// ParsePdfFileToTree parses a text pdf file and returns text in tree hierarchical order.
// It takes a file name, and returns the output in tree format, from the TextNode struct.
func ParsePdfFileToTree(filename string) *TextNode {
	f, r, err := pdf.Open(filename)
	catch(err)

	defer f.Close()

	return readFromReader(r)
}

func readFromReader(r *pdf.Reader) *TextNode {
	data := getPDF(r)
	titleFonts := findTitles(data)
	parent := hierarchizeText(data, titleFonts)

	return parent
}

func convertToJSON(parent *TextNode) []byte {
	file, err := json.Marshal(parent.Children)
	catch(err)

	return file
}
