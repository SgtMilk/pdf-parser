package pdfparser

const countCutoffDivisor int = 4

type font struct {
	name  string
	size  float64
	width float64
}

/*
Finds and returns a list of fonts in the PDF document that are title fonts.
*/
func findTitles(texts []section) []font {
	fonts := findFonts(texts)

	characterCount := 0
	fontSum := 0.

	for f, v := range fonts {
		characterCount += v
		fontSum += float64(v) * f.size
	}

	avgFontSize := int(fontSum / float64(characterCount))
	countCutoff := characterCount / countCutoffDivisor

	isTitle := make(map[font]bool)
	for font, count := range fonts {
		isTitle[font] = count < countCutoff && font.size >= float64(avgFontSize)
	}

	var titleFonts []font

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
func findFonts(texts []section) map[font]int {
	mLength := make(map[font]int)
	mWidth := make(map[font]float64)

	for _, text := range texts {
		curFont := font{
			name: text.text.Font,
			size: text.text.FontSize,
		}
		if val, ok := mLength[curFont]; ok {
			mLength[curFont] = val + len(text.text.S)
			mWidth[curFont] += text.text.W * float64(len(text.text.S))
		} else {
			mLength[curFont] = len(text.text.S)
			mWidth[curFont] = text.text.W * float64(len(text.text.S))
		}
	}

	m := make(map[font]int)

	for k, v := range mLength {
		k.width = mWidth[k] / float64(v)
		m[k] = v
	}

	return m
}
