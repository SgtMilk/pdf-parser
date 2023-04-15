package pdfparser

import (
	"encoding/json"
	"pdf-parser/utils"
)

func ParsePDF(filename string) []byte{
	data := getPDF(filename)

	titleFonts := findTitles(data)
	parent := hierarchizeText(data, titleFonts)
	file, err := json.Marshal(parent.Children)
	utils.Catch(err)
	
	return file
}