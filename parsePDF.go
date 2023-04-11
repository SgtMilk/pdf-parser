package main

import "encoding/json"

func ParsePDF(filename string) []byte{
	data, err := getPDF(filename)
	catch(err)

	texts, titleFonts := findTitles(data)
	parent := hierarchizeText(texts, titleFonts)
	file, err := json.Marshal(parent.Children)

	catch(err)
	return file
}