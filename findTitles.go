package main

import (
	"github.com/ledongthuc/pdf"
)

type TextType struct{
	text pdf.Text
	isTitle bool
}

type Font struct {
	name string
	size float64
}

func findTitles(texts []pdf.Text) ([]TextType, []Font) {
	fonts := findFonts(texts)

	characterCount := 0
	fontSum := 0.
	for f, v := range fonts{
		characterCount += v
		fontSum += float64(v) * f.size
	}
	avgFontSize := fontSum / float64(characterCount)
	countCutoff := characterCount / 4

	isTitle := make(map[Font]bool)
	for font, count := range fonts{
		isTitle[font] = count < countCutoff && font.size >= avgFontSize
	}

	stringCutoff := 50
	categorizedTexts := make([]TextType, len(texts))
	for i, text := range texts{
		font := Font{
			name: text.Font,
			size: text.FontSize,
		}
		textType := TextType{
			text: text,
			isTitle: isTitle[font] && len(text.S) < stringCutoff,
		}
		categorizedTexts[i] = textType


	}

    var titleFonts []Font
    for k, v := range isTitle {
        if(v){titleFonts = append(titleFonts, k)}
    }

	return categorizedTexts, titleFonts
}

func findFonts(texts []pdf.Text) map[Font]int {
	m := make(map[Font]int)

	for _, text := range texts{
		font := Font{
			name: text.Font,
			size: text.FontSize,
		}
		if val, ok := m[font]; ok {
			m[font] = val + len(text.S)
		}else{
			m[font] = len(text.S)
		}
	}
	return m
}
