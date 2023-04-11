package main

import (
	"encoding/json"
	"io/ioutil"
)

func main() {
	data, err := getPDF("./source/resume-en.pdf")
	catch(err)

	texts, titleFonts := findTitles(data)
	parent := hierarchizeText(texts, titleFonts)
	file, _ := json.MarshalIndent(parent, "", " ")
 
	_ = ioutil.WriteFile("test.json", file, 0644)
}

func catch(e error) {
	if e != nil {
		panic(e)
	}
}