package pdfparser

const countCutoffDivisor int = 4

type Font struct {
	name  string
	size  float64
	width float64
}

/*
Finds and returns a list of fonts in the PDF document that are title fonts.
*/
func findTitles(texts []Section) []Font {
	fonts := findFonts(texts)

	characterCount := 0
	fontSum := 0.

	for f, v := range fonts {
		characterCount += v
		fontSum += float64(v) * f.size
	}

	avgFontSize := int(fontSum / float64(characterCount))
	countCutoff := characterCount / countCutoffDivisor

	isTitle := make(map[Font]bool)
	for font, count := range fonts {
		isTitle[font] = count < countCutoff && font.size >= float64(avgFontSize)
	}

	var titleFonts []Font

	for k, v := range isTitle {
		if v {
			titleFonts = append(titleFonts, k)
		}
	}

	return titleFonts
}

/*
From an array of texts, finds all the fonts included and stores them in a map.
Also stores in the map the number of characters for every font.
*/
func findFonts(texts []Section) map[Font]int {
	mLength := make(map[Font]int)
	mWidth := make(map[Font]float64)

	for _, text := range texts {
		font := Font{
			name: text.text.Font,
			size: text.text.FontSize,
		}
		if val, ok := mLength[font]; ok {
			mLength[font] = val + len(text.text.S)
			mWidth[font] += text.text.W * float64(len(text.text.S))
		} else {
			mLength[font] = len(text.text.S)
			mWidth[font] = text.text.W * float64(len(text.text.S))
		}
	}

	m := make(map[Font]int)

	for k, v := range mLength {
		k.width = mWidth[k] / float64(v)
		m[k] = v
	}

	return m
}
