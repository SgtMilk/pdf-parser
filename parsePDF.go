package main

import "encoding/json"

func ParsePDF(filename string) string{
	data, err := getPDF("./source/resume-en.pdf")
	catch(err)

	texts, titleFonts := findTitles(data)
	parent := hierarchizeText(texts, titleFonts)
	file, err := json.MarshalIndent(parent, "", " ")

	catch(err)
	return string(file)
}