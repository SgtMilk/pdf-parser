package pdfparser

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"
)

const minColumnEntries int = 4
const fontMultiplicatorForSameLine float64 = 3
const fontMultiplicatorForNextLine float64 = 8
const maxLengthForSubtitle = 40

type pageHeight struct {
	min float64
	max float64
}

type rect struct {
	Top    float64
	Bottom float64
	Left   float64
	Right  float64
}

type section struct {
	text     pdf.Text
	position rect
}

func getPDF(r *pdf.Reader) []section {
	sections := []section{}

	// fetching all the content once and getting the actual height of each page
	totalPage := r.NumPage()
	pageContents := make([][]pdf.Text, totalPage)
	pageHeights := make([]pageHeight, totalPage)

	for pageIndex := 0; pageIndex < totalPage; pageIndex++ {
		p := r.Page(pageIndex + 1)
		if p.V.IsNull() {
			continue
		}

		content := p.Content().Text
		pageContents[pageIndex] = content

		min := 100000.
		max := 0.

		for _, t := range content {
			if t.Y < min {
				min = t.Y
			}

			if t.Y > max {
				max = t.Y
			}
		}

		pageHeights[pageIndex] = pageHeight{min, max}
	}

	// parsing the content and separating in sections
	for page, texts := range pageContents {
		var lastTextStyle pdf.Text

		sum := 0.
		length := 0.

		var top, right, bottom, left float64 = 0, 0, 10000, 10000

		for _, text := range texts {
			if isSameSentence(lastTextStyle, text) {
				if text.W != 0 {
					sum += text.W
					length++
				}

				// updating coords
				if top < text.Y {
					top = text.Y
				}

				if bottom > text.Y {
					bottom = text.Y
				}

				if right < text.X {
					right = text.X
				}

				if left > text.X {
					left = text.X
				}

				lastTextStyle.S += text.S
			} else {
				lastTextStyle.W = sum / length

				sections = addString(lastTextStyle, sections, rect{
					Top:    calculateY(top, page, pageHeights),
					Bottom: calculateY(bottom, page, pageHeights),
					Right:  right,
					Left:   left,
				})
				lastTextStyle = text

				top, right, bottom, left = text.Y, text.X, text.Y, text.X
				if text.W != 0 {
					sum = text.W
					length = 1
				}
			}
		}

		sections = addString(lastTextStyle, sections, rect{
			Top:    calculateY(top, page, pageHeights),
			Bottom: calculateY(bottom, page, pageHeights),
			Right:  right,
			Left:   left,
		})
	}

	// sorting the sections because pdf format is trash
	sortSections(sections)

	return sections
}

func isSameSentence(prev, cur pdf.Text) bool {
	if prev.S == "" {
		return false
	}

	styleCheck := math.Abs(prev.FontSize-cur.FontSize) < 1 && prev.Font == cur.Font
	heightCheck := math.Abs(prev.Y-cur.Y) < prev.FontSize*fontMultiplicatorForNextLine

	return styleCheck && heightCheck
}

func addString(cur pdf.Text, sections []section, rect rect) []section {
	cur.S = cleanString(cur.S)
	if cur.S == "" {
		return sections
	}

	if isUpper(cur.S) {
		cur.Font += "-CAPS"
	}

	sections = append(sections, section{
		text:     cur,
		position: rect,
	})

	return sections
}

func cleanString(s string) string {
	re := regexp.MustCompile(`[^[:print:]À-ÿ]`)
	s = re.ReplaceAllLiteralString(s, "")
	re = regexp.MustCompile(`/s+`)
	s = re.ReplaceAllLiteralString(s, " ")
	s = strings.ReplaceAll(s, "\n", "")

	if strings.TrimSpace(s) == "" {
		return ""
	}

	if s != "" && s[0] == ' ' {
		return s[1:]
	}

	return s
}

func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func calculateY(y float64, curPage int, pageHeights []pageHeight) float64 {
	offset := 25
	newY := y - pageHeights[curPage].min + float64((offset * ((len(pageHeights)) - curPage)))

	for i := curPage + 1; i < len(pageHeights); i++ {
		newY += pageHeights[i].max - pageHeights[i].min
	}

	return newY
}

func sortSections(sections []section) {
	sort.SliceStable(sections, func(i, j int) bool {
		if math.Abs(sections[i].position.Top-sections[j].position.Top) < sections[i].text.FontSize/2 {
			return sections[i].position.Left < sections[j].position.Left
		}
		return sections[i].position.Top > sections[j].position.Top
	})

	fixColumns(sections)
	fixSubtitleBeforeTitle(sections)
}

func fixColumns(sections []section) {
	for i := 1; i < len(sections)-1; i++ {
		if !(math.Abs(sections[i].text.Y-sections[i+1].text.Y) < sections[i+1].text.FontSize/2) {
			continue
		}

		j := i + 1
		for j < len(sections) &&
			math.Abs(sections[j].position.Right-sections[j].position.Left) < math.Abs(sections[i+1].position.Right-sections[i].position.Left) &&
			math.Abs(sections[j].text.Y-sections[j-1].text.Y) < sections[j].text.FontSize*3 {
			j++
		}

		if j-i < minColumnEntries {
			continue
		}

		colArr := sections[i:j]

		sort.SliceStable(colArr, func(i, j int) bool {
			if math.Abs(colArr[i].text.X-colArr[j].text.X) < colArr[i].text.FontSize/2 {
				return colArr[i].text.Y > colArr[j].text.Y
			}
			return colArr[i].text.X < colArr[j].text.X
		})

		for k := i; k < j; k++ {
			sections[k] = colArr[k-i]
		}

		i = j
	}
}

func fixSubtitleBeforeTitle(sections []section) {
	for i := 1; i < len(sections); i++ {
		if subtitleIsSameParagraph(sections[i].text, sections[i-1].text) {
			sections[i], sections[i-1] = sections[i-1], sections[i]
			i++
		}
	}
}

func subtitleIsSameParagraph(cur, prev pdf.Text) bool {
	sizeCond := cur.FontSize-prev.FontSize > 1
	positionalCond := math.Abs(prev.Y-cur.Y) < prev.FontSize*fontMultiplicatorForSameLine
	lengthCond := len(prev.S) < maxLengthForSubtitle

	return positionalCond && lengthCond && sizeCond
}
